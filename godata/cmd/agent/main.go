package main

import (
	"log"

	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/logger"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("agent")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix Agent standalone server",
		zap.Int("port", cfg.Server.Port),
		zap.String("version", cfg.Monitor.ServiceVersion),
	)

	// TODO: Initialize and start the agent service
	// This server runs agents independently without the full API stack.
	// It will connect to the gRPC service mesh for inter-service communication.

	zap.L().Info("agent server started")
	select {}
}
