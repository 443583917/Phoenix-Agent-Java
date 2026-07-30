package runtime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/phoenix-agent-go/internal/model"

	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
)

// AgentManager manages the lifecycle of tRPC-Agent-Go agents and provides
// streaming call functionality via the runner.
//
// Each StreamCall creates a fresh tRPC-Agent-Go agent and runner instance
// for request isolation.
type AgentManager struct {
	registry *AgentRegistry
}

// NewAgentManager creates a new AgentManager with the given agent registry.
func NewAgentManager(registry *AgentRegistry) *AgentManager {
	return &AgentManager{registry: registry}
}

// Registry returns the underlying agent registry for external configuration.
func (m *AgentManager) Registry() *AgentRegistry {
	return m.registry
}

// StreamCall executes a streaming agent call and returns a channel of SSE events.
//
// It creates a fresh tRPC-Agent-Go agent and runner for each call, sends the
// user message, and converts runner events into SSEEvent values on the returned
// channel. The channel is closed when the run completes.
func (m *AgentManager) StreamCall(
	ctx context.Context,
	req model.StreamRequest,
) (<-chan model.SSEEvent, error) {
	agentCfg, err := m.registry.Get(req.AgentSN)
	if err != nil {
		return nil, err
	}

	// Build model instance with provider-specific options.
	modelOpts := []openai.Option{
		openai.WithVariant(openai.VariantDeepSeek),
	}
	if agentCfg.BaseURL != "" {
		modelOpts = append(modelOpts, openai.WithBaseURL(agentCfg.BaseURL))
	}
	if agentCfg.APIKey != "" {
		modelOpts = append(modelOpts, openai.WithAPIKey(agentCfg.APIKey))
	}
	modelInstance := openai.New(agentCfg.ModelName, modelOpts...)

	// Build generation config.
	genConfig := tmodel.GenerationConfig{
		Stream: true,
	}
	if agentCfg.MaxTokens > 0 {
		maxTokens := agentCfg.MaxTokens
		genConfig.MaxTokens = &maxTokens
	}
	if agentCfg.Temperature > 0 {
		temp := agentCfg.Temperature
		genConfig.Temperature = &temp
	}

	// Build llmagent with tools if configured.
	agentOpts := []llmagent.Option{
		llmagent.WithModel(modelInstance),
		llmagent.WithGenerationConfig(genConfig),
	}
	if len(agentCfg.Tools) > 0 {
		agentOpts = append(agentOpts, llmagent.WithTools(agentCfg.Tools))
	}
	llmAgent := llmagent.New(req.AgentSN, agentOpts...)

	// Create runner and execute.
	r := runner.NewRunner(req.AgentSN+"-app", llmAgent)
	events, err := r.Run(
		ctx, req.UserID, req.SessionID, tmodel.NewUserMessage(req.Message),
	)
	if err != nil {
		return nil, fmt.Errorf("agent run: %w", err)
	}

	// Convert runner event stream to SSE event stream.
	ch := make(chan model.SSEEvent, 10)
	go func() {
		defer close(ch)
		for evt := range events {
			sseEvent := convertRunnerEvent(evt)
			ch <- sseEvent
		}
		// Signal end of stream.
		ch <- model.SSEEvent{
			Event: "content",
			Data: model.ContentEvent{
				Content: "",
				End:     true,
			},
		}
	}()

	return ch, nil
}

// convertRunnerEvent converts a tRPC-Agent-Go runner Event to an SSEEvent.
//
// Handles:
//   - Chat completion chunks with text content → ContentEvent
//   - Tool call events → ToolCallEvent
//   - Error responses → ErrorEvent
func convertRunnerEvent(evt *event.Event) model.SSEEvent {
	// Handle API-level errors in the response.
	if evt.Response != nil && evt.Response.Error != nil {
		code := 500
		if evt.Response.Error.Code != nil {
			code = errorCodeToInt(*evt.Response.Error.Code)
		}
		return model.SSEEvent{
			Event: "error",
			Data: model.ErrorEvent{
				Message: evt.Response.Error.Message,
				Code:    code,
			},
		}
	}

	// Handle chat completion chunks.
	if evt.Response != nil && len(evt.Response.Choices) > 0 {
		choice := evt.Response.Choices[0]

		// Check for tool calls in the delta.
		if len(choice.Delta.ToolCalls) > 0 {
			tc := choice.Delta.ToolCalls[0]
			args := make(map[string]interface{})
			if len(tc.Function.Arguments) > 0 {
				_ = json.Unmarshal(tc.Function.Arguments, &args)
			}
			return model.SSEEvent{
				Event: "tool_call",
				Data: model.ToolCallEvent{
					ToolName: tc.Function.Name,
					Args:     args,
				},
			}
		}

		// Handle text content chunks.
		if choice.Delta.Content != "" {
			return model.SSEEvent{
				Event: "content",
				Data: model.ContentEvent{
					Content: choice.Delta.Content,
					End:     false,
				},
			}
		}
	}

	// Default: skip unrecognized events (empty SSEEvent is discarded by callers).
	return model.SSEEvent{}
}

// errorCodeToInt converts an API error code string to an integer.
// If the string is a numeric code it parses it; otherwise returns 500.
func errorCodeToInt(code string) int {
	if code == "" {
		return 500
	}
	var n int
	if _, err := fmt.Sscanf(code, "%d", &n); err == nil {
		return n
	}
	return 500
}
