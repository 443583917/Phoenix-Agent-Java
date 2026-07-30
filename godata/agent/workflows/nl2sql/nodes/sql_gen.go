package nodes

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SqlGenerateNode generates SQL from the execution plan.
type SqlGenerateNode struct{}

// Name returns the unique node identifier.
func (n *SqlGenerateNode) Name() string {
	return "sql_generate"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
func (n *SqlGenerateNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)
	nl2state.SQLQuery = "SELECT 'stub - Phase 5 NL2SQL'"
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
