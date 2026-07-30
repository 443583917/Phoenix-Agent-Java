package nodes

import (
	"context"
	"fmt"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// PlannerNode generates an execution plan for the query.
type PlannerNode struct{}

// Name returns the unique node identifier.
func (n *PlannerNode) Name() string {
	return "planner"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// If the plan was previously rejected, the conditional edge router
// (defined in graph.go) will route to "end".
func (n *PlannerNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	nl2state.Plan = &nl2sql.ExecutionPlan{
		Steps: []nl2sql.PlanStep{
			{
				Description: fmt.Sprintf("Execute query for: %s", nl2state.Query),
			},
		},
		Reasoning: "Stub plan - Phase 5",
	}
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
