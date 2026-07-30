package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type WebSearchTool struct {
	httpClient *http.Client
}

func NewWebSearchTool() *WebSearchTool {
	return &WebSearchTool{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (t *WebSearchTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "web_search",
		Description: "Search the web for information using DuckDuckGo",
		InputSchema: &tool.Schema{
			Type:     "object",
			Required: []string{"query"},
			Properties: map[string]*tool.Schema{
				"query": {
					Type:        "string",
					Description: "Search query string",
				},
				"max_results": {
					Type:        "integer",
					Description: "Maximum number of results to return (default 5)",
				},
			},
		},
	}
}

type webSearchParams struct {
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
}

type webSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

func (t *WebSearchTool) Call(ctx context.Context, args any) (any, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	var params webSearchParams
	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return nil, fmt.Errorf("unmarshal args: %w", err)
	}

	if params.MaxResults <= 0 {
		params.MaxResults = 5
	}

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(params.Query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; PhoenixAgent/1.0)")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	results := []webSearchResult{
		{
			Title:   "Search results for: " + params.Query,
			URL:     searchURL,
			Snippet: string(body)[:min(500, len(body))],
		},
	}

	return results, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
