package nodes

import (
	"context"
	"strconv"

	"github.com/phoenix-agent-go/agent/workflows/kg/types"
	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type GraphMergeNode struct {
	EntityRepo   repository.KGEntityRepository
	RelationRepo repository.KGRelationRepository
}

func (n *GraphMergeNode) Name() string {
	return "graph_merge"
}

func (n *GraphMergeNode) Execute(ctx context.Context, state graph.State) (any, error) {
	kgState := getKGState(state)

	entityIDMap := make(map[string]string)
	entitiesCreated := 0

	if kgState.EntityOutput != nil {
		for _, e := range kgState.EntityOutput.Entities {
			entityID := strconv.FormatUint(id.MustGenerateID(), 10)
			entityIDMap[e.Name] = entityID

			entity := &model.KGEntity{
				ID: entityID,
				Name:        e.Name,
				Type:        e.Type,
				Description: e.Description,
				DomainId:    kgState.DomainID,
				Source:      "llm_extract",
			}
			if err := n.EntityRepo.Create(ctx, entity); err == nil {
				entitiesCreated++
			}
		}
	}

	relationsCreated := 0
	if kgState.RelationOutput != nil {
		for _, r := range kgState.RelationOutput.Relations {
			sourceID, srcOK := entityIDMap[r.Source]
			targetID, tgtOK := entityIDMap[r.Target]
			if !srcOK || !tgtOK {
				continue
			}

			relation := &model.KGRelation{
				ID: strconv.FormatUint(id.MustGenerateID(), 10),
				SourceEntityId: sourceID,
				TargetEntityId: targetID,
				RelationType:   r.RelationType,
				Properties:     r.Description,
				Weight:         1.0,
			}
			if err := n.RelationRepo.Create(ctx, relation); err == nil {
				relationsCreated++
			}
		}
	}

	output := &types.GraphMergeOutput{
		EntitiesCreated:  entitiesCreated,
		RelationsCreated: relationsCreated,
	}
	kgState.GraphMergeOutput = output
	kgState.CurrentNode = n.Name()

	return graph.State{
		"kg_state":                  kgState,
		types.StateKeyGraphMergeOutput: output,
	}, nil
}
