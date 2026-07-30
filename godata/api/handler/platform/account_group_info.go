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

// AccountGroupInfoHandler handles account-group-association CRUD HTTP requests.
type AccountGroupInfoHandler struct {
	svc *service.PlatformService
}

// NewAccountGroupInfoHandler creates a new AccountGroupInfoHandler.
func NewAccountGroupInfoHandler(svc *service.PlatformService) *AccountGroupInfoHandler {
	return &AccountGroupInfoHandler{svc: svc}
}

// Page returns a paginated list of account-group associations.
// GET /platform/account-group-info/page
func (h *AccountGroupInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.AccountGroupInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAccountGroupInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single account-group association by its primary key.
// GET /platform/account-group-info/:id
func (h *AccountGroupInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAccountGroupInfoByID(c.Request.Context(), c.Param("id"))
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

// Create creates a new account-group association.
// POST /platform/account-group-info
func (h *AccountGroupInfoHandler) Create(c *gin.Context) {
	var entity model.AccountGroupInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateAccountGroupInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing account-group association.
// PUT /platform/account-group-info
func (h *AccountGroupInfoHandler) Update(c *gin.Context) {
	var entity model.AccountGroupInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAccountGroupInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an account-group association by its ID.
// DELETE /platform/account-group-info/:id
func (h *AccountGroupInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAccountGroupInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
