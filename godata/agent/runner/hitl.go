package runner

import (
	"context"
	"fmt"
	"sync"
)

// HitlHandler manages Human-In-The-Loop confirm/reject state for tool calls
// emitted by agents.
//
// When an agent emits a tool_call event, the caller should register the
// pending confirmation via RegisterPending and await external approval
// (e.g. from a UI button click) via HandleConfirm.
//
// NOTE: The tRPC-Agent-Go framework provides tool callbacks via
// tool.Callbacks and model.Callbacks. For server-side HITL patterns,
// the framework's interrupt/resume mechanism (graph.Interrupt) or
// BeforeTool callbacks can intercept tool calls before execution.
// This custom handler complements those patterns for Phoenix's
// web-based confirmation flow where approvals arrive asynchronously
// from the UI via HTTP endpoints.
//
// Framework integration: When an agent is configured with tool
// callbacks, the BeforeTool callback can call RegisterPending to
// pause execution until the user confirms via HandleConfirm:
//
//	toolCallbacks := tool.NewCallbacks()
//	toolCallbacks.RegisterBeforeTool(func(ctx, args) (any, error) {
//	    approved := <-handler.RegisterPending(confirmID, sessionID)
//	    if !approved { return nil, errors.New("rejected by user") }
//	    return nil, nil
//	})
type HitlHandler struct {
	mu      sync.RWMutex
	pending map[string]*hitlEntry
}

type hitlEntry struct {
	sessionID string
	approved  chan bool // receives true (approved) or false (rejected)
}

// NewHitlHandler creates a new HITL handler.
func NewHitlHandler() *HitlHandler {
	return &HitlHandler{
		pending: make(map[string]*hitlEntry),
	}
}

// RegisterPending registers a tool call as awaiting human confirmation.
// Returns a channel that receives true when approved or false when rejected.
//
// The confirmID should be unique per tool call (e.g. a UUID).
func (h *HitlHandler) RegisterPending(confirmID, sessionID string) <-chan bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	approved := make(chan bool, 1)
	h.pending[confirmID] = &hitlEntry{
		sessionID: sessionID,
		approved:  approved,
	}
	return approved
}

// HandleConfirm resolves a pending tool confirmation.
//
// It looks up the pending entry by confirmID and sends the allowed value
// on the approval channel. Returns an error if no matching pending entry
// is found.
func (h *HitlHandler) HandleConfirm(ctx context.Context, confirmID string, allowed bool) error {
	h.mu.Lock()
	entry, ok := h.pending[confirmID]
	if !ok {
		h.mu.Unlock()
		return fmt.Errorf("no pending confirmation for confirmID %q", confirmID)
	}
	delete(h.pending, confirmID)
	h.mu.Unlock()

	select {
	case entry.approved <- allowed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PendingCount returns the number of tool calls awaiting confirmation.
func (h *HitlHandler) PendingCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.pending)
}

// PendingIDs returns all confirmIDs currently awaiting confirmation.
func (h *HitlHandler) PendingIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]string, 0, len(h.pending))
	for id := range h.pending {
		ids = append(ids, id)
	}
	return ids
}
