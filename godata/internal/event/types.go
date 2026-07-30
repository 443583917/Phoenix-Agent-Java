package event

import (
	"encoding/json"
	"time"
)

const (
	TypeUserCreated     = "privilege.user.created"
	TypeLoginSuccess    = "privilege.login.success"
	TypeRoleUpdated     = "privilege.role.updated"
	TypeAgentCalled     = "agent.action.recorded"
	TypeSessionCreated  = "chat.session.created"
)

type Event struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
	Source    string          `json:"source,omitempty"`
}

func NewEvent(eventType string, payload interface{}) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Event{
		Type:      eventType,
		Payload:   data,
		Timestamp: time.Now(),
	}, nil
}

type UserCreatedPayload struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type LoginSuccessPayload struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	IP       string `json:"ip"`
}

type AgentActionPayload struct {
	AgentID   string `json:"agentId"`
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Action    string `json:"action"`
}

type SessionCreatedPayload struct {
	SessionID string `json:"sessionId"`
	AgentID   string `json:"agentId"`
	UserID    string `json:"userId"`
	Title     string `json:"title"`
}
