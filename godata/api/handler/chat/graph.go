package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/agent/tools/datasource"
	"github.com/phoenix-agent-go/agent/workflows/nl2sql"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"github.com/phoenix-agent-go/internal/service"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/graph"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// GraphHandler handles SSE streaming endpoints for graph-based NL2SQL search.
type GraphHandler struct {
	svc       *service.DataService
	llm       model.Model
	dbManager *datasource.DatasourceManager
}

// NewGraphHandler creates a new GraphHandler.
// Pass nil for model/dbManager if not yet configured — the handler will
// report an error when StreamSearch is called without a model.
func NewGraphHandler(svc *service.DataService, llm model.Model, dbManager *datasource.DatasourceManager) *GraphHandler {
	return &GraphHandler{
		svc:       svc,
		llm:       llm,
		dbManager: dbManager,
	}
}

// StreamSearch streams NL2SQL graph search results via SSE.
// GET /stream/search?agentId=&threadId=&query=&humanFeedback=&humanFeedbackContent=&rejectedPlan=&nl2sqlOnly=
func (h *GraphHandler) StreamSearch(c *gin.Context) {
	agentID := c.Query("agentId")
	threadID := c.Query("threadId")
	query := c.Query("query")
	humanFeedback := c.Query("humanFeedback")
	humanFeedbackContent := c.Query("humanFeedbackContent")
	nl2sqlOnly := c.Query("nl2sqlOnly")

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	// Send initial event
	sendSSE(c.Writer, flusher, "start", map[string]interface{}{
		"query":    query,
		"agentId":  agentID,
		"threadId": threadID,
	})

	// Check if LLM is configured
	if h.llm == nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": "LLM model not configured. The NL2SQL graph needs a model to execute.",
		})
		return
	}

	// Build NL2SQL graph
	graphBuilder := nl2sql.NewNL2SQLGraph(h.llm, nil, h.dbManager)
	if humanFeedback == "true" {
		graphBuilder.WithHumanReview(true)
	}

	g, err := graphBuilder.Build()
	if err != nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Failed to build graph: %s", err.Error()),
		})
		return
	}

	// Prepare initial state
	nl2State := &nl2sqltypes.NL2SQLState{
		Input:              query,
		AgentID:            agentID,
		HumanReviewEnabled: humanFeedback == "true",
		HumanFeedbackData:  humanFeedbackContent,
		IsOnlyNL2SQL:       nl2sqlOnly == "true",
	}

	initialState := graph.State{
		"user_input":   query,
		"nl2sql_state": nl2State,
	}

	// Create executor
	executor, err := graph.NewExecutor(g)
	if err != nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Failed to create executor: %s", err.Error()),
		})
		return
	}

	// Create invocation
	inv := agent.NewInvocation(
		agent.WithInvocationID(agentID),
		agent.WithInvocationModel(h.llm),
	)

	sendSSE(c.Writer, flusher, "progress", map[string]interface{}{
		"message": "Starting NL2SQL graph execution",
		"node":    "start",
	})

	// Run the graph
	eventChan, err := executor.Execute(c.Request.Context(), initialState, inv)
	if err != nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Graph execution error: %s", err.Error()),
		})
		return
	}

	// Process graph events
	for evt := range eventChan {
		eventData := serializeEvent(evt)
		sendSSE(c.Writer, flusher, "node", eventData)

		// Check for state data in response events
		if evt != nil && evt.StateDelta != nil {
			sendSSE(c.Writer, flusher, "state_update", map[string]interface{}{
				"state_delta": evt.StateDelta,
			})
		}
	}

	// Send completion event (channel closed = graph finished)
	sendSSE(c.Writer, flusher, "complete", map[string]interface{}{
		"message": "Graph execution completed",
	})
}

// sendSSE sends a server-sent event.
func sendSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonBytes)
	flusher.Flush()
}

// serializeEvent converts a graph event to a serializable map.
func serializeEvent(evt *event.Event) map[string]interface{} {
	result := map[string]interface{}{
		"type": "graph_event",
	}

	if evt != nil {
		if evt.Author != "" {
			result["author"] = evt.Author
		}
		if evt.InvocationID != "" {
			result["invocation_id"] = evt.InvocationID
		}
		if evt.StateDelta != nil {
			result["state_delta"] = evt.StateDelta
		}
	}

	return result
}
