package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/internal/service"
)

type SessionEventHandler struct {
	svc *service.DataService
}

func NewSessionEventHandler(svc *service.DataService) *SessionEventHandler {
	return &SessionEventHandler{svc: svc}
}

func (h *SessionEventHandler) StreamSessions(c *gin.Context) {
	agentIDStr := c.Param("id")
	agentID, _ := strconv.Atoi(agentIDStr)
	userID := getUserID(c)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	writeEvent := func(eventType string, data interface{}) {
		payload := map[string]interface{}{
			"type": eventType,
			"data": data,
		}
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return
		}
		fmt.Fprintf(c.Writer, "data: %s\n\n", jsonBytes)
		flusher.Flush()
	}

	list, err := h.svc.ListChatSessions(c.Request.Context(), agentID, userID)
	if err == nil && list != nil {
		writeEvent("session_update", map[string]interface{}{
			"sessions": list,
			"count":    len(list),
		})
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			writeEvent("heartbeat", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
		}
	}
}
