package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/agent/tools/datasource"
	"github.com/phoenix-agent-go/agent/workflows/nl2sql"
	nl2sqltypes "github.com/phoenix-agent-go/agent/workflows/nl2sql/types"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/service/tracing"
	"go.opentelemetry.io/otel/trace"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/graph"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

type GraphHandler struct {
	svc       *service.DataService
	llm       model.Model
	dbManager *datasource.DatasourceManager
	retriever *knowledge.Retriever
	mtcm      *runtime.MultiTurnContextManager
	tracing   *tracing.TracingService
	mu        sync.Mutex
	cached    *graph.Graph
}

func NewGraphHandler(svc *service.DataService, llm model.Model, dbManager *datasource.DatasourceManager, retriever *knowledge.Retriever, mtcm *runtime.MultiTurnContextManager, tracing *tracing.TracingService) *GraphHandler {
	return &GraphHandler{
		svc:       svc,
		llm:       llm,
		dbManager: dbManager,
		retriever: retriever,
		mtcm:      mtcm,
		tracing:   tracing,
	}
}

func (h *GraphHandler) StreamSearch(c *gin.Context) {
	agentID := c.Query("agentId")
	threadID := c.Query("threadId")
	query := c.Query("query")
	humanFeedback := c.Query("humanFeedback")
	humanFeedbackContent := c.Query("humanFeedbackContent")
	nl2sqlOnly := c.Query("nl2sqlOnly")
	rejectedPlan := c.Query("rejectedPlan")

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	sendSSE(c.Writer, flusher, "start", map[string]interface{}{
		"query":    query,
		"agentId":  agentID,
		"threadId": threadID,
	})

	if h.llm == nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": "LLM model not configured. The NL2SQL graph needs a model to execute.",
		})
		return
	}

	if h.mtcm != nil {
		if humanFeedback == "true" {
			h.mtcm.RestartLastTurn(threadID)
		} else {
			h.mtcm.BeginTurn(threadID)
		}
	}

	g, err := h.getOrBuildGraph(humanFeedback == "true")
	if err != nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Failed to build graph: %s", err.Error()),
		})
		return
	}

	nl2State := &nl2sqltypes.NL2SQLState{
		Input:              query,
		AgentID:            agentID,
		HumanReviewEnabled: humanFeedback == "true",
		HumanFeedbackData:  humanFeedbackContent,
		IsOnlyNL2SQL:       nl2sqlOnly == "true",
	}

	if rejectedPlan != "" {
		nl2State.PlanRepairCount = 1
		nl2State.PlanValidationError = rejectedPlan
	}

	if h.mtcm != nil {
		if ctxStr := h.mtcm.BuildContext(threadID); ctxStr != "" {
			nl2State.MultiTurnContext = ctxStr
		}
	}

	initialState := graph.State{
		"user_input":   query,
		"nl2sql_state": nl2State,
	}

	executor, err := graph.NewExecutor(g)
	if err != nil {
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Failed to create executor: %s", err.Error()),
		})
		return
	}

	inv := agent.NewInvocation(
		agent.WithInvocationID(agentID),
		agent.WithInvocationModel(h.llm),
	)

	sendSSE(c.Writer, flusher, "progress", map[string]interface{}{
		"message": "Starting NL2SQL graph execution",
		"node":    "start",
	})

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var span trace.Span
	if h.tracing != nil {
		ctx, span = h.tracing.StartGraphSpan(ctx, threadID, query)
	}

	eventChan, err := executor.Execute(ctx, initialState, inv)
	if err != nil {
		if h.tracing != nil {
			h.tracing.RecordError(span, err)
			span.End()
		}
		if h.mtcm != nil {
			h.mtcm.DiscardPending(threadID)
		}
		sendSSE(c.Writer, flusher, "error", map[string]interface{}{
			"message": fmt.Sprintf("Graph execution error: %s", err.Error()),
		})
		return
	}

	var lastResponse string
	for {
		select {
		case <-ctx.Done():
			if h.mtcm != nil {
				h.mtcm.DiscardPending(threadID)
			}
			if h.tracing != nil {
				span.End()
			}
			sendSSE(c.Writer, flusher, "complete", map[string]interface{}{
				"message": "Client disconnected",
			})
			return
		case evt, ok := <-eventChan:
			if !ok {
				if h.mtcm != nil {
					h.mtcm.FinishTurn(threadID, query, lastResponse)
				}
				if h.tracing != nil {
					h.tracing.EndSpan(span, lastResponse)
				}
				sendSSE(c.Writer, flusher, "complete", map[string]interface{}{
					"message": "Graph execution completed",
				})
				return
			}
			eventData := serializeEvent(evt)
			sendSSE(c.Writer, flusher, "node", eventData)

			if evt != nil && evt.StateDelta != nil {
				sendSSE(c.Writer, flusher, "state_update", map[string]interface{}{
					"state_delta": evt.StateDelta,
				})
				if raw, ok := evt.StateDelta["nl2sql_state"]; ok && len(raw) > 0 {
					var delta nl2sqltypes.NL2SQLState
					if err := json.Unmarshal(raw, &delta); err == nil {
						if delta.Result != "" {
							lastResponse = delta.Result
						} else if delta.SQLExecuteOutput != "" {
							lastResponse = delta.SQLExecuteOutput
						} else if delta.PythonAnalyzeOutput != "" {
							lastResponse = delta.PythonAnalyzeOutput
						}
					}
				}
			}
		}
	}
}

func (h *GraphHandler) getOrBuildGraph(humanReview bool) (*graph.Graph, error) {
	if !humanReview {
		h.mu.Lock()
		if h.cached != nil {
			g := h.cached
			h.mu.Unlock()
			return g, nil
		}
		h.mu.Unlock()
	}

	graphBuilder := nl2sql.NewNL2SQLGraph(h.llm, h.retriever, h.dbManager)
	if humanReview {
		graphBuilder.WithHumanReview(true)
	}

	g, err := graphBuilder.Build()
	if err != nil {
		return nil, err
	}

	if !humanReview {
		h.mu.Lock()
		h.cached = g
		h.mu.Unlock()
	}

	return g, nil
}

func sendSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonBytes)
	flusher.Flush()
}

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
