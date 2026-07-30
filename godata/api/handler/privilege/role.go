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

// RoleHandler handles role CRUD and ACL HTTP requests.
type RoleHandler struct {
	svc *service.PrivilegeService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(svc *service.PrivilegeService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// Page returns a paginated list of roles.
// POST /api/privilege/role/page
func (h *RoleHandler) Page(c *gin.Context) {
	var query model.PrivilegeRoleQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	list, total, err := h.svc.PageRoles(c.Request.Context(), query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, query.Page, query.Size)
}

// GetByID returns a single role by its primary key.
// GET /api/privilege/role/:id
func (h *RoleHandler) GetByID(c *gin.Context) {
	role, err := h.svc.GetRoleByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, role)
}

// GetByCompany returns roles for the specified company.
// GET /api/privilege/role/company/:companyId
func (h *RoleHandler) GetByCompany(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	roles, err := h.svc.GetRolesByCompanyID(c.Request.Context(), companyID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, roles)
}

// Create creates a new role.
// POST /api/privilege/role
func (h *RoleHandler) Create(c *gin.Context) {
	var dto model.PrivilegeRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateRole(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing role. The role ID comes from the JSON body.
// PUT /api/privilege/role
func (h *RoleHandler) Update(c *gin.Context) {
	var dto model.PrivilegeRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateRole(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a role by its ID.
// DELETE /api/privilege/role/:id
func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteRole(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// GetAcls returns the ACL entries for a role.
// GET /api/privilege/role/:id/acls
func (h *RoleHandler) GetAcls(c *gin.Context) {
	acls, err := h.svc.GetRoleAcls(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, acls)
}
