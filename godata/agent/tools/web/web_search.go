package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
		Description: "Search the web for information. Returns structured results with title, URL, and snippet.",
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

type SearchResult struct {
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

	results, err := t.search(ctx, params.Query, params.MaxResults)
	if err != nil {
		return []SearchResult{}, nil
	}
	return results, nil
}

var (
	titleRe   = regexp.MustCompile(`<a[^>]*class="result__a"[^>]*>(.*?)</a>`)
	snippetRe = regexp.MustCompile(`<a[^>]*class="result__snippet"[^>]*>(.*?)</a>`)
	urlRe     = regexp.MustCompile(`<a[^>]*class="result__url"[^>]*href="([^"]*)"`)
	tagRe     = regexp.MustCompile(`<[^>]+>`)
)

func (t *WebSearchTool) search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)

	titles := titleRe.FindAllStringSubmatch(html, maxResults)
	snippets := snippetRe.FindAllStringSubmatch(html, maxResults)
	urls := urlRe.FindAllStringSubmatch(html, maxResults)

	var results []SearchResult
	for i := 0; i < maxResults && i < len(titles); i++ {
		title := stripTags(titles[i][1])
		snippet := ""
		if i < len(snippets) {
			snippet = stripTags(snippets[i][1])
		}
		resultURL := ""
		if i < len(urls) {
			resultURL = urls[i][1]
			if decoded, err := url.QueryUnescape(resultURL); err == nil && strings.Contains(decoded, "uddg=") {
				if parsed, perr := url.Parse(decoded); perr == nil {
					if u := parsed.Query().Get("uddg"); u != "" {
						resultURL = u
					}
				}
			}
		}
		results = append(results, SearchResult{
			Title:   title,
			URL:     resultURL,
			Snippet: snippet,
		})
	}
	return results, nil
}

func stripTags(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	return strings.TrimSpace(tagRe.ReplaceAllString(s, ""))
}
