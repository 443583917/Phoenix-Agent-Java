package platform

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// AccountLoginHandler handles platform authentication HTTP requests.
type AccountLoginHandler struct {
	svc        *service.PlatformService
	jwtManager *jwt.JWTManager
}

// NewAccountLoginHandler creates a new AccountLoginHandler.
func NewAccountLoginHandler(svc *service.PlatformService, jwtManager *jwt.JWTManager) *AccountLoginHandler {
	return &AccountLoginHandler{svc: svc, jwtManager: jwtManager}
}

// Login authenticates an account and returns a JWT token.
// POST /auth/login
func (h *AccountLoginHandler) Login(c *gin.Context) {
	var dto model.AccountLoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	account, err := h.svc.Login(c.Request.Context(), dto)
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}

	userID, _ := strconv.ParseUint(account.ID, 10, 64)
	token, _ := h.jwtManager.GenerateToken(userID, account.Username)

	response.Success(c, gin.H{
		"token":     token,
		"accountId": account.ID,
		"username":  account.Username,
		"realName":  account.RealName,
		"nickName":  account.NickName,
		"avatarUrl": account.AvatarURL,
	})
}

// Logout handles logout (no-op for stateless JWT).
// POST /auth/logout
func (h *AccountLoginHandler) Logout(c *gin.Context) {
	response.Success(c, "退出成功")
}

// ThirdLogin authenticates via a third-party identity.
// POST /auth/thirdLogin
func (h *AccountLoginHandler) ThirdLogin(c *gin.Context) {
	var dto model.ThirdPartyLoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	account, err := h.svc.ThirdPartyLogin(c.Request.Context(), dto)
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}

	userID, _ := strconv.ParseUint(account.ID, 10, 64)
	token, _ := h.jwtManager.GenerateToken(userID, account.Username)

	response.Success(c, gin.H{
		"token":     token,
		"accountId": account.ID,
		"username":  account.Username,
		"realName":  account.RealName,
		"nickName":  account.NickName,
		"avatarUrl": account.AvatarURL,
	})
}

// UpdatePassword changes the current user's password after verifying the old one.
// PUT /auth/updatePassword
func (h *AccountLoginHandler) UpdatePassword(c *gin.Context) {
	var dto model.UpdatePwdDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	// Use the user ID from JWT if not provided in the body.
	if dto.UserID == "" {
		uid, _ := c.Get("user_id")
		dto.UserID = strconv.FormatUint(uid.(uint64), 10)
	}

	if err := h.svc.UpdatePassword(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, "密码修改成功")
}
