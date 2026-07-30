package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/agent/workflows/rag/types"
	"trpc.group/trpc-go/trpc-agent-go/graph"
)

type RetrieveNode struct {
	Retriever *knowledge.Retriever
}

func (n *RetrieveNode) Name() string { return "retrieve" }

func (n *RetrieveNode) Execute(ctx context.Context, state graph.State) (any, error) {
	ragState := getRAGState(state)
	topK := ragState.TopK
	if topK <= 0 {
		topK = 5
	}

	var docs []knowledge.Document
	if n.Retriever != nil {
		var err error
		docs, err = n.Retriever.Search(ctx, ragState.Input, topK)
		if err != nil {
			docs = nil
		}
	}

	ragState.Documents = docs
	ragState.Context = assembleContext(docs)

	return graph.State{
		"rag_state": ragState,
	}, nil
}

type AssembleNode struct{}

func (n *AssembleNode) Name() string { return "assemble" }

func (n *AssembleNode) Execute(ctx context.Context, state graph.State) (any, error) {
	ragState := getRAGState(state)
	if ragState.Context == "" && len(ragState.Documents) > 0 {
		ragState.Context = assembleContext(ragState.Documents)
	}

	return graph.State{
		"rag_state": ragState,
	}, nil
}

func getRAGState(state graph.State) *types.RagState {
	if s, ok := state["rag_state"].(*types.RagState); ok && s != nil {
		return s
	}
	return &types.RagState{}
}

func assembleContext(docs []knowledge.Document) string {
	var parts []string
	for i, doc := range docs {
		parts = append(parts, fmt.Sprintf("[文档 %d] %s", i+1, doc.Content))
	}
	return strings.Join(parts, "\n\n")
}
