package knowledge

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// AgentKnowledgeHandler handles agent knowledge CRUD, recall toggle, and embedding retry.
type AgentKnowledgeHandler struct {
	svc *service.DataService
}

// NewAgentKnowledgeHandler creates a new AgentKnowledgeHandler.
func NewAgentKnowledgeHandler(svc *service.DataService) *AgentKnowledgeHandler {
	return &AgentKnowledgeHandler{svc: svc}
}

// Page returns a paginated list of knowledge entries.
// GET /api/agent-knowledge/page
func (h *AgentKnowledgeHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.AgentKnowledge
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAgentKnowledge(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single knowledge entry by its primary key.
// GET /api/agent-knowledge/:id
func (h *AgentKnowledgeHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAgentKnowledgeByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new knowledge entry (supports multipart for file upload).
// POST /api/agent-knowledge
func (h *AgentKnowledgeHandler) Create(c *gin.Context) {
	var entity model.AgentKnowledge

	// Try JSON first, then multipart form.
	contentType := c.ContentType()
	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&entity); err != nil {
			response.Error(c, errcode.InvalidParams)
			return
		}
	} else {
		// Multipart form.
		agentIDStr := c.PostForm("agentId")
		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			response.Error(c, errcode.InvalidParams)
			return
		}
		entity.AgentId = agentID
		entity.Title = c.PostForm("title")
		entity.Type = c.PostForm("type")
		entity.Question = c.PostForm("question")
		entity.Content = c.PostForm("content")
		entity.SplitterType = c.DefaultPostForm("splitterType", "token")

		// Handle file upload.
		file, header, err := c.Request.FormFile("file")
		if err == nil {
			defer file.Close()
			entity.SourceFilename = header.Filename
			entity.FileType = header.Header.Get("Content-Type")
			entity.FileSize = header.Size
			// In a real implementation, save file to storage and set entity.FilePath.
		}
	}

	if _, err := h.svc.CreateAgentKnowledge(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing knowledge entry.
// PUT /api/agent-knowledge
func (h *AgentKnowledgeHandler) Update(c *gin.Context) {
	var entity model.AgentKnowledge
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAgentKnowledge(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a knowledge entry by its ID.
// DELETE /api/agent-knowledge/:id
func (h *AgentKnowledgeHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAgentKnowledge(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ToggleRecall toggles the is_recall flag.
// PUT /api/agent-knowledge/:id/recall-toggle
func (h *AgentKnowledgeHandler) ToggleRecall(c *gin.Context) {
	if err := h.svc.ToggleAgentKnowledgeRecall(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// RetryEmbedding resets embedding status to "pending" for retry.
// POST /api/agent-knowledge/:id/retry-embedding
func (h *AgentKnowledgeHandler) RetryEmbedding(c *gin.Context) {
	if err := h.svc.RetryAgentKnowledgeEmbedding(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// QueryPage returns a paginated list using POST body for complex filter conditions.
// POST /api/agent-knowledge/query/page
func (h *AgentKnowledgeHandler) QueryPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	var query model.AgentKnowledge
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	list, total, err := h.svc.PageAgentKnowledge(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}
