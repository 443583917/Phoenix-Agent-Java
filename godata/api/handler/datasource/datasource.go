package datasource

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type DatasourceHandler struct {
	svc *service.DataService
}

func NewDatasourceHandler(svc *service.DataService) *DatasourceHandler {
	return &DatasourceHandler{svc: svc}
}

func (h *DatasourceHandler) GetTypes(c *gin.Context) {
	types := []string{"MySQL", "PostgreSQL", "Oracle", "SQLServer", "ClickHouse", "MongoDB", "Redis"}
	response.Success(c, types)
}

func (h *DatasourceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.Datasource
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageDatasource(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *DatasourceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	entity, err := h.svc.GetDatasourceByID(c.Request.Context(), id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, entity)
}

func (h *DatasourceHandler) GetTables(c *gin.Context) {
	id := c.Param("id")
	tables, err := h.svc.GetDatasourceTables(c.Request.Context(), id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, tables)
}

func (h *DatasourceHandler) GetColumns(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Param("tableName")
	columns, err := h.svc.GetDatasourceColumns(c.Request.Context(), id, tableName)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, columns)
}

func (h *DatasourceHandler) Create(c *gin.Context) {
	var entity model.Datasource
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	result, err := h.svc.CreateDatasource(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *DatasourceHandler) Update(c *gin.Context) {
	var entity model.Datasource
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.ID = c.Param("id")
	if err := h.svc.UpdateDatasource(c.Request.Context(), &entity); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *DatasourceHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteDatasource(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *DatasourceHandler) TestConnection(c *gin.Context) {
	err := h.svc.TestDatasourceConnection(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Success(c, gin.H{"success": false, "message": err.Error()})
		return
	}
	response.Success(c, gin.H{"success": true, "message": "连接成功"})
}

func (h *DatasourceHandler) ListLogicalRelations(c *gin.Context) {
	dsID, _ := strconv.Atoi(c.Param("id"))
	list, err := h.svc.ListLogicalRelations(c.Request.Context(), dsID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *DatasourceHandler) CreateLogicalRelation(c *gin.Context) {
	var entity model.LogicalRelation
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	dsID, _ := strconv.Atoi(c.Param("id"))
	entity.DatasourceId = dsID
	result, err := h.svc.CreateLogicalRelation(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *DatasourceHandler) UpdateLogicalRelation(c *gin.Context) {
	dsID, _ := strconv.Atoi(c.Param("id"))
	var relations []*model.LogicalRelation
	if err := c.ShouldBindJSON(&relations); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdateLogicalRelations(c.Request.Context(), dsID, relations); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *DatasourceHandler) DeleteLogicalRelation(c *gin.Context) {
	dsID, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteLogicalRelations(c.Request.Context(), dsID); err != nil {
		response.Error(c, errcode.InternalError)
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
