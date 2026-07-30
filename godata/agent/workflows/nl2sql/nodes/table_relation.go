package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// TableRelationNode analyzes table and column documents to determine
// relationships and generate a semantic model prompt.
type TableRelationNode struct {
	LLM *LLMService
}

// Name returns the unique node identifier.
func (n *TableRelationNode) Name() string {
	return "table_relation"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It takes table/column documents from state, uses the LLM to analyze
// schema relationships, and generates a semantic model prompt with retry support.
func (n *TableRelationNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	// Build schema context from table/column documents
	tableDocs := nl2state.TableDocumentsForSchema
	columnDocs := nl2state.ColumnDocumentsForSchema

	var schemaBuf strings.Builder
	schemaBuf.WriteString("数据库表结构：\n\n")
	for i, doc := range tableDocs {
		schemaBuf.WriteString(fmt.Sprintf("表 %d:\n%s\n\n", i+1, doc))
	}
	schemaBuf.WriteString("字段信息：\n\n")
	for i, doc := range columnDocs {
		schemaBuf.WriteString(fmt.Sprintf("字段 %d:\n%s\n\n", i+1, doc))
	}

	schemaInfo := schemaBuf.String()
	if schemaInfo == "" {
		schemaInfo = "无表结构信息"
	}

	dialect := nl2state.DBDialectType
	if dialect == "" {
		dialect = "mysql"
	}

	// Use LLM to analyze table relations
	prompt := `你是一个数据库表关系分析助手。请分析以下表结构，确定表之间的关系（外键关系、逻辑关联），并生成一个语义模型描述。

请输出JSON格式：
{
  "table_relation_description": "表关系描述，描述主外键关联、逻辑关系等",
  "semantic_model_prompt": "用于SQL生成的语义模型提示，包含表名、字段名、类型、业务含义、关联关系等"
}

数据库方言：` + dialect + `

用户查询：` + query + `

` + schemaInfo

	resp, err := n.LLM.Call(ctx, prompt, "请分析上述表结构并生成JSON输出。")

	var tableRelation string
	var semanticPrompt string

	if err == nil {
		jsonStr := extractJSON(resp)
		if jsonStr != "" {
			var result struct {
				TableRelationDescription string `json:"table_relation_description"`
				SemanticModelPrompt      string `json:"semantic_model_prompt"`
			}
			if json.Unmarshal([]byte(jsonStr), &result) == nil {
				tableRelation = result.TableRelationDescription
				semanticPrompt = result.SemanticModelPrompt
			}
		}
	}

	// Fallback
	if tableRelation == "" {
		tableRelation = fmt.Sprintf("表结构: %s", schemaInfo)
	}
	if semanticPrompt == "" {
		semanticPrompt = tableRelation
	}

	nl2state.TableRelationOutput = tableRelation
	nl2state.GeneratedSemanticModelPrompt = semanticPrompt
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state": nl2state,
		nl2sqltypes.StateKeyTableRelationOutput:          tableRelation,
		nl2sqltypes.StateKeyGeneratedSemanticModelPrompt: semanticPrompt,
	}, nil
}
