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

// AccountInfoHandler handles account-info CRUD and management HTTP requests.
type AccountInfoHandler struct {
	svc *service.PlatformService
}

// NewAccountInfoHandler creates a new AccountInfoHandler.
func NewAccountInfoHandler(svc *service.PlatformService) *AccountInfoHandler {
	return &AccountInfoHandler{svc: svc}
}

// Page returns a paginated list of accounts.
// GET /platform/account-info/page
func (h *AccountInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.AccountInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageAccountInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single account by its primary key.
// GET /platform/account-info/:id
func (h *AccountInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetAccountInfoByID(c.Request.Context(), c.Param("id"))
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

// GetByUsername returns an account by username.
// GET /platform/account-info/username/:username
func (h *AccountInfoHandler) GetByUsername(c *gin.Context) {
	entity, err := h.svc.GetAccountInfoByUsername(c.Request.Context(), c.Param("username"))
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

// GetByCode returns an account by employee code.
// GET /platform/account-info/code/:code
func (h *AccountInfoHandler) GetByCode(c *gin.Context) {
	entity, err := h.svc.GetAccountInfoByCode(c.Request.Context(), c.Param("code"))
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

// GetByThirdPartyID returns an account by third-party ID.
// GET /platform/account-info/third-party/:thirdPartyId
func (h *AccountInfoHandler) GetByThirdPartyID(c *gin.Context) {
	entity, err := h.svc.GetAccountInfoByThirdPartyID(c.Request.Context(), c.Param("thirdPartyId"))
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

// List returns all accounts.
// GET /platform/account-info/list
func (h *AccountInfoHandler) List(c *gin.Context) {
	list, err := h.svc.ListAccountInfo(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetMyAgents returns the agents accessible to the current account.
// GET /platform/account-info/getMyAgents
func (h *AccountInfoHandler) GetMyAgents(c *gin.Context) {
	userID, _ := c.Get("user_id")
	accountID := strconv.FormatUint(userID.(uint64), 10)

	list, err := h.svc.GetMyAgents(c.Request.Context(), accountID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetUnGroupPageByGroupId returns accounts not assigned to the specified group.
// GET /platform/account-info/getUnGroupPageByGroupId
func (h *AccountInfoHandler) GetUnGroupPageByGroupId(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	groupID := c.Query("groupId")

	var query model.AccountInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.GetUnGroupPage(c.Request.Context(), page, size, groupID, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// Create creates a new account.
// POST /platform/account-info
func (h *AccountInfoHandler) Create(c *gin.Context) {
	var entity model.AccountInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateAccountInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing account.
// PUT /platform/account-info
func (h *AccountInfoHandler) Update(c *gin.Context) {
	var entity model.AccountInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateAccountInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an account by its ID.
// DELETE /platform/account-info/:id
func (h *AccountInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteAccountInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// BatchStatus updates the status of multiple accounts.
// PUT /platform/account-info/batch-status
func (h *AccountInfoHandler) BatchStatus(c *gin.Context) {
	var req struct {
		IDs    []string `json:"ids" binding:"required"`
		Status string   `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.BatchUpdateAccountStatus(c.Request.Context(), req.IDs, req.Status); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// GetByStatus returns accounts filtered by status.
// GET /platform/account-info/status/:status
func (h *AccountInfoHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")
	var query model.AccountInfo
	query.Status = status
	list, _, err := h.svc.PageAccountInfo(c.Request.Context(), 1, 1000, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}
