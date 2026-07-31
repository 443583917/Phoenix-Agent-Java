package config

type CodeExecutorConfig struct {
	Type             string `mapstructure:"type"`
	ImageName        string `mapstructure:"image_name"`
	ContainerPrefix  string `mapstructure:"container_prefix"`
	LimitMemoryMB    int    `mapstructure:"limit_memory_mb"`
	CPUCores         int    `mapstructure:"cpu_cores"`
	CodeTimeout      int    `mapstructure:"code_timeout_seconds"`
	ContainerTimeout int    `mapstructure:"container_timeout_seconds"`
	NetworkMode      string `mapstructure:"network_mode"`
}

func DefaultCodeExecutorConfig() CodeExecutorConfig {
	return CodeExecutorConfig{
		Type:             "simulation",
		ImageName:        "python:3.11-slim",
		ContainerPrefix:  "phoenix-python-exec-",
		LimitMemoryMB:    500,
		CPUCores:         1,
		CodeTimeout:      60,
		ContainerTimeout: 300,
		NetworkMode:      "none",
	}
}
