package nodes

import (
	"context"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type ReportGeneratorNode struct {
	LLM *LLMService
}

func (n *ReportGeneratorNode) Name() string {
	return "report_generate"
}

func (n *ReportGeneratorNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	sql := nl2state.SQLGenerateOutput
	sqlResult := nl2state.SQLExecuteOutput
	pythonAnalysis := nl2state.PythonAnalyzeOutput

	if pythonAnalysis == "" {
		pythonAnalysis = sqlResult
	}

	prompt := GetPrompt("report_generate")
	data := map[string]string{
		"Query":          query,
		"SQL":            sql,
		"SQLResult":      sqlResult,
		"PythonAnalysis": pythonAnalysis,
		"Continuation":   "",
	}

	report, err := n.LLM.CallWithPrompt(ctx, prompt, data)
	if err != nil {
		report = generateFallbackReport(query, sql, sqlResult, pythonAnalysis)
	}

	report = strings.TrimSpace(report)
	if report == "" {
		report = generateFallbackReport(query, sql, sqlResult, pythonAnalysis)
	}

	nl2state.Result = report
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":              nl2state,
		nl2sqltypes.StateKeyResult:  report,
		"last_response":             report,
	}, nil
}

// generateFallbackReport creates a simple report when LLM generation fails.
func generateFallbackReport(query, sql, sqlResult, pythonAnalysis string) string {
	var b strings.Builder
	b.WriteString("# 数据分析报告\n\n")
	b.WriteString("## 查询说明\n\n")
	b.WriteString("**用户需求**: " + query + "\n\n")

	if sql != "" {
		b.WriteString("## 执行SQL\n\n")
		b.WriteString("```sql\n")
		b.WriteString(sql + "\n")
		b.WriteString("```\n\n")
	}

	if sqlResult != "" {
		b.WriteString("## 查询结果\n\n")
		b.WriteString(sqlResult + "\n\n")
	}

	if pythonAnalysis != "" {
		b.WriteString("## 分析结果\n\n")
		b.WriteString(pythonAnalysis + "\n\n")
	}

	b.WriteString("---\n\n")
	b.WriteString("*报告由 NL2SQL 工作流自动生成*\n")

	return b.String()
}
