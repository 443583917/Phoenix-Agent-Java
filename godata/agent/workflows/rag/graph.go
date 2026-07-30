package rag

import "context"

// RagWorkflow represents a RAG retrieval workflow.
// This is a stub that will be filled in with actual workflow logic.
type RagWorkflow struct {
	name string
}

func NewRagWorkflow() *RagWorkflow {
	return &RagWorkflow{name: "rag-workflow"}
}

func (w *RagWorkflow) Name() string {
	return w.name
}

// Run executes the RAG workflow.
// This is a stub — actual implementation will integrate with the agent runtime.
func (w *RagWorkflow) Run(ctx context.Context, query string, opts ...Option) (string, error) {
	return "", nil
}

// Option is a functional option for RagWorkflow.
type Option func(*options)

type options struct {
	TopK        int
	CategoryIDs []string
	FileTypes   []string
}

func WithTopK(k int) Option {
	return func(o *options) {
		o.TopK = k
	}
}

func WithCategoryIDs(ids ...string) Option {
	return func(o *options) {
		o.CategoryIDs = ids
	}
}

func WithFileTypes(types ...string) Option {
	return func(o *options) {
		o.FileTypes = types
	}
}
