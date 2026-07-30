package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// PlatformInfoHandler handles platform-info CRUD and query HTTP requests.
type PlatformInfoHandler struct {
	svc *service.PlatformService
}

// NewPlatformInfoHandler creates a new PlatformInfoHandler.
func NewPlatformInfoHandler(svc *service.PlatformService) *PlatformInfoHandler {
	return &PlatformInfoHandler{svc: svc}
}

// Page returns a paginated list of platforms.
// GET /platform/platform-info/page
func (h *PlatformInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.PlatformInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PagePlatformInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single platform by its primary key.
// GET /platform/platform-info/:id
func (h *PlatformInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetPlatformInfoByID(c.Request.Context(), c.Param("id"))
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

// GetByType returns platforms of a specific type.
// GET /platform/platform-info/type/:type
func (h *PlatformInfoHandler) GetByType(c *gin.Context) {
	list, err := h.svc.FindPlatformInfoByType(c.Request.Context(), c.Param("type"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByTypeEnabled returns enabled platforms of a specific type.
// GET /platform/platform-info/type/:type/enabled
func (h *PlatformInfoHandler) GetByTypeEnabled(c *gin.Context) {
	list, err := h.svc.FindPlatformInfoEnabledByType(c.Request.Context(), c.Param("type"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetEnabledPlatform returns all enabled platforms.
// GET /platform/platform-info/getEnabledPlatform
func (h *PlatformInfoHandler) GetEnabledPlatform(c *gin.Context) {
	// Query all enabled platforms via a broad page query filtered by status="0".
	// Use a large page size to effectively return all records.
	list, _, err := h.svc.PagePlatformInfo(c.Request.Context(), 1, 9999, &model.PlatformInfo{Status: "0"})
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// Create creates a new platform.
// POST /platform/platform-info
func (h *PlatformInfoHandler) Create(c *gin.Context) {
	var entity model.PlatformInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreatePlatformInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing platform.
// PUT /platform/platform-info
func (h *PlatformInfoHandler) Update(c *gin.Context) {
	var entity model.PlatformInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdatePlatformInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a platform by its ID.
// DELETE /platform/platform-info/:id
func (h *PlatformInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeletePlatformInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
