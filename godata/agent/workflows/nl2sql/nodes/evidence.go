package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// EvidenceRecallNode retrieves relevant evidence/knowledge for the query.
type EvidenceRecallNode struct{}

// Name returns the unique node identifier.
func (n *EvidenceRecallNode) Name() string {
	return "evidence_recall"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
func (n *EvidenceRecallNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)
	nl2state.EvidenceContexts = append(nl2state.EvidenceContexts, nl2sql.EvidenceContext{
		Content: "Stub evidence - Phase 5",
		Score:   1.0,
		Source:  "stub",
	})
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
