package nodes

import (
	"context"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SemanticConsistencyNode validates generated SQL against the user's query
// and schema to ensure semantic correctness.
type SemanticConsistencyNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *SemanticConsistencyNode) Name() string {
	return "semantic_consistency"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It calls the LLM with the semantic_consistency prompt to analyze whether
// the generated SQL accurately reflects the user's intent. The result is
// either "通过" (pass) or "不通过。reason" (fail with reason).
func (n *SemanticConsistencyNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	sql := nl2state.SQLGenerateOutput
	evidence := nl2state.Evidence
	schemaInfo := nl2state.GeneratedSemanticModelPrompt
	if schemaInfo == "" {
		schemaInfo = nl2state.TableRelationOutput
	}

	// If no SQL was generated, validation fails immediately.
	if sql == "" || sql == "-- SQL generation failed, no valid SQL produced" {
		nl2state.SemanticConsistencyOutput = "不通过。SQL为空或生成失败"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state": nl2state,
			nl2sqltypes.StateKeySemanticConsistencyOutput: "不通过。SQL为空或生成失败",
		}, nil
	}

	prompt := GetPrompt("semantic_consistency")
	data := map[string]string{
		"Query":      query,
		"SQL":        sql,
		"Evidence":   evidence,
		"SchemaInfo": schemaInfo,
	}

	var result struct {
		Result string `json:"result"`
		Reason string `json:"reason"`
	}

	validationResult := "通过"
	if err := n.LLM.CallJSON(ctx, prompt, data, &result); err == nil {
		if strings.Contains(result.Result, "不通过") {
			validationResult = "不通过。" + result.Reason
		} else {
			validationResult = "通过"
		}
	}

	nl2state.SemanticConsistencyOutput = validationResult
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
		nl2sqltypes.StateKeySemanticConsistencyOutput: validationResult,
	}, nil
}
