package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/logger"
	"github.com/phoenix-agent-go/internal/job"
	"github.com/phoenix-agent-go/internal/job/jobs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	internalConfig "github.com/phoenix-agent-go/internal/config"
)

func main() {
	cfg, err := config.Load("job")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := logger.Init(&cfg.Monitor); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	zap.L().Info("starting Phoenix job server",
		zap.String("version", cfg.Monitor.ServiceVersion),
	)

	database := initDB(&cfg.DB)
	if database == nil {
		zap.L().Fatal("database connection required for job server")
	}

	scheduler := job.NewScheduler()

	scheduler.AddJob(jobs.NewAgentStatisticsJob(database), 1*time.Hour)
	scheduler.AddJob(jobs.NewKnowledgeEmbeddingRetryJob(database), 30*time.Minute)
	scheduler.AddJob(jobs.NewSessionCleanupJob(database), 24*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	zap.L().Info("job server started, waiting for signals")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("job server shutting down")
	scheduler.Stop()
	cancel()
}

func initDB(cfg *internalConfig.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		zap.L().Fatal("database connection failed", zap.Error(err))
		return nil
	}
	return conn
}
