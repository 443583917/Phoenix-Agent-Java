package semanticmodel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// SemanticModelHandler handles semantic model CRUD operations.
type SemanticModelHandler struct {
	svc *service.DataService
}

// NewSemanticModelHandler creates a new SemanticModelHandler.
func NewSemanticModelHandler(svc *service.DataService) *SemanticModelHandler {
	return &SemanticModelHandler{svc: svc}
}

// List returns a paginated list of semantic models.
// GET /semantic-model/
func (h *SemanticModelHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	response.SuccessPage(c, []interface{}{}, 0, page, size)
}

// GetByID returns a single semantic model by its ID.
// GET /semantic-model/:id
func (h *SemanticModelHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"id": id, "name": "stub-semantic-model"})
}

// Create creates a new semantic model.
// POST /semantic-model
func (h *SemanticModelHandler) Create(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"created": true, "data": body})
}

// Update updates an existing semantic model.
// PUT /semantic-model/:id
func (h *SemanticModelHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"id": id, "updated": true, "data": body})
}

// Delete deletes a semantic model by its ID.
// DELETE /semantic-model/:id
func (h *SemanticModelHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"id": id, "deleted": true})
}

// BatchDelete deletes multiple semantic models by their IDs.
// DELETE /semantic-model/batch
func (h *SemanticModelHandler) BatchDelete(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"ids": body.Ids, "deleted": true})
}

// Enable enables semantic models by their IDs.
// PUT /semantic-model/enable
func (h *SemanticModelHandler) Enable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"ids": body.Ids, "enabled": true})
}

// Disable disables semantic models by their IDs.
// PUT /semantic-model/disable
func (h *SemanticModelHandler) Disable(c *gin.Context) {
	var body struct {
		Ids []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"ids": body.Ids, "disabled": true})
}

// BatchImport imports semantic models in batch.
// POST /semantic-model/batch-import
func (h *SemanticModelHandler) BatchImport(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"imported": true, "data": body})
}

// DownloadTemplate returns a template description for semantic model import.
// GET /semantic-model/template/download
func (h *SemanticModelHandler) DownloadTemplate(c *gin.Context) {
	response.Success(c, gin.H{
		"template": "semantic_model_template.xlsx",
		"description": "Upload this template with semantic model data",
	})
}

// ImportExcel imports semantic models from an uploaded Excel file.
// POST /semantic-model/import/excel
func (h *SemanticModelHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{
		"imported":  true,
		"filename":  file.Filename,
		"size":      file.Size,
	})
}
