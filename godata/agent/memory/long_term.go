package memory

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/memory"
	"trpc.group/trpc-go/trpc-agent-go/memory/inmemory"
)

// Memory represents a single long-term memory entry retrieved from the
// framework memory service.
type Memory struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Topics   []string               `json:"topics,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score"`
}

// LongTermMemory wraps the tRPC-Agent-Go memory.Service for long-term
// memory persistence and vector search. The framework handles embedding,
// storage, and similarity search internally.
//
// Previous version used custom Milvus HTTP calls; now delegates to the
// framework's inmemory.MemoryService with built-in memory management.
type LongTermMemory struct {
	svc     memory.Service
	appName string
}

// Embedder abstracts the embedding model call so LongTermMemory is decoupled
// from any specific embedding implementation.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

// NewLongTermMemory creates a new long-term memory instance backed by the
// tRPC-Agent-Go memory service.
//
// The service is created as in-memory by default. For production deployments,
// consider using a database-backed service (postgres, mysql, tencentdb).
func NewLongTermMemory(appName string) *LongTermMemory {
	if appName == "" {
		appName = "phoenix"
	}
	return &LongTermMemory{
		svc:     inmemory.NewMemoryService(),
		appName: appName,
	}
}

// NewLongTermMemoryWithService creates a LongTermMemory using a
// pre-configured framework memory.Service (e.g., postgres-backed).
func NewLongTermMemoryWithService(svc memory.Service, appName string) *LongTermMemory {
	if appName == "" {
		appName = "phoenix"
	}
	return &LongTermMemory{
		svc:     svc,
		appName: appName,
	}
}

// SearchMemories searches for relevant long-term memories for a given user
// using the framework's vector similarity search.
//
// Returns results sorted by relevance score (higher is better).
func (l *LongTermMemory) SearchMemories(
	ctx context.Context,
	userID string,
	query string,
	topK int,
) ([]Memory, error) {
	if topK <= 0 {
		topK = 5
	}

	userKey := memory.UserKey{
		AppName: l.appName,
		UserID:  userID,
	}

	entries, err := l.svc.SearchMemories(ctx, userKey, query,
		memory.WithSearchOptions(memory.SearchOptions{
			Query:      query,
			MaxResults: topK,
		}),
	)
	if err != nil {
		return nil, err
	}

	memories := make([]Memory, 0, len(entries))
	for _, entry := range entries {
		m := Memory{
			ID:    entry.ID,
			Score: entry.Score,
		}
		if entry.Memory != nil {
			m.Content = entry.Memory.Memory
			m.Topics = entry.Memory.Topics
		}
		memories = append(memories, m)
	}
	return memories, nil
}

// SaveMemory stores a new memory entry for a user via the framework
// memory service.
func (l *LongTermMemory) SaveMemory(
	ctx context.Context,
	userID string,
	content string,
	topics []string,
) error {
	userKey := memory.UserKey{
		AppName: l.appName,
		UserID:  userID,
	}
	return l.svc.AddMemory(ctx, userKey, content, topics)
}

// GetService returns the underlying framework memory.Service.
func (l *LongTermMemory) GetService() memory.Service {
	return l.svc
}
