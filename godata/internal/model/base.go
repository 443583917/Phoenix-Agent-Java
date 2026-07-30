package model

import "time"

type BaseModel struct {
	ID         string         `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	CreateBy   *string        `gorm:"column:create_by;type:varchar(64)" json:"createBy,omitempty"`
	UpdateBy   *string        `gorm:"column:update_by;type:varchar(64)" json:"updateBy,omitempty"`
	DelFlag    int            `gorm:"column:del_flag;default:0" json:"delFlag"`
}
