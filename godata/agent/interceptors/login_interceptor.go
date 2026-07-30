package interceptors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/phoenix-agent-go/agent/memory"
	"github.com/phoenix-agent-go/internal/repository"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

const (
	// redisDedupKeyPrefix is the Redis key prefix for per-(user, session) dedup.
	redisDedupKeyPrefix = "phoenix:agent:dedup:"

	// dedupWindow is the deduplication window duration.
	dedupWindow = 1 * time.Minute
)

// LoginUserAgentInterceptor is a pre-processing interceptor that runs before
// each agent conversation turn. It performs:
//
//  1. Redis-based deduplication per (userID, sessionID) within a 1-minute window.
//  2. Asynchronous recording of agent usage (fire-and-forget).
//  3. History memory injection: searches the vector store for relevant memories
//     and injects them as a system message.
type LoginUserAgentInterceptor struct {
	redisClient       *redis.Client
	userAgentInfoRepo repository.UserAgentInfoRepository
	longTermMemory    *memory.LongTermMemory
}

// NewLoginUserAgentInterceptor creates a new LoginUserAgentInterceptor.
// redisClient may be nil; if so, dedup and memory injection are skipped.
func NewLoginUserAgentInterceptor(
	redisClient *redis.Client,
	userAgentInfoRepo repository.UserAgentInfoRepository,
	longTermMemory *memory.LongTermMemory,
) *LoginUserAgentInterceptor {
	return &LoginUserAgentInterceptor{
		redisClient:       redisClient,
		userAgentInfoRepo: userAgentInfoRepo,
		longTermMemory:    longTermMemory,
	}
}

// DedupResult indicates whether the current request is a duplicate.
type DedupResult struct {
	IsDuplicate bool
}

// CheckDedup checks Redis for a recent dedup entry for the given user and session.
// Returns true if this request is a duplicate within the dedup window.
func (i *LoginUserAgentInterceptor) CheckDedup(
	ctx context.Context,
	userID, sessionID string,
) bool {
	if i.redisClient == nil {
		return false
	}

	key := redisDedupKeyPrefix + userID + ":" + sessionID
	exists, err := i.redisClient.Exists(ctx, key).Result()
	if err != nil {
		zap.L().Warn("login interceptor: redis dedup check failed",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return false
	}

	if exists > 0 {
		return true
	}

	// Set the dedup key with TTL.
	if err := i.redisClient.Set(ctx, key, "1", dedupWindow).Err(); err != nil {
		zap.L().Warn("login interceptor: failed to set dedup key",
			zap.String("userID", userID),
			zap.Error(err),
		)
	}

	return false
}

// RecordUsage asynchronously records the user's agent usage action.
// This is designed to be called via "go" (fire-and-forget).
func (i *LoginUserAgentInterceptor) RecordUsage(userID, agentSN string) {
	if i.userAgentInfoRepo == nil {
		return
	}

	// Use a background context with timeout for the async operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := i.userAgentInfoRepo.RecordAction(ctx, userID, agentSN); err != nil {
		zap.L().Warn("login interceptor: failed to record agent usage",
			zap.String("userID", userID),
			zap.String("agentSN", agentSN),
			zap.Error(err),
		)
	}
}

// InjectHistoryMemories searches the vector store for memories relevant to
// the user's current query and returns them formatted as a system message.
//
// Returns an empty string if no relevant memories are found or if the
// vector store is not configured.
func (i *LoginUserAgentInterceptor) InjectHistoryMemories(
	ctx context.Context,
	userID, query string,
) string {
	if i.longTermMemory == nil {
		return ""
	}

	memories, err := i.longTermMemory.SearchMemories(ctx, userID, query, 5)
	if err != nil {
		zap.L().Warn("login interceptor: failed to search memories",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return ""
	}

	if len(memories) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("【相关历史记忆】\n")
	for _, mem := range memories {
		builder.WriteString(fmt.Sprintf("- %s\n", mem.Content))
	}

	return builder.String()
}
