package main

import (
	"fmt"
	"log"
	"net"

	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/logger"
	internalConfig "github.com/phoenix-agent-go/internal/config"
	"github.com/phoenix-agent-go/internal/dao/db"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
	pb "github.com/phoenix-agent-go/rpc/proto"
	"github.com/phoenix-agent-go/rpc/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("rpc")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix gRPC server",
		zap.Int("port", cfg.Server.Port),
	)

	database := initDB(&cfg.DB)

	var privilegeSvc *service.PrivilegeService
	var dataSvc *service.DataService
	if database != nil {
		privilegeSvc = buildPrivilegeService(database)
		dataSvc = buildDataService(database)
	} else {
		zap.L().Warn("database unavailable, gRPC services will not be registered")
	}

	srv := grpc.NewServer()
	reflection.Register(srv)

	if privilegeSvc != nil {
		pb.RegisterPrivilegeServiceServer(srv, server.NewPrivilegeServer(privilegeSvc))
	}
	if dataSvc != nil {
		pb.RegisterAgentServiceServer(srv, server.NewAgentServer(dataSvc))
		pb.RegisterDataServiceServer(srv, server.NewDataServer(dataSvc))
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Fatal("failed to listen", zap.String("addr", addr), zap.Error(err))
	}

	zap.L().Info("gRPC server listening", zap.String("addr", addr))
	if err := srv.Serve(lis); err != nil {
		zap.L().Fatal("gRPC server error", zap.Error(err))
	}
}

func initDB(cfg *internalConfig.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		zap.L().Warn("database connection failed", zap.Error(err))
		return nil
	}
	return conn
}

func buildPrivilegeService(database *gorm.DB) *service.PrivilegeService {
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

	uc := usecase.NewPrivilegeUsecase(
		userRepo, roleRepo, userRoleRepo, moduleRepo, aclRepo,
		deptRepo, companyRepo, employeeRepo, dictRepo, pvalueRepo, loginLogRepo, nil,
	)

	return service.NewPrivilegeService(uc)
}

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
