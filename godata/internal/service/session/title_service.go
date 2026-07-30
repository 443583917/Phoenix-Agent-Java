package session

import (
	"context"

	"go.uber.org/zap"
)

type TitleService struct {
	logger *zap.Logger
}

func NewTitleService() *TitleService {
	return &TitleService{logger: zap.L().Named("session.title")}
}

func (s *TitleService) GenerateTitle(ctx context.Context, sessionID string, firstMessage string) (string, error) {
	if firstMessage == "" {
		return "New Chat", nil
	}

	if len(firstMessage) <= 30 {
		return firstMessage, nil
	}

	s.logger.Info("generating session title",
		zap.String("sessionId", sessionID),
		zap.Int("msgLen", len(firstMessage)),
	)

	return firstMessage[:30] + "...", nil
}
