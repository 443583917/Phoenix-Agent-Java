package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/internal/service"
)

// GraphHandler handles SSE streaming endpoints for graph search.
type GraphHandler struct {
	svc *service.DataService
}

// NewGraphHandler creates a new GraphHandler.
func NewGraphHandler(svc *service.DataService) *GraphHandler {
	return &GraphHandler{svc: svc}
}

// StreamSearch streams NL2SQL graph search results via SSE.
// GET /stream/search?agentId=&threadId=&query=&humanFeedback=&humanFeedbackContent=&rejectedPlan=&nl2sqlOnly=
func (h *GraphHandler) StreamSearch(c *gin.Context) {
	agentId := c.Query("agentId")
	threadId := c.Query("threadId")
	query := c.Query("query")
	humanFeedback := c.Query("humanFeedback")
	humanFeedbackContent := c.Query("humanFeedbackContent")
	rejectedPlan := c.Query("rejectedPlan")
	nl2sqlOnly := c.Query("nl2sqlOnly")

	_ = agentId
	_ = threadId
	_ = humanFeedback
	_ = humanFeedbackContent
	_ = rejectedPlan
	_ = nl2sqlOnly

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
		"type": "info",
		"data": map[string]interface{}{
			"message": "NL2SQL graph search streaming - Phase 5 stub",
			"query":   query,
		},
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(c.Writer, "data: %s\n\n", jsonBytes)
	flusher.Flush()
}
