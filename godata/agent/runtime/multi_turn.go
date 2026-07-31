package runtime

import (
	"fmt"
	"strings"
	"sync"
)

type MultiTurnContextManager struct {
	mu       sync.RWMutex
	contexts map[string]*TurnContext
	maxTurns int
}

type TurnContext struct {
	ThreadID      string
	History       []TurnEntry
	PendingQuery  string
	PendingChunks strings.Builder
}

type TurnEntry struct {
	Query    string
	Response string
}

func NewMultiTurnContextManager(maxTurns int) *MultiTurnContextManager {
	if maxTurns <= 0 {
		maxTurns = 5
	}
	return &MultiTurnContextManager{
		contexts: make(map[string]*TurnContext),
		maxTurns: maxTurns,
	}
}

func (m *MultiTurnContextManager) BeginTurn(threadID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, ok := m.contexts[threadID]
	if !ok {
		ctx = &TurnContext{ThreadID: threadID}
		m.contexts[threadID] = ctx
	}
	ctx.PendingQuery = ""
	ctx.PendingChunks.Reset()
}

func (m *MultiTurnContextManager) FinishTurn(threadID, query, response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, ok := m.contexts[threadID]
	if !ok {
		return
	}
	ctx.History = append(ctx.History, TurnEntry{Query: query, Response: response})
	if len(ctx.History) > m.maxTurns {
		ctx.History = ctx.History[len(ctx.History)-m.maxTurns:]
	}
	ctx.PendingQuery = ""
	ctx.PendingChunks.Reset()
}

func (m *MultiTurnContextManager) DiscardPending(threadID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, ok := m.contexts[threadID]
	if !ok {
		return
	}
	ctx.PendingQuery = ""
	ctx.PendingChunks.Reset()
}

func (m *MultiTurnContextManager) RestartLastTurn(threadID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, ok := m.contexts[threadID]
	if !ok || len(ctx.History) == 0 {
		return
	}
	last := ctx.History[len(ctx.History)-1]
	ctx.PendingQuery = last.Query
	ctx.History = ctx.History[:len(ctx.History)-1]
}

func (m *MultiTurnContextManager) BuildContext(threadID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ctx, ok := m.contexts[threadID]
	if !ok || len(ctx.History) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, entry := range ctx.History {
		sb.WriteString(fmt.Sprintf("[第%d轮] 用户: %s\n助手: %s\n", i+1, entry.Query, entry.Response))
	}
	return sb.String()
}

func (m *MultiTurnContextManager) AppendPlannerChunk(threadID, chunk string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, ok := m.contexts[threadID]
	if !ok {
		return
	}
	ctx.PendingChunks.WriteString(chunk)
}
