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

// DictionaryHandler handles dictionary entry CRUD HTTP requests.
type DictionaryHandler struct {
	svc *service.PrivilegeService
}

// NewDictionaryHandler creates a new DictionaryHandler.
func NewDictionaryHandler(svc *service.PrivilegeService) *DictionaryHandler {
	return &DictionaryHandler{svc: svc}
}

// Page returns a paginated list of dictionary entries.
// GET /api/privilege/dictionary/page
func (h *DictionaryHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.svc.PageDictionaries(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single dictionary entry by its primary key.
// GET /api/privilege/dictionary/:id
func (h *DictionaryHandler) GetByID(c *gin.Context) {
	dict, err := h.svc.GetDictionaryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, dict)
}

// GetBySystem returns all dictionary entries for the given system SN.
// GET /api/privilege/dictionary/system/:systemSn
func (h *DictionaryHandler) GetBySystem(c *gin.Context) {
	list, err := h.svc.GetDictionariesBySystemSN(c.Request.Context(), c.Param("systemSn"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByPCode returns all dictionary entries with the given parent code.
// GET /api/privilege/dictionary/pcode/:pcode
func (h *DictionaryHandler) GetByPCode(c *gin.Context) {
	list, err := h.svc.GetDictionariesByPCode(c.Request.Context(), c.Param("pcode"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// Create creates a new dictionary entry.
// POST /api/privilege/dictionary
func (h *DictionaryHandler) Create(c *gin.Context) {
	var dto model.PrivilegeDictionaryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateDictionary(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing dictionary entry. The ID comes from the JSON body.
// PUT /api/privilege/dictionary
func (h *DictionaryHandler) Update(c *gin.Context) {
	var dto model.PrivilegeDictionaryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateDictionary(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a dictionary entry by its ID.
// DELETE /api/privilege/dictionary/:id
func (h *DictionaryHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteDictionary(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
