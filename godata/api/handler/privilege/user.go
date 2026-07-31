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

// UserHandler handles user CRUD and password-management HTTP requests.
type UserHandler struct {
	svc *service.PrivilegeService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc *service.PrivilegeService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Page returns a paginated list of users.
// GET /api/privilege/user/page
func (h *UserHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.PrivilegeUserPageQuery
	_ = c.ShouldBindQuery(&query)
	query.Page = page
	query.Size = size

	list, total, err := h.svc.PageUsers(c.Request.Context(), query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single user by its primary key.
// GET /api/privilege/user/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.svc.GetUserByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, user)
}

// GetByCode returns a single user by its employee code.
// GET /api/privilege/user/code/:code
func (h *UserHandler) GetByCode(c *gin.Context) {
	user, err := h.svc.GetUserByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, user)
}

// GetByUsername returns a single user by its username.
// GET /api/privilege/user/username/:username
func (h *UserHandler) GetByUsername(c *gin.Context) {
	user, err := h.svc.GetUserByUsername(c.Request.Context(), c.Param("username"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, user)
}

// Create creates a new user with a default password.
// POST /api/privilege/user
func (h *UserHandler) Create(c *gin.Context) {
	var dto model.PrivilegeUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateUser(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing user. The user ID comes from the JSON body.
// PUT /api/privilege/user
func (h *UserHandler) Update(c *gin.Context) {
	var dto model.PrivilegeUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateUser(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a user by its ID.
// DELETE /api/privilege/user/:id
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// UpdatePassword changes the current user's password after verifying the old one.
// PUT /api/privilege/user/password
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var dto model.PasswordUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
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

// SetPassword sets a user's password directly (admin operation).
// PUT /api/privilege/user/setPassword
func (h *UserHandler) SetPassword(c *gin.Context) {
	var dto model.PasswordUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdatePassword(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ResetPassword resets a user's password to a random value and returns the plaintext.
// PUT /api/privilege/user/reset-password/:id
func (h *UserHandler) ResetPassword(c *gin.Context) {
	newPwd, err := h.svc.ResetPassword(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.ErrorWithMsg(c, errcode.NotFound, "用户不存在")
		return
	}
	response.Success(c, newPwd)
}
