package vectorstore

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type VectorStoreType string

const (
	VectorStoreTypePgVector VectorStoreType = "pgvector"
	VectorStoreTypeMilvus   VectorStoreType = "milvus"
)

type Document struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Score    float64           `json:"score,omitempty"`
}

type AgentVectorStoreService struct {
	logger *zap.Logger
}

func NewAgentVectorStoreService() *AgentVectorStoreService {
	return &AgentVectorStoreService{logger: zap.L().Named("vectorstore.agent")}
}

func (s *AgentVectorStoreService) Search(ctx context.Context, agentID int64, query string, topK int) ([]Document, error) {
	s.logger.Info("vector search",
		zap.Int64("agentId", agentID),
		zap.String("query", query),
		zap.Int("topK", topK),
	)
	return []Document{}, nil
}

func (s *AgentVectorStoreService) AddDocuments(ctx context.Context, agentID int64, docs []Document) error {
	s.logger.Info("add documents",
		zap.Int64("agentId", agentID),
		zap.Int("count", len(docs)),
	)
	return nil
}

func (s *AgentVectorStoreService) DeleteDocuments(ctx context.Context, agentID int64, docIDs []string) error {
	s.logger.Info("delete documents",
		zap.Int64("agentId", agentID),
		zap.Int("count", len(docIDs)),
	)
	return nil
}

func (s *AgentVectorStoreService) DeleteByMetadata(ctx context.Context, agentID int64, metadata map[string]string) error {
	s.logger.Info("delete by metadata",
		zap.Int64("agentId", agentID),
	)
	return nil
}

type DynamicFilter struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

func BuildDynamicFilter(filters []DynamicFilter) string {
	return fmt.Sprintf("filters: %v", filters)
}
