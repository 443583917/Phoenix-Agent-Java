package types

const (
	StateKeyInput             = "input"
	StateKeyEntityExtract     = "entity_extract_output"
	StateKeyRelationExtract   = "relation_extract_output"
	StateKeyGraphMergeOutput  = "graph_merge_output"
	StateKeyDomainID          = "domain_id"
	StateKeyCurrentNode       = "current_node"
)

type KGState struct {
	Input            string                 `json:"input"`
	DomainID         string                 `json:"domainId"`
	EntityOutput     *EntityExtractOutput   `json:"entityOutput"`
	RelationOutput   *RelationExtractOutput `json:"relationOutput"`
	GraphMergeOutput *GraphMergeOutput      `json:"graphMergeOutput"`
	CurrentNode      string                 `json:"currentNode"`
}

type ExtractedEntity struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ExtractedRelation struct {
	Source       string `json:"source"`
	Target       string `json:"target"`
	RelationType string `json:"relationType"`
	Description  string `json:"description"`
}

type EntityExtractOutput struct {
	Entities []ExtractedEntity `json:"entities"`
}

type RelationExtractOutput struct {
	Relations []ExtractedRelation `json:"relations"`
}

type GraphMergeOutput struct {
	EntitiesCreated  int `json:"entitiesCreated"`
	RelationsCreated int `json:"relationsCreated"`
}
