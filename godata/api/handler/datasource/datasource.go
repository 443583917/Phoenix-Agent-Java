package datasource

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// DatasourceHandler handles datasource CRUD, connection test, and logical relations.
// All methods are stubs for Phase 5D; full implementations deferred to later phases.
type DatasourceHandler struct {
	svc *service.DataService
}

// NewDatasourceHandler creates a new DatasourceHandler.
func NewDatasourceHandler(svc *service.DataService) *DatasourceHandler {
	return &DatasourceHandler{svc: svc}
}

// GetTypes returns the list of supported datasource types (stub).
// GET /datasource/types
func (h *DatasourceHandler) GetTypes(c *gin.Context) {
	types := []string{"MySQL", "PostgreSQL", "Oracle", "SQLServer", "ClickHouse", "MongoDB", "Redis"}
	response.Success(c, types)
}

// List returns a paginated list of datasources (stub).
// GET /datasource/
func (h *DatasourceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	response.SuccessPage(c, []interface{}{}, 0, page, size)
}

// GetByID returns a stub datasource entity by ID.
// GET /datasource/:id
func (h *DatasourceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, map[string]interface{}{
		"id":   id,
		"name": "stub-datasource",
		"type": "MySQL",
	})
}

// GetTables returns a stub table list for the given datasource.
// GET /datasource/:id/tables
func (h *DatasourceHandler) GetTables(c *gin.Context) {
	response.Success(c, []map[string]interface{}{
		{"tableName": "stub_table_1", "comment": "Stub table 1 - Phase 5"},
		{"tableName": "stub_table_2", "comment": "Stub table 2 - Phase 5"},
	})
}

// GetColumns returns a stub column list for the given table of a datasource.
// GET /datasource/:id/tables/:tableName/columns
func (h *DatasourceHandler) GetColumns(c *gin.Context) {
	tableName := c.Param("tableName")
	response.Success(c, []map[string]interface{}{
		{"columnName": "id", "dataType": "bigint", "tableName": tableName},
		{"columnName": "name", "dataType": "varchar", "tableName": tableName},
	})
}

// Create creates a new datasource (stub).
// POST /datasource
func (h *DatasourceHandler) Create(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"id":      "stub-generated-id",
		"message": "Datasource create stub - Phase 5",
		"body":    body,
	})
}

// Update updates an existing datasource (stub).
// PUT /datasource/:id
func (h *DatasourceHandler) Update(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	body["id"] = c.Param("id")
	response.Success(c, body)
}

// Delete soft-deletes a datasource by ID (stub).
// DELETE /datasource/:id
func (h *DatasourceHandler) Delete(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Datasource delete stub - Phase 5",
	})
}

// TestConnection tests the connection for a given datasource (stub).
// POST /datasource/:id/test
func (h *DatasourceHandler) TestConnection(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"success": true,
		"message": "Connection test stub - Phase 5",
	})
}

// ListLogicalRelations returns an empty list of logical relations (stub).
// GET /datasource/:id/logical-relations
func (h *DatasourceHandler) ListLogicalRelations(c *gin.Context) {
	response.Success(c, []interface{}{})
}

// CreateLogicalRelation creates a new logical relation (stub).
// POST /datasource/:id/logical-relations
func (h *DatasourceHandler) CreateLogicalRelation(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	body["datasourceId"] = c.Param("id")
	response.Success(c, body)
}

// UpdateLogicalRelation updates logical relations in bulk (stub).
// PUT /datasource/:id/logical-relations
func (h *DatasourceHandler) UpdateLogicalRelation(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, body)
}

// DeleteLogicalRelation deletes a logical relation (stub).
// DELETE /datasource/:id/logical-relations
func (h *DatasourceHandler) DeleteLogicalRelation(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"message": "Logical relation delete stub - Phase 5",
	})
}
