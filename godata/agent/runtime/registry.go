package runtime

import (
	"fmt"
	"sync"

	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// AgentConfig holds the static configuration for an agent instance.
// The AgentManager uses this to create fresh tRPC-Agent-Go agent instances
// for each streaming call.
type AgentConfig struct {
	// SN is the unique serial number identifying this agent configuration.
	SN string
	// ModelName is the LLM model name (e.g. "deepseek-chat", "gpt-4o").
	ModelName string
	// BaseURL is the API base URL for the model provider.
	BaseURL string
	// APIKey is the authentication key for the model provider.
	APIKey string
	// Tools is the list of tools available to this agent.
	Tools []tool.Tool
	// Stream enables streaming output mode.
	Stream bool
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int
	// Temperature controls randomness (0.0 to 2.0).
	Temperature float64
	// AgentType indicates the agent type ("react", "assistant", "workflow").
	AgentType string
}

// AgentRegistry stores agent configurations by their serial number.
// Thread-safe for concurrent access.
type AgentRegistry struct {
	mu     sync.RWMutex
	agents map[string]*AgentConfig
}

// NewAgentRegistry creates a new empty agent registry.
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*AgentConfig),
	}
}

// Register adds an agent configuration to the registry.
func (r *AgentRegistry) Register(cfg *AgentConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[cfg.SN] = cfg
}

// Get retrieves an agent configuration by serial number.
func (r *AgentRegistry) Get(sn string) (*AgentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.agents[sn]
	if !ok {
		return nil, fmt.Errorf("agent %q not found in registry", sn)
	}
	return cfg, nil
}

// List returns all registered agent serial numbers.
func (r *AgentRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sns := make([]string, 0, len(r.agents))
	for sn := range r.agents {
		sns = append(sns, sn)
	}
	return sns
}
