package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := Load("api")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 8066, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.DB.Host)
}
