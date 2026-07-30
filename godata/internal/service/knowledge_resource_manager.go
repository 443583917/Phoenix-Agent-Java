package service

import (
	"context"

	"go.uber.org/zap"
)

type KnowledgeResourceManager struct {
	logger *zap.Logger
}

func NewKnowledgeResourceManager() *KnowledgeResourceManager {
	return &KnowledgeResourceManager{logger: zap.L().Named("knowledge.resource")}
}

func (m *KnowledgeResourceManager) Embed(ctx context.Context, knowledgeID string) error {
	m.logger.Info("embedding knowledge", zap.String("id", knowledgeID))
	return nil
}

func (m *KnowledgeResourceManager) DeleteResources(ctx context.Context, knowledgeID string) error {
	m.logger.Info("deleting knowledge resources", zap.String("id", knowledgeID))
	return nil
}

func (m *KnowledgeResourceManager) CleanupOrphaned(ctx context.Context) (int, error) {
	m.logger.Info("cleaning up orphaned knowledge resources")
	return 0, nil
}
