package nodes

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// ReportGeneratorNode generates the final HTML report.
type ReportGeneratorNode struct{}

// Name returns the unique node identifier.
func (n *ReportGeneratorNode) Name() string {
	return "report_generate"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
func (n *ReportGeneratorNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)
	nl2state.ReportContent = "<html><body><h1>NL2SQL Report Stub</h1><p>Phase 5 placeholder report</p></body></html>"
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
