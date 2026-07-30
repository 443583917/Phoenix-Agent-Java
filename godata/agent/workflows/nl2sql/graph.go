package nl2sql

import (
	"github.com/phoenix-agent-go/agent/workflows/nl2sql/nodes"
	"github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
)

// StateGraph is a simplified graph structure that holds nodes and edges.
type StateGraph struct {
	Nodes map[string]types.Node
	Edges map[string]string // from node -> to node (simple linear graph)
}

// NL2SQLGraph builds the NL2SQL workflow graph.
type NL2SQLGraph struct{}

// NewNL2SQLGraph creates a new NL2SQLGraph builder.
func NewNL2SQLGraph() *NL2SQLGraph {
	return &NL2SQLGraph{}
}

// Build instantiates all node stubs and wires them into a StateGraph.
func (g *NL2SQLGraph) Build() *StateGraph {
	graph := &StateGraph{
		Nodes: make(map[string]types.Node),
		Edges: make(map[string]string),
	}

	// Node stubs
	intentNode := &nodes.IntentRecognitionNode{}
	evidenceNode := &nodes.EvidenceRecallNode{}
	schemaNode := &nodes.SchemaRecallNode{}
	plannerNode := &nodes.PlannerNode{}
	sqlGenNode := &nodes.SqlGenerateNode{}
	pythonExecNode := &nodes.PythonExecuteNode{}
	reportNode := &nodes.ReportGeneratorNode{}

	// Register nodes
	graph.Nodes[intentNode.Name()] = intentNode
	graph.Nodes[evidenceNode.Name()] = evidenceNode
	graph.Nodes[schemaNode.Name()] = schemaNode
	graph.Nodes[plannerNode.Name()] = plannerNode
	graph.Nodes[sqlGenNode.Name()] = sqlGenNode
	graph.Nodes[pythonExecNode.Name()] = pythonExecNode
	graph.Nodes[reportNode.Name()] = reportNode

	// Wire edges in order
	graph.Edges["intent_recognition"] = "evidence_recall"
	graph.Edges["evidence_recall"] = "schema_recall"
	graph.Edges["schema_recall"] = "planner"
	graph.Edges["planner"] = "sql_generate"
	graph.Edges["sql_generate"] = "python_execute"
	graph.Edges["python_execute"] = "report_generate"
	graph.Edges["report_generate"] = "end"

	return graph
}
