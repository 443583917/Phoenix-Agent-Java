package nodes

import (
	"context"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// QueryEnhanceNode enhances the user's query for better retrieval and SQL generation.
type QueryEnhanceNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *QueryEnhanceNode) Name() string {
	return "query_enhance"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It calls the LLM to enhance/rewrite the query and stores the result in state.
func (n *QueryEnhanceNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}
	rewrittenQuery := nl2state.Evidence
	if rewrittenQuery == "" {
		rewrittenQuery = query
	}

	prompt := GetPrompt("query_enhancement")
	data := map[string]string{
		"Query":          query,
		"RewrittenQuery": rewrittenQuery,
	}

	var output nl2sqltypes.QueryEnhanceOutput
	if err := n.LLM.CallJSON(ctx, prompt, data, &output); err != nil {
		output = nl2sqltypes.QueryEnhanceOutput{
			CanonicalQuery: query,
		}
	}

	if output.CanonicalQuery == "" {
		output.CanonicalQuery = query
	}

	nl2state.QueryEnhanceOutput = &output
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                        nl2state,
		nl2sqltypes.StateKeyQueryEnhanceOutput: &output,
		"user_input":                          output.CanonicalQuery,
	}, nil
}
