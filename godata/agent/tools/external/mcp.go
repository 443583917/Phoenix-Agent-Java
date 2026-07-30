package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type MCPServerConfig struct {
	Name    string            `json:"name" yaml:"name"`
	URL     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type MCPToolSource struct {
	servers []MCPServerConfig
}

func NewMCPToolSource(servers []MCPServerConfig) *MCPToolSource {
	return &MCPToolSource{servers: servers}
}

type mcpTool struct {
	name        string
	description string
	inputSchema *tool.Schema
	serverURL   string
	headers     map[string]string
}

func (t *mcpTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        t.name,
		Description: t.description,
		InputSchema: t.inputSchema,
	}
}

func (t *mcpTool) Call(ctx context.Context, args any) (any, error) {
	body, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"params": map[string]any{
			"name":      t.name,
			"arguments": args,
		},
		"id": 1,
	})
	if err != nil {
		return nil, fmt.Errorf("mcp marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.serverURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("mcp request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mcp call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mcp read: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("mcp unmarshal: %w", err)
	}
	return result, nil
}

type mcpListToolsResponse struct {
	Result struct {
		Tools []struct {
			Name        string       `json:"name"`
			Description string       `json:"description"`
			InputSchema *tool.Schema `json:"inputSchema"`
		} `json:"tools"`
	} `json:"result"`
}

func (m *MCPToolSource) DiscoverTools(ctx context.Context) ([]tool.Tool, error) {
	var allTools []tool.Tool

	for _, server := range m.servers {
		tools, err := m.discoverFromServer(ctx, server)
		if err != nil {
			continue
		}
		allTools = append(allTools, tools...)
	}

	return allTools, nil
}

func (m *MCPToolSource) discoverFromServer(ctx context.Context, server MCPServerConfig) ([]tool.Tool, error) {
	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "tools/list",
		"id":      1,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range server.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResp mcpListToolsResponse
	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return nil, err
	}

	var tools []tool.Tool
	for _, t := range listResp.Result.Tools {
		tools = append(tools, &mcpTool{
			name:        t.Name,
			description: t.Description,
			inputSchema: t.InputSchema,
			serverURL:   server.URL,
			headers:     server.Headers,
		})
	}

	return tools, nil
}

func (m *MCPToolSource) RegisterAll(ctx context.Context, registry interface{ Register(tool.Tool) }) error {
	tools, err := m.DiscoverTools(ctx)
	if err != nil {
		return err
	}
	for _, t := range tools {
		registry.Register(t)
	}
	return nil
}
