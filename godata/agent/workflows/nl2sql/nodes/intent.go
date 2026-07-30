package nodes

import (
	"context"

	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// IntentRecognitionNode classifies the user's intent.
type IntentRecognitionNode struct{}

// Name returns the unique node identifier.
func (n *IntentRecognitionNode) Name() string {
	return "intent_recognition"
}

// Execute classifies intent and sets stub values.
func (n *IntentRecognitionNode) Execute(ctx context.Context, state *nl2sql.NL2SQLState) (*nl2sql.NodeOutput, error) {
	state.Intent = "sql"
	state.IntentConfidence = 0.95
	return &nl2sql.NodeOutput{NextNode: "evidence_recall"}, nil
}
