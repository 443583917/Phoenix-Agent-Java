package model

import "time"

// UserAgentInfo — tbl_agent_user_agent_info
// 用户智能体信息表
type UserAgentInfo struct {
	BaseModel
	UserID      string     `gorm:"column:user_id;type:varchar(64)" json:"userId"`
	AgentSn     string     `gorm:"column:agent_sn;type:varchar(64)" json:"agentSn"`
	ActionCount *int64     `gorm:"column:action_count;default:0" json:"actionCount"`
	LastDate    *time.Time `gorm:"column:last_date" json:"lastDate"`
}

func (UserAgentInfo) TableName() string { return "tbl_agent_user_agent_info" }

// UserMemoryInfo — tbl_agent_user_memory_info
// 用户记忆信息表
type UserMemoryInfo struct {
	BaseModel
	UserID     string `gorm:"column:user_id;type:varchar(64)" json:"userId"`
	AgentSn    string `gorm:"column:agent_sn;type:varchar(64)" json:"agentSn"`
	MemoryType string `gorm:"column:memory_type;type:varchar(32)" json:"memoryType"`
	Content    string `gorm:"column:content;type:text" json:"content"`
}

func (UserMemoryInfo) TableName() string { return "tbl_agent_user_memory_info" }

// UserProfileInfo — tbl_agent_user_profile_info
// 用户画像信息表
type UserProfileInfo struct {
	BaseModel
	UserID      string `gorm:"column:user_id;type:varchar(64)" json:"userId"`
	AgentSn     string `gorm:"column:agent_sn;type:varchar(64)" json:"agentSn"`
	ProfileData string `gorm:"column:profile_data;type:text" json:"profileData"`
}

func (UserProfileInfo) TableName() string { return "tbl_agent_user_profile_info" }

// CombinedStore — tbl_vector_store_combined
// 向量存储组合表
type CombinedStore struct {
	BaseModel
	Namespace string `gorm:"column:namespace;type:varchar(128)" json:"namespace"`
	KeyName   string `gorm:"column:key_name;type:varchar(256)" json:"keyName"`
	ValueJson string `gorm:"column:value_json;type:text" json:"valueJson"`
}

func (CombinedStore) TableName() string { return "tbl_vector_store_combined" }
