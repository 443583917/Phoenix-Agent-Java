package privilege

import (
	"context"
	"encoding/json"

	"github.com/phoenix-agent-go/internal/event"
	"go.uber.org/zap"
)

type UserCreatedHandler struct {
	logger *zap.Logger
}

func NewUserCreatedHandler() *UserCreatedHandler {
	return &UserCreatedHandler{
		logger: zap.L().Named("event.handler.user_created"),
	}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, e *event.Event) error {
	var payload event.UserCreatedPayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		h.logger.Warn("failed to unmarshal user created payload", zap.Error(err))
		return nil
	}

	h.logger.Info("user created event processed",
		zap.String("userId", payload.UserID),
		zap.String("username", payload.Username),
	)

	return nil
}
