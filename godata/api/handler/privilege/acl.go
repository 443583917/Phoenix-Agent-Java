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

// ACLHandler handles ACL-related HTTP requests.
type ACLHandler struct {
	svc *service.PrivilegeService
}

// NewACLHandler creates a new ACLHandler.
func NewACLHandler(svc *service.PrivilegeService) *ACLHandler {
	return &ACLHandler{svc: svc}
}

// Page returns a paginated list of ACLs.
// GET /api/privilege/acl/page
func (h *ACLHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.svc.PageACLs(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single ACL by its primary key.
// GET /api/privilege/acl/:id
func (h *ACLHandler) GetByID(c *gin.Context) {
	acl, err := h.svc.GetACLByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, acl)
}

// GetByRelease returns all ACLs for the specified release.
// GET /api/privilege/acl/release/:releaseId
func (h *ACLHandler) GetByRelease(c *gin.Context) {
	acls, err := h.svc.GetACLsByReleaseID(c.Request.Context(), c.Param("releaseId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, acls)
}

// GetByReleaseAndModule returns ACLs for the specified release and module.
// GET /api/privilege/acl/release/module/:releaseId/:moduleId
func (h *ACLHandler) GetByReleaseAndModule(c *gin.Context) {
	acls, err := h.svc.GetACLsByReleaseIDAndModuleID(c.Request.Context(), c.Param("releaseId"), c.Param("moduleId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, acls)
}

// Update updates an existing ACL.
// PUT /api/privilege/acl
func (h *ACLHandler) Update(c *gin.Context) {
	var dto model.PrivilegeAclDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateACL(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an ACL by its ID.
// DELETE /api/privilege/acl/:id
func (h *ACLHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteACL(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// SaveAll saves all ACL entries for a release with a given check status.
// POST /api/privilege/acl/saveAll/:releaseId/:checkStatus
func (h *ACLHandler) SaveAll(c *gin.Context) {
	releaseID := c.Param("releaseId")
	checkStatus, _ := strconv.Atoi(c.Param("checkStatus"))

	var dtos []model.PrivilegeAclDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SaveAllACLs(c.Request.Context(), dtos, releaseID, checkStatus); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// SaveModule saves a single ACL module entry.
// POST /api/privilege/acl/saveModule
func (h *ACLHandler) SaveModule(c *gin.Context) {
	var dto model.PrivilegeAclDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SaveModuleACL(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
