package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// SqlGenerateNode generates SQL from the execution plan.
type SqlGenerateNode struct{}

// Name returns the unique node identifier.
func (n *SqlGenerateNode) Name() string {
	return "sql_generate"
}

// Execute sets a stub SQL query. In NL2SQL-only mode, skips Python execution
// and proceeds directly to report generation.
func (n *SqlGenerateNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.SQLQuery = "SELECT 'stub - Phase 5 NL2SQL'"

	if state.NL2SQLOnly {
		return &nl2sql.NodeOutput{NextNode: "report_generate"}, nil
	}
	return &nl2sql.NodeOutput{NextNode: "python_execute"}, nil
}
