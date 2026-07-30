package agent

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// AgentCategoryHandler handles agent category CRUD and tree operations.
type AgentCategoryHandler struct {
	svc *service.DataService
}

// NewAgentCategoryHandler creates a new AgentCategoryHandler.
func NewAgentCategoryHandler(svc *service.DataService) *AgentCategoryHandler {
	return &AgentCategoryHandler{svc: svc}
}

// Page returns a paginated list of categories.
// GET /api/agent-category/page
func (h *AgentCategoryHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.AgentCategory
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAgentCategory(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single category by its primary key.
// GET /api/agent-category/:id
func (h *AgentCategoryHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAgentCategoryByID(c.Request.Context(), c.Param("id"))
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

// Tree returns the full category tree.
// GET /api/agent-category/tree
func (h *AgentCategoryHandler) Tree(c *gin.Context) {
	tree, err := h.svc.TreeAgentCategory(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tree)
}

// Create creates a new category.
// POST /api/agent-category
func (h *AgentCategoryHandler) Create(c *gin.Context) {
	var entity model.AgentCategory
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateAgentCategory(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing category.
// PUT /api/agent-category
func (h *AgentCategoryHandler) Update(c *gin.Context) {
	var entity model.AgentCategory
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAgentCategory(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a category by its ID.
// DELETE /api/agent-category/:id
func (h *AgentCategoryHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAgentCategory(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
