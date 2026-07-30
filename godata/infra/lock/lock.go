package lock

import (
    "context"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
)

type Locker interface {
    Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Unlock(ctx context.Context, key string) error
}

type redisLocker struct {
    client *redis.Client
    mu     sync.Mutex
    tokens map[string]string
}

var unlockScript = redis.NewScript(`
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
`)

func NewRedisLocker(client *redis.Client) Locker {
    return &redisLocker{
        client: client,
        tokens: make(map[string]string),
    }
}

func (l *redisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    token := uuid.New().String()
    fullKey := "lock:" + key
    ok, err := l.client.SetNX(ctx, fullKey, token, ttl).Result()
    if err != nil {
        return false, err
    }
    if ok {
        l.mu.Lock()
        l.tokens[key] = token
        l.mu.Unlock()
    }
    return ok, nil
}

func (l *redisLocker) Unlock(ctx context.Context, key string) error {
    l.mu.Lock()
    token, exists := l.tokens[key]
    if exists {
        delete(l.tokens, key)
    }
    l.mu.Unlock()
    if !exists {
        return nil
    }
    fullKey := "lock:" + key
    return unlockScript.Run(ctx, l.client, []string{fullKey}, token).Err()
}
