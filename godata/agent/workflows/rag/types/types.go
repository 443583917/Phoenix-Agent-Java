package types

import "github.com/phoenix-agent-go/agent/knowledge"

type RagState struct {
	Input      string               `json:"input"`
	TopK       int                  `json:"topK"`
	Categories []string             `json:"categories,omitempty"`
	FileTypes  []string             `json:"fileTypes,omitempty"`
	Documents  []knowledge.Document `json:"documents,omitempty"`
	Context    string               `json:"context"`
}
