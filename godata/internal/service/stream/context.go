package stream

import (
	"context"
	"strings"
	"sync"
	"time"
)

type Context struct {
	ThreadID        string
	AgentID         string
	StartTime       time.Time
	CollectedOutput strings.Builder
	TextType        string
	CancelFunc      context.CancelFunc
}

type ContextManager struct {
	mu       sync.RWMutex
	contexts map[string]*Context
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[string]*Context),
	}
}

func (m *ContextManager) Register(threadID, agentID string, cancel context.CancelFunc) *Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx := &Context{
		ThreadID:   threadID,
		AgentID:    agentID,
		StartTime:  time.Now(),
		CancelFunc: cancel,
	}
	m.contexts[threadID] = ctx
	return ctx
}

func (m *ContextManager) Get(threadID string) *Context {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.contexts[threadID]
}

func (m *ContextManager) Unregister(threadID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ctx, ok := m.contexts[threadID]; ok {
		if ctx.CancelFunc != nil {
			ctx.CancelFunc()
		}
		delete(m.contexts, threadID)
	}
}

func (m *ContextManager) List() []*Context {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Context, 0, len(m.contexts))
	for _, ctx := range m.contexts {
		result = append(result, ctx)
	}
	return result
}
