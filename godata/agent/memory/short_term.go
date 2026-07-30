package memory

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Message represents a single conversation turn.
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ShortTermConfig configures the short-term memory store.
type ShortTermConfig struct {
	// MaxMessages is the maximum number of messages retained per session.
	MaxMessages int
	// RedisClient, if set, enables Redis-backed persistence in addition to
	// the in-memory cache. When nil, only in-memory storage is used.
	RedisClient *redis.Client
	// TTL is the expiration duration for session history in Redis.
	// Defaults to 1 hour when RedisClient is set.
	TTL time.Duration
}

// defaultMaxMessages is used when config.MaxMessages is zero or negative.
const defaultMaxMessages = 50

// ShortTermMemory manages per-session conversation windows.
//
// Messages are stored in an in-memory buffer and, optionally, persisted to
// Redis for durability across restarts. Oldest messages are evicted when the
// per-session limit is exceeded (FIFO).
type ShortTermMemory struct {
	mu     sync.RWMutex
	store  map[string][]Message
	cfg    ShortTermConfig
}

// NewShortTermMemory creates a new short-term memory store.
func NewShortTermMemory(cfg ShortTermConfig) *ShortTermMemory {
	if cfg.MaxMessages <= 0 {
		cfg.MaxMessages = defaultMaxMessages
	}
	if cfg.RedisClient != nil && cfg.TTL <= 0 {
		cfg.TTL = 1 * time.Hour
	}
	return &ShortTermMemory{
		store: make(map[string][]Message),
		cfg:   cfg,
	}
}

// AddMessage appends a message to the session history. Messages exceeding the
// per-session limit are trimmed from the front (FIFO).
func (m *ShortTermMemory) AddMessage(sessionID string, msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	msg.Timestamp = time.Now()
	history := m.store[sessionID]
	history = append(history, msg)
	if len(history) > m.cfg.MaxMessages {
		history = history[len(history)-m.cfg.MaxMessages:]
	}
	m.store[sessionID] = history

	// Best-effort Redis persistence.
	if m.cfg.RedisClient != nil {
		data, _ := json.Marshal(history)
		m.cfg.RedisClient.Set(context.Background(), redisKey(sessionID), data, m.cfg.TTL)
	}
}

// GetHistory returns the full conversation history for a session.
// Returns an empty slice if no history exists.
func (m *ShortTermMemory) GetHistory(sessionID string) []Message {
	m.mu.RLock()
	history, ok := m.store[sessionID]
	m.mu.RUnlock()
	if ok {
		return copyMessages(history)
	}

	// Try Redis fallback.
	if m.cfg.RedisClient != nil {
		data, err := m.cfg.RedisClient.Get(context.Background(), redisKey(sessionID)).Bytes()
		if err != nil {
			return nil
		}
		var msgs []Message
		if err := json.Unmarshal(data, &msgs); err != nil {
			return nil
		}
		// Populate in-memory cache.
		m.mu.Lock()
		m.store[sessionID] = msgs
		m.mu.Unlock()
		return copyMessages(msgs)
	}
	return nil
}

// Clear removes all history for a session.
func (m *ShortTermMemory) Clear(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, sessionID)
	if m.cfg.RedisClient != nil {
		m.cfg.RedisClient.Del(context.Background(), redisKey(sessionID))
	}
}

// Len returns the number of messages stored for a session.
func (m *ShortTermMemory) Len(sessionID string) int {
	return len(m.GetHistory(sessionID))
}

func redisKey(sessionID string) string {
	return "memory:short_term:" + sessionID
}

func copyMessages(src []Message) []Message {
	dst := make([]Message, len(src))
	copy(dst, src)
	return dst
}
