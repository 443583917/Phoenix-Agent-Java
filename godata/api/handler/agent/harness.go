package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/agent/runner"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
)

// HarnessHandler handles harness/agent chat SSE streaming with
// Human-in-the-Loop (HITL) tool confirmation support.
type HarnessHandler struct {
	agentManager *runtime.AgentManager
	hitlHandler  *runner.HitlHandler
}

// NewHarnessHandler creates a new HarnessHandler.
func NewHarnessHandler(agentManager *runtime.AgentManager, hitlHandler *runner.HitlHandler) *HarnessHandler {
	return &HarnessHandler{
		agentManager: agentManager,
		hitlHandler:  hitlHandler,
	}
}

// harnessChatEvent is the SSE event payload for harness chat responses,
// matching the Java controller's output format.
type harnessChatEvent struct {
	Content     string          `json:"content"`
	End         bool            `json:"end"`
	NeedConfirm bool            `json:"needConfirm,omitempty"`
	ToolCalls   json.RawMessage `json:"toolCalls,omitempty"`
	Buttons     []button        `json:"buttons,omitempty"`
}

type button struct {
	Text   string `json:"text"`
	Action string `json:"action"`
	Type   string `json:"type"`
}

// confirmResultEvent is the SSE event payload for confirm results.
type confirmResultEvent struct {
	Content string `json:"content"`
	End     bool   `json:"end"`
}

// Chat streams harness chat responses via SSE, including tool confirmation
// buttons when the agent requires human approval.
//
// POST /api/admin/harness/chat
//
// SSE event format:
//   - Content: {"content":"...","end":false}
//   - Tool confirm: {"content":"","end":false,"needConfirm":true,"toolCalls":[...],"buttons":[...]}
//   - Complete: {"content":"","end":true}
func (h *HarnessHandler) Chat(c *gin.Context) {
	var req model.HarnessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)

	streamReq := model.StreamRequest{
		AgentSN:   req.HarnessSn,
		UserID:    userIDStr,
		SessionID: req.SessionID,
		Message:   req.Message,
	}

	events, err := h.agentManager.StreamCall(c.Request.Context(), streamReq)
	if err != nil {
		response.Error(c, errcode.ModelError)
		return
	}

	// Set SSE headers and start streaming.
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

			var sseData []byte
			var marshalErr error

			switch evt.Event {
			case "tool_call":
				// Convert tool_call event to harness chat SSE with confirm buttons.
				tc, ok := evt.Data.(model.ToolCallEvent)
				if !ok {
					continue
				}
				// Generate a unique confirm ID and register with HITL handler.
				confirmID := strconv.FormatUint(id.MustGenerateID(), 36)
				tc.ConfirmID = confirmID
				tc.UserID = userIDStr
				tc.SessionID = req.SessionID

				_ = h.hitlHandler.RegisterPending(confirmID, req.SessionID)

				// Build tool calls JSON for the SSE response.
				toolCallList := []map[string]interface{}{
					{
						"toolName": tc.ToolName,
						"args":     tc.Args,
					},
				}
				tcJSON, _ := json.Marshal(toolCallList)

				event := harnessChatEvent{
					Content:     "",
					End:         false,
					NeedConfirm: true,
					ToolCalls:   tcJSON,
					Buttons: []button{
						{Text: "确认", Action: "confirm", Type: "primary"},
						{Text: "取消", Action: "cancel", Type: "danger"},
					},
				}
				sseData, marshalErr = json.Marshal(event)

			case "content":
				// Stream content event as-is.
				ce, ok := evt.Data.(model.ContentEvent)
				if !ok {
					// Attempt a generic marshal.
					sseData, marshalErr = json.Marshal(evt.Data)
				} else {
					event := harnessChatEvent{
						Content: ce.Content,
						End:     ce.End,
					}
					sseData, marshalErr = json.Marshal(event)
				}

			case "error":
				// Forward error events with the "error" SSE event type.
				sseData, marshalErr = json.Marshal(evt.Data)
				if marshalErr == nil {
					fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", sseData)
					flusher.Flush()
				}
				continue

			default:
				// Skip unrecognized event types.
				continue
			}

			if marshalErr != nil {
				continue
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", sseData)
			flusher.Flush()
		}
	}
}

// Confirm processes a human-in-the-loop confirmation request.
//
// POST /api/admin/harness/confirm
//
// Accepts a JSON ConfirmRequest body. On success, returns a single SSE event
// indicating the confirmation result.
func (h *HarnessHandler) Confirm(c *gin.Context) {
	var req model.ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)
	req.UserID = userIDStr

	// Resolve the pending confirmation via HITL handler.
	if err := h.hitlHandler.HandleConfirm(c.Request.Context(), req.SessionID, req.Confirmed); err != nil {
		response.ErrorWithMsg(c, errcode.InternalError, err.Error())
		return
	}

	// Return a single SSE event with the confirmation result.
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	result := confirmResultEvent{
		Content: "",
		End:     true,
	}
	data, _ := json.Marshal(result)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()
}
