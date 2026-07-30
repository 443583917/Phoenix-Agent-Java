package model

// ===== KG Module GORM Entities =====
//
// Each entity embeds BaseModel for common fields:
// ID (string PK), CreateTime, UpdateTime, CreateBy, UpdateBy, DelFlag.

// ----------------------------------------------------------------
// 1. KGEntity — tbl_kg_entity
// Represents a node/entity in the knowledge graph.
// ----------------------------------------------------------------

type KGEntity struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(256)" json:"name"`
	Type        string `gorm:"column:type;type:varchar(64)" json:"type"`
	Description string `gorm:"column:description;type:text" json:"description"`
	Properties  string `gorm:"column:properties;type:text" json:"properties"`
	Source      string `gorm:"column:source;type:varchar(64)" json:"source"`
	DomainId    string `gorm:"column:domain_id;type:varchar(64)" json:"domainId"`
}

func (KGEntity) TableName() string { return "tbl_kg_entity" }

// ----------------------------------------------------------------
// 2. KGRelation — tbl_kg_relation
// Represents an edge/relation between two KG entities.
// ----------------------------------------------------------------

type KGRelation struct {
	BaseModel
	SourceEntityId string  `gorm:"column:source_entity_id;type:varchar(64)" json:"sourceEntityId"`
	TargetEntityId string  `gorm:"column:target_entity_id;type:varchar(64)" json:"targetEntityId"`
	RelationType   string  `gorm:"column:relation_type;type:varchar(64)" json:"relationType"`
	Properties     string  `gorm:"column:properties;type:text" json:"properties"`
	Weight         float64 `gorm:"column:weight;type:decimal(10,4);default:1.0" json:"weight"`
}

func (KGRelation) TableName() string { return "tbl_kg_relation" }

// ----------------------------------------------------------------
// 3. KGDomain — tbl_kg_domain
// Domain/category for organizing KG entities.
// ----------------------------------------------------------------

type KGDomain struct {
	BaseModel
	Code        string `gorm:"column:code;type:varchar(64)" json:"code"`
	Name        string `gorm:"column:name;type:varchar(128)" json:"name"`
	Description string `gorm:"column:description;type:varchar(256)" json:"description"`
}

func (KGDomain) TableName() string { return "tbl_kg_domain" }
