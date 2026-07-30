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

// AgentHandler handles agent CRUD, publish/offline, and API key management.
type AgentHandler struct {
	svc *service.DataService
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(svc *service.DataService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

// Page returns a paginated list of agents.
// GET /api/agent/page
func (h *AgentHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.Agent
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAgent(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	// Mask API keys in response.
	for _, a := range list {
		if a.ApiKey != "" {
			a.ApiKey = maskKey(a.ApiKey)
		}
	}
	response.SuccessPage(c, list, total, page, size)
}

// List returns all agents.
// GET /api/agent/list
func (h *AgentHandler) List(c *gin.Context) {
	list, err := h.svc.ListAgent(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	for _, a := range list {
		if a.ApiKey != "" {
			a.ApiKey = maskKey(a.ApiKey)
		}
	}
	response.Success(c, list)
}

// GetByID returns a single agent by its primary key.
// GET /api/agent/:id
func (h *AgentHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAgentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	if entity.ApiKey != "" {
		entity.ApiKey = maskKey(entity.ApiKey)
	}
	response.Success(c, entity)
}

// Create creates a new agent.
// POST /api/agent
func (h *AgentHandler) Create(c *gin.Context) {
	var entity model.Agent
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateAgent(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing agent.
// PUT /api/agent
func (h *AgentHandler) Update(c *gin.Context) {
	var entity model.Agent
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAgent(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an agent by its ID.
// DELETE /api/agent/:id
func (h *AgentHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAgent(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Publish sets agent status to "published".
// POST /api/agent/:id/publish
func (h *AgentHandler) Publish(c *gin.Context) {
	if err := h.svc.PublishAgent(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Offline sets agent status to "offline".
// POST /api/agent/:id/offline
func (h *AgentHandler) Offline(c *gin.Context) {
	if err := h.svc.OfflineAgent(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ──────────────────── API Key management ────────────────────

// GenerateAPIKey generates a new API key for the agent.
// POST /api/agent/:id/api-key/generate
func (h *AgentHandler) GenerateAPIKey(c *gin.Context) {
	rawKey, err := h.svc.GenerateAgentAPIKey(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"apiKey": rawKey})
}

// ResetAPIKey replaces the agent's API key.
// POST /api/agent/:id/api-key/reset
func (h *AgentHandler) ResetAPIKey(c *gin.Context) {
	rawKey, err := h.svc.ResetAgentAPIKey(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"apiKey": rawKey})
}

// DeleteAPIKey removes the API key from the agent.
// DELETE /api/agent/:id/api-key
func (h *AgentHandler) DeleteAPIKey(c *gin.Context) {
	if err := h.svc.DeleteAgentAPIKey(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ToggleAPIKeyEnabled toggles the API key enabled flag.
// PUT /api/agent/:id/api-key/toggle
func (h *AgentHandler) ToggleAPIKeyEnabled(c *gin.Context) {
	if err := h.svc.ToggleAgentAPIKeyEnabled(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// GetAPIKeyMasked returns the masked API key.
// GET /api/agent/:id/api-key
func (h *AgentHandler) GetAPIKeyMasked(c *gin.Context) {
	masked, err := h.svc.GetAgentAPIKeyMasked(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"apiKey": masked})
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
