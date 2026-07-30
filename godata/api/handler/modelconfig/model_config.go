package modelconfig

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type ModelConfigHandler struct {
	svc *service.DataService
}

func NewModelConfigHandler(svc *service.DataService) *ModelConfigHandler {
	return &ModelConfigHandler{svc: svc}
}

func (h *ModelConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := h.svc.PageModelConfig(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *ModelConfigHandler) Add(c *gin.Context) {
	var entity model.ModelConfig
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	result, err := h.svc.CreateModelConfig(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *ModelConfigHandler) Update(c *gin.Context) {
	var entity model.ModelConfig
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdateModelConfig(c.Request.Context(), &entity); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *ModelConfigHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteModelConfig(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *ModelConfigHandler) Activate(c *gin.Context) {
	if err := h.svc.ActivateModelConfig(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *ModelConfigHandler) Test(c *gin.Context) {
	var entity model.ModelConfig
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	start := time.Now()
	result, err := h.svc.TestModelConfig(c.Request.Context(), &entity)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		response.Success(c, gin.H{"success": false, "message": err.Error(), "latency": latency})
		return
	}
	response.Success(c, gin.H{"success": result, "latency": latency})
}

func (h *ModelConfigHandler) CheckReady(c *gin.Context) {
	ready, err := h.svc.CheckModelConfigReady(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"ready": ready})
}

func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
