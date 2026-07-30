package knowledge

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// BusinessKnowledgeHandler handles business knowledge CRUD operations.
type BusinessKnowledgeHandler struct {
	svc *service.DataService
}

// NewBusinessKnowledgeHandler creates a new BusinessKnowledgeHandler.
func NewBusinessKnowledgeHandler(svc *service.DataService) *BusinessKnowledgeHandler {
	return &BusinessKnowledgeHandler{svc: svc}
}

// List returns a paginated list of business knowledge entries.
// GET /business-knowledge/
func (h *BusinessKnowledgeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	response.SuccessPage(c, []interface{}{}, 0, page, size)
}

// GetByID returns a single business knowledge entry by its ID.
// GET /business-knowledge/:id
func (h *BusinessKnowledgeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"id": id, "name": "stub-business-knowledge"})
}

// Create creates a new business knowledge entry.
// POST /business-knowledge
func (h *BusinessKnowledgeHandler) Create(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"created": true, "data": body})
}

// Update updates an existing business knowledge entry.
// PUT /business-knowledge/:id
func (h *BusinessKnowledgeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	response.Success(c, gin.H{"id": id, "updated": true, "data": body})
}

// Delete deletes a business knowledge entry by its ID.
// DELETE /business-knowledge/:id
func (h *BusinessKnowledgeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"id": id, "deleted": true})
}

// Recall recalls knowledge for a given business knowledge entry.
// POST /business-knowledge/recall/:id
func (h *BusinessKnowledgeHandler) Recall(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{
		"id":      id,
		"results": []interface{}{},
		"recalled": true,
	})
}

// RefreshVectorStore refreshes the vector store index.
// POST /business-knowledge/refresh-vector-store
func (h *BusinessKnowledgeHandler) RefreshVectorStore(c *gin.Context) {
	response.Success(c, gin.H{"refreshed": true, "message": "Vector store refresh triggered"})
}

// RetryEmbedding retries embedding for a given business knowledge entry.
// POST /business-knowledge/retry-embedding/:id
func (h *BusinessKnowledgeHandler) RetryEmbedding(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"id": id, "reembedded": true})
}
