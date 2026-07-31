package chat

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/service/session"
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

	if entity.Role == "user" {
		sess, sErr := h.svc.GetChatSessionByID(c.Request.Context(), sessionID)
		if sErr == nil && sess != nil && sess.Title == "" {
			titleSvc := session.NewTitleService()
			title, _ := titleSvc.GenerateTitle(c.Request.Context(), sessionID, entity.Content)
			_ = h.svc.RenameSession(c.Request.Context(), sessionID, title)
		}
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
	html := session.BuildReportHTML(messages, sessionID)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="report_%s.html"`, sessionID))
	c.String(http.StatusOK, html)
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
