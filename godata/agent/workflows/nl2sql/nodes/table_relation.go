package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/agent/workflows/nl2sql/prompts"
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
	prompt := prompts.TableRelationPrompt
	data := map[string]string{
		"SchemaInfo": schemaInfo,
		"DBDialect":  dialect,
	}

	var result struct {
		TableRelationDescription string `json:"tableRelationDescription"`
		SemanticModelPrompt      string `json:"semanticModelPrompt"`
	}

	var tableRelation string
	var semanticPrompt string

	if err := n.LLM.CallJSON(ctx, prompt, data, &result); err == nil {
		tableRelation = result.TableRelationDescription
		semanticPrompt = result.SemanticModelPrompt
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
