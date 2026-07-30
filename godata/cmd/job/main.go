package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/logger"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("job")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix job server",
		zap.String("version", cfg.Monitor.ServiceVersion),
	)

	// TODO: Register cron jobs here
	// e.g.:
	// - Vector embedding refresh
	// - Knowledge graph sync
	// - Data cleanup
	// - Report generation

	zap.L().Info("job server started, waiting for signals")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("job server shutting down")
}
