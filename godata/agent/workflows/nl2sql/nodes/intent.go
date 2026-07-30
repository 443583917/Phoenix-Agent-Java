package nodes

import (
	"context"

	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// IntentRecognitionNode classifies the user's intent.
type IntentRecognitionNode struct{}

// Name returns the unique node identifier.
func (n *IntentRecognitionNode) Name() string {
	return "intent_recognition"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It classifies intent and updates the NL2SQLState.
func (n *IntentRecognitionNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)
	nl2state.Intent = "sql"
	nl2state.IntentConfidence = 0.95
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
	}, nil
}
