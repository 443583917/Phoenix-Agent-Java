package modelconfig

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// ModelConfigHandler handles LLM model configuration CRUD and testing.
// All methods are stubs for Phase 5D; full implementations deferred to later phases.
type ModelConfigHandler struct {
	svc *service.DataService
}

// NewModelConfigHandler creates a new ModelConfigHandler.
func NewModelConfigHandler(svc *service.DataService) *ModelConfigHandler {
	return &ModelConfigHandler{svc: svc}
}

// List returns a paginated list of model configurations (stub).
// GET /model-config/list
func (h *ModelConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	response.SuccessPage(c, []interface{}{}, 0, page, size)
}

// Add creates a new model configuration (stub).
// POST /model-config/add
func (h *ModelConfigHandler) Add(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"id":      "stub-model-config-id",
		"message": "Model config add stub - Phase 5",
		"body":    body,
	})
}

// Update updates an existing model configuration (stub).
// PUT /model-config/update
func (h *ModelConfigHandler) Update(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, body)
}

// Delete soft-deletes a model configuration by ID (stub).
// DELETE /model-config/:id
func (h *ModelConfigHandler) Delete(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Model config delete stub - Phase 5",
	})
}

// Activate activates a model configuration by ID (stub).
// POST /model-config/activate/:id
func (h *ModelConfigHandler) Activate(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Model config activate stub - Phase 5",
	})
}

// Test tests a model configuration connection (stub).
// POST /model-config/test
func (h *ModelConfigHandler) Test(c *gin.Context) {
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)
	response.Success(c, map[string]interface{}{
		"success": true,
		"latency": "100ms",
	})
}

// CheckReady checks whether the model configuration service is ready (stub).
// GET /model-config/check-ready
func (h *ModelConfigHandler) CheckReady(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"ready": true,
	})
}
