package privilege

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// LoginLogHandler handles login log query and management HTTP requests.
type LoginLogHandler struct {
	svc *service.PrivilegeService
}

// NewLoginLogHandler creates a new LoginLogHandler.
func NewLoginLogHandler(svc *service.PrivilegeService) *LoginLogHandler {
	return &LoginLogHandler{svc: svc}
}

// Page returns a paginated list of login logs.
// GET /api/privilege/login-log/page
func (h *LoginLogHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.svc.PageLoginLogs(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single login log by its primary key.
// GET /api/privilege/login-log/:id
func (h *LoginLogHandler) GetByID(c *gin.Context) {
	log, err := h.svc.GetLoginLogByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, log)
}

// Create creates a new login log entry.
// POST /api/privilege/login-log
func (h *LoginLogHandler) Create(c *gin.Context) {
	var dto model.PrivilegeLoginLogDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateLoginLog(c.Request.Context(), dto); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a login log by its ID.
// DELETE /api/privilege/login-log/:id
func (h *LoginLogHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteLoginLog(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
