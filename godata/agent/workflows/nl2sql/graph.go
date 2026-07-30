package nl2sql

import (
	"context"
	"reflect"

	"github.com/phoenix-agent-go/agent/workflows/nl2sql/nodes"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"

	"trpc.group/trpc-go/trpc-agent-go/graph"
)

// NL2SQLGraph builds the NL2SQL workflow graph using the tRPC-Agent-Go
// graph.StateGraph builder. The workflow consists of 7 nodes connected
// sequentially:
//
//	intent_recognition -> evidence_recall -> schema_recall -> planner
//	  -> sql_generate -> python_execute -> report_generate -> end
//
// Conditional routing at planner: if the plan was rejected, skip to end.
type NL2SQLGraph struct{}

// NewNL2SQLGraph creates a new NL2SQLGraph builder.
func NewNL2SQLGraph() *NL2SQLGraph {
	return &NL2SQLGraph{}
}

// Build instantiates all node stubs and wires them into a tRPC-Agent-Go
// StateGraph. Returns the compiled graph ready for execution.
func (g *NL2SQLGraph) Build() (*graph.Graph, error) {
	// Use the framework's message-oriented state schema.
	schema := graph.MessagesStateSchema()

	// Add custom field for NL2SQL state.
	schema.AddField("nl2sql_state", graph.StateField{
		Type:    reflect.TypeOf(&nl2sqltypes.NL2SQLState{}),
		Reducer: graph.DefaultReducer,
		Default: func() any { return &nl2sqltypes.NL2SQLState{} },
	})

	sg := graph.NewStateGraph(schema)

	// Node stubs — each implements graph.NodeFunc signature:
	//   func(ctx context.Context, state graph.State) (any, error)
	intentNode := &nodes.IntentRecognitionNode{}
	evidenceNode := &nodes.EvidenceRecallNode{}
	schemaNode := &nodes.SchemaRecallNode{}
	plannerNode := &nodes.PlannerNode{}
	sqlGenNode := &nodes.SqlGenerateNode{}
	pythonExecNode := &nodes.PythonExecuteNode{}
	reportNode := &nodes.ReportGeneratorNode{}

	// Register nodes with graph.NodeFunc wrappers.
	sg.AddNode(intentNode.Name(), intentNode.Execute)
	sg.AddNode(evidenceNode.Name(), evidenceNode.Execute)
	sg.AddNode(schemaNode.Name(), schemaNode.Execute)
	sg.AddNode(plannerNode.Name(), plannerNode.Execute)
	sg.AddNode(sqlGenNode.Name(), sqlGenNode.Execute)
	sg.AddNode(pythonExecNode.Name(), pythonExecNode.Execute)
	sg.AddNode(reportNode.Name(), reportNode.Execute)

	// Wire linear edges.
	sg.AddEdge(intentNode.Name(), evidenceNode.Name())
	sg.AddEdge(evidenceNode.Name(), schemaNode.Name())
	sg.AddEdge(schemaNode.Name(), plannerNode.Name())

	// Conditional edge from planner: go to sql_generate or end if rejected.
	sg.AddConditionalEdges(plannerNode.Name(),
		func(ctx context.Context, state graph.State) (string, error) {
			nl2state, ok := state["nl2sql_state"].(*nl2sqltypes.NL2SQLState)
			if !ok || nl2state == nil {
				return sqlGenNode.Name(), nil
			}
			if nl2state.RejectedPlan {
				return graph.End, nil
			}
			return sqlGenNode.Name(), nil
		},
		map[string]string{
			sqlGenNode.Name(): sqlGenNode.Name(),
			graph.End:         graph.End,
		},
	)
	sg.AddEdge(sqlGenNode.Name(), pythonExecNode.Name())
	sg.AddEdge(pythonExecNode.Name(), reportNode.Name())
	sg.AddEdge(reportNode.Name(), graph.End)

	// Set entry and finish points.
	sg.SetEntryPoint(intentNode.Name())
	sg.SetFinishPoint(graph.End)

	return sg.Compile()
}
