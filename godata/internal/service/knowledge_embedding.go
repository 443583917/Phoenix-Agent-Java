package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/repository"
	"go.uber.org/zap"
)

type EmbeddingService struct {
	repo   repository.AgentKnowledgeRepository
	logger *zap.Logger
}

func NewEmbeddingService(repo repository.AgentKnowledgeRepository) *EmbeddingService {
	return &EmbeddingService{
		repo:   repo,
		logger: zap.L().Named("embedding"),
	}
}

func (s *EmbeddingService) TriggerEmbedding(ctx context.Context, knowledgeID string) error {
	s.logger.Info("triggering embedding", zap.String("knowledgeId", knowledgeID))
	entity, err := s.repo.FindByID(ctx, knowledgeID)
	if err != nil {
		return err
	}
	if entity == nil {
		return nil
	}
	entity.EmbeddingStatus = "processing"
	if err := s.repo.Update(ctx, entity); err != nil {
		return err
	}
	// TODO: integrate with actual vector store and embedding model
	entity.EmbeddingStatus = "complete"
	return s.repo.Update(ctx, entity)
}

func (s *EmbeddingService) TriggerDeletion(ctx context.Context, knowledgeID string) error {
	s.logger.Info("triggering knowledge resource deletion", zap.String("knowledgeId", knowledgeID))
	// TODO: delete from vector store
	return nil
}

func (s *EmbeddingService) ProcessPendingBatch(ctx context.Context, batchSize int) (int, error) {
	s.logger.Info("processing pending embeddings", zap.Int("batchSize", batchSize))
	// TODO: query pending items and process them
	return 0, nil
}
