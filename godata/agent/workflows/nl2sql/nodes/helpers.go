package nodes

import (
	nl2sql "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// getOrCreateState retrieves the NL2SQLState from the graph.State or creates
// a new one if not present. This is used by all node implementations to
// access the shared workflow state.
func getOrCreateState(state graph.State) *nl2sql.NL2SQLState {
	if s, ok := state["nl2sql_state"].(*nl2sql.NL2SQLState); ok && s != nil {
		return s
	}
	return &nl2sql.NL2SQLState{}
}
