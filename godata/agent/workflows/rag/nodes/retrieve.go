package nodes

import "context"

// RetrieveNode is a stub for the document retrieval step in the RAG workflow.
type RetrieveNode struct {
	name string
}

func NewRetrieveNode() *RetrieveNode {
	return &RetrieveNode{name: "retrieve"}
}

func (n *RetrieveNode) Name() string {
	return n.name
}

// Execute retrieves relevant documents based on the query.
// This is a stub — actual implementation will use vector search (Milvus) and
// keyword search to find relevant documents.
func (n *RetrieveNode) Execute(ctx context.Context, query string, topK int) ([]Document, error) {
	return nil, nil
}

// Document represents a retrieved document chunk.
type Document struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Score    float64 `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
