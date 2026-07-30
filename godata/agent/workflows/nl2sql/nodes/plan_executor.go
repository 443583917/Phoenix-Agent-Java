package nodes

import (
	"context"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// PlanExecutorNode validates and executes the next step in the execution plan.
// It routes to the correct tool node based on the step's ToolToUse field
// and updates plan_current_step / plan_next_node in state.
// No LLM call is required for this node.
type PlanExecutorNode struct{}

// Name returns the unique node identifier.
func (n *PlanExecutorNode) Name() string {
	return "plan_executor"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It reads the current execution step from the plan, maps ToolToUse to
// the corresponding graph node name, and advances the step counter.
// When all steps are completed, it routes to "report_generate".
func (n *PlanExecutorNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	plan := nl2state.PlannerOutput
	if plan == nil || len(plan.ExecutionPlan) == 0 {
		// No plan, default to sql_generate
		nl2state.PlanNextNode = "sql_generate"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                    nl2state,
			nl2sqltypes.StateKeyPlanNextNode: "sql_generate",
		}, nil
	}

	currentStep := nl2state.PlanCurrentStep
	if currentStep >= len(plan.ExecutionPlan) {
		// All steps completed, go to report
		nl2state.PlanNextNode = "report_generate"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                    nl2state,
			nl2sqltypes.StateKeyPlanNextNode: "report_generate",
		}, nil
	}

	step := plan.ExecutionPlan[currentStep]

	// Map tool_to_use to actual node name
	var nextNode string
	switch step.ToolToUse {
	case "sql_generate":
		nextNode = "sql_generate"
	case "python_execute":
		nextNode = "python_generate_node"
	default:
		nextNode = "sql_generate"
	}

	// Increment step counter
	nl2state.PlanCurrentStep = currentStep + 1
	nl2state.PlanNextNode = nextNode
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                       nl2state,
		nl2sqltypes.StateKeyPlanCurrentStep:  currentStep + 1,
		nl2sqltypes.StateKeyPlanNextNode:     nextNode,
	}, nil
}
