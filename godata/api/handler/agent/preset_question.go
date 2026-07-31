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

// AgentPresetQuestionHandler handles preset question CRUD scoped to an agent.
type AgentPresetQuestionHandler struct {
	svc *service.DataService
}

// NewAgentPresetQuestionHandler creates a new AgentPresetQuestionHandler.
func NewAgentPresetQuestionHandler(svc *service.DataService) *AgentPresetQuestionHandler {
	return &AgentPresetQuestionHandler{svc: svc}
}

// Page returns a paginated list of preset questions for an agent.
// GET /api/agent/:agentId/preset-question/page
func (h *AgentPresetQuestionHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var query model.AgentPresetQuestion
	query.AgentId = agentID
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAgentPresetQuestion(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single preset question by its ID.
// GET /api/agent/:agentId/preset-question/:id
func (h *AgentPresetQuestionHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAgentPresetQuestionByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new preset question for an agent.
// POST /api/agent/:agentId/preset-question
func (h *AgentPresetQuestionHandler) Create(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var entity model.AgentPresetQuestion
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.AgentId = agentID

	if _, err := h.svc.CreateAgentPresetQuestion(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing preset question.
// PUT /api/agent/:agentId/preset-question
func (h *AgentPresetQuestionHandler) Update(c *gin.Context) {
	var entity model.AgentPresetQuestion
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAgentPresetQuestion(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a preset question by its ID.
// DELETE /api/agent/:agentId/preset-question/:id
func (h *AgentPresetQuestionHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAgentPresetQuestion(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// List returns all preset questions for an agent (non-paginated).
// GET /api/agent/:agentId/preset-question
func (h *AgentPresetQuestionHandler) List(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var query model.AgentPresetQuestion
	query.AgentId = agentID
	_ = c.ShouldBindQuery(&query)

	list, _, err := h.svc.PageAgentPresetQuestion(c.Request.Context(), 1, 1000, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// ListByAccount returns preset questions filtered by account ID.
// GET /api/agent/:agentId/:accountId/preset-questions
func (h *AgentPresetQuestionHandler) ListByAccount(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var query model.AgentPresetQuestion
	query.AgentId = agentID
	query.AccountId = c.Param("accountId")

	list, _, err := h.svc.PageAgentPresetQuestion(c.Request.Context(), 1, 1000, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// BatchSave creates or updates multiple preset questions in one call.
// POST /api/agent/:agentId/preset-questions
func (h *AgentPresetQuestionHandler) BatchSave(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var entities []*model.AgentPresetQuestion
	if err := c.ShouldBindJSON(&entities); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	var created []*model.AgentPresetQuestion
	for _, entity := range entities {
		entity.AgentId = agentID
		result, err := h.svc.CreateAgentPresetQuestion(c.Request.Context(), entity)
		if err != nil {
			continue
		}
		created = append(created, result)
	}
	response.Success(c, created)
}
