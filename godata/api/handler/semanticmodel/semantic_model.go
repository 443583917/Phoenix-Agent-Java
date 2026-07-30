package semanticmodel

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

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
	var entities []*model.SemanticModel
	if err := c.ShouldBindJSON(&entities); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	count, err := h.svc.BatchCreateSemanticModels(c.Request.Context(), entities)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"imported": count})
}

func (h *SemanticModelHandler) DownloadTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="semantic_model_template.csv"`)

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"agentId", "datasourceId", "tableName", "columnName", "businessName", "synonyms", "businessDescription", "columnComment", "dataType"})
	_ = w.Write([]string{"1", "1", "users", "name", "用户姓名", "姓名;用户名", "存储用户的姓名信息", "", "varchar(128)"})
	w.Flush()
}

func (h *SemanticModelHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	f, err := file.Open()
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		response.ErrorWithMsg(c, errcode.InvalidParams, fmt.Sprintf("CSV解析失败: %v", err))
		return
	}

	var entities []*model.SemanticModel
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 9 {
			continue
		}
		agentID, _ := strconv.ParseInt(strings.TrimSpace(row[0]), 10, 64)
		dsID, _ := strconv.Atoi(strings.TrimSpace(row[1]))
		entities = append(entities, &model.SemanticModel{
			AgentId:             agentID,
			DatasourceId:        dsID,
			Table:               strings.TrimSpace(row[2]),
			ColumnName:          strings.TrimSpace(row[3]),
			BusinessName:        strings.TrimSpace(row[4]),
			Synonyms:            strings.TrimSpace(row[5]),
			BusinessDescription: strings.TrimSpace(row[6]),
			ColumnComment:       strings.TrimSpace(row[7]),
			DataType:            strings.TrimSpace(row[8]),
		})
	}

	count, err := h.svc.BatchCreateSemanticModels(c.Request.Context(), entities)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, gin.H{"imported": count, "filename": file.Filename})
}

func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
