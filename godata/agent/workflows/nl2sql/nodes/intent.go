package nodes

import (
	"context"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// IntentRecognitionNode classifies the user's intent.
type IntentRecognitionNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *IntentRecognitionNode) Name() string {
	return "intent_recognition"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It classifies intent by calling the LLM and updates the NL2SQLState.
func (n *IntentRecognitionNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}
	contextStr := getStateStr(state, nl2sqltypes.StateKeyMultiTurnContext)

	prompt := GetPrompt("intent_recognition")
	data := map[string]string{
		"Query":   query,
		"Context": contextStr,
	}

	var output nl2sqltypes.IntentRecognitionOutput
	if err := n.LLM.CallJSON(ctx, prompt, data, &output); err != nil {
		// Default to data analysis on error
		output = nl2sqltypes.IntentRecognitionOutput{
			Classification: "可能的数据分析请求",
			Confidence:     0.5,
			Reasoning:      "LLM call failed, defaulting to data analysis",
		}
	}

	if output.Classification == "" {
		output.Classification = "可能的数据分析请求"
	}

	nl2state.IntentOutput = &output
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                   nl2state,
		nl2sqltypes.StateKeyIntentOutput: nl2state.IntentOutput,
	}, nil
}
