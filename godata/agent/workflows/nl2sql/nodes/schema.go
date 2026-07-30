package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SchemaRecallNode retrieves relevant schema information for the query.
type SchemaRecallNode struct{}

// Name returns the unique node identifier.
func (n *SchemaRecallNode) Name() string {
	return "schema_recall"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
func (n *SchemaRecallNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)
	nl2state.SchemaContext = append(nl2state.SchemaContext, nl2sql.SchemaContext{
		TableName:    "stub_table",
		ColumnName:   "stub_column",
		DataType:     "VARCHAR",
		BusinessName: "Stub Business Name",
		Description:  "Stub schema - Phase 5",
	})
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
