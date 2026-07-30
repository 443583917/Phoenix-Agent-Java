package platform

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// TenantInfoHandler handles tenant-info CRUD HTTP requests.
type TenantInfoHandler struct {
	svc *service.PlatformService
}

// NewTenantInfoHandler creates a new TenantInfoHandler.
func NewTenantInfoHandler(svc *service.PlatformService) *TenantInfoHandler {
	return &TenantInfoHandler{svc: svc}
}

// Page returns a paginated list of tenants.
// GET /platform/tenant-info/page
func (h *TenantInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.TenantInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageTenantInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single tenant by its primary key.
// GET /platform/tenant-info/:id
func (h *TenantInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetTenantInfoByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, entity)
}

// Create creates a new tenant.
// POST /platform/tenant-info
func (h *TenantInfoHandler) Create(c *gin.Context) {
	var entity model.TenantInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateTenantInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing tenant.
// PUT /platform/tenant-info
func (h *TenantInfoHandler) Update(c *gin.Context) {
	var entity model.TenantInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateTenantInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a tenant by its ID.
// DELETE /platform/tenant-info/:id
func (h *TenantInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteTenantInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
