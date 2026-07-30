package privilege

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	svc        *service.PrivilegeService
	jwtManager *jwt.JWTManager
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *service.PrivilegeService, jwtManager *jwt.JWTManager) *AuthHandler {
	return &AuthHandler{svc: svc, jwtManager: jwtManager}
}

// Captcha returns a captcha key and a base64-encoded image.
// GET /api/privilege/auth/captcha
func (h *AuthHandler) Captcha(c *gin.Context) {
	captchaKey := "captcha_" + strconv.FormatUint(id.MustGenerateID(), 10)
	c.SetCookie("captcha_key", captchaKey, 300, "/", "", false, true)
	response.Success(c, gin.H{"captchaKey": captchaKey, "image": "placeholder"})
}

// Login authenticates the user and returns a JWT token with user info.
// POST /api/privilege/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var dto model.LoginInfoDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}

	ip := c.ClientIP()
	userInfo, err := h.svc.Login(c.Request.Context(), dto, ip)
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}

	userID, _ := strconv.ParseUint(userInfo.UserID, 10, 64)
	token, _ := h.jwtManager.GenerateToken(userID, userInfo.Username)
	userInfo.Token = token
	response.Success(c, userInfo)
}

// Logout handles logout (no-op for stateless JWT).
// POST /api/privilege/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, "退出成功")
}

// Menus returns the current user's menu tree.
// GET /api/privilege/auth/menus
func (h *AuthHandler) Menus(c *gin.Context) {
	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}
	userID, _ := c.Get("user_id")
	userIDStr := strconv.FormatUint(userID.(uint64), 10)

	menus, err := h.svc.GetUserMenus(c.Request.Context(), userIDStr)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, menus)
}

// GetLoginUserInfo returns the current user's profile.
// GET /api/privilege/auth/getLoginUserInfo
func (h *AuthHandler) GetLoginUserInfo(c *gin.Context) {
	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}
	userID, _ := c.Get("user_id")
	userIDStr := strconv.FormatUint(userID.(uint64), 10)

	user, err := h.svc.GetUserByID(c.Request.Context(), userIDStr)
	if err != nil {
		response.Error(c, errcode.Unauthorized)
		return
	}
	response.Success(c, user)
}
