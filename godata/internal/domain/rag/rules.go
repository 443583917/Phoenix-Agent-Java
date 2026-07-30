package rag

const (
	CategoryTypeFile = "file"
	CategoryTypeQA   = "qa"
	CategoryTypeFAQ  = "faq"
)

func IsValidCategoryType(categoryType string) bool {
	switch categoryType {
	case CategoryTypeFile, CategoryTypeQA, CategoryTypeFAQ:
		return true
	}
	return false
}
