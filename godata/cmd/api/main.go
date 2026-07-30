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
	"github.com/phoenix-agent-go/api"
	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/logger"
	"github.com/phoenix-agent-go/infra/monitoring"
	"go.uber.org/zap"
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

	// 设置路由
	router := api.SetupRouter(cfg, jwtManager, enforcer)

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
