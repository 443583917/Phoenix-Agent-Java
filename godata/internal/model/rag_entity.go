package model

import "time"

// ===== RAG Module GORM Entities =====
// Audit columns follow Java com.phoenix.common.model.BaseModel convention:
// creator / updator (not create_by / update_by).

// ----------------------------------------------------------------
// 1. RagCategory — tbl_rag_category
// ----------------------------------------------------------------

type RagCategory struct {
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

func (RagCategory) TableName() string { return "tbl_rag_category" }

// ----------------------------------------------------------------
// 2. RagFileInfo — tbl_rag_file_info
// ----------------------------------------------------------------

type RagFileInfo struct {
	ID               string    `gorm:"primaryKey;column:id;type:varchar(64)" json:"id"`
	FileType         string    `gorm:"column:file_type;type:varchar(32)" json:"fileType"`
	CategoryId       string    `gorm:"column:category_id;type:varchar(64)" json:"categoryId"`
	Name             string    `gorm:"column:name;type:varchar(256)" json:"name"`
	Title            string    `gorm:"column:title;type:varchar(256)" json:"title"`
	Description      string    `gorm:"column:description;type:text" json:"description"`
	Path             string    `gorm:"column:path;type:varchar(512)" json:"path"`
	PdfType          string    `gorm:"column:pdf_type;type:varchar(32)" json:"pdfType"`
	PageTopMargin    int       `gorm:"column:page_top_margin" json:"pageTopMargin"`
	PagesPerDocument int       `gorm:"column:pages_per_document" json:"pagesPerDocument"`
	TextSplitter     bool      `gorm:"column:text_splitter;default:false" json:"textSplitter"`
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Creator          *string   `gorm:"column:creator;type:varchar(255)" json:"creator,omitempty"`
	UpdateTime       time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Updator          *string   `gorm:"column:updator;type:varchar(255)" json:"updator,omitempty"`
	DelFlag          int       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (RagFileInfo) TableName() string { return "tbl_rag_file_info" }
