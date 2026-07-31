package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/phoenix-agent-go/agent/hooks"
	"github.com/phoenix-agent-go/agent/interceptors"
	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/agent/memory"
	"github.com/phoenix-agent-go/agent/runner"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/agent/tools/datasource"
	"github.com/phoenix-agent-go/api"
	infraCache "github.com/phoenix-agent-go/infra/cache"
	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/jwt"
	internalConfig "github.com/phoenix-agent-go/internal/config"
	"github.com/phoenix-agent-go/infra/logger"
	"github.com/phoenix-agent-go/infra/monitoring"
	cachePkg "github.com/phoenix-agent-go/internal/dao/cache"
	"github.com/phoenix-agent-go/internal/dao/db"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/service/tracing"
	"github.com/phoenix-agent-go/internal/usecase"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"

	"github.com/phoenix-agent-go/internal/dao/queue"
	"github.com/phoenix-agent-go/internal/dao/vectorstore"
)

func main() {
	// 加载配置
	cfg, err := config.Load("api")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 兜底: viper mapstructure 可能无法解析 time.Duration 字符串
	if cfg.Auth.Expire == 0 {
		cfg.Auth.Expire = 24 * time.Hour
	}

	// 初始化日志
	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix API server",
		zap.Int("port", cfg.Server.Port),
		zap.String("version", cfg.Monitor.ServiceVersion),
	)

	// 初始化 RabbitMQ 消费者
	consumer, err := queue.NewConsumer(&cfg.RabbitMQ)
	if err != nil {
		zap.L().Warn("failed to init RabbitMQ consumer, continuing without event bus",
			zap.Error(err),
		)
	}

	// 初始化 OpenTelemetry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tp, err := monitoring.InitTracer(ctx, &cfg.Monitor)
	if err != nil {
		zap.L().Warn("failed to init tracer", zap.Error(err))
	}
	if tp != nil {
		defer func() {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			tp.Shutdown(shutdownCtx)
		}()
	}

	// 初始化 JWT 管理器
	jwtManager := jwt.NewJWTManager(cfg.Auth.Secret, cfg.Auth.Expire)

	// 初始化 Casbin 权限管理器
	modelPath := filepath.Join("internal", "config", "casbin_model.conf")
	policyPath := filepath.Join("internal", "config", "casbin_policy.csv")
	
	var enforcer *casbin.Enforcer
	if _, err := os.Stat(policyPath); err == nil {
		enforcer, err = casbin.NewEnforcer(modelPath, policyPath)
		if err != nil {
			zap.L().Fatal("failed to init casbin enforcer with policy", zap.Error(err))
		}
		enforcer.EnableEnforce(true)
		zap.L().Info("casbin enforcement enabled")
	} else {
		enforcer, err = casbin.NewEnforcer(modelPath)
		if err != nil {
			zap.L().Fatal("failed to init casbin enforcer", zap.Error(err))
		}
		enforcer.EnableEnforce(false)
			zap.L().Info("no casbin policy file found, enforcement disabled")
	}
	defer enforcer.EnableEnforce(false)

	// 初始化数据库
	database := initDB(&cfg.DB)

	// 初始化 Redis (shared client for privilege cache, HITL cache, login interceptor)
	var redisClient *redis.Client
	if rdb, err := infraCache.InitRedis(&cfg.Redis); err == nil {
		redisClient = rdb
		zap.L().Info("redis connected")
	} else {
		zap.L().Warn("redis unavailable, some features will be degraded", zap.Error(err))
	}

	// 构建依赖链
	var privilegeSvc *service.PrivilegeService
	var platformSvc *service.PlatformService
	var dataSvc *service.DataService
	var ragSvc *service.RagService
	var kgSvc *service.KgService
	if database != nil {
		privilegeSvc = buildPrivilegeService(database, cfg, redisClient)
		platformSvc = buildPlatformService(database)
		dataSvc = buildDataService(database)
		ragSvc = buildRagService(database)
		kgSvc = buildKgService(database)

		// Wire cache invalidation for privilege role updates into the consumer.
		if consumer != nil {
			consumer.SetCacheInvalidator(func(ctx context.Context, eventType string, payload []byte) {
				zap.L().Info("consuming cache invalidation",
					zap.String("eventType", eventType),
				)
				// The privilege service exposes InvalidateCache for this purpose.
				// For now we rely on the PrivilegeCache invalidation in the service layer;
				// future: call privilegeSvc.InvalidateRoleCache() when that method exists.
			})
		}
	} else {
		zap.L().Warn("database unavailable, privilege and platform endpoints will return errors")
	}

	// 初始化 Agent 框架 (session service, agent manager, HITL handler, hooks, interceptors)
	sessSvc := inmemory.NewSessionService()
	agentManager := initAgentManager(cfg, sessSvc, database, redisClient)
	hitlHandler := runner.NewHitlHandler()

	var llmModel tmodel.Model
	if cfg.Agent.Model.APIKey != "" {
		keySet := cfg.Agent.Model.APIKey[:8] + "..."
		llmModel = openai.New(
			cfg.Agent.Model.Model,
			openai.WithBaseURL(cfg.Agent.Model.BaseURL),
			openai.WithAPIKey(cfg.Agent.Model.APIKey),
		)
		zap.L().Info("LLM model initialized",
			zap.String("model", cfg.Agent.Model.Model),
			zap.String("apiKey", keySet))
	}

	dbManager := datasource.NewDatasourceManager()
	mtcm := runtime.NewMultiTurnContextManager(5)
	tracingSvc := tracing.NewTracingService()

	var retriever *knowledge.Retriever

	var embeddingSvc *service.EmbeddingService
	if database != nil {
		agentKnowledgeRepo := db.NewAgentKnowledgeRepository(database)
		embeddingSvc = service.NewEmbeddingService(agentKnowledgeRepo)
	}

	if cfg.Milvus.Addr != "" {
		milvusStore, storeErr := vectorstore.NewMilvusStore(cfg.Milvus.Addr, internalConfig.DefaultVectorStoreConfig())
		if storeErr != nil {
			zap.L().Warn("milvus connect failed, vector search will be disabled", zap.Error(storeErr))
		} else {
			zap.L().Info("milvus connected", zap.String("addr", cfg.Milvus.Addr))
			retriever = knowledge.NewRetrieverWithVectorStore(&milvusVectorStoreAdapter{store: milvusStore})
			defer milvusStore.Close()
		}
	}

	// 设置路由
	router := api.SetupRouter(cfg, jwtManager, enforcer, privilegeSvc, platformSvc, dataSvc, ragSvc, kgSvc, agentManager, hitlHandler, redisClient, llmModel, dbManager, retriever, mtcm, tracingSvc, embeddingSvc)

	// 启动服务
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// 启动消息消费者
	if consumer != nil {
		go func() {
			zap.L().Info("starting RabbitMQ consumer")
			if err := consumer.Start(context.Background()); err != nil {
				zap.L().Error("consumer stopped with error", zap.Error(err))
			}
		}()
	}

	// 优雅关闭
	go func() {
		zap.L().Info("server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if consumer != nil {
		zap.L().Info("stopping RabbitMQ consumer")
		if err := consumer.Stop(); err != nil {
			zap.L().Error("failed to stop consumer", zap.Error(err))
		}
	}

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Fatal("server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("server stopped")
}

// initAgentManager creates and configures the AgentManager with hooks,
// interceptors, and the async memory pipeline.
func initAgentManager(
	cfg *config.AppConfig,
	sessSvc *inmemory.SessionService,
	database *gorm.DB,
	redisClient *redis.Client,
) *runtime.AgentManager {
	registry := runtime.NewAgentRegistry()

	// Register a default agent from config.
	defaultAgent := &runtime.AgentConfig{
		SN:          "default",
		ModelName:   cfg.Agent.Model.Model,
		BaseURL:     cfg.Agent.Model.BaseURL,
		APIKey:      cfg.Agent.Model.APIKey,
		Stream:      cfg.Agent.Stream,
		MaxTokens:   cfg.Agent.MaxTokens,
		Temperature: cfg.Agent.Temperature,
		AgentType:   "react",
	}
	registry.Register(defaultAgent)

	// Build manager options from available components.
	var managerOpts []runtime.AgentManagerOption

	if database != nil {
		// Create repositories for agent memory infrastructure.
		profileRepo := db.NewUserProfileRepository(database)
		memoryRepo := db.NewUserMemoryRepository(database)
		userAgentInfoRepo := db.NewUserAgentInfoRepository(database)

		// Create the extraction model for memory pipeline and summarization hook.
		extractionModel := openai.New(
			cfg.Agent.Model.Model,
			openai.WithVariant(openai.VariantDeepSeek),
			openai.WithBaseURL(cfg.Agent.Model.BaseURL),
			openai.WithAPIKey(cfg.Agent.Model.APIKey),
		)

		// Create long-term memory (vector store) for both memory pipeline and login interceptor.
		longTermMemory := memory.NewLongTermMemory("phoenix")

		// Wire profile injection hook.
		profileHook := hooks.NewProfileInjectionHook(profileRepo)
		managerOpts = append(managerOpts, runtime.WithProfileHook(profileHook))

		// Wire model call limit hook (default 15 calls).
		limitHook := hooks.NewModelCallLimitHook(15)
		managerOpts = append(managerOpts, runtime.WithLimitHook(limitHook))

		// Wire summarization hook (threshold 20 messages).
		summarizationHook := hooks.NewSummarizationHook(extractionModel, 20)
		managerOpts = append(managerOpts, runtime.WithSummarizationHook(summarizationHook))

		// Wire login interceptor (Redis dedup, async usage recording, history memory injection).
		loginInterceptor := interceptors.NewLoginUserAgentInterceptor(
			redisClient, userAgentInfoRepo, longTermMemory,
		)
		managerOpts = append(managerOpts, runtime.WithLoginInterceptor(loginInterceptor))

		// Wire memory pipeline for async memory extraction after each turn.
		memoryPipeline := memory.NewMemoryPipeline(
			extractionModel, profileRepo, memoryRepo, longTermMemory,
		)
		managerOpts = append(managerOpts, runtime.WithMemoryPipeline(memoryPipeline))

		zap.L().Info("agent infrastructure initialized",
			zap.Bool("profileHook", true),
			zap.Bool("limitHook", true),
			zap.Bool("summarizationHook", true),
			zap.Bool("loginInterceptor", true),
			zap.Bool("memoryPipeline", true),
		)
	} else {
		zap.L().Warn("database unavailable, agent infrastructure hooks disabled")
	}

	// Register a default agent from config — already done above, but log the model info.
	zap.L().Info("agent manager initialized",
		zap.String("model", defaultAgent.ModelName),
		zap.String("base_url", defaultAgent.BaseURL),
		zap.Int("max_tokens", defaultAgent.MaxTokens),
		zap.Float64("temperature", defaultAgent.Temperature),
	)

	return runtime.NewAgentManager(registry, sessSvc, managerOpts...)
}

// initDB attempts to connect to PostgreSQL. If the connection fails, it logs
// a warning and returns nil so the server can still start (e.g. for /echo).
func initDB(cfg *internalConfig.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		zap.L().Warn("database connection failed, privilege module will be unavailable", zap.Error(err))
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Warn("failed to get underlying sql.DB", zap.Error(err))
		return nil
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	zap.L().Info("database connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("dbname", cfg.Name),
	)
	return db
}

// buildPlatformService constructs the full platform dependency chain:
// repos -> usecase -> service.
func buildPlatformService(database *gorm.DB) *service.PlatformService {
	groupInfoRepo := db.NewGroupInfoRepository(database)
	groupAgentInfoRepo := db.NewGroupAgentInfoRepository(database)
	accountInfoRepo := db.NewAccountInfoRepository(database)
	accountGroupInfoRepo := db.NewAccountGroupInfoRepository(database)
	accountTenantInfoRepo := db.NewAccountTenantInfoRepository(database)
	tenantInfoRepo := db.NewTenantInfoRepository(database)
	platformInfoRepo := db.NewPlatformInfoRepository(database)

	uc := usecase.NewPlatformUsecase(
		groupInfoRepo, groupAgentInfoRepo, accountInfoRepo,
		accountGroupInfoRepo, accountTenantInfoRepo, tenantInfoRepo, platformInfoRepo,
	)

	return service.NewPlatformService(uc)
}

// buildDataService constructs the full data management dependency chain:
// repos -> usecase -> service.
func buildDataService(database *gorm.DB) *service.DataService {
	agentRepo := db.NewAgentRepository(database)
	agentCategoryRepo := db.NewAgentCategoryRepository(database)
	agentDatasourceRepo := db.NewAgentDatasourceRepository(database)
	agentKnowledgeRepo := db.NewAgentKnowledgeRepository(database)
	agentPresetQuestionRepo := db.NewAgentPresetQuestionRepository(database)
	agentDatasourceTablesRepo := db.NewAgentDatasourceTablesRepository(database)
	chatSessionRepo := db.NewChatSessionRepository(database)
	chatMessageRepo := db.NewChatMessageRepository(database)
	datasourceRepo := db.NewDatasourceRepository(database)
	logicalRelationRepo := db.NewLogicalRelationRepository(database)
	modelConfigRepo := db.NewModelConfigRepository(database)
	userPromptConfigRepo := db.NewUserPromptConfigRepository(database)
	semanticModelRepo := db.NewSemanticModelRepository(database)
	businessKnowledgeRepo := db.NewBusinessKnowledgeRepository(database)
	dsAccessor := db.NewDatasourceAccessor()

	uc := usecase.NewDataUsecase(
		agentRepo, agentCategoryRepo, agentDatasourceRepo,
		agentKnowledgeRepo, agentPresetQuestionRepo, agentDatasourceTablesRepo,
		chatSessionRepo, chatMessageRepo, datasourceRepo,
		logicalRelationRepo, modelConfigRepo, userPromptConfigRepo,
		semanticModelRepo, businessKnowledgeRepo, dsAccessor,
	)

	return service.NewDataService(uc)
}

// buildRagService constructs the RAG dependency chain:
// repos -> usecase -> service.
func buildRagService(database *gorm.DB) *service.RagService {
	categoryRepo := db.NewRagCategoryRepository(database)
	fileRepo := db.NewRagFileInfoRepository(database)

	uc := usecase.NewRagUsecase(categoryRepo, fileRepo)

	return service.NewRagService(uc)
}

// buildKgService constructs the KG dependency chain:
// repos -> usecase -> service.
func buildKgService(database *gorm.DB) *service.KgService {
	entityRepo := db.NewKGEntityRepository(database)
	relationRepo := db.NewKGRelationRepository(database)
	domainRepo := db.NewKGDomainRepository(database)

	uc := usecase.NewKgUsecase(entityRepo, relationRepo, domainRepo)

	return service.NewKgService(uc)
}

// buildPrivilegeService constructs the full privilege dependency chain:
// repos -> cache -> usecase -> service.
func buildPrivilegeService(database *gorm.DB, cfg *config.AppConfig, redisClient *redis.Client) *service.PrivilegeService {
	// Initialize cache (Redis + BigCache)
	var privilegeCache *cachePkg.PrivilegeCache
	if redisClient != nil {
		local, err := infraCache.InitBigCache()
		if err != nil {
			zap.L().Warn("failed to init BigCache, using memory-only cache", zap.Error(err))
			local = nil
		}
		privilegeCache = cachePkg.NewPrivilegeCache(redisClient, local)
		zap.L().Info("cache initialized (Redis + BigCache)")
	} else {
		zap.L().Warn("redis unavailable, caching disabled")
		// Still create a cache without Redis/BigCache — operations degrade gracefully.
		local, _ := infraCache.InitBigCache()
		if local != nil {
			privilegeCache = cachePkg.NewPrivilegeCache(nil, local)
		}
	}

	// Initialize all repositories
	userRepo := db.NewUserRepository(database)
	roleRepo := db.NewRoleRepository(database)
	userRoleRepo := db.NewUserRoleRepository(database)
	moduleRepo := db.NewModuleRepository(database)
	aclRepo := db.NewACLRepository(database)
	deptRepo := db.NewDepartmentRepository(database)
	companyRepo := db.NewCompanyRepository(database)
	employeeRepo := db.NewEmployeeRepository(database)
	dictRepo := db.NewDictionaryRepository(database)
	pvalueRepo := db.NewPvalueRepository(database)
	loginLogRepo := db.NewLoginLogRepository(database)

	// Initialize usecase
	uc := usecase.NewPrivilegeUsecase(
		userRepo, roleRepo, userRoleRepo, moduleRepo, aclRepo,
		deptRepo, companyRepo, employeeRepo, dictRepo, pvalueRepo, loginLogRepo,
		privilegeCache,
	)

	// Initialize service
	return service.NewPrivilegeService(uc)
}

type milvusVectorStoreAdapter struct {
	store *vectorstore.MilvusStore
}

func (a *milvusVectorStoreAdapter) Search(ctx context.Context, query string, embedding []float64, topK int) ([]knowledge.Document, error) {
	results, err := a.store.Search(ctx, query, embedding, topK)
	if err != nil {
		return nil, err
	}
	docs := make([]knowledge.Document, len(results))
	for i, r := range results {
		docs[i] = knowledge.Document{
			ID:      r.ID,
			Content: r.Content,
			Score:   r.Score,
		}
	}
	return docs, nil
}
