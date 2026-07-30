package privilege

import (
	"context"
	"encoding/json"

	"github.com/phoenix-agent-go/internal/event"
	"go.uber.org/zap"
)

type LoginSuccessHandler struct {
	logger *zap.Logger
}

func NewLoginSuccessHandler() *LoginSuccessHandler {
	return &LoginSuccessHandler{
		logger: zap.L().Named("event.handler.login_success"),
	}
}

func (h *LoginSuccessHandler) Handle(ctx context.Context, e *event.Event) error {
	var payload event.LoginSuccessPayload
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		h.logger.Warn("failed to unmarshal login success payload", zap.Error(err))
		return nil
	}

	h.logger.Info("login success event processed",
		zap.String("userId", payload.UserID),
		zap.String("username", payload.Username),
		zap.String("ip", payload.IP),
	)

	return nil
}
