package model

// ===== RAG Module GORM Entities =====
//
// Each entity embeds BaseModel for common fields:
// ID (string PK), CreateTime, UpdateTime, CreateBy, UpdateBy, DelFlag.
// Table names match the Java MyBatis Flex @Table annotations exactly.
// Field types and JSON names match the Java entity sources.

// ----------------------------------------------------------------
// 1. RagCategory — tbl_rag_category
// ----------------------------------------------------------------

type RagCategory struct {
	BaseModel
	Code        string `gorm:"column:code;type:varchar(64)" json:"code"`
	Name        string `gorm:"column:name;type:varchar(128)" json:"name"`
	Description string `gorm:"column:description;type:varchar(256)" json:"description"`
}

func (RagCategory) TableName() string { return "tbl_rag_category" }

// ----------------------------------------------------------------
// 2. RagFileInfo — tbl_rag_file_info
// ----------------------------------------------------------------

type RagFileInfo struct {
	BaseModel
	FileType          string `gorm:"column:file_type;type:varchar(32)" json:"fileType"`
	CategoryId        string `gorm:"column:category_id;type:varchar(64)" json:"categoryId"`
	Name              string `gorm:"column:name;type:varchar(256)" json:"name"`
	Title             string `gorm:"column:title;type:varchar(256)" json:"title"`
	Description       string `gorm:"column:description;type:text" json:"description"`
	Path              string `gorm:"column:path;type:varchar(512)" json:"path"`
	PdfType           string `gorm:"column:pdf_type;type:varchar(32)" json:"pdfType"`
	PageTopMargin     int    `gorm:"column:page_top_margin" json:"pageTopMargin"`
	PagesPerDocument  int    `gorm:"column:pages_per_document" json:"pagesPerDocument"`
	TextSplitter      bool   `gorm:"column:text_splitter;default:false" json:"textSplitter"`
}

func (RagFileInfo) TableName() string { return "tbl_rag_file_info" }
