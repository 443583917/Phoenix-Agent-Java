package config

type MonitorConfig struct {
	OTELEndpoint    string  `mapstructure:"otel_endpoint"`
	ServiceName     string  `mapstructure:"service_name"`
	ServiceVersion  string  `mapstructure:"service_version"`
	TraceSampleRate float64 `mapstructure:"trace_sample_rate"`
	LogLevel        string  `mapstructure:"log_level"`
	LogFormat       string  `mapstructure:"log_format"`
}
