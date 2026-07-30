package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/redis/go-redis/v9"
)

// PrivilegeCache is a two-level cache (L1: BigCache, L2: Redis) for privilege
// domain entities: users, roles, and user-role associations.
type PrivilegeCache struct {
	redis *redis.Client
	local *bigcache.BigCache
	ttl   time.Duration
}

// NewPrivilegeCache creates a new PrivilegeCache with the given Redis client and
// BigCache instance. The default TTL for Redis entries is 10 minutes.
func NewPrivilegeCache(redis *redis.Client, local *bigcache.BigCache) *PrivilegeCache {
	return &PrivilegeCache{redis: redis, local: local, ttl: 10 * time.Minute}
}

func (c *PrivilegeCache) key(prefix, id string) string {
	return "privilege:" + prefix + ":" + id
}

// ──────────────────────────── User ────────────────────────────

// GetUser retrieves a PrivilegeUser by ID, checking L1 (BigCache) first then
// L2 (Redis). A hit in L2 backfills L1.
func (c *PrivilegeCache) GetUser(ctx context.Context, id string) (*model.PrivilegeUser, error) {
	if data, err := c.local.Get(c.key("user", id)); err == nil {
		var user model.PrivilegeUser
		if err := json.Unmarshal(data, &user); err == nil {
			return &user, nil
		}
	}
	data, err := c.redis.Get(ctx, c.key("user", id)).Bytes()
	if err != nil {
		return nil, err
	}
	var user model.PrivilegeUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	_ = c.local.Set(c.key("user", id), data) // backfill L1
	return &user, nil
}

// SetUser stores a PrivilegeUser in both L1 and L2 caches.
func (c *PrivilegeCache) SetUser(ctx context.Context, user *model.PrivilegeUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if err := c.local.Set(c.key("user", user.ID), data); err != nil {
		return err
	}
	return c.redis.Set(ctx, c.key("user", user.ID), data, c.ttl).Err()
}

// InvalidateUser removes a user from both L1 and L2 caches.
func (c *PrivilegeCache) InvalidateUser(ctx context.Context, id string) error {
	if err := c.local.Delete(c.key("user", id)); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}
	return c.redis.Del(ctx, c.key("user", id)).Err()
}

// ──────────────────────────── Role ────────────────────────────

// GetRole retrieves a PrivilegeRole by ID, checking L1 first then L2.
func (c *PrivilegeCache) GetRole(ctx context.Context, id string) (*model.PrivilegeRole, error) {
	if data, err := c.local.Get(c.key("role", id)); err == nil {
		var role model.PrivilegeRole
		if err := json.Unmarshal(data, &role); err == nil {
			return &role, nil
		}
	}
	data, err := c.redis.Get(ctx, c.key("role", id)).Bytes()
	if err != nil {
		return nil, err
	}
	var role model.PrivilegeRole
	if err := json.Unmarshal(data, &role); err != nil {
		return nil, err
	}
	_ = c.local.Set(c.key("role", id), data)
	return &role, nil
}

// SetRole stores a PrivilegeRole in both L1 and L2 caches.
func (c *PrivilegeCache) SetRole(ctx context.Context, role *model.PrivilegeRole) error {
	data, err := json.Marshal(role)
	if err != nil {
		return err
	}
	if err := c.local.Set(c.key("role", role.ID), data); err != nil {
		return err
	}
	return c.redis.Set(ctx, c.key("role", role.ID), data, c.ttl).Err()
}

// InvalidateRole removes a role from both L1 and L2 caches.
func (c *PrivilegeCache) InvalidateRole(ctx context.Context, id string) error {
	if err := c.local.Delete(c.key("role", id)); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}
	return c.redis.Del(ctx, c.key("role", id)).Err()
}

// ──────────────────────────── User-Roles ────────────────────────────

// GetUserRoles retrieves the list of role IDs assigned to a user.
func (c *PrivilegeCache) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	if data, err := c.local.Get(c.key("user_roles", userID)); err == nil {
		var roleIDs []string
		if err := json.Unmarshal(data, &roleIDs); err == nil {
			return roleIDs, nil
		}
	}
	data, err := c.redis.Get(ctx, c.key("user_roles", userID)).Bytes()
	if err != nil {
		return nil, err
	}
	var roleIDs []string
	if err := json.Unmarshal(data, &roleIDs); err != nil {
		return nil, err
	}
	_ = c.local.Set(c.key("user_roles", userID), data)
	return roleIDs, nil
}

// SetUserRoles stores the list of role IDs for a user in both caches.
func (c *PrivilegeCache) SetUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	data, err := json.Marshal(roleIDs)
	if err != nil {
		return err
	}
	if err := c.local.Set(c.key("user_roles", userID), data); err != nil {
		return err
	}
	return c.redis.Set(ctx, c.key("user_roles", userID), data, c.ttl).Err()
}

// InvalidateUserRoles removes the cached role IDs for a user.
func (c *PrivilegeCache) InvalidateUserRoles(ctx context.Context, userID string) error {
	if err := c.local.Delete(c.key("user_roles", userID)); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}
	return c.redis.Del(ctx, c.key("user_roles", userID)).Err()
}
