package agent

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
)

// ReactAgentHandler handles agent-related SSE streaming endpoints.
type ReactAgentHandler struct {
	agentManager *runtime.AgentManager
}

// NewReactAgentHandler creates a new ReactAgentHandler.
func NewReactAgentHandler(agentManager *runtime.AgentManager) *ReactAgentHandler {
	return &ReactAgentHandler{agentManager: agentManager}
}

// Chat streams agent chat responses via SSE.
//
// POST /api/admin/agent/chat
//
// Accepts a JSON ChatModelRequest body and returns a text/event-stream of
// ContentEvent objects: {"content":"...","end":false} / {"content":"","end":true}.
func (h *ReactAgentHandler) Chat(c *gin.Context) {
	var req model.ChatModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)

	streamReq := model.StreamRequest{
		AgentSN:   req.AgentSn,
		UserID:    userIDStr,
		SessionID: req.SessionID,
		Message:   req.Content,
	}

	events, err := h.agentManager.StreamCall(c.Request.Context(), streamReq)
	if err != nil {
		response.Error(c, errcode.ModelError)
		return
	}

	streamSSE(c, events)
}

// StreamChatSQL is a stub SSE endpoint for graph search streaming.
//
// GET /api/admin/agent/stream/chatsql
//
// Full implementation will be provided in Phase 5.
func (h *ReactAgentHandler) StreamChatSQL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	event := model.ContentEvent{
		Content: "Graph search streaming coming in Phase 5",
		End:     true,
	}
	data, _ := json.Marshal(event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()
}

// streamSSE writes model.SSEEvent values from the channel to the Gin response
// as SSE data frames.
func streamSSE(c *gin.Context, events <-chan model.SSEEvent) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			if evt.Event == "" {
				continue
			}
			data, err := json.Marshal(evt.Data)
			if err != nil {
				continue
			}

			if evt.Event == "error" {
				fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", data)
			} else {
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			}
			flusher.Flush()
		}
	}
}
