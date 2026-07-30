package runtime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/phoenix-agent-go/agent/hooks"
	"github.com/phoenix-agent-go/agent/interceptors"
	"github.com/phoenix-agent-go/agent/memory"
	"github.com/phoenix-agent-go/internal/model"

	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// AgentManager manages the lifecycle of tRPC-Agent-Go agents and provides
// streaming call functionality via the framework runner.
//
// Each StreamCall creates a fresh tRPC-Agent-Go agent and runner instance
// for request isolation. The session service is passed to the framework
// Runner so conversation history is automatically managed.
//
// The manager also wires agent hooks (profile injection, call limiting,
// summarization), an interceptor (login dedup + memory injection), and
// the async memory pipeline for post-turn memory extraction.
type AgentManager struct {
	registry          *AgentRegistry
	sessSvc           *inmemory.SessionService
	memoryPipeline    *memory.MemoryPipeline
	profileHook       *hooks.ProfileInjectionHook
	limitHook         *hooks.ModelCallLimitHook
	summarizationHook *hooks.SummarizationHook
	loginInterceptor  *interceptors.LoginUserAgentInterceptor
}

// AgentManagerOption configures an AgentManager.
type AgentManagerOption func(*AgentManager)

// WithMemoryPipeline sets the memory pipeline for async memory extraction.
func WithMemoryPipeline(p *memory.MemoryPipeline) AgentManagerOption {
	return func(m *AgentManager) {
		m.memoryPipeline = p
	}
}

// WithProfileHook sets the profile injection hook.
func WithProfileHook(h *hooks.ProfileInjectionHook) AgentManagerOption {
	return func(m *AgentManager) {
		m.profileHook = h
	}
}

// WithLimitHook sets the model call limit hook.
func WithLimitHook(h *hooks.ModelCallLimitHook) AgentManagerOption {
	return func(m *AgentManager) {
		m.limitHook = h
	}
}

// WithSummarizationHook sets the summarization hook.
func WithSummarizationHook(h *hooks.SummarizationHook) AgentManagerOption {
	return func(m *AgentManager) {
		m.summarizationHook = h
	}
}

// WithLoginInterceptor sets the login interceptor.
func WithLoginInterceptor(i *interceptors.LoginUserAgentInterceptor) AgentManagerOption {
	return func(m *AgentManager) {
		m.loginInterceptor = i
	}
}

// NewAgentManager creates a new AgentManager with the given agent registry.
// sessSvc is the tRPC-Agent-Go session service used for persisting
// conversation events across runs.
func NewAgentManager(registry *AgentRegistry, sessSvc *inmemory.SessionService, opts ...AgentManagerOption) *AgentManager {
	m := &AgentManager{registry: registry, sessSvc: sessSvc}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Registry returns the underlying agent registry for external configuration.
func (m *AgentManager) Registry() *AgentRegistry {
	return m.registry
}

// GetSessionService returns the session service used by this manager.
func (m *AgentManager) GetSessionService() *inmemory.SessionService {
	return m.sessSvc
}

// StreamCall executes a streaming agent call and returns a channel of SSE events.
//
// It creates a fresh tRPC-Agent-Go agent and runner for each call, sends the
// user message, and converts runner events into SSEEvent values on the returned
// channel. The channel is closed when the run completes.
//
// The session service is passed to the framework Runner via
// runner.WithSessionService so events are automatically persisted.
//
// After the run completes, the memory pipeline is invoked asynchronously
// (fire-and-forget) to extract and persist user memories.
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

	// Build model callbacks from hooks.
	modelCallbacks := m.buildModelCallbacks()

	// Build llmagent with tools, callbacks, and hooks if configured.
	agentOpts := []llmagent.Option{
		llmagent.WithModel(modelInstance),
		llmagent.WithGenerationConfig(genConfig),
	}
	if len(agentCfg.Tools) > 0 {
		agentOpts = append(agentOpts, llmagent.WithTools(agentCfg.Tools))
	}
	if modelCallbacks != nil {
		agentOpts = append(agentOpts, llmagent.WithModelCallbacks(modelCallbacks))
	}
	llmAgent := llmagent.New(req.AgentSN, agentOpts...)

	// Enrich context with userID, sessionID, and agentSN for hooks/interceptors.
	ctx = context.WithValue(ctx, "userID", req.UserID)
	ctx = context.WithValue(ctx, "sessionID", req.SessionID)
	ctx = context.WithValue(ctx, "agentSN", req.AgentSN)

	// Apply login interceptor: inject history memories into the user message.
	message := req.Message
	if m.loginInterceptor != nil {
		// Asynchronously record agent usage.
		go m.loginInterceptor.RecordUsage(req.UserID, req.AgentSN)

		// Inject relevant history memories if available.
		if historyText := m.loginInterceptor.InjectHistoryMemories(ctx, req.UserID, req.Message); historyText != "" {
			message = historyText + "\n\n【当前问题】\n" + req.Message
		}
	}

	// Build runner options with session service for automatic event persistence.
	runnerOpts := []runner.Option{}
	if m.sessSvc != nil {
		runnerOpts = append(runnerOpts, runner.WithSessionService(m.sessSvc))
	}

	// Create runner and execute.
	r := runner.NewRunner(req.AgentSN+"-app", llmAgent, runnerOpts...)
	events, err := r.Run(
		ctx, req.UserID, req.SessionID, tmodel.NewUserMessage(message),
	)
	if err != nil {
		return nil, fmt.Errorf("agent run: %w", err)
	}

	// Convert runner event stream to SSE event stream.
	ch := make(chan model.SSEEvent, 10)
	go func() {
		defer close(ch)
		var userQuery string
		for evt := range events {
			sseEvent := convertRunnerEvent(evt)

			// Accumulate user query or assistant response for memory extraction.
			if evt.Response != nil && len(evt.Response.Choices) > 0 {
				choice := evt.Response.Choices[0]
				if choice.Delta.Content != "" {
					userQuery += choice.Delta.Content
				}
			}

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

		// Fire-and-forget: trigger async memory extraction.
		if m.memoryPipeline != nil {
			go m.memoryPipeline.ProcessAndExtractMemory(
				context.Background(),
				req.UserID,
				req.AgentSN,
				req.Message,
			)
		}

		// Reset hook state for this session after the run.
		if m.limitHook != nil {
			m.limitHook.Reset(req.SessionID)
		}
		if m.summarizationHook != nil {
			m.summarizationHook.Reset(req.SessionID)
		}
	}()

	return ch, nil
}

// buildModelCallbacks creates tRPC-Agent-Go model callbacks from the configured hooks.
// Returns nil if no hooks are configured.
func (m *AgentManager) buildModelCallbacks() *tmodel.Callbacks {
	hasHooks := m.profileHook != nil || m.limitHook != nil || m.summarizationHook != nil
	if !hasHooks {
		return nil
	}

	callbacks := tmodel.NewCallbacks()

	if m.profileHook != nil {
		callbacks.RegisterBeforeModel(m.profileHook.BeforeModel)
	}
	if m.limitHook != nil {
		callbacks.RegisterBeforeModel(m.limitHook.BeforeModel)
	}
	if m.summarizationHook != nil {
		callbacks.RegisterBeforeModel(m.summarizationHook.BeforeModel)
	}

	return callbacks
}

// convertRunnerEvent converts a tRPC-Agent-Go runner Event to an SSEEvent.
//
// Handles:
//   - Chat completion chunks with text content -> ContentEvent
//   - Tool call events -> ToolCallEvent
//   - Error responses -> ErrorEvent
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
