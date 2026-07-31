package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appcfg "github.com/phoenix-agent-go/internal/config"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Server   appcfg.ServerConfig
	DB       appcfg.DBConfig
	Redis    appcfg.RedisConfig
	Milvus   appcfg.MilvusConfig
	RabbitMQ appcfg.RabbitMQConfig
	Monitor  appcfg.MonitorConfig
	Agent    appcfg.AgentConfig
	RPC      appcfg.RPCConfig
	Cors     appcfg.CorsConfig
	Auth     appcfg.AuthConfig
	Graph    appcfg.GraphConfig
}

func Load(serviceName string) (*AppConfig, error) {
	v := viper.New()

	// 配置搜索路径
	v.SetConfigName("db")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")

	// 环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("PHOENIX")

	if err := loadConfigFile(v, "db"); err != nil {
		return nil, err
	}
	if err := loadConfigFile(v, "redis"); err != nil {
		return nil, err
	}
	if err := loadConfigFile(v, "milvus"); err != nil {
		return nil, err
	}
	if err := loadConfigFile(v, "rabbitmq"); err != nil {
		return nil, err
	}
	if err := loadConfigFile(v, "monitor"); err != nil {
		return nil, err
	}

	// 服务专用配置
	// 尝试多个相对路径（测试和工作目录可能不同）
	for _, base := range []string{"../../configs", "../configs", "configs"} {
		serviceConfigPath := filepath.Join(base, serviceName, "app.yaml")
		if _, err := os.Stat(serviceConfigPath); err == nil {
			v.SetConfigFile(serviceConfigPath)
			if err := v.MergeInConfig(); err != nil {
				return nil, fmt.Errorf("failed to load service config %s: %w", serviceConfigPath, err)
			}
			break
		}
	}

	// 展开 ${VAR} 环境变量
	// 注意：不能用 v.Set("a.b.c", val) 写扁平 key，那样会在 override 中创建
	// 不完整的嵌套结构，导致后续 UnmarshalKey("a") 时覆盖掉 yaml 中的其他字段。
	// 改为按顶层 key 取出完整嵌套结构、递归展开后再整体 Set。
	seen := map[string]bool{}
	for _, key := range v.AllKeys() {
		top := key
		if idx := strings.Index(key, "."); idx > 0 {
			top = key[:idx]
		}
		if seen[top] {
			continue
		}
		seen[top] = true
		v.Set(top, expandEnvVars(v.Get(top)))
	}

	cfg := &AppConfig{}
	if err := v.UnmarshalKey("database", &cfg.DB); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("milvus", &cfg.Milvus); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("rabbitmq", &cfg.RabbitMQ); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("monitor", &cfg.Monitor); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("server", &cfg.Server); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("cors", &cfg.Cors); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("agent", &cfg.Agent); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("rpc", &cfg.RPC); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return nil, err
	}
	if err := v.UnmarshalKey("graph", &cfg.Graph); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadConfigFile(v *viper.Viper, name string) error {
	v.SetConfigName(name)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to load %s.yaml: %w", name, err)
		}
		// 配置文件不存在可接受
	}
	return nil
}

// expandEnvVars 递归展开配置值中的 ${VAR} 环境变量占位符。
func expandEnvVars(val interface{}) interface{} {
	switch t := val.(type) {
	case string:
		if strings.HasPrefix(t, "${") && strings.HasSuffix(t, "}") {
			if vv, ok := os.LookupEnv(t[2 : len(t)-1]); ok {
				return vv
			}
		}
		return t
	case map[string]interface{}:
		m := make(map[string]interface{}, len(t))
		for k, v := range t {
			m[k] = expandEnvVars(v)
		}
		return m
	case []interface{}:
		s := make([]interface{}, len(t))
		for i, v := range t {
			s[i] = expandEnvVars(v)
		}
		return s
	default:
		return val
	}
}
