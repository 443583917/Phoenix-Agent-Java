package prompt

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// PromptConfigHandler handles user-defined prompt template configuration CRUD.
// All methods are stubs for Phase 5D; full implementations deferred to later phases.
type PromptConfigHandler struct {
	svc *service.DataService
}

// NewPromptConfigHandler creates a new PromptConfigHandler.
func NewPromptConfigHandler(svc *service.DataService) *PromptConfigHandler {
	return &PromptConfigHandler{svc: svc}
}

// Save creates or updates a prompt configuration (stub).
// POST /prompt-config/save
func (h *PromptConfigHandler) Save(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"id":      "stub-prompt-config-id",
		"message": "Prompt config save stub - Phase 5",
		"body":    body,
	})
}

// GetByID returns a single prompt configuration by ID (stub).
// GET /prompt-config/:id
func (h *PromptConfigHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, map[string]interface{}{
		"id":          id,
		"name":        "stub-prompt",
		"promptType":  "system",
		"description": "Stub prompt config - Phase 5",
	})
}

// List returns a paginated list of prompt configurations (stub).
// GET /prompt-config/list
func (h *PromptConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	response.SuccessPage(c, []interface{}{}, 0, page, size)
}

// ListByType returns a paginated list of prompt configurations filtered by type (stub).
// GET /prompt-config/list-by-type/:type
func (h *PromptConfigHandler) ListByType(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	promptType := c.Param("type")
	response.SuccessPage(c, []interface{}{}, 0, page, size)
	_ = promptType // placeholder: will be used for filtering in real implementation
}

// GetActiveByType returns the single active prompt configuration for a type (stub).
// GET /prompt-config/active/:type
func (h *PromptConfigHandler) GetActiveByType(c *gin.Context) {
	promptType := c.Param("type")
	response.Success(c, map[string]interface{}{
		"id":          "stub-active-config",
		"name":        "Active " + promptType,
		"promptType":  promptType,
		"enabled":     true,
		"description": "Stub active prompt config - Phase 5",
	})
}

// GetActiveAllByType returns all active prompt configurations for a type (stub).
// GET /prompt-config/active-all/:type
func (h *PromptConfigHandler) GetActiveAllByType(c *gin.Context) {
	promptType := c.Param("type")
	response.Success(c, []map[string]interface{}{
		{
			"id":         "stub-active-1",
			"name":       "Active " + promptType + " #1",
			"promptType": promptType,
			"enabled":    true,
		},
	})
}

// Delete soft-deletes a prompt configuration by ID (stub).
// DELETE /prompt-config/:id
func (h *PromptConfigHandler) Delete(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Prompt config delete stub - Phase 5",
	})
}

// Enable enables a prompt configuration by ID (stub).
// POST /prompt-config/:id/enable
func (h *PromptConfigHandler) Enable(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Prompt config enable stub - Phase 5",
	})
}

// Disable disables a prompt configuration by ID (stub).
// POST /prompt-config/:id/disable
func (h *PromptConfigHandler) Disable(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Prompt config disable stub - Phase 5",
	})
}

// GetTypes returns the list of supported prompt types (stub).
// GET /prompt-config/types
func (h *PromptConfigHandler) GetTypes(c *gin.Context) {
	response.Success(c, []string{"system", "user", "assistant", "tool", "reasoning"})
}

// BatchEnable enables multiple prompt configurations by IDs (stub).
// POST /prompt-config/batch-enable
func (h *PromptConfigHandler) BatchEnable(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"message": "Batch enable stub - Phase 5",
		"body":    body,
	})
}

// BatchDisable disables multiple prompt configurations by IDs (stub).
// POST /prompt-config/batch-disable
func (h *PromptConfigHandler) BatchDisable(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"message": "Batch disable stub - Phase 5",
		"body":    body,
	})
}

// SetPriority sets the priority of a prompt configuration (stub).
// POST /prompt-config/:id/priority
func (h *PromptConfigHandler) SetPriority(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Set priority stub - Phase 5",
		"body":    body,
	})
}

// SetDisplayOrder sets the display order of a prompt configuration (stub).
// POST /prompt-config/:id/display-order
func (h *PromptConfigHandler) SetDisplayOrder(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, map[string]interface{}{
		"id":      c.Param("id"),
		"message": "Set display order stub - Phase 5",
		"body":    body,
	})
}
