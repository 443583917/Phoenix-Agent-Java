package nodes

import (
	"context"
	"fmt"

	"github.com/phoenix-agent-go/agent/workflows/kg/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

const entitySystemPrompt = `You are a knowledge graph entity extraction assistant. Extract all entities from the text. Return a JSON object with an "entities" array. Each entity has: "name", "type" (Person/Organization/Location/Concept/Product), "description".`

type EntityExtractNode struct {
	LLM *LLMService
}

func (n *EntityExtractNode) Name() string {
	return "entity_extract"
}

func (n *EntityExtractNode) Execute(ctx context.Context, state graph.State) (any, error) {
	kgState := getKGState(state)

	input := ""
	if v, ok := state[types.StateKeyInput].(string); ok {
		input = v
	}
	if input == "" {
		input = kgState.Input
	}

	var output types.EntityExtractOutput
	if err := n.LLM.CallJSON(ctx, entitySystemPrompt, fmt.Sprintf("Text:\n%s", input), &output); err != nil {
		return graph.State{
			"kg_state":              kgState,
			types.StateKeyEntityExtract: &types.EntityExtractOutput{},
		}, nil
	}

	kgState.EntityOutput = &output
	kgState.CurrentNode = n.Name()

	return graph.State{
		"kg_state":              kgState,
		types.StateKeyEntityExtract: &output,
	}, nil
}
