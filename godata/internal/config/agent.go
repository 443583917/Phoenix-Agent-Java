package config

import "time"

type AgentModelConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"`
}

type AgentConfig struct {
	Model          AgentModelConfig `mapstructure:"model"`
	Stream         bool             `mapstructure:"stream"`
	MaxTokens      int              `mapstructure:"max_tokens"`
	Temperature    float64          `mapstructure:"temperature"`
	SessionTimeout time.Duration    `mapstructure:"session_timeout"`
}
