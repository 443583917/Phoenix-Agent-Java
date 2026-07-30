package privilege

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// UserRoleHandler handles user-role association HTTP requests.
type UserRoleHandler struct {
	svc *service.PrivilegeService
}

// NewUserRoleHandler creates a new UserRoleHandler.
func NewUserRoleHandler(svc *service.PrivilegeService) *UserRoleHandler {
	return &UserRoleHandler{svc: svc}
}

// GetByUser returns all role associations for the given user.
// GET /api/privilege/user-role/user/:userId
func (h *UserRoleHandler) GetByUser(c *gin.Context) {
	list, err := h.svc.GetUserRolesByUserID(c.Request.Context(), c.Param("userId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByRole returns all user associations for the given role.
// GET /api/privilege/user-role/role/:roleId
func (h *UserRoleHandler) GetByRole(c *gin.Context) {
	list, err := h.svc.GetUserRolesByRoleID(c.Request.Context(), c.Param("roleId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// Create creates a single user-role association.
// POST /api/privilege/user-role
func (h *UserRoleHandler) Create(c *gin.Context) {
	var dto model.PrivilegeUserRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SaveUserRoles(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete removes a single user-role association.
// DELETE /api/privilege/user-role/:id?userId=xxx
func (h *UserRoleHandler) Delete(c *gin.Context) {
	dto := model.PrivilegeUserRoleDTO{
		UserID: c.Query("userId"),
		RoleID: c.Param("id"),
	}
	if dto.UserID == "" || dto.RoleID == "" {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.DeleteUserRoles(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// BatchSave replaces all roles for a user with the given list.
// POST /api/privilege/user-role/batch-save
func (h *UserRoleHandler) BatchSave(c *gin.Context) {
	var dto model.UserRoleBatchDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.BatchSaveUserRoles(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// BatchRemove removes a set of roles from a user.
// DELETE /api/privilege/user-role/batch-remove
func (h *UserRoleHandler) BatchRemove(c *gin.Context) {
	var dto model.UserRoleBatchDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.BatchDeleteUserRoles(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
