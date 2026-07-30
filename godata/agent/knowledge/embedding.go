package knowledge

import (
	"context"

	openaiembed "trpc.group/trpc-go/trpc-agent-go/knowledge/embedder/openai"
)

// EmbeddingConfig holds the configuration for an OpenAI-compatible embedding API.
type EmbeddingConfig struct {
	// BaseURL is the API endpoint base (e.g. "https://api.openai.com/v1").
	BaseURL string
	// APIKey is the authentication key.
	APIKey string
	// Model is the embedding model name (e.g. "text-embedding-3-small").
	Model string
}

// EmbeddingModel wraps the tRPC-Agent-Go openaiembed.Embedder to provide a
// simplified interface for the Phoenix agent layer. The framework handles
// HTTP transport, retries, authentication, and tracing internally.
type EmbeddingModel struct {
	cfg      EmbeddingConfig
	embedder *openaiembed.Embedder
}

// NewEmbeddingModel creates a new embedding model instance backed by the
// tRPC-Agent-Go OpenAI embedder.
func NewEmbeddingModel(cfg EmbeddingConfig) *EmbeddingModel {
	if cfg.Model == "" {
		cfg.Model = "text-embedding-3-small"
	}
	opts := []openaiembed.Option{
		openaiembed.WithModel(cfg.Model),
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openaiembed.WithBaseURL(cfg.BaseURL))
	}
	if cfg.APIKey != "" {
		opts = append(opts, openaiembed.WithAPIKey(cfg.APIKey))
	}
	return &EmbeddingModel{
		cfg:      cfg,
		embedder: openaiembed.New(opts...),
	}
}

// Embed generates an embedding vector for the given text.
// Returns float64 slice to match the framework embedder.Embedder interface.
func (e *EmbeddingModel) Embed(ctx context.Context, text string) ([]float64, error) {
	if e.cfg.BaseURL == "" {
		// Delegate to framework default (uses OPENAI_API_KEY env var or defaults).
	}
	return e.embedder.GetEmbedding(ctx, text)
}

// EmbedBatch generates embeddings for multiple texts by issuing individual
// calls. The framework openaiembed currently does not expose a batch endpoint,
// but we iterate to preserve the existing batch API surface.
func (e *EmbeddingModel) EmbedBatch(ctx context.Context, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	vecs := make([][]float64, 0, len(texts))
	for _, t := range texts {
		vec, err := e.embedder.GetEmbedding(ctx, t)
		if err != nil {
			return nil, err
		}
		vecs = append(vecs, vec)
	}
	return vecs, nil
}

// GetEmbedder returns the underlying framework embedder for use with
// knowledge.New(WithEmbedder(...)).
func (e *EmbeddingModel) GetEmbedder() *openaiembed.Embedder {
	return e.embedder
}

// Ensure EmbeddingModel satisfies the local Embedder interface.
// The framework embedder returns []float64 which is compatible.
var _ Embedder = (*EmbeddingModel)(nil)
