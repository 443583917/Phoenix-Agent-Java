package logger

import (
	"sync"

	"github.com/phoenix-agent-go/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Init(cfg *config.MonitorConfig) error {
	var err error
	once.Do(func() {
		level := zapcore.InfoLevel
		if cfg != nil {
			_ = level.UnmarshalText([]byte(cfg.LogLevel))
		}

		var zapCfg zap.Config
		if cfg != nil && cfg.LogFormat == "json" {
			zapCfg = zap.NewProductionConfig()
		} else {
			zapCfg = zap.NewDevelopmentConfig()
		}
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, err = zapCfg.Build(zap.AddCallerSkip(1))
		if err == nil {
			zap.ReplaceGlobals(logger)
		}
	})
	return err
}

func L() *zap.Logger {
	if logger == nil {
		l, _ := zap.NewDevelopment()
		return l
	}
	return logger
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}
