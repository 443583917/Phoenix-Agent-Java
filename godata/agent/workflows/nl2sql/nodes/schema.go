package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// SchemaRecallNode retrieves relevant schema information for the query.
type SchemaRecallNode struct{}

// Name returns the unique node identifier.
func (n *SchemaRecallNode) Name() string {
	return "schema_recall"
}

// Execute adds a stub schema context to the state.
func (n *SchemaRecallNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.SchemaContext = append(state.SchemaContext, nl2sql.SchemaContext{
		TableName:    "stub_table",
		ColumnName:   "stub_column",
		DataType:     "VARCHAR",
		BusinessName: "Stub Business Name",
		Description:  "Stub schema - Phase 5",
	})
	return &nl2sql.NodeOutput{NextNode: "planner"}, nil
}
