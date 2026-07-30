package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// PythonExecuteNode executes Python code for data processing.
type PythonExecuteNode struct{}

// Name returns the unique node identifier.
func (n *PythonExecuteNode) Name() string {
	return "python_execute"
}

// Execute sets a stub execution result.
func (n *PythonExecuteNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.ExecutionResult = map[string]interface{}{
		"stub":    true,
		"message": "Python execution stub - Phase 5",
	}
	return &nl2sql.NodeOutput{NextNode: "report_generate"}, nil
}
