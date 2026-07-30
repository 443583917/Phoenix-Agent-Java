package nodes

import (
	"context"
	"fmt"

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

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	pyOutput := nl2state.PythonExecuteOutput
	sqlResult := nl2state.SQLExecuteOutput
	sql := nl2state.SQLGenerateOutput

	if !nl2state.PythonIsSuccess {
		nl2state.PythonAnalyzeOutput = "Python执行失败，跳过分析"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                       nl2state,
			nl2sqltypes.StateKeyPythonAnalyzeOutput: "Python执行失败，跳过分析",
		}, nil
	}

	prompt := fmt.Sprintf(`你是一个数据分析结果解读助手。请根据以下信息分析数据并给出结论。

用户需求：%s

SQL查询：%s

SQL查询结果：%s

Python执行输出：
%s

请分析这些结果，输出JSON格式：
{
  "analysis": "数据分析结论，包括关键发现、趋势、异常等",
  "key_insights": ["洞察1", "洞察2", "洞察3"],
  "recommendation": "基于数据的建议"
}`, query, sql, sqlResult, pyOutput)

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
