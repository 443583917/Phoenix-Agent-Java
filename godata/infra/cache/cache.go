package cache

import (
    "context"
    "time"

    "github.com/allegro/bigcache/v3"
    "github.com/phoenix-agent-go/internal/config"
    "github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return rdb, nil
}

func InitBigCache() (*bigcache.BigCache, error) {
    return bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
}
