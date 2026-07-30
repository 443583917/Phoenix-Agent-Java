package agents

import (
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// NewAssistantAgent creates a simple conversational assistant agent
// without any tools. It uses streaming output by default.
//
// This is the simplest agent type, suitable for general Q&A and
// conversational use cases that don't require tool calling.
func NewAssistantAgent(name string, m model.Model, opts ...llmagent.Option) *llmagent.LLMAgent {
	defaultOpts := []llmagent.Option{
		llmagent.WithModel(m),
		llmagent.WithGenerationConfig(model.GenerationConfig{
			Stream: true,
		}),
	}
	allOpts := append(defaultOpts, opts...)
	return llmagent.New(name, allOpts...)
}
