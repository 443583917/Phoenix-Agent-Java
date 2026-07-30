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
	"github.com/phoenix-agent-go/agent/runner"
	"github.com/phoenix-agent-go/agent/runtime"
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
	"github.com/phoenix-agent-go/internal/usecase"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("api")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
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
	enforcer, err := casbin.NewEnforcer(modelPath)
	if err != nil {
		zap.L().Fatal("failed to init casbin enforcer", zap.Error(err))
	}
	defer enforcer.EnableEnforce(false)

	// 初始化数据库
	database := initDB(&cfg.DB)

	// 构建特权模块依赖链
	var privilegeSvc *service.PrivilegeService
	var platformSvc *service.PlatformService
	if database != nil {
		privilegeSvc = buildPrivilegeService(database, cfg)
		platformSvc = buildPlatformService(database)
	} else {
		zap.L().Warn("database unavailable, privilege and platform endpoints will return errors")
	}

	// 初始化 Agent 框架
	agentManager := initAgentManager(cfg)
	hitlHandler := runner.NewHitlHandler()

	// 设置路由
	router := api.SetupRouter(cfg, jwtManager, enforcer, privilegeSvc, platformSvc, agentManager, hitlHandler)

	// 启动服务
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Fatal("server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("server stopped")
}

// initAgentManager creates and configures the AgentManager with a default agent.
func initAgentManager(cfg *config.AppConfig) *runtime.AgentManager {
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

	zap.L().Info("agent manager initialized",
		zap.String("model", defaultAgent.ModelName),
		zap.String("base_url", defaultAgent.BaseURL),
		zap.Int("max_tokens", defaultAgent.MaxTokens),
		zap.Float64("temperature", defaultAgent.Temperature),
	)

	return runtime.NewAgentManager(registry)
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
// repos → usecase → service.
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

// buildPrivilegeService constructs the full privilege dependency chain:
// repos → cache → usecase → service.
func buildPrivilegeService(database *gorm.DB, cfg *config.AppConfig) *service.PrivilegeService {
	// Initialize cache (Redis + BigCache)
	var privilegeCache *cachePkg.PrivilegeCache
	if rdb, err := infraCache.InitRedis(&cfg.Redis); err == nil {
		local, err := infraCache.InitBigCache()
		if err != nil {
			zap.L().Warn("failed to init BigCache, using memory-only cache", zap.Error(err))
			local = nil
		}
		privilegeCache = cachePkg.NewPrivilegeCache(rdb, local)
		zap.L().Info("cache initialized (Redis + BigCache)")
	} else {
		zap.L().Warn("redis unavailable, caching disabled", zap.Error(err))
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
