package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// ReportGeneratorNode generates the final HTML report.
type ReportGeneratorNode struct{}

// Name returns the unique node identifier.
func (n *ReportGeneratorNode) Name() string {
	return "report_generate"
}

// Execute sets a stub HTML report.
func (n *ReportGeneratorNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.ReportContent = "<html><body><h1>NL2SQL Report Stub</h1><p>Phase 5 placeholder report</p></body></html>"
	return &nl2sql.NodeOutput{NextNode: "end"}, nil
}
