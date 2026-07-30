package rag

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// RagCategoryHandler handles RagCategory CRUD operations.
type RagCategoryHandler struct {
	svc *service.RagService
}

func NewRagCategoryHandler(svc *service.RagService) *RagCategoryHandler {
	return &RagCategoryHandler{svc: svc}
}

// Page returns a paginated list of RAG categories.
// GET /api/rag/category/page
func (h *RagCategoryHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.RagCategory
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageRagCategory(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// List returns all RAG categories.
// GET /api/rag/category/list
func (h *RagCategoryHandler) List(c *gin.Context) {
	list, err := h.svc.ListRagCategory(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByID returns a single RAG category by its primary key.
// GET /api/rag/category/:id
func (h *RagCategoryHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetRagCategoryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, entity)
}

// Create creates a new RAG category.
// POST /api/rag/category
func (h *RagCategoryHandler) Create(c *gin.Context) {
	var entity model.RagCategory
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateRagCategory(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing RAG category.
// PUT /api/rag/category
func (h *RagCategoryHandler) Update(c *gin.Context) {
	var entity model.RagCategory
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateRagCategory(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a RAG category by its ID.
// DELETE /api/rag/category/:id
func (h *RagCategoryHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteRagCategory(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
