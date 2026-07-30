package config

import "time"

type MilvusConfig struct {
	Addr       string        `mapstructure:"addr"`
	Collection string        `mapstructure:"collection"`
	Dim        int           `mapstructure:"dim"`
	IndexType  string        `mapstructure:"index_type"`
	MetricType string        `mapstructure:"metric_type"`
	Timeout    time.Duration `mapstructure:"timeout"`
}
