package config

import "time"

type RabbitMQConfig struct {
	Addr                 string        `mapstructure:"addr"`
	Exchange             string        `mapstructure:"exchange"`
	PrefetchCount        int           `mapstructure:"prefetch_count"`
	ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
	MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
}
