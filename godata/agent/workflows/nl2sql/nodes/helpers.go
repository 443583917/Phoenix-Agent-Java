package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/agent/workflows/nl2sql/prompts"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ──────────────────────────── State Helpers ────────────────────────────

// getOrCreateState retrieves the NL2SQLState from the graph.State or creates
// a new one if not present.
func getOrCreateState(state graph.State) *nl2sqltypes.NL2SQLState {
	if s, ok := state["nl2sql_state"].(*nl2sqltypes.NL2SQLState); ok && s != nil {
		return s
	}
	return &nl2sqltypes.NL2SQLState{}
}

// getStateStr retrieves a string value from graph state by key.
func getStateStr(state graph.State, key string) string {
	if v, ok := state[key].(string); ok {
		return v
	}
	return ""
}

// ──────────────────────────── LLM Helper ────────────────────────────

// LLMService wraps the model.Model for use by NL2SQL nodes.
// It provides convenience methods for structured LLM calls.
type LLMService struct {
	model.Model
}

// NewLLMService creates a new LLMService.
func NewLLMService(m model.Model) *LLMService {
	return &LLMService{Model: m}
}

// LLMResponse is the raw response from a model call.
type LLMResponse struct {
	Content string
	Usage   *struct {
		PromptTokens     int `json:"prompt_tokens,omitempty"`
		CompletionTokens int `json:"completion_tokens,omitempty"`
	}
}

// Call sends a prompt to the LLM and returns the text response.
func float64Ptr(v float64) *float64 {
	return &v
}

func (s *LLMService) Call(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := []model.Message{
		model.NewSystemMessage(systemPrompt),
		model.NewUserMessage(userPrompt),
	}

	req := &model.Request{
		Messages: messages,
		GenerationConfig: model.GenerationConfig{
			Temperature: float64Ptr(0.1),
		},
	}

	respChan, err := s.GenerateContent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	var fullContent strings.Builder
	for resp := range respChan {
		if resp.Error != nil {
			return "", fmt.Errorf("LLM response error: %s", resp.Error.Message)
		}
		for _, choice := range resp.Choices {
			// 流式响应内容在 Delta，非流式响应内容在 Message，两者都读取
			fullContent.WriteString(choice.Delta.Content)
			fullContent.WriteString(choice.Message.Content)
		}
	}

	return fullContent.String(), nil
}

// CallWithPrompt renders a template and calls the LLM with system/user messages.
func (s *LLMService) CallWithPrompt(ctx context.Context, promptTemplate string, data map[string]string) (string, error) {
	// Simple template rendering (Go text/template replacement)
	systemMsg := extractSystemFromPrompt(promptTemplate)
	userMsg := renderTemplate(promptTemplate, data)

	return s.Call(ctx, systemMsg, userMsg)
}

// CallJSON calls the LLM and parses the response as JSON.
func (s *LLMService) CallJSON(ctx context.Context, promptTemplate string, data map[string]string, output any) error {
	resp, err := s.CallWithPrompt(ctx, promptTemplate, data)
	if err != nil {
		return err
	}

	// Extract JSON from the response (handle markdown code blocks)
	jsonStr := extractJSON(resp)
	if jsonStr == "" {
		return fmt.Errorf("no JSON found in LLM response: %s", truncate(resp, 200))
	}

	if err := json.Unmarshal([]byte(jsonStr), output); err != nil {
		return fmt.Errorf("failed to parse LLM JSON response: %w\nResponse: %s", err, truncate(resp, 200))
	}

	return nil
}

// ──────────────────────────── Template Helpers ────────────────────────────

// renderTemplate performs simple {{.key}} replacement in the template string.
func renderTemplate(tmpl string, data map[string]string) string {
	result := tmpl
	for k, v := range data {
		placeholder := "{{." + k + "}}"
		result = strings.ReplaceAll(result, placeholder, v)
	}
	return result
}

// extractSystemFromPrompt extracts the system message from a prompt string.
// If the prompt starts with a system-like instruction block, returns it as system message.
// For simplicity, returns the full prompt as user-facing text.
func extractSystemFromPrompt(prompt string) string {
	// For NL2SQL prompts, the entire prompt is treated as the system instruction.
	// The user's specific query is passed separately via the template data.
	return prompt
}

// extractJSON extracts a JSON object or array from a string that may contain
// markdown code fences, explanatory text, etc.
func extractJSON(s string) string {
	// Try to find JSON between code fences first
	if idx := strings.Index(s, "```json"); idx >= 0 {
		rest := s[idx+7:]
		if end := strings.Index(rest, "```"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
	}
	if idx := strings.Index(s, "```"); idx >= 0 {
		rest := s[idx+3:]
		if end := strings.Index(rest, "```"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
	}

	// Try finding a { ... } block
	if start := strings.Index(s, "{"); start >= 0 {
		depth := 0
		for i := start; i < len(s); i++ {
			switch s[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					return s[start : i+1]
				}
			}
		}
	}

	return ""
}

// truncate truncates a string to the given length.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// ──────────────────────────── Prompt Lookup ────────────────────────────

// GetPrompt returns the named prompt template.
func GetPrompt(name string) string {
	switch name {
	case "intent_recognition":
		return prompts.IntentRecognitionPrompt
	case "evidence_rewrite":
		return prompts.EvidenceRewritePrompt
	case "query_enhancement":
		return prompts.QueryEnhancementPrompt
	case "feasibility_assessment":
		return prompts.FeasibilityAssessmentPrompt
	case "planner":
		return prompts.PlannerPrompt
	case "sql_generate":
		return prompts.SQLGeneratePrompt
	case "semantic_consistency":
		return prompts.SemanticConsistencyPrompt
	case "python_generate":
		return prompts.PythonGeneratePrompt
	case "report_generate":
		return prompts.ReportGeneratePrompt
	case "table_relation":
		return prompts.TableRelationPrompt
	case "python_analyze":
		return prompts.PythonAnalyzePrompt
	default:
		return ""
	}
}
