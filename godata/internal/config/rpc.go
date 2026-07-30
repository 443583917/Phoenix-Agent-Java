package config

type RPCConfig struct {
	Port       int  `mapstructure:"port"`
	Reflection bool `mapstructure:"reflection"`
}
