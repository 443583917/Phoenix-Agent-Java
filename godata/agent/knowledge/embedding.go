package knowledge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EmbeddingConfig holds the configuration for an OpenAI-compatible embedding API.
type EmbeddingConfig struct {
	// BaseURL is the API endpoint base (e.g. "https://api.openai.com/v1").
	BaseURL string
	// APIKey is the authentication key.
	APIKey string
	// Model is the embedding model name (e.g. "text-embedding-ada-002").
	Model string
}

// EmbeddingModel calls an OpenAI-compatible embedding API to produce vector
// representations of text.
type EmbeddingModel struct {
	cfg    EmbeddingConfig
	client *http.Client
}

// NewEmbeddingModel creates a new embedding model instance.
func NewEmbeddingModel(cfg EmbeddingConfig) *EmbeddingModel {
	if cfg.Model == "" {
		cfg.Model = "text-embedding-ada-002"
	}
	return &EmbeddingModel{
		cfg:    cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Embed generates an embedding vector for the given text.
func (e *EmbeddingModel) Embed(ctx context.Context, text string) ([]float32, error) {
	if e.cfg.BaseURL == "" {
		return nil, fmt.Errorf("embedding base URL not configured")
	}

	type embedReq struct {
		Input string `json:"input"`
		Model string `json:"model"`
	}
	type embedData struct {
		Embedding []float32 `json:"embedding"`
	}
	type embedResp struct {
		Data  []embedData `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	reqBody := embedReq{
		Input: text,
		Model: e.cfg.Model,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	url := e.cfg.BaseURL + "/embeddings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if e.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding http %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result embedResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("embedding api error: %s", result.Error.Message)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Data[0].Embedding, nil
}

// EmbedBatch generates embeddings for multiple texts in a single API call.
func (e *EmbeddingModel) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if e.cfg.BaseURL == "" {
		return nil, fmt.Errorf("embedding base URL not configured")
	}
	if len(texts) == 0 {
		return nil, nil
	}

	type embedReq struct {
		Input []string `json:"input"`
		Model string   `json:"model"`
	}
	type embedData struct {
		Embedding []float32 `json:"embedding"`
	}
	type embedResp struct {
		Data  []embedData `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	reqBody := embedReq{Input: texts, Model: e.cfg.Model}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	url := e.cfg.BaseURL + "/embeddings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if e.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding http %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result embedResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("embedding api error: %s", result.Error.Message)
	}

	vecs := make([][]float32, len(result.Data))
	for i, d := range result.Data {
		vecs[i] = d.Embedding
	}
	return vecs, nil
}
