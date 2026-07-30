package prompt

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type PromptConfigHandler struct {
	svc *service.DataService
}

func NewPromptConfigHandler(svc *service.DataService) *PromptConfigHandler {
	return &PromptConfigHandler{svc: svc}
}

func (h *PromptConfigHandler) Save(c *gin.Context) {
	var entity model.UserPromptConfig
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	result, err := h.svc.SavePromptConfig(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *PromptConfigHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetPromptConfigByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, entity)
}

func (h *PromptConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := h.svc.PagePromptConfig(c.Request.Context(), page, size, "")
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *PromptConfigHandler) ListByType(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	promptType := c.Param("type")
	list, total, err := h.svc.ListPromptConfigByType(c.Request.Context(), promptType, page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *PromptConfigHandler) GetActiveByType(c *gin.Context) {
	entity, err := h.svc.GetActivePromptConfigByType(c.Request.Context(), c.Param("type"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

func (h *PromptConfigHandler) GetActiveAllByType(c *gin.Context) {
	list, err := h.svc.GetActiveAllPromptConfigByType(c.Request.Context(), c.Param("type"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *PromptConfigHandler) Delete(c *gin.Context) {
	if err := h.svc.DeletePromptConfig(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) Enable(c *gin.Context) {
	if err := h.svc.EnablePromptConfig(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) Disable(c *gin.Context) {
	if err := h.svc.DisablePromptConfig(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) GetTypes(c *gin.Context) {
	response.Success(c, []string{"system", "user", "assistant", "tool", "reasoning"})
}

func (h *PromptConfigHandler) BatchEnable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.BatchEnablePromptConfig(c.Request.Context(), body.Ids); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) BatchDisable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.BatchDisablePromptConfig(c.Request.Context(), body.Ids); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) SetPriority(c *gin.Context) {
	var body struct {
		Priority int `json:"priority"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.SetPromptConfigPriority(c.Request.Context(), c.Param("id"), body.Priority); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *PromptConfigHandler) SetDisplayOrder(c *gin.Context) {
	var body struct {
		DisplayOrder int `json:"displayOrder"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.SetPromptConfigDisplayOrder(c.Request.Context(), c.Param("id"), body.DisplayOrder); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
