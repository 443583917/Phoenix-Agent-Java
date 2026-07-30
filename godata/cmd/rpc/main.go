package main

import (
	"fmt"
	"log"
	"net"

	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 加载配置
	cfg, err := config.Load("rpc")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix gRPC server",
		zap.Int("port", cfg.Server.Port),
	)

	// 创建 gRPC 服务器
	srv := grpc.NewServer()
	reflection.Register(srv)

	// TODO: Register gRPC services here when defined

	// 启动监听
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
