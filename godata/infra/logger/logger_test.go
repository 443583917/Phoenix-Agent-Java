package logger

import (
	"testing"

	"github.com/phoenix-agent-go/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	cfg := &config.MonitorConfig{
		LogLevel:  "debug",
		LogFormat: "console",
	}
	err := Init(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, L())

	L().Info("test log message",
		zap.String("key", "value"),
	)
	Sync()
}
