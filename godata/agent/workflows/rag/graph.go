package rag

import (
	"reflect"

	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/agent/workflows/rag/nodes"
	"github.com/phoenix-agent-go/agent/workflows/rag/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

const (
	NodeRetrieve = "retrieve"
	NodeAssemble = "assemble"
)

type RagGraph struct {
	retriever *knowledge.Retriever
}

func NewRagGraph(retriever *knowledge.Retriever) *RagGraph {
	return &RagGraph{retriever: retriever}
}

func (g *RagGraph) Build() (*graph.Graph, error) {
	schema := graph.NewStateSchema()
	schema.AddField("rag_state", graph.StateField{
		Type:    reflectType[*types.RagState](),
		Reducer: graph.DefaultReducer,
		Default: func() any { return &types.RagState{} },
	})

	sg := graph.NewStateGraph(schema)

	retrieveNode := &nodes.RetrieveNode{Retriever: g.retriever}
	assembleNode := &nodes.AssembleNode{}

	sg.AddNode(retrieveNode.Name(), retrieveNode.Execute)
	sg.AddNode(assembleNode.Name(), assembleNode.Execute)

	sg.SetEntryPoint(NodeRetrieve)
	sg.AddEdge(NodeRetrieve, NodeAssemble)
	sg.AddEdge(NodeAssemble, graph.End)

	return sg.Compile()
}

func reflectType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
