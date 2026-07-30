package model

// PlatformInfo — tbl_platform_platform_info
// 三方平台信息
type PlatformInfo struct {
	PlatformBaseModel
	ID         string `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	Type       string `gorm:"column:type;type:varchar(32)" json:"type"`
	Name       string `gorm:"column:name;type:varchar(128)" json:"name"`
	Status     string `gorm:"column:status;type:varchar(16);default:0" json:"status"`
	CorpID     string `gorm:"column:corpid;type:varchar(128)" json:"corpid"`
	CorpSecret string `gorm:"column:corpsecret;type:varchar(256)" json:"corpsecret"`
	AgentID    string `gorm:"column:agentid;type:varchar(128)" json:"agentid"`
	AppKey     string `gorm:"column:app_key;type:varchar(128)" json:"appKey"`
}

func (PlatformInfo) TableName() string { return "tbl_platform_platform_info" }
