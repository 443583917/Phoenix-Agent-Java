package nodes

import (
	"context"
	"fmt"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// PlannerNode generates an execution plan for the query.
type PlannerNode struct{}

// Name returns the unique node identifier.
func (n *PlannerNode) Name() string {
	return "planner"
}

// Execute creates a stub execution plan. If the plan was previously rejected,
// the graph terminates at this node.
func (n *PlannerNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	if state.RejectedPlan {
		return &nl2sql.NodeOutput{NextNode: "end"}, nil
	}

	state.Plan = &nl2sql.ExecutionPlan{
		Steps: []nl2sql.PlanStep{
			{
				Description: fmt.Sprintf("Execute query for: %s", state.Query),
			},
		},
		Reasoning: "Stub plan - Phase 5",
	}
	return &nl2sql.NodeOutput{NextNode: "sql_generate"}, nil
}
