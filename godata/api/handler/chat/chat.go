package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// ChatHandler handles session and message CRUD endpoints.
type ChatHandler struct {
	svc *service.DataService
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(svc *service.DataService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// ListSessions returns a stub session list for an agent.
// GET /agent/:id/sessions
func (h *ChatHandler) ListSessions(c *gin.Context) {
	agentID := c.Param("id")
	_ = agentID
	response.Success(c, []interface{}{})
}

// CreateSession returns a stub created session.
// POST /agent/:id/sessions
func (h *ChatHandler) CreateSession(c *gin.Context) {
	agentID := c.Param("id")
	response.Success(c, gin.H{
		"id":      "stub-session-id",
		"agentId": agentID,
		"title":   "New Session",
	})
}

// DeleteAllSessions returns a stub deletion confirmation.
// DELETE /agent/:id/sessions
func (h *ChatHandler) DeleteAllSessions(c *gin.Context) {
	agentID := c.Param("id")
	_ = agentID
	response.Success(c, gin.H{"deleted": true})
}

// GetMessages returns a stub message list for a session.
// GET /sessions/:id/messages
func (h *ChatHandler) GetMessages(c *gin.Context) {
	sessionID := c.Param("id")
	_ = sessionID
	response.Success(c, []interface{}{})
}

// CreateMessage returns a stub created message.
// POST /sessions/:id/messages
func (h *ChatHandler) CreateMessage(c *gin.Context) {
	sessionID := c.Param("id")
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)
	response.Success(c, gin.H{
		"id":        "stub-message-id",
		"sessionId": sessionID,
		"role":      "user",
		"content":   "stub message content",
	})
}

// PinSession returns a stub pin toggle response.
// PUT /sessions/:id/pin?isPinned=
func (h *ChatHandler) PinSession(c *gin.Context) {
	sessionID := c.Param("id")
	isPinned := c.Query("isPinned")
	response.Success(c, gin.H{
		"id":       sessionID,
		"isPinned": isPinned,
	})
}

// RenameSession returns a stub rename response.
// PUT /sessions/:id/rename?title=
func (h *ChatHandler) RenameSession(c *gin.Context) {
	sessionID := c.Param("id")
	title := c.Query("title")
	response.Success(c, gin.H{
		"id":    sessionID,
		"title": title,
	})
}

// DeleteSession returns a stub deletion confirmation.
// DELETE /sessions/:id
func (h *ChatHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	_ = sessionID
	response.Success(c, gin.H{"deleted": true})
}

// GenerateReportHTML returns a stub HTML report.
// POST /sessions/:id/reports/html
func (h *ChatHandler) GenerateReportHTML(c *gin.Context) {
	sessionID := c.Param("id")
	_ = sessionID
	response.Success(c, gin.H{
		"html": "<html><body><h1>Report - Phase 5 Stub</h1></body></html>",
	})
}
