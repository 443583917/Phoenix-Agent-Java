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

// AccountTenantInfoHandler handles account-tenant-association CRUD HTTP requests.
type AccountTenantInfoHandler struct {
	svc *service.PlatformService
}

// NewAccountTenantInfoHandler creates a new AccountTenantInfoHandler.
func NewAccountTenantInfoHandler(svc *service.PlatformService) *AccountTenantInfoHandler {
	return &AccountTenantInfoHandler{svc: svc}
}

// Page returns a paginated list of account-tenant associations.
// GET /platform/account-tenant-info/page
func (h *AccountTenantInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.AccountTenantInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAccountTenantInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single account-tenant association by its primary key.
// GET /platform/account-tenant-info/:id
func (h *AccountTenantInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAccountTenantInfoByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new account-tenant association.
// POST /platform/account-tenant-info
func (h *AccountTenantInfoHandler) Create(c *gin.Context) {
	var entity model.AccountTenantInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateAccountTenantInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing account-tenant association.
// PUT /platform/account-tenant-info
func (h *AccountTenantInfoHandler) Update(c *gin.Context) {
	var entity model.AccountTenantInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAccountTenantInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an account-tenant association by its ID.
// DELETE /platform/account-tenant-info/:id
func (h *AccountTenantInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAccountTenantInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
