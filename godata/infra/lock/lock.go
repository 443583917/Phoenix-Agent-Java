package lock

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type Locker interface {
    Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Unlock(ctx context.Context, key string) error
}

type redisLocker struct {
    client *redis.Client
}

func NewRedisLocker(client *redis.Client) Locker {
    return &redisLocker{client: client}
}

func (l *redisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return l.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

func (l *redisLocker) Unlock(ctx context.Context, key string) error {
    return l.client.Del(ctx, "lock:"+key).Err()
}
