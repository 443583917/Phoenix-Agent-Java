package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	tools  []ToolDefinition
}

type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

func NewServer() *Server {
	return &Server{
		logger: zap.L().Named("mcp.server"),
		tools:  defaultTools(),
	}
}

func defaultTools() []ToolDefinition {
	agentListSchema, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]string{"type": "string", "description": "Filter by agent status"},
		},
	})

	nl2sqlSchema, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]string{"type": "string", "description": "Natural language query"},
		},
		"required": []string{"query"},
	})

	return []ToolDefinition{
		{
			Name:        "list_agents",
			Description: "List all available agents with optional status filter",
			InputSchema: agentListSchema,
		},
		{
			Name:        "nl2sql_convert",
			Description: "Convert a natural language query to SQL using the NL2SQL engine",
			InputSchema: nl2sqlSchema,
		},
	}
}

func (s *Server) ListTools(ctx context.Context) ([]ToolDefinition, error) {
	return s.tools, nil
}

func (s *Server) CallTool(ctx context.Context, name string, args json.RawMessage) (json.RawMessage, error) {
	s.logger.Info("mcp tool call", zap.String("name", name))

	switch name {
	case "list_agents":
		result, _ := json.Marshal(map[string]interface{}{
			"agents": []interface{}{},
		})
		return result, nil
	case "nl2sql_convert":
		result, _ := json.Marshal(map[string]interface{}{
			"sql": "-- NL2SQL conversion not yet wired",
		})
		return result, nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}
