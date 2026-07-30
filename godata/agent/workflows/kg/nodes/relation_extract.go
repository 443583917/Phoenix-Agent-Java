package nodes

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/workflows/kg/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

const relationSystemPrompt = `You are a knowledge graph relation extraction assistant. Given the text and extracted entities, identify relationships between entities. Return a JSON object with a "relations" array. Each relation has: "source" (entity name), "target" (entity name), "relationType", "description".`

type RelationExtractNode struct {
	LLM *LLMService
}

func (n *RelationExtractNode) Name() string {
	return "relation_extract"
}

func (n *RelationExtractNode) Execute(ctx context.Context, state graph.State) (any, error) {
	kgState := getKGState(state)

	input := ""
	if v, ok := state[types.StateKeyInput].(string); ok {
		input = v
	}
	if input == "" {
		input = kgState.Input
	}

	entityContext := ""
	if kgState.EntityOutput != nil {
		for _, e := range kgState.EntityOutput.Entities {
			entityContext += fmt.Sprintf("- %s (%s): %s\n", e.Name, e.Type, e.Description)
		}
	}

	userPrompt := fmt.Sprintf("Text:\n%s\n\nExtracted entities:\n%s\n\nIdentify relationships between these entities.", input, entityContext)

	var output types.RelationExtractOutput
	if err := n.LLM.CallJSON(ctx, relationSystemPrompt, userPrompt, &output); err != nil {
		return graph.State{
			"kg_state":                kgState,
			types.StateKeyRelationExtract: &types.RelationExtractOutput{},
		}, nil
	}

	kgState.RelationOutput = &output
	kgState.CurrentNode = n.Name()

	return graph.State{
		"kg_state":                kgState,
		types.StateKeyRelationExtract: &output,
	}, nil
}
