package nodes

import (
	"context"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// FeasibilityAssessmentNode classifies the user query as one of:
// "数据分析", "需要澄清", or "自由闲聊".
type FeasibilityAssessmentNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *FeasibilityAssessmentNode) Name() string {
	return "feasibility_assessment"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It calls the LLM with the feasibility_assessment prompt to classify
// the user's request and stores the result in state.
func (n *FeasibilityAssessmentNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	schemaInfo := nl2state.TableRelationOutput
	if schemaInfo == "" {
		schemaInfo = "无表结构信息"
	}

	prompt := GetPrompt("feasibility_assessment")
	data := map[string]string{
		"Query":      query,
		"SchemaInfo": schemaInfo,
	}

	var result struct {
		Classification string `json:"classification"`
		Reasoning      string `json:"reasoning"`
	}

	if err := n.LLM.CallJSON(ctx, prompt, data, &result); err != nil {
		result.Classification = "数据分析"
		result.Reasoning = "默认处理为数据分析请求"
	}

	if result.Classification == "" {
		result.Classification = "数据分析"
	}

	classification := strings.TrimSpace(result.Classification)

	nl2state.FeasibilityOutput = classification
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                        nl2state,
		nl2sqltypes.StateKeyFeasibilityOutput: classification,
	}, nil
}
