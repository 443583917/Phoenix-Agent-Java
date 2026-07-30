package chat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type ChatHandler struct {
	svc *service.DataService
}

func NewChatHandler(svc *service.DataService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

func (h *ChatHandler) ListSessions(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Param("id"))
	userID := getUserID(c)
	list, err := h.svc.ListChatSessions(c.Request.Context(), agentID, userID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *ChatHandler) CreateSession(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Param("id"))
	var entity model.ChatSession
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.AgentId = agentID
	entity.UserId = getUserID(c)
	result, err := h.svc.CreateChatSession(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *ChatHandler) DeleteAllSessions(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteAllSessions(c.Request.Context(), agentID); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	sessionID := c.Param("id")
	list, err := h.svc.GetSessionMessages(c.Request.Context(), sessionID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *ChatHandler) CreateMessage(c *gin.Context) {
	sessionID := c.Param("id")
	var entity model.ChatMessage
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.SessionId = sessionID
	result, err := h.svc.AddChatMessage(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *ChatHandler) PinSession(c *gin.Context) {
	sessionID := c.Param("id")
	isPinned := c.Query("isPinned") == "true"
	if err := h.svc.PinSession(c.Request.Context(), sessionID, isPinned); err != nil {
		handleUsecaseErr(c, err)
		return
	}
	response.Success(c, gin.H{"id": sessionID, "isPinned": isPinned})
}

func (h *ChatHandler) RenameSession(c *gin.Context) {
	sessionID := c.Param("id")
	title := c.Query("title")
	if err := h.svc.RenameSession(c.Request.Context(), sessionID, title); err != nil {
		handleUsecaseErr(c, err)
		return
	}
	response.Success(c, gin.H{"id": sessionID, "title": title})
}

func (h *ChatHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	if err := h.svc.DeleteSession(c.Request.Context(), sessionID); err != nil {
		handleUsecaseErr(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *ChatHandler) GenerateReportHTML(c *gin.Context) {
	sessionID := c.Param("id")
	messages, err := h.svc.GetSessionMessages(c.Request.Context(), sessionID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}

	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>Session Report</title>`)
	sb.WriteString(`<style>body{font-family:sans-serif;max-width:800px;margin:0 auto;padding:20px}
		.msg{margin:10px 0;padding:10px;border-radius:8px}
		.user{background:#e3f2fd}.assistant{background:#f5f5f5}
		h1{color:#333;border-bottom:2px solid #2196f3;padding-bottom:10px}</style></head><body>`)
	sb.WriteString(`<h1>Session Report</h1>`)

	if len(messages) == 0 {
		sb.WriteString(`<p>No messages in this session.</p>`)
	} else {
		for _, msg := range messages {
			role := msg.Role
			if role == "" {
				role = "user"
			}
			sb.WriteString(fmt.Sprintf(`<div class="msg %s"><strong>%s:</strong><br>%s</div>`, role, role, msg.Content))
		}
	}

	sb.WriteString(fmt.Sprintf(`<hr><small>Generated at %s | Session: %s</small>`, time.Now().Format("2006-01-02 15:04:05"), sessionID))
	sb.WriteString(`</body></html>`)

	response.Success(c, gin.H{"html": sb.String()})
}

func getUserID(c *gin.Context) string {
	if uid, exists := c.Get("user_id"); exists {
		return strconv.FormatUint(uid.(uint64), 10)
	}
	return ""
}

func handleUsecaseErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
