package chat

import (
	"context"
	"encoding/json"

	"github.com/phoenix-agent-go/internal/event"
	"go.uber.org/zap"
)

type SessionCreatedHandler struct {
	logger *zap.Logger
}

func NewSessionCreatedHandler() *SessionCreatedHandler {
	return &SessionCreatedHandler{
		logger: zap.L().Named("event.handler.session_created"),
	}
}

func (h *SessionCreatedHandler) Handle(ctx context.Context, e *event.Event) error {
	var payload event.SessionCreatedPayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		h.logger.Warn("failed to unmarshal session created payload", zap.Error(err))
		return nil
	}

	h.logger.Info("session created event processed",
		zap.String("sessionId", payload.SessionID),
		zap.String("agentId", payload.AgentID),
		zap.String("userId", payload.UserID),
		zap.String("title", payload.Title),
	)

	return nil
}
