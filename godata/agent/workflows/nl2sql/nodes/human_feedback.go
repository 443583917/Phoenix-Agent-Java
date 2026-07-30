package nodes

import (
	"context"
	"fmt"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type HumanFeedbackNode struct{}

func (n *HumanFeedbackNode) Name() string {
	return "human_feedback"
}

func (n *HumanFeedbackNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	// Check if human review is enabled
	humanReviewEnabled := nl2state.HumanReviewEnabled

	if !humanReviewEnabled {
		// Skip human feedback - not enabled
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state": nl2state,
		}, nil
	}

	// Check if feedback data is already available (resume from interrupt)
	feedbackData := nl2state.HumanFeedbackData
	if feedbackData != "" {
		// Feedback already provided, continue
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                    nl2state,
			nl2sqltypes.StateKeyHumanFeedbackData: feedbackData,
		}, nil
	}

	// No feedback yet, interrupt execution
	// Build context for the human
	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	interruptMsg := fmt.Sprintf("需要人工确认:\n用户请求: %s\n请确认是否继续执行?", query)

	// Return a Command with interrupt signal
	return &graph.Command{
		GoTo: graph.End,
		Update: graph.State{
			"nl2sql_state":                    nl2state,
			nl2sqltypes.StateKeyHumanFeedbackData: "",
		},
		Resume: interruptMsg,
	}, nil
}
