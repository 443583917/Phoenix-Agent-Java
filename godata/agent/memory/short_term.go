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

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ShortTermConfig struct {
	AppName     string
	MaxMessages int
	SessionTTL  time.Duration
}

const defaultMaxMessages = 50

type ShortTermMemory struct {
	mu      sync.RWMutex
	sessSvc *inmemory.SessionService
	cfg     ShortTermConfig
}

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

func (m *ShortTermMemory) AddMessage(sessionID string, msg Message) {
	msg.Timestamp = time.Now()

	key := session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	}

	sess, err := m.sessSvc.GetSession(context.Background(), key)
	if err != nil {
		sess, err = m.sessSvc.CreateSession(context.Background(), key, session.StateMap{})
		if err != nil {
			return
		}
	}

	data, _ := json.Marshal(msg)

	evt := event.Event{
		Timestamp: msg.Timestamp,
		Version:   event.CurrentVersion,
		Author:    msg.Role,
		StateDelta: map[string][]byte{
			"message": data,
		},
	}

	_ = m.sessSvc.AppendEvent(context.Background(), sess, &evt)
}

func (m *ShortTermMemory) GetHistory(sessionID string) []Message {
	key := session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	}

	sess, err := m.sessSvc.GetSession(context.Background(), key)
	if err != nil || sess == nil {
		return nil
	}

	sess.EventMu.RLock()
	defer sess.EventMu.RUnlock()

	events := sess.Events
	if m.cfg.MaxMessages > 0 && len(events) > m.cfg.MaxMessages {
		events = events[len(events)-m.cfg.MaxMessages:]
	}

	msgs := make([]Message, 0, len(events))
	for _, evt := range events {
		var msg Message
		if data, ok := evt.StateDelta["message"]; ok {
			if err := json.Unmarshal(data, &msg); err == nil {
				msgs = append(msgs, msg)
				continue
			}
		}
		if evt.Author != "" {
			msgs = append(msgs, Message{
				Role:      evt.Author,
				Content:   "",
				Timestamp: evt.Timestamp,
			})
		}
	}
	return msgs
}

func (m *ShortTermMemory) Clear(sessionID string) {
	_ = m.sessSvc.DeleteSession(context.Background(), session.Key{
		AppName:   m.cfg.AppName,
		UserID:    "system",
		SessionID: sessionID,
	})
}

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

func (m *ShortTermMemory) GetSessionService() *inmemory.SessionService {
	return m.sessSvc
}
