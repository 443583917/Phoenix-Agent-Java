package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// EvidenceRecallNode retrieves relevant evidence/knowledge for the query.
type EvidenceRecallNode struct{}

// Name returns the unique node identifier.
func (n *EvidenceRecallNode) Name() string {
	return "evidence_recall"
}

// Execute adds a stub evidence context to the state.
func (n *EvidenceRecallNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.EvidenceContexts = append(state.EvidenceContexts, nl2sql.EvidenceContext{
		Content: "Stub evidence - Phase 5",
		Score:   1.0,
		Source:  "stub",
	})
	return &nl2sql.NodeOutput{NextNode: "schema_recall"}, nil
}
