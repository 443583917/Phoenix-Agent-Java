package agents

import (
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

// ReactAgentBuilder builds a ReAct-pattern (Reasoning + Acting) LLM agent
// with tool-calling capabilities and streaming enabled by default.
//
// Usage:
//
//	builder := NewReactAgentBuilder("my-agent", modelInstance).
//	    WithTools(tools).
//	    WithMaxTokens(4096)
//	agent := builder.Build()
type ReactAgentBuilder struct {
	name        string
	agentModel  model.Model
	tools       []tool.Tool
	genConfig   model.GenerationConfig
}

// NewReactAgentBuilder creates a new builder for a ReAct agent.
// Streaming is enabled by default.
func NewReactAgentBuilder(name string, m model.Model) *ReactAgentBuilder {
	return &ReactAgentBuilder{
		name:       name,
		agentModel: m,
		genConfig: model.GenerationConfig{
			Stream: true,
		},
	}
}

// WithTools sets the tools available to the agent.
func (b *ReactAgentBuilder) WithTools(tools []tool.Tool) *ReactAgentBuilder {
	b.tools = tools
	return b
}

// WithGenerationConfig sets a custom generation configuration.
func (b *ReactAgentBuilder) WithGenerationConfig(cfg model.GenerationConfig) *ReactAgentBuilder {
	b.genConfig = cfg
	return b
}

// WithMaxTokens sets the maximum number of tokens to generate.
func (b *ReactAgentBuilder) WithMaxTokens(n int) *ReactAgentBuilder {
	b.genConfig.MaxTokens = &n
	return b
}

// WithTemperature sets the sampling temperature.
func (b *ReactAgentBuilder) WithTemperature(t float64) *ReactAgentBuilder {
	b.genConfig.Temperature = &t
	return b
}

// Build creates the llmagent.LLMAgent instance with all configured options.
func (b *ReactAgentBuilder) Build() *llmagent.LLMAgent {
	opts := []llmagent.Option{
		llmagent.WithModel(b.agentModel),
		llmagent.WithGenerationConfig(b.genConfig),
	}
	if len(b.tools) > 0 {
		opts = append(opts, llmagent.WithTools(b.tools))
	}
	return llmagent.New(b.name, opts...)
}
