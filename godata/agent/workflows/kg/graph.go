package kg

import "context"

// KgWorkflow represents a Knowledge Graph query workflow.
// This is a stub that will be filled in with actual workflow logic.
type KgWorkflow struct {
	name string
}

func NewKgWorkflow() *KgWorkflow {
	return &KgWorkflow{name: "kg-workflow"}
}

func (w *KgWorkflow) Name() string {
	return w.name
}

// Run executes the KG query workflow.
// This is a stub — actual implementation will traverse the knowledge graph
// to find relevant entities and relationships based on the query.
func (w *KgWorkflow) Run(ctx context.Context, query string, opts ...Option) (string, error) {
	return "", nil
}

// Option is a functional option for KgWorkflow.
type Option func(*options)

type options struct {
	MaxDepth      int
	EntityTypes   []string
	RelationTypes []string
}

func WithMaxDepth(d int) Option {
	return func(o *options) {
		o.MaxDepth = d
	}
}

func WithEntityTypes(types ...string) Option {
	return func(o *options) {
		o.EntityTypes = types
	}
}

func WithRelationTypes(types ...string) Option {
	return func(o *options) {
		o.RelationTypes = types
	}
}
