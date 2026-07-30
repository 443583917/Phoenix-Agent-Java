package nodes

import (
	"context"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SqlGenerateNode generates SQL from the execution plan using the LLM.
type SqlGenerateNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *SqlGenerateNode) Name() string {
	return "sql_generate"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It retrieves retry context from state, calls the LLM with the sql_generate
// prompt, parses and cleans the generated SQL, and updates the state.
func (n *SqlGenerateNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	schemaInfo := nl2state.GeneratedSemanticModelPrompt
	if schemaInfo == "" {
		schemaInfo = nl2state.TableRelationOutput
	}
	evidence := nl2state.Evidence
	dialect := nl2state.DBDialectType
	if dialect == "" {
		dialect = "mysql"
	}

	// Get retry reason from state (for regeneration context).
	retryReason := ""
	if nl2state.SQLRegenReason != nil {
		retryReason = nl2state.SQLRegenReason.Reason
	}

	// Increment generate count.
	generateCount := nl2state.SQLGenerateCount + 1

	prompt := GetPrompt("sql_generate")
	data := map[string]string{
		"Query":      query,
		"SchemaInfo": schemaInfo,
		"Evidence":   evidence,
		"Dialect":    dialect,
	}

	if retryReason != "" {
		data["RetryReason"] = retryReason
	}

	var result struct {
		SQL         string `json:"sql"`
		Explanation string `json:"explanation"`
	}

	sql := ""
	if err := n.LLM.CallJSON(ctx, prompt, data, &result); err == nil && result.SQL != "" {
		sql = cleanSQL(result.SQL)
	}

	// Fallback if SQL generation failed or produced no valid SQL.
	if sql == "" {
		sql = "-- SQL generation failed, no valid SQL produced"
	}

	nl2state.SQLGenerateOutput = sql
	nl2state.SQLGenerateCount = generateCount
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                      nl2state,
		nl2sqltypes.StateKeySQLGenerateOutput: sql,
		nl2sqltypes.StateKeySQLGenerateCount:  generateCount,
	}, nil
}

// cleanSQL removes markdown code fences and trims whitespace from SQL.
func cleanSQL(sql string) string {
	sql = strings.TrimSpace(sql)
	sql = strings.TrimPrefix(sql, "```sql")
	sql = strings.TrimPrefix(sql, "```")
	sql = strings.TrimSuffix(sql, "```")
	sql = strings.TrimSpace(sql)
	return sql
}
