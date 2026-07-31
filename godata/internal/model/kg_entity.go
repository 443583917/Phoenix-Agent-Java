package model

import "time"

// ===== KG Module GORM Entities =====
// Audit columns use creator / updator convention (matching RAG tables).

// ----------------------------------------------------------------
// 1. KGEntity — tbl_kg_entity
// ----------------------------------------------------------------

type KGEntity struct {
	ID          string    `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(256)" json:"name"`
	Type        string    `gorm:"column:type;type:varchar(64)" json:"type"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Properties  string    `gorm:"column:properties;type:text" json:"properties"`
	Source      string    `gorm:"column:source;type:varchar(64)" json:"source"`
	DomainId    string    `gorm:"column:domain_id;type:varchar(64)" json:"domainId"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Creator     *string   `gorm:"column:creator;type:varchar(255)" json:"creator,omitempty"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Updator     *string   `gorm:"column:updator;type:varchar(255)" json:"updator,omitempty"`
	DelFlag     int       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (KGEntity) TableName() string { return "tbl_kg_entity" }

// ----------------------------------------------------------------
// 2. KGRelation — tbl_kg_relation
// ----------------------------------------------------------------

type KGRelation struct {
	ID             string    `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	SourceEntityId string    `gorm:"column:source_entity_id;type:varchar(64)" json:"sourceEntityId"`
	TargetEntityId string    `gorm:"column:target_entity_id;type:varchar(64)" json:"targetEntityId"`
	RelationType   string    `gorm:"column:relation_type;type:varchar(64)" json:"relationType"`
	Properties     string    `gorm:"column:properties;type:text" json:"properties"`
	Weight         float64   `gorm:"column:weight;type:decimal(10,4);default:1.0" json:"weight"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Creator        *string   `gorm:"column:creator;type:varchar(255)" json:"creator,omitempty"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Updator        *string   `gorm:"column:updator;type:varchar(255)" json:"updator,omitempty"`
	DelFlag        int       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (KGRelation) TableName() string { return "tbl_kg_relation" }

// ----------------------------------------------------------------
// 3. KGDomain — tbl_kg_domain
// ----------------------------------------------------------------

type KGDomain struct {
	ID          string    `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	Code        string    `gorm:"column:code;type:varchar(64)" json:"code"`
	Name        string    `gorm:"column:name;type:varchar(128)" json:"name"`
	Description string    `gorm:"column:description;type:varchar(256)" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Creator     *string   `gorm:"column:creator;type:varchar(255)" json:"creator,omitempty"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Updator     *string   `gorm:"column:updator;type:varchar(255)" json:"updator,omitempty"`
	DelFlag     int       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (KGDomain) TableName() string { return "tbl_kg_domain" }
