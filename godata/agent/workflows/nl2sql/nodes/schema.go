package nodes

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/knowledge"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// SchemaRecallNode retrieves relevant schema information for the query.
type SchemaRecallNode struct {
	Retriever *knowledge.Retriever
	LLM       *LLMService
}

// Name returns the unique node identifier.
func (n *SchemaRecallNode) Name() string {
	return "schema_recall"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It searches for table and column documents using the retriever
// and stores the results in state.
func (n *SchemaRecallNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := nl2state.QueryEnhanceOutput.CanonicalQuery
	if query == "" {
		query = getStateStr(state, nl2sqltypes.StateKeyInput)
	}
	if query == "" {
		query, _ = state["user_input"].(string)
	}

	var tableDocs []string
	var columnDocs []string

	if n.Retriever != nil {
		// Search for table schemas
		tableResults, err := n.Retriever.Search(ctx, fmt.Sprintf("table schema %s", query), 10)
		if err == nil {
			for _, doc := range tableResults {
				if doc.Content != "" {
					tableDocs = append(tableDocs, doc.Content)
				}
			}
		}

		// Search for column details
		columnResults, err := n.Retriever.Search(ctx, fmt.Sprintf("column definition %s", query), 15)
		if err == nil {
			for _, doc := range columnResults {
				if doc.Content != "" {
					columnDocs = append(columnDocs, doc.Content)
				}
			}
		}
	}

	nl2state.TableDocumentsForSchema = tableDocs
	nl2state.ColumnDocumentsForSchema = columnDocs
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":                              nl2state,
		nl2sqltypes.StateKeyTableDocumentsForSchema:  tableDocs,
		nl2sqltypes.StateKeyColumnDocumentsForSchema: columnDocs,
	}, nil
}
