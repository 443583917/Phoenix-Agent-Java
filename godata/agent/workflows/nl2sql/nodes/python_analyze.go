package nodes

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/workflows/nl2sql/prompts"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type PythonAnalyzeNode struct {
	LLM *LLMService
}

func (n *PythonAnalyzeNode) Name() string {
	return "python_analyze"
}

func (n *PythonAnalyzeNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	pyOutput := nl2state.PythonExecuteOutput
	sqlResult := nl2state.SQLExecuteOutput

	if !nl2state.PythonIsSuccess {
		nl2state.PythonAnalyzeOutput = "Python执行失败，跳过分析"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                       nl2state,
			nl2sqltypes.StateKeyPythonAnalyzeOutput: "Python执行失败，跳过分析",
		}, nil
	}

	prompt := fmt.Sprintf(prompts.PythonAnalyzePrompt, sqlResult, pyOutput)

	var result struct {
		Analysis       string   `json:"analysis"`
		KeyInsights    []string `json:"key_insights"`
		Recommendation string   `json:"recommendation"`
	}

	analysis := ""
	if err := n.LLM.CallJSON(ctx, prompt, map[string]string{}, &result); err == nil && result.Analysis != "" {
		analysis = result.Analysis
		if result.Recommendation != "" {
			analysis += "\n\n建议：" + result.Recommendation
		}
	}

	if analysis == "" {
		analysis = fmt.Sprintf("Python执行结果:\n%s", pyOutput)
	}

	nl2state.PythonAnalyzeOutput = analysis
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                       nl2state,
		nl2sqltypes.StateKeyPythonAnalyzeOutput: analysis,
	}, nil
}
