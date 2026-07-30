package agent

import (
	"context"
	"encoding/json"

	"github.com/phoenix-agent-go/internal/event"
	"go.uber.org/zap"
)

type ActionRecordedHandler struct {
	logger *zap.Logger
}

func NewActionRecordedHandler() *ActionRecordedHandler {
	return &ActionRecordedHandler{
		logger: zap.L().Named("event.handler.action_recorded"),
	}
}

func (h *ActionRecordedHandler) Handle(ctx context.Context, e *event.Event) error {
	var payload event.AgentActionPayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		h.logger.Warn("failed to unmarshal agent action payload", zap.Error(err))
		return nil
	}

	h.logger.Info("agent action recorded",
		zap.String("agentId", payload.AgentID),
		zap.String("userId", payload.UserID),
		zap.String("sessionId", payload.SessionID),
		zap.String("action", payload.Action),
	)

	return nil
}
