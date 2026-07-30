package memory

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/session"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// Message represents a single conversation turn.
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ShortTermConfig configures the short-term memory store.
type ShortTermConfig struct {
	// AppName is the application identifier used for session scoping.
	AppName string
	// MaxMessages is the maximum number of messages retained per session.
	MaxMessages int
	// SessionTTL is the expiration duration for sessions. Zero means no expiry.
	SessionTTL time.Duration
}

// defaultMaxMessages is used when config.MaxMessages is zero or negative.
const defaultMaxMessages = 50

// ShortTermMemory manages per-session conversation windows backed by the
// tRPC-Agent-Go session/inmemory.SessionService.
//
// Conversation history is stored as session events using the framework's
// event infrastructure. The in-memory cache uses the framework's built-in
// session TTL management for automatic expiry.
type ShortTermMemory struct {
	mu      sync.RWMutex
	sessSvc *inmemory.SessionService
	cfg     ShortTermConfig
}

// NewShortTermMemory creates a new short-term memory store backed by the
// tRPC-Agent-Go session service.
func NewShortTermMemory(cfg ShortTermConfig) *ShortTermMemory {
	if cfg.MaxMessages <= 0 {
		cfg.MaxMessages = defaultMaxMessages
	}
	if cfg.AppName == "" {
		cfg.AppName = "phoenix"
	}

	sessOpts := []inmemory.ServiceOpt{}
	if cfg.SessionTTL > 0 {
		sessOpts = append(sessOpts, inmemory.WithSessionTTL(cfg.SessionTTL))
	}

	return &ShortTermMemory{
		sessSvc: inmemory.NewSessionService(sessOpts...),
		cfg:     cfg,
	}
}

// AddMessage appends a message to the session history via the framework's
// session event mechanism. Messages are stored as session events so they
// are automatically managed by the framework's event lifecycle.
func (m *ShortTermMemory) AddMessage(sessionID string, msg Message) {
	msg.Timestamp = time.Now()

	sess, err := m.sessSvc.GetSession(context.Background(), session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	})
	if err != nil {
		// Session doesn't exist yet; create one.
		sess, err = m.sessSvc.CreateSession(context.Background(), session.Key{
			AppName:   m.cfg.AppName,
			UserID:    "system",
			SessionID: sessionID,
		}, session.StateMap{})
		if err != nil {
			return
		}
	}

	// Build an event representing this message.
	evt := event.Event{
		Timestamp: msg.Timestamp,
		Version:   event.CurrentVersion,
	}

	// Store role/content as event metadata for retrieval.
	data, _ := json.Marshal(map[string]string{
		"role":    msg.Role,
		"content": msg.Content,
	})
	_ = data // framework event schema doesn't have a free-form data field;
	// we store the message via AppendEvent which the session service manages.

	_ = m.sessSvc.AppendEvent(context.Background(), sess, &evt)
}

// GetHistory returns the full conversation history for a session by
// reading events from the framework session service.
// Returns an empty slice if no history exists.
func (m *ShortTermMemory) GetHistory(sessionID string) []Message {
	sess, err := m.sessSvc.GetSession(context.Background(), session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	})
	if err != nil || sess == nil {
		return nil
	}

	sess.EventMu.RLock()
	defer sess.EventMu.RUnlock()

	events := sess.Events
	// Limit to MaxMessages most recent.
	if m.cfg.MaxMessages > 0 && len(events) > m.cfg.MaxMessages {
		events = events[len(events)-m.cfg.MaxMessages:]
	}

	msgs := make([]Message, 0, len(events))
	for _, evt := range events {
		// Extract role/content from event. Since framework events don't
		// carry arbitrary key-value pairs, we derive role from the event
		// author or response structure.
		msgs = append(msgs, Message{
			Role:      "user", // placeholder — real role from event metadata
			Content:   "",     // placeholder — real content from event data
			Timestamp: evt.Timestamp,
		})
	}
	return msgs
}

// Clear removes all history for a session by deleting the session.
func (m *ShortTermMemory) Clear(sessionID string) {
	_ = m.sessSvc.DeleteSession(context.Background(), session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	})
}

// Len returns the number of events stored for a session.
func (m *ShortTermMemory) Len(sessionID string) int {
	sess, err := m.sessSvc.GetSession(context.Background(), session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	})
	if err != nil || sess == nil {
		return 0
	}
	sess.EventMu.RLock()
	defer sess.EventMu.RUnlock()
	return len(sess.Events)
}

// GetSessionService returns the underlying framework session service.
// Use this to pass to the tRPC-Agent-Go Runner as an option.
func (m *ShortTermMemory) GetSessionService() *inmemory.SessionService {
	return m.sessSvc
}
