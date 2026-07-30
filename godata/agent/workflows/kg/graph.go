package kg

import (
	"reflect"

	"github.com/phoenix-agent-go/agent/workflows/kg/nodes"
	"github.com/phoenix-agent-go/agent/workflows/kg/types"
	"github.com/phoenix-agent-go/internal/repository"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

const (
	NodeEntityExtract   = "entity_extract"
	NodeRelationExtract = "relation_extract"
	NodeGraphMerge      = "graph_merge"
)

type KGGraph struct {
	llm          *nodes.LLMService
	entityRepo   repository.KGEntityRepository
	relationRepo repository.KGRelationRepository
}

func NewKGGraph(llm *nodes.LLMService, entityRepo repository.KGEntityRepository, relationRepo repository.KGRelationRepository) *KGGraph {
	return &KGGraph{
		llm:          llm,
		entityRepo:   entityRepo,
		relationRepo: relationRepo,
	}
}

func (g *KGGraph) Build() (*graph.Graph, error) {
	schema := graph.NewStateSchema()
	schema.AddField("kg_state", graph.StateField{
		Type:    reflectType[*types.KGState](),
		Reducer: graph.DefaultReducer,
		Default: func() any { return &types.KGState{} },
	})

	sg := graph.NewStateGraph(schema)

	entityNode := &nodes.EntityExtractNode{LLM: g.llm}
	relationNode := &nodes.RelationExtractNode{LLM: g.llm}
	mergeNode := &nodes.GraphMergeNode{
		EntityRepo:   g.entityRepo,
		RelationRepo: g.relationRepo,
	}

	sg.AddNode(entityNode.Name(), entityNode.Execute)
	sg.AddNode(relationNode.Name(), relationNode.Execute)
	sg.AddNode(mergeNode.Name(), mergeNode.Execute)

	sg.SetEntryPoint(NodeEntityExtract)
	sg.AddEdge(NodeEntityExtract, NodeRelationExtract)
	sg.AddEdge(NodeRelationExtract, NodeGraphMerge)
	sg.AddEdge(NodeGraphMerge, graph.End)

	return sg.Compile()
}

func reflectType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

var _ = (*KGGraph)(nil)
