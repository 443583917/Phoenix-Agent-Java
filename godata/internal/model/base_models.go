package model

import "time"

// PlatformBaseModel is the base embedded struct for all platform entities.
//
// Fields: createTime, creator, updateTime, updator, delFlag, keyword.
// Unlike the privilege BaseModel (which has CreateBy/UpdateBy), PlatformBaseModel
// uses creator/updator to match the Java com.phoenix.common.model.BaseModel.
// It does NOT include an ID field — each platform entity declares its own
// primary key (string ID).
type PlatformBaseModel struct {
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Creator    *string   `gorm:"column:creator;type:varchar(64)" json:"creator,omitempty"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Updator    *string   `gorm:"column:updator;type:varchar(64)" json:"updator,omitempty"`
	DelFlag    int       `gorm:"column:del_flag;default:0" json:"delFlag"`
	Keyword    string    `gorm:"-" json:"keyword,omitempty"`
}
