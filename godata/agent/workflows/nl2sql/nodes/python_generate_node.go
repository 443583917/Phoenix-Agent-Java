package nodes

import (
	"context"
	"fmt"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// PythonGenerateNode generates Python data-analysis code based on SQL query
// results by calling the LLM with the python_generate prompt.
type PythonGenerateNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *PythonGenerateNode) Name() string {
	return "python_generate_node"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It builds a prompt from the user query, SQL result, instruction, and
// optional retry info, calls the LLM, and stores the generated Python
// code in the graph state.
func (n *PythonGenerateNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	sqlResult := nl2state.SQLExecuteOutput
	if sqlResult == "" {
		sqlResult = "无SQL查询结果"
	}

	// Get instruction from current plan step.
	instruction := "分析数据并生成可视化图表"
	if nl2state.PlannerOutput != nil && nl2state.PlanCurrentStep > 0 {
		stepIdx := nl2state.PlanCurrentStep - 1
		if stepIdx < len(nl2state.PlannerOutput.ExecutionPlan) {
			step := nl2state.PlannerOutput.ExecutionPlan[stepIdx]
			if params, ok := step.ToolParameters["instruction"]; ok && params != "" {
				instruction = params
			}
		}
	}

	// Retry info if regenerating.
	retryInfo := ""
	if nl2state.PythonTriesCount > 0 {
		retryInfo = fmt.Sprintf("第%d次重试", nl2state.PythonTriesCount+1)
	}

	prompt := GetPrompt("python_generate")
	data := map[string]string{
		"Query":       query,
		"SQLResult":   sqlResult,
		"Instruction": instruction,
		"RetryInfo":   retryInfo,
	}

	var result struct {
		Code        string `json:"code"`
		Explanation string `json:"explanation"`
	}

	code := ""
	if err := n.LLM.CallJSON(ctx, prompt, data, &result); err == nil && result.Code != "" {
		code = result.Code
	}

	if code == "" {
		code = "# Python code generation failed"
	}

	nl2state.PythonGenerateOutput = code
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                      nl2state,
		nl2sqltypes.StateKeyPythonGenerateOutput: code,
	}, nil
}
