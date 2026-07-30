package runner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/phoenix-agent-go/internal/model"
)

// WriteSSE streams SSE events to the given http.ResponseWriter.
//
// It sets the appropriate headers (Content-Type: text/event-stream,
// Cache-Control: no-cache, Connection: keep-alive), then writes each event
// from the channel as a standard SSE data frame. When the channel closes,
// WriteSSE returns.
//
// Usage in an HTTP handler:
//
//	events, err := runner.Run(ctx, req)
//	if err != nil { ... }
//	runner.WriteSSE(w, events)
func WriteSSE(w http.ResponseWriter, events <-chan model.SSEEvent) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Set SSE headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	for evt := range events {
		// Skip empty events.
		if evt.Event == "" {
			continue
		}

		data, err := json.Marshal(evt)
		if err != nil {
			continue
		}

		// Write SSE frame: "event: <type>\ndata: <json>\n\n"
		// When the event type is "content" we also send a named event for
		// client-side dispatch.
		if evt.Event == "content" {
			fmt.Fprintf(w, "event: content\ndata: %s\n\n", data)
		} else {
			fmt.Fprintf(w, "data: %s\n\n", data)
		}
		flusher.Flush()
	}
}
