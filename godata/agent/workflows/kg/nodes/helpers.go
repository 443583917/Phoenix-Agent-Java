package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/agent/workflows/kg/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type LLMService struct {
	model.Model
}

func NewLLMService(m model.Model) *LLMService {
	return &LLMService{Model: m}
}

func (s *LLMService) Call(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := []model.Message{
		model.NewSystemMessage(systemPrompt),
		model.NewUserMessage(userPrompt),
	}
	temp := 0.1
	req := &model.Request{
		Messages: messages,
		GenerationConfig: model.GenerationConfig{
			Temperature: &temp,
		},
	}
	respChan, err := s.GenerateContent(ctx, req)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for resp := range respChan {
		if resp.Error != nil {
			return "", fmt.Errorf("LLM error: %s", resp.Error.Message)
		}
		for _, choice := range resp.Choices {
			sb.WriteString(choice.Delta.Content)
		}
	}
	return sb.String(), nil
}

func (s *LLMService) CallJSON(ctx context.Context, systemPrompt, userPrompt string, output any) error {
	resp, err := s.Call(ctx, systemPrompt, userPrompt)
	if err != nil {
		return err
	}
	jsonStr := extractJSON(resp)
	if jsonStr == "" {
		return fmt.Errorf("no JSON in LLM response")
	}
	return json.Unmarshal([]byte(jsonStr), output)
}

func extractJSON(s string) string {
	if idx := strings.Index(s, "```json"); idx >= 0 {
		rest := s[idx+7:]
		if end := strings.Index(rest, "```"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
	}
	if start := strings.Index(s, "{"); start >= 0 {
		depth := 0
		for i := start; i < len(s); i++ {
			if s[i] == '{' {
				depth++
			} else if s[i] == '}' {
				depth--
				if depth == 0 {
					return s[start : i+1]
				}
			}
		}
	}
	return ""
}

func getKGState(state graph.State) *types.KGState {
	if s, ok := state["kg_state"].(*types.KGState); ok && s != nil {
		return s
	}
	return &types.KGState{}
}
