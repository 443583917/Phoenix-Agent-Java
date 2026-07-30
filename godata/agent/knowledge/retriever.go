package knowledge

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Document represents a single retrieved knowledge document.
type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score"`
}

// VectorStore abstracts a vector-search backend so Retriever can work with
// Milvus, Qdrant, Weaviate, or any compatible store.
type VectorStore interface {
	Search(ctx context.Context, vector []float32, topK int, filter string) ([]Document, error)
}

// Retriever performs hybrid search combining vector (semantic) and keyword
// (lexical) retrieval with reciprocal rank fusion.
type Retriever struct {
	vectorStore VectorStore
	embedder    Embedder
}

// Embedder is the interface an embedding model must satisfy.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// NewRetriever creates a new hybrid retriever.
func NewRetriever(vs VectorStore, emb Embedder) *Retriever {
	return &Retriever{
		vectorStore: vs,
		embedder:    emb,
	}
}

// Search performs hybrid retrieval: it runs vector search via the configured
// VectorStore and a keyword-based BM25-like placeholder search, then merges
// results with reciprocal rank fusion.
//
// When the VectorStore is nil (e.g. no Milvus configured), only keyword
// search runs.
func (r *Retriever) Search(ctx context.Context, query string, topK int) ([]Document, error) {
	if topK <= 0 {
		topK = 5
	}

	if r.vectorStore == nil && r.embedder == nil {
		return r.keywordSearch(query, topK), nil
	}

	if r.vectorStore == nil {
		return r.keywordSearch(query, topK), nil
	}

	// Vector search.
	vec, err := r.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embed: %w", err)
	}
	vectorDocs, err := r.vectorStore.Search(ctx, vec, topK*2, "")
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	// Keyword search.
	keywordDocs := r.keywordSearch(query, topK)

	// Merge with reciprocal rank fusion.
	return fuseResults(vectorDocs, keywordDocs, topK), nil
}

// keywordSearch performs a simple keyword-based search against local or
// in-memory document storage. This is a placeholder for a full-text index
// (e.g. Elasticsearch or PostgreSQL tsvector).
func (r *Retriever) keywordSearch(query string, topK int) []Document {
	// Placeholder: returns empty. A real implementation would search an
	// inverted index or full-text search backend.
	_ = query
	_ = topK
	return nil
}

// fuseResults merges vector and keyword results using reciprocal rank fusion.
func fuseResults(vectorDocs, keywordDocs []Document, topK int) []Document {
	const k = 60.0

	scoreMap := make(map[string]float64)
	contentMap := make(map[string]Document)

	// Accumulate reciprocal rank scores.
	for i, doc := range vectorDocs {
		scoreMap[doc.ID] += 1.0 / (k + float64(i+1))
		if _, ok := contentMap[doc.ID]; !ok {
			contentMap[doc.ID] = doc
		}
	}
	for i, doc := range keywordDocs {
		scoreMap[doc.ID] += 1.0 / (k + float64(i+1))
		if _, ok := contentMap[doc.ID]; !ok {
			contentMap[doc.ID] = doc
		}
	}

	// Sort by fused score descending.
	type scoredDoc struct {
		id    string
		score float64
	}
	ranked := make([]scoredDoc, 0, len(scoreMap))
	for id, score := range scoreMap {
		ranked = append(ranked, scoredDoc{id, score})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].score > ranked[j].score
	})

	limit := topK
	if limit > len(ranked) {
		limit = len(ranked)
	}
	result := make([]Document, 0, limit)
	for i := 0; i < limit; i++ {
		doc := contentMap[ranked[i].id]
		doc.Score = ranked[i].score
		result = append(result, doc)
	}
	return result
}

// IndexDocuments indexes a batch of documents into the vector store.
// Each document is embedded and upserted.
func (r *Retriever) IndexDocuments(ctx context.Context, docs []Document) error {
	if r.vectorStore == nil || r.embedder == nil {
		return fmt.Errorf("vector store or embedder not configured")
	}
	// NOTE: This is a placeholder — a real implementation depends on the
	// concrete VectorStore supporting an upsert/insert operation. The Milvus
	// client in agent/memory/long_term.go implements SaveMemory for this.
	for _, doc := range docs {
		_, err := r.embedder.Embed(ctx, doc.Content)
		if err != nil {
			return fmt.Errorf("embed doc %s: %w", doc.ID, err)
		}
		// Store the embedding + document in the vector store.
		_ = doc // placeholder — actual insert depends on VectorStore contract
	}
	return nil
}

// ContainsKeyword is a helper that checks whether any document content
// contains the given keyword (case-insensitive).
func ContainsKeyword(docs []Document, keyword string) []Document {
	if len(docs) == 0 || keyword == "" {
		return nil
	}
	lower := strings.ToLower(keyword)
	var matched []Document
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Content), lower) {
			matched = append(matched, doc)
		}
	}
	return matched
}
