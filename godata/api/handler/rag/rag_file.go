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

// RagFileInfoHandler handles RagFileInfo CRUD operations.
type RagFileInfoHandler struct {
	svc *service.RagService
}

func NewRagFileInfoHandler(svc *service.RagService) *RagFileInfoHandler {
	return &RagFileInfoHandler{svc: svc}
}

// Page returns a paginated list of RAG files.
// GET /api/rag/file/page
func (h *RagFileInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.RagFileInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageRagFile(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// List returns all RAG files.
// GET /api/rag/file/list
func (h *RagFileInfoHandler) List(c *gin.Context) {
	list, err := h.svc.ListRagFile(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByID returns a single RAG file by its primary key.
// GET /api/rag/file/:id
func (h *RagFileInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetRagFileByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new RAG file.
// POST /api/rag/file
func (h *RagFileInfoHandler) Create(c *gin.Context) {
	var entity model.RagFileInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateRagFile(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing RAG file.
// PUT /api/rag/file
func (h *RagFileInfoHandler) Update(c *gin.Context) {
	var entity model.RagFileInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateRagFile(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a RAG file by its ID.
// DELETE /api/rag/file/:id
func (h *RagFileInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteRagFile(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
