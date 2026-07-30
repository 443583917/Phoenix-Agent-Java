package tools

import (
	"fmt"
	"sync"

	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// ToolRegistry stores function tools with their names and descriptions.
// Thread-safe for concurrent registration and lookup.
type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]tool.Tool
}

// NewToolRegistry creates a new empty tool registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]tool.Tool),
	}
}

// Register adds a tool to the registry keyed by its declaration name.
func (r *ToolRegistry) Register(t tool.Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := t.Declaration().Name
	r.tools[name] = t
}

// Get retrieves a tool by name.
func (r *ToolRegistry) Get(name string) (tool.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool %q not found in registry", name)
	}
	return t, nil
}

// List returns the names of all registered tools.
func (r *ToolRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetAll returns all registered tools as a slice.
func (r *ToolRegistry) GetAll() []tool.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tools := make([]tool.Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// RegisterAll registers multiple tools at once.
func (r *ToolRegistry) RegisterAll(tools []tool.Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range tools {
		name := t.Declaration().Name
		r.tools[name] = t
	}
}
