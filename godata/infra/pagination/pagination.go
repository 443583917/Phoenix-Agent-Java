package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageQuery struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

const (
	DefaultPage = 1
	DefaultSize = 10
	MaxSize     = 100
)

func ParsePageQuery(c *gin.Context) PageQuery {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = DefaultPage
	}
	if size < 1 {
		size = DefaultSize
	}
	if size > MaxSize {
		size = MaxSize
	}

	return PageQuery{Page: page, Size: size}
}

func (p PageQuery) Offset() int {
	return (p.Page - 1) * p.Size
}
