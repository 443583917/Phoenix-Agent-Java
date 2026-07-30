package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Memory represents a single long-term memory entry retrieved from the
// vector store (Milvus).
type Memory struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Embedding []float32              `json:"embedding,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Score     float64                `json:"score"`
}

// MilvusConfig holds connection parameters for the Milvus vector database.
type MilvusConfig struct {
	// Address is the Milvus RESTful API endpoint (e.g. "http://localhost:19530").
	Address string
	// Collection is the Milvus collection name for memory storage.
	Collection string
	// User is the optional authentication username.
	User string
	// Password is the optional authentication password.
	Password string
}

// LongTermMemory performs vector search and persistence via the Milvus RESTful API.
//
// It embeds text via the provided Embedder interface and stores/retrieves
// vectors in the configured Milvus collection.
type LongTermMemory struct {
	cfg      MilvusConfig
	embedder Embedder
	client   *http.Client
}

// Embedder abstracts the embedding model call so LongTermMemory is decoupled
// from any specific embedding implementation.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// NewLongTermMemory creates a new long-term memory instance.
func NewLongTermMemory(cfg MilvusConfig, embedder Embedder) *LongTermMemory {
	if cfg.Collection == "" {
		cfg.Collection = "agent_memory"
	}
	return &LongTermMemory{
		cfg:      cfg,
		embedder: embedder,
		client:   &http.Client{},
	}
}

// SearchMemories searches for relevant long-term memories for a given user
// using vector similarity search in Milvus.
//
// Returns results sorted by relevance score (higher is better).
func (l *LongTermMemory) SearchMemories(
	ctx context.Context,
	userID string,
	query string,
	topK int,
) ([]Memory, error) {
	if l.cfg.Address == "" {
		return nil, fmt.Errorf("milvus address not configured")
	}
	if topK <= 0 {
		topK = 5
	}

	vec, err := l.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	return l.searchMilvus(ctx, vec, topK, fmt.Sprintf("user_id == \"%s\"", userID))
}

// SaveMemory stores a new memory entry in Milvus with its embedding vector.
func (l *LongTermMemory) SaveMemory(
	ctx context.Context,
	userID string,
	content string,
	metadata map[string]interface{},
) error {
	if l.cfg.Address == "" {
		return fmt.Errorf("milvus address not configured")
	}

	vec, err := l.embedder.Embed(ctx, content)
	if err != nil {
		return fmt.Errorf("embed content: %w", err)
	}

	return l.insertMilvus(ctx, vec, content, userID, metadata)
}

// ──────────────────────────── Milvus REST helpers ────────────────────────────

type milvusSearchReq struct {
	CollectionName string    `json:"collectionName"`
	Vector         []float32 `json:"vector"`
	Limit          int       `json:"limit"`
	Filter         string    `json:"filter,omitempty"`
	OutputFields   []string  `json:"outputFields"`
}

type milvusSearchResp struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    []milvusSearchResult  `json:"data"`
}

type milvusSearchResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Fields   map[string]interface{} `json:"fields"`
	Vector   []float32              `json:"vector"`
}

func (l *LongTermMemory) searchMilvus(
	ctx context.Context,
	vector []float32,
	topK int,
	filter string,
) ([]Memory, error) {
	reqBody := milvusSearchReq{
		CollectionName: l.cfg.Collection,
		Vector:         vector,
		Limit:          topK,
		Filter:         filter,
		OutputFields:   []string{"content", "metadata"},
	}

	var resp milvusSearchResp
	if err := l.doPost(ctx, "/v2/vectordb/entities/search", reqBody, &resp); err != nil {
		return nil, fmt.Errorf("milvus search: %w", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("milvus search error (code=%d): %s", resp.Code, resp.Message)
	}

	memories := make([]Memory, 0, len(resp.Data))
	for _, r := range resp.Data {
		content, _ := r.Fields["content"].(string)
		rawMeta, _ := r.Fields["metadata"].(map[string]interface{})
		memories = append(memories, Memory{
			ID:       r.ID,
			Content:  content,
			Score:    r.Score,
			Metadata: rawMeta,
		})
	}
	return memories, nil
}

type milvusInsertReq struct {
	CollectionName string               `json:"collectionName"`
	Data           []milvusInsertEntity `json:"data"`
}

type milvusInsertEntity struct {
	Vector   []float32              `json:"vector"`
	Content  string                 `json:"content"`
	UserID   string                 `json:"user_id"`
	Metadata map[string]interface{} `json:"metadata"`
}

type milvusInsertResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (l *LongTermMemory) insertMilvus(
	ctx context.Context,
	vector []float32,
	content, userID string,
	metadata map[string]interface{},
) error {
	reqBody := milvusInsertReq{
		CollectionName: l.cfg.Collection,
		Data: []milvusInsertEntity{{
			Vector:   vector,
			Content:  content,
			UserID:   userID,
			Metadata: metadata,
		}},
	}

	var resp milvusInsertResp
	if err := l.doPost(ctx, "/v2/vectordb/entities/insert", reqBody, &resp); err != nil {
		return fmt.Errorf("milvus insert: %w", err)
	}
	if resp.Code != 0 {
		return fmt.Errorf("milvus insert error (code=%d): %s", resp.Code, resp.Message)
	}
	return nil
}

func (l *LongTermMemory) doPost(ctx context.Context, path string, body, into interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	url := l.cfg.Address + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if l.cfg.User != "" {
		req.SetBasicAuth(l.cfg.User, l.cfg.Password)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("milvus http %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if into != nil {
		if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
			return fmt.Errorf("decode: %w", err)
		}
	}
	return nil
}
