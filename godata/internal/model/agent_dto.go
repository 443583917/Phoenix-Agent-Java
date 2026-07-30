package model

// ChatModelRequest — request for agent chat
type ChatModelRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Content   string `json:"content" binding:"required"`
	AgentSn   string `json:"agentSn" binding:"required"`
}

// ConfirmRequest — human-in-the-loop confirm/reject request
type ConfirmRequest struct {
	AgentSn   string `json:"agentSn" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
	SessionID string `json:"sessionId" binding:"required"`
	Confirmed bool   `json:"confirmed"`
}

// HarnessRequest — harness chat request
type HarnessRequest struct {
	HarnessSn string `json:"harnessSn" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
	SessionID string `json:"sessionId" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

// StreamRequest — internal stream call request to agent runtime
type StreamRequest struct {
	AgentSN   string
	UserID    string
	SessionID string
	Message   string
}

// SSEEvent — generic SSE event envelope
type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ContentEvent — streaming text content chunk
type ContentEvent struct {
	Content string `json:"content"`
	End     bool   `json:"end"`
}

// ToolCallEvent — tool call requiring human confirmation
type ToolCallEvent struct {
	ToolName  string                 `json:"toolName"`
	Args      map[string]interface{} `json:"args"`
	ConfirmID string                 `json:"confirmId"`
	UserID    string                 `json:"userId"`
	SessionID string                 `json:"sessionId"`
}

// ConfirmResultEvent — result of a confirm/reject action
type ConfirmResultEvent struct {
	ConfirmID string `json:"confirmId"`
	Approved  bool   `json:"approved"`
}

// ErrorEvent — SSE error payload
type ErrorEvent struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
