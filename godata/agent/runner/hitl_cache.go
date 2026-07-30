package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	// hitlKeyPrefix is the Redis key prefix for HITL pending confirmations.
	// Key format: phoenix:agentscope:session:{sessionID}:pending_confirm
	hitlKeyPrefix = "phoenix:agentscope:session:"

	// hitlKeySuffix is appended to the session ID to form the full key.
	hitlKeySuffix = ":pending_confirm"

	// hitlTTL is the default TTL for HITL confirmation entries.
	hitlTTL = 2 * time.Hour
)

// HitlCacheService provides Redis-backed storage for Human-In-The-Loop
// tool confirmation events. It replaces the in-memory channel-based
// HitlHandler with a durable, distributed cache that survives process
// restarts and works across multiple server instances.
type HitlCacheService struct {
	redis *redis.Client
}

// NewHitlCacheService creates a new HitlCacheService backed by Redis.
// redis may be nil; if so, operations are no-ops that return errors.
func NewHitlCacheService(redis *redis.Client) *HitlCacheService {
	return &HitlCacheService{redis: redis}
}

// SavePendingConfirm stores a ToolCallEvent in Redis as a pending
// confirmation entry. The entry is keyed by sessionID and has a 2-hour TTL.
func (c *HitlCacheService) SavePendingConfirm(
	ctx context.Context,
	sessionID string,
	event model.ToolCallEvent,
) error {
	if c.redis == nil {
		return fmt.Errorf("hitl cache: redis not configured")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("hitl cache: marshal event: %w", err)
	}

	key := c.buildKey(sessionID)
	if err := c.redis.Set(ctx, key, string(data), hitlTTL).Err(); err != nil {
		return fmt.Errorf("hitl cache: save pending confirm: %w", err)
	}

	return nil
}

// GetAndRemovePendingConfirm retrieves and atomically removes a pending
// confirmation entry from Redis. Returns nil if no entry exists for the
// given session ID.
func (c *HitlCacheService) GetAndRemovePendingConfirm(
	ctx context.Context,
	sessionID string,
) (*model.ToolCallEvent, error) {
	if c.redis == nil {
		return nil, fmt.Errorf("hitl cache: redis not configured")
	}

	key := c.buildKey(sessionID)

	// Get the value and delete the key atomically via pipeline.
	pipe := c.redis.Pipeline()
	getCmd := pipe.Get(ctx, key)
	pipe.Del(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("hitl cache: get and remove: %w", err)
	}

	val, err := getCmd.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("hitl cache: get pending confirm: %w", err)
	}

	var event model.ToolCallEvent
	if err := json.Unmarshal([]byte(val), &event); err != nil {
		return nil, fmt.Errorf("hitl cache: unmarshal event: %w", err)
	}

	return &event, nil
}

// HasPendingConfirm checks whether a pending confirmation exists for the
// given session ID without removing it.
func (c *HitlCacheService) HasPendingConfirm(
	ctx context.Context,
	sessionID string,
) (bool, error) {
	if c.redis == nil {
		return false, nil
	}

	key := c.buildKey(sessionID)
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("hitl cache: check pending: %w", err)
	}

	return exists > 0, nil
}

// RemovePendingConfirm deletes a pending confirmation entry without
// retrieving it first.
func (c *HitlCacheService) RemovePendingConfirm(
	ctx context.Context,
	sessionID string,
) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildKey(sessionID)
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("hitl cache: remove pending: %w", err)
	}

	return nil
}

// buildKey constructs the Redis key for a pending confirmation entry.
func (c *HitlCacheService) buildKey(sessionID string) string {
	return hitlKeyPrefix + sessionID + hitlKeySuffix
}
