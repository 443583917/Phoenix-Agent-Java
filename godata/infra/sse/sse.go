package sse

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func Stream(c *gin.Context, events <-chan Event) {
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
		case event, ok := <-events:
			if !ok {
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
