package hooks

import (
	"context"
	"sync"

	"go.uber.org/zap"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

// ModelCallLimitHook tracks the number of LLM calls made during a single
// agent run and forces a final output when the limit is exceeded.
//
// When the limit is reached, subsequent BeforeModel calls inject a system
// message instructing the model to produce a final summary instead of
// continuing the conversation.
type ModelCallLimitHook struct {
	mu      sync.Mutex
	counts  map[string]int // keyed by sessionID extracted from context
	maxCall int
}

// NewModelCallLimitHook creates a new ModelCallLimitHook.
// maxCall is the maximum number of LLM calls allowed per agent run (default 15).
func NewModelCallLimitHook(maxCall int) *ModelCallLimitHook {
	if maxCall <= 0 {
		maxCall = 15
	}
	return &ModelCallLimitHook{
		counts:  make(map[string]int),
		maxCall: maxCall,
	}
}

// BeforeModel is the tRPC-Agent-Go BeforeModelCallbackStructured callback.
// It increments the call counter for the given session and, if the limit is
// exceeded, injects a termination system message to force a final output.
func (h *ModelCallLimitHook) BeforeModel(
	ctx context.Context,
	args *tmodel.BeforeModelArgs,
) (*tmodel.BeforeModelResult, error) {
	// Extract session identifier from context.
	sessionID, _ := ctx.Value("sessionID").(string)
	if sessionID == "" {
		return nil, nil
	}

	h.mu.Lock()
	count := h.counts[sessionID]
	count++
	h.counts[sessionID] = count
	h.mu.Unlock()

	if count <= h.maxCall {
		return nil, nil
	}

	// Limit exceeded: inject a termination prompt.
	zap.L().Warn("model call limit reached, forcing final output",
		zap.String("sessionID", sessionID),
		zap.Int("count", count),
		zap.Int("maxCall", h.maxCall),
	)

	terminationMsg := tmodel.NewSystemMessage(
		"你已经进行了太多次调用，请基于已有的信息直接给出最终答案，不要再调用工具或提出更多问题。",
	)
	args.Request.Messages = append(args.Request.Messages, terminationMsg)

	return nil, nil
}

// Reset clears the call counter for a given session.
func (h *ModelCallLimitHook) Reset(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.counts, sessionID)
}

// Count returns the current call count for a session.
func (h *ModelCallLimitHook) Count(sessionID string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.counts[sessionID]
}
