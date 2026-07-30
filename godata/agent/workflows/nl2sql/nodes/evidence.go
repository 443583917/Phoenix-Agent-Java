package nodes

import (
	"context"
	"strings"

	"github.com/phoenix-agent-go/agent/knowledge"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// EvidenceRecallNode retrieves relevant evidence/knowledge for the query.
type EvidenceRecallNode struct {
	LLM       *LLMService
	Retriever *knowledge.Retriever
}

// Name returns the unique node identifier.
func (n *EvidenceRecallNode) Name() string {
	return "evidence_recall"
}

// Execute implements graph.NodeFunc for the tRPC-Agent-Go StateGraph.
// It rewrites the query for retrieval, searches for relevant evidence,
// and stores the results in state.
func (n *EvidenceRecallNode) Execute(ctx context.Context, state graph.State) (any, error) {
	nl2state := getOrCreateState(state)

	query := getStateStr(state, nl2sqltypes.StateKeyInput)
	if query == "" {
		query, _ = state["user_input"].(string)
	}
	contextStr := getStateStr(state, nl2sqltypes.StateKeyMultiTurnContext)

	// Step 1: Rewrite query for retrieval
	rewritePrompt := GetPrompt("evidence_rewrite")
	rewriteData := map[string]string{
		"Query":   query,
		"Context": contextStr,
	}

	var rewriteResult struct {
		CanonicalQuery string `json:"canonical_query"`
	}
	rewrittenQuery := query
	if err := n.LLM.CallJSON(ctx, rewritePrompt, rewriteData, &rewriteResult); err == nil && rewriteResult.CanonicalQuery != "" {
		rewrittenQuery = rewriteResult.CanonicalQuery
	}

	// Step 2: Search for evidence
	var evidenceParts []string
	if n.Retriever != nil {
		docs, err := n.Retriever.Search(ctx, rewrittenQuery, 5)
		if err == nil {
			for _, doc := range docs {
				if doc.Content != "" {
					evidenceParts = append(evidenceParts, doc.Content)
				}
			}
		}
	}

	evidence := ""
	if len(evidenceParts) > 0 {
		evidence = strings.Join(evidenceParts, "\n\n")
	}

	nl2state.Evidence = evidence
	nl2state.CurrentNode = n.Name()

	return graph.State{
		"nl2sql_state":             nl2state,
		nl2sqltypes.StateKeyEvidence: evidence,
	}, nil
}
