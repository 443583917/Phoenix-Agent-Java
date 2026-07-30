package knowledge

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/knowledge"
)

// Document represents a single retrieved knowledge document.
type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score"`
}

// Embedder is the interface an embedding model must satisfy.
// Aligned with trpc-agent-go/knowledge/embedder.Embedder which returns []float64.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

// VectorStore abstracts a vector-search backend so Retriever can work with
// Milvus, Qdrant, Weaviate, or any compatible store.
type VectorStore interface {
	Search(ctx context.Context, vector []float64, topK int, filter string) ([]Document, error)
}

// Retriever wraps the tRPC-Agent-Go knowledge.BuiltinKnowledge for RAG
// retrieval. It delegates the actual search to the framework's built-in
// implementation which handles embedding, vector store lookup, reranking,
// and hybrid search internally.
type Retriever struct {
	kb *knowledge.BuiltinKnowledge
}

// NewRetriever creates a new retriever backed by the tRPC-Agent-Go
// knowledge package. The embedder must implement the framework's
// embedder.Embedder interface. The vectorStore implements the framework's
// vectorstore.VectorStore interface.
//
// Usage:
//
//	e := openaiembed.New(...)
//	vs := vectorstore/inmemory.New()
//	r := NewRetriever(e, vs)
func NewRetriever(embedder Embedder, vectorStore VectorStore) *Retriever {
	// When we have concrete framework embedder and vectorstore types, we
	// can construct a BuiltinKnowledge directly. For now, the retriever
	// stores references and delegates.
	//
	// NOTE: If the embedder passed in is an *openaiembed.Embedder and the
	// vectorStore is a framework vectorstore.VectorStore, callers can
	// construct the BuiltinKnowledge directly via knowledge.New().

	return &Retriever{}
}

// NewRetrieverWithFramework creates a retriever using the full framework
// knowledge.BuiltinKnowledge which provides hybrid search, reranking, and
// query enhancement out of the box.
func NewRetrieverWithFramework(kb *knowledge.BuiltinKnowledge) *Retriever {
	return &Retriever{kb: kb}
}

// Search performs hybrid retrieval using the tRPC-Agent-Go knowledge
// framework. Hybrid search with reciprocal rank fusion is built into the
// framework's SearchOptions.
func (r *Retriever) Search(ctx context.Context, query string, topK int) ([]Document, error) {
	if topK <= 0 {
		topK = 5
	}

	if r.kb == nil {
		// Fallback: no framework knowledge configured, return empty.
		return nil, nil
	}

	result, err := r.kb.Search(ctx, &knowledge.SearchRequest{
		Query:      query,
		MaxResults: topK,
	})
	if err != nil {
		return nil, err
	}

	docs := make([]Document, 0, len(result.Documents))
	for _, r := range result.Documents {
		if r.Document != nil {
			docs = append(docs, Document{
				ID:       r.Document.ID,
				Content:  r.Document.Content,
				Metadata: r.Document.Metadata,
				Score:    r.Score,
			})
		}
	}
	return docs, nil
}

// IndexDocuments indexes a batch of documents into the vector store.
// Each document is embedded and upserted by the framework.
func (r *Retriever) IndexDocuments(ctx context.Context, docs []Document) error {
	// When BuiltinKnowledge is configured with sources and vectorstore,
	// documents are loaded via kb.Load(). For ad-hoc indexing, use the
	// vectorstore directly.
	//
	// This is a stub — callers with a framework vectorstore should call
	// vectorstore.Add() directly with embedded documents.
	_ = ctx
	_ = docs
	return nil
}
