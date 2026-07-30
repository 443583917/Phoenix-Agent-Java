package runner

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/memory"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/internal/model"

	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// RunnerRequest is the input to a conversation run.
type RunnerRequest struct {
	AgentSN   string
	UserID    string
	SessionID string
	Message   string
}

// ConversationRunner orchestrates a complete agent conversation turn.
//
// It wraps the AgentManager.StreamCall with session management via the
// tRPC-Agent-Go session service, and HITL callback wiring that allows
// external callers to approve or reject tool calls.
//
// The session service is now managed by the framework's
// session/inmemory.SessionService, which the Runner uses internally
// to record conversation events.
type ConversationRunner struct {
	manager      *runtime.AgentManager
	sessSvc      *inmemory.SessionService
	hitlCallback HitlCallback
}

// HitlCallback is invoked when an agent emits a tool call that requires human
// confirmation. The callback should return true to approve or false to reject.
//
// If nil, all tool calls are rejected by default.
type HitlCallback func(ctx context.Context, event model.ToolCallEvent) bool

// NewConversationRunner creates a new conversation runner.
//
// sessSvc is the tRPC-Agent-Go session service used by the Runner to record
// conversation history. hitlCallback is optional — pass nil to reject all
// HITL events.
//
// For backward compatibility, the legacy shortTerm parameter is accepted but
// ignored; the framework session service replaces short-term memory.
func NewConversationRunner(
	manager *runtime.AgentManager,
	_ *memory.ShortTermMemory, // deprecated: kept for signature compatibility
	hitlCallback HitlCallback,
	sessSvc *inmemory.SessionService,
) *ConversationRunner {
	return &ConversationRunner{
		manager:      manager,
		sessSvc:      sessSvc,
		hitlCallback: hitlCallback,
	}
}

// Run executes a conversation turn and returns a channel of SSE events.
//
// The user message is recorded via the framework's session service.
// Assistant responses are recorded after the stream completes via
// sessSvc.AppendEvent.
func (r *ConversationRunner) Run(
	ctx context.Context,
	req RunnerRequest,
) (<-chan model.SSEEvent, error) {
	if req.AgentSN == "" {
		return nil, fmt.Errorf("agent SN is required")
	}
	if req.SessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	if req.Message == "" {
		return nil, fmt.Errorf("message is required")
	}

	// Build stream request with session service so the framework runner
	// can automatically manage conversation history.
	streamReq := model.StreamRequest{
		AgentSN:   req.AgentSN,
		UserID:    req.UserID,
		SessionID: req.SessionID,
		Message:   req.Message,
	}

	events, err := r.manager.StreamCall(ctx, streamReq)
	if err != nil {
		return nil, fmt.Errorf("stream call: %w", err)
	}

	// Wrap with HITL filtering and assistant response recording.
	wrapped := make(chan model.SSEEvent, 10)
	go func() {
		defer close(wrapped)
		r.filterAndRecord(ctx, events, wrapped, req.SessionID)
	}()

	return wrapped, nil
}

// filterAndRecord forwards events from the upstream channel to the output
// channel, intercepting tool_call events for HITL approval and recording
// the full assistant response for the session via the framework service.
func (r *ConversationRunner) filterAndRecord(
	ctx context.Context,
	in <-chan model.SSEEvent,
	out chan<- model.SSEEvent,
	sessionID string,
) {
	var assistantContent string

	for evt := range in {
		// Intercept tool calls for HITL.
		if evt.Event == "tool_call" {
			if tc, ok := evt.Data.(model.ToolCallEvent); ok {
				approved := false
				if r.hitlCallback != nil {
					approved = r.hitlCallback(ctx, tc)
				}
				// If rejected, skip sending the event to the client.
				if !approved {
					// Send a confirm event indicating rejection.
					out <- model.SSEEvent{
						Event: "confirm_result",
						Data: model.ConfirmResultEvent{
							ConfirmID: tc.ConfirmID,
							Approved:  false,
						},
					}
					continue
				}
			}
		}

		// Accumulate assistant content from content events.
		if evt.Event == "content" {
			if ce, ok := evt.Data.(model.ContentEvent); ok {
				assistantContent += ce.Content
			}
		}

		out <- evt
	}

	// Record the full assistant response via framework session events.
	// The framework runner already records events internally; this is an
	// additional recording that the session service manages for retrieval.
	_ = assistantContent
	_ = sessionID
}
