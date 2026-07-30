package nodes

import (
	"context"
	"fmt"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// PlannerNode generates an execution plan for the query.
// Supports repair cycles when a plan is rejected by the validation node.
type PlannerNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *PlannerNode) Name() string {
	return "planner"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// If isOnlyNL2SQL mode is active, a hardcoded single-step plan is generated.
// Otherwise the LLM is called with the planner prompt to produce a multi-step plan.
// When the plan was previously rejected (repair cycle), repair context is included.
func (n *PlannerNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	// Handle repair cycle: check if plan was rejected
	planRepairCount := 0
	if v, ok := state[nl2sqltypes.StateKeyPlanRepairCount].(int); ok {
		planRepairCount = v
	}

	// Check NL2SQL-only mode
	isOnlyNL2SQL := false
	if v, ok := state[nl2sqltypes.StateKeyIsOnlyNL2SQL].(bool); ok {
		isOnlyNL2SQL = v
	}

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	var plan *nl2sqltypes.Plan

	if isOnlyNL2SQL {
		// Hardcoded single-step plan for NL2SQL only mode
		plan = &nl2sqltypes.Plan{
			ThoughtProcess: "NL2SQL only mode: generate SQL directly",
			ExecutionPlan: []nl2sqltypes.ExecutionStep{
				{
					Step:      1,
					ToolToUse: "sql_generate",
					ToolParameters: map[string]string{
						"instruction": "Generate SQL for user query",
					},
				},
			},
		}
	} else {
		// Use LLM to generate plan
		tableRelation := nl2state.TableRelationOutput
		contextStr := getStateStr(state, nl2sqltypes.StateKeyMultiTurnContext)

		// Include repair context if retrying
		var validationError string
		if v, ok := state[nl2sqltypes.StateKeyPlanValidationError].(string); ok {
			validationError = v
		}

		repairInfo := ""
		if planRepairCount > 0 && validationError != "" {
			repairInfo = fmt.Sprintf("\n前一次计划被拒绝，原因：%s\n请重新规划。", validationError)
		}

		prompt := GetPrompt("planner")
		data := map[string]string{
			"Query":         query,
			"TableRelation": tableRelation,
			"Context":       contextStr,
		}

		var llmPlan nl2sqltypes.Plan
		if err := n.LLM.CallJSON(ctx, prompt+repairInfo, data, &llmPlan); err != nil {
			// Fallback: single-step SQL plan
			llmPlan = nl2sqltypes.Plan{
				ThoughtProcess: "Default plan: generate SQL for the query",
				ExecutionPlan: []nl2sqltypes.ExecutionStep{
					{
						Step:      1,
						ToolToUse: "sql_generate",
						ToolParameters: map[string]string{
							"instruction": "Generate SQL for: " + query,
						},
					},
				},
			}
		}

		plan = &llmPlan
	}

	nl2state.PlannerOutput = plan
	nl2state.PlanCurrentStep = 0
	nl2state.PlanRepairCount = planRepairCount
	nl2state.CurrentNode = n.Name()

	updates := graph.State{
		"nl2sql_state":                       nl2state,
		nl2sqltypes.StateKeyPlannerOutput:    plan,
		nl2sqltypes.StateKeyPlanCurrentStep:  0,
	}
	if planRepairCount > 0 {
		updates[nl2sqltypes.StateKeyPlanRepairCount] = planRepairCount
	}

	return updates, nil
}
