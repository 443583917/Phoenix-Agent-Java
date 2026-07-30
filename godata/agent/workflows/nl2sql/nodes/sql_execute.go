package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/agent/tools/datasource"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SqlExecuteNode executes the generated SQL against a configured datasource
// and captures the results as formatted strings for downstream nodes.
type SqlExecuteNode struct {
	DBManager *datasource.DatasourceManager
}

// Name returns the unique node identifier.
func (n *SqlExecuteNode) Name() string {
	return "sql_execute"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It runs the generated SQL against the datasource, formats the output,
// and stores both the formatted text and a row-level list in state.
func (n *SqlExecuteNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	sql := nl2state.SQLGenerateOutput
	if sql == "" || strings.HasPrefix(sql, "--") {
		nl2state.SQLExecuteOutput = "SQL为空或无效，无法执行"
		nl2state.CurrentNode = n.Name()
		return graph.State{
			"nl2sql_state":                   nl2state,
			nl2sqltypes.StateKeySQLExecuteOutput: "SQL为空或无效，无法执行",
		}, nil
	}

	// Get datasource name from state or use a sensible default.
	datasourceName := getStateStr(state, nl2sqltypes.StateKeyAgentID)
	if datasourceName == "" {
		datasourceName = "default"
	}

	var output string
	var resultList []string

	if n.DBManager != nil {
		result, err := n.DBManager.QueryRows(ctx, datasourceName, sql)
		if err != nil {
			output = fmt.Sprintf("SQL执行失败: %s", err.Error())
		} else {
			// Format the result as a human-readable table.
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("查询结果 (%d行):\n", result.RowCount))
			buf.WriteString(strings.Join(result.Columns, " | "))
			buf.WriteString("\n")

			for _, row := range result.Rows {
				vals := make([]string, len(row))
				for i, v := range row {
					vals[i] = fmt.Sprintf("%v", v)
				}
				buf.WriteString(strings.Join(vals, " | "))
				buf.WriteString("\n")
			}

			if result.Truncated {
				buf.WriteString(fmt.Sprintf("\n(结果被截断，仅显示前%d行)", result.RowCount))
			}

			output = buf.String()

			// Store each row as a tab-separated string in the list memory.
			for _, row := range result.Rows {
				vals := make([]string, len(row))
				for i, v := range row {
					vals[i] = fmt.Sprintf("%v", v)
				}
				resultList = append(resultList, strings.Join(vals, "\t"))
			}
		}
	} else {
		output = "数据库管理器未配置，无法执行SQL"
	}

	nl2state.SQLExecuteOutput = output
	nl2state.SQLResultListMemory = resultList
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                       nl2state,
		nl2sqltypes.StateKeySQLExecuteOutput:    output,
		nl2sqltypes.StateKeySQLResultListMemory: resultList,
	}, nil
}
