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

// DataBaseModel for data tables (tbl_data_*): uses created_time/updated_time,
// no create_by/update_by/del_flag columns.
type DataBaseModel struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	CreateTime time.Time `gorm:"column:created_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:updated_time;autoUpdateTime" json:"updateTime"`
}
