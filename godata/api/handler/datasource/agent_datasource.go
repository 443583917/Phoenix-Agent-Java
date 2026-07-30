package datasource

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// AgentDatasourceHandler handles agent-datasource link CRUD, toggle, and tables.
type AgentDatasourceHandler struct {
	svc *service.DataService
}

// NewAgentDatasourceHandler creates a new AgentDatasourceHandler.
func NewAgentDatasourceHandler(svc *service.DataService) *AgentDatasourceHandler {
	return &AgentDatasourceHandler{svc: svc}
}

// Page returns a paginated list of datasource-agent links for a given agent.
// GET /api/agent/:agentId/datasource/page
func (h *AgentDatasourceHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	agentIDStr := c.Param("agentId")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var query model.AgentDatasource
	query.AgentId = agentID
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAgentDatasource(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single datasource-agent link by its ID.
// GET /api/agent/:agentId/datasource/:id
func (h *AgentDatasourceHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAgentDatasourceByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new datasource-agent link.
// POST /api/agent/:agentId/datasource
func (h *AgentDatasourceHandler) Create(c *gin.Context) {
	agentIDStr := c.Param("agentId")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var entity model.AgentDatasource
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.AgentId = agentID

	if _, err := h.svc.CreateAgentDatasource(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing datasource-agent link.
// PUT /api/agent/:agentId/datasource
func (h *AgentDatasourceHandler) Update(c *gin.Context) {
	var entity model.AgentDatasource
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAgentDatasource(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a datasource-agent link by its ID.
// DELETE /api/agent/:agentId/datasource/:id
func (h *AgentDatasourceHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAgentDatasource(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ToggleActive toggles the is_active flag.
// PUT /api/agent/:agentId/datasource/:id/toggle
func (h *AgentDatasourceHandler) ToggleActive(c *gin.Context) {
	if err := h.svc.ToggleAgentDatasourceActive(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// GetTables returns the table list for a specific datasource-agent link.
// GET /api/agent/:agentId/datasource/:id/tables
func (h *AgentDatasourceHandler) GetTables(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	tables, err := h.svc.GetAgentDatasourceTables(c.Request.Context(), adID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tables)
}

// SaveTables replaces the table selection for a datasource-agent link.
// POST /api/agent/:agentId/datasource/:id/tables
func (h *AgentDatasourceHandler) SaveTables(c *gin.Context) {
	idStr := c.Param("id")
	adID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var req struct {
		Tables []string `json:"tables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SaveAgentDatasourceTables(c.Request.Context(), adID, req.Tables); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
