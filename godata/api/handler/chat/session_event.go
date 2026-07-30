package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/internal/service"
)

// SessionEventHandler handles SSE streaming for session events.
type SessionEventHandler struct {
	svc *service.DataService
}

// NewSessionEventHandler creates a new SessionEventHandler.
func NewSessionEventHandler(svc *service.DataService) *SessionEventHandler {
	return &SessionEventHandler{svc: svc}
}

// StreamSessions streams session update events via SSE.
// GET /agent/:agentId/sessions/stream
func (h *SessionEventHandler) StreamSessions(c *gin.Context) {
	agentId := c.Param("agentId")

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	data := map[string]interface{}{
		"type": "session_update",
		"data": map[string]interface{}{
			"message": "Session event stream - Phase 5 stub",
			"agentId": agentId,
		},
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(c.Writer, "data: %s\n\n", jsonBytes)
	flusher.Flush()
}
