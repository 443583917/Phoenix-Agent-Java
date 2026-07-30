package runner

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/memory"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/internal/model"
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
// It wraps the AgentManager.StreamCall with session management, short-term
// memory recording, and HITL callback wiring that allows external callers to
// approve or reject tool calls.
type ConversationRunner struct {
	manager        *runtime.AgentManager
	shortTerm      *memory.ShortTermMemory
	hitlCallback   HitlCallback
}

// HitlCallback is invoked when an agent emits a tool call that requires human
// confirmation. The callback should return true to approve or false to reject.
//
// If nil, all tool calls are rejected by default.
type HitlCallback func(ctx context.Context, event model.ToolCallEvent) bool

// NewConversationRunner creates a new conversation runner.
//
// shortTerm and hitlCallback are optional — pass nil to skip short-term
// memory recording or to reject all HITL events.
func NewConversationRunner(
	manager *runtime.AgentManager,
	shortTerm *memory.ShortTermMemory,
	hitlCallback HitlCallback,
) *ConversationRunner {
	return &ConversationRunner{
		manager:      manager,
		shortTerm:    shortTerm,
		hitlCallback: hitlCallback,
	}
}

// Run executes a conversation turn and returns a channel of SSE events.
//
// The user message is recorded in short-term memory before the call.
// Assistant responses are recorded after the stream completes.
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

	// Record user message in short-term memory.
	if r.shortTerm != nil {
		r.shortTerm.AddMessage(req.SessionID, memory.Message{
			Role:    "user",
			Content: req.Message,
		})
	}

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
// the full assistant response for short-term memory.
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

	// Record the full assistant response in short-term memory.
	if r.shortTerm != nil && assistantContent != "" {
		r.shortTerm.AddMessage(sessionID, memory.Message{
			Role:    "assistant",
			Content: assistantContent,
		})
	}
}
