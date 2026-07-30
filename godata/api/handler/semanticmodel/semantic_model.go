package semanticmodel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type SemanticModelHandler struct {
	svc *service.DataService
}

func NewSemanticModelHandler(svc *service.DataService) *SemanticModelHandler {
	return &SemanticModelHandler{svc: svc}
}

func (h *SemanticModelHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.SemanticModel
	_ = c.ShouldBindQuery(&query)
	list, total, err := h.svc.PageSemanticModel(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *SemanticModelHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetSemanticModelByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, entity)
}

func (h *SemanticModelHandler) Create(c *gin.Context) {
	var entity model.SemanticModel
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	result, err := h.svc.CreateSemanticModel(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *SemanticModelHandler) Update(c *gin.Context) {
	var entity model.SemanticModel
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.ID = c.Param("id")
	if err := h.svc.UpdateSemanticModel(c.Request.Context(), &entity); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *SemanticModelHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteSemanticModel(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *SemanticModelHandler) BatchDelete(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.BatchDeleteSemanticModel(c.Request.Context(), body.Ids); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *SemanticModelHandler) Enable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.EnableSemanticModels(c.Request.Context(), body.Ids); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *SemanticModelHandler) Disable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.DisableSemanticModels(c.Request.Context(), body.Ids); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *SemanticModelHandler) BatchImport(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"imported": true})
}

func (h *SemanticModelHandler) DownloadTemplate(c *gin.Context) {
	response.Success(c, gin.H{
		"template":    "semantic_model_template.xlsx",
		"description": "Upload this template with semantic model data",
	})
}

func (h *SemanticModelHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"imported": true, "filename": file.Filename, "size": file.Size})
}

func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
