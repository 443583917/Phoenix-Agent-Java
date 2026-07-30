package common

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

// PlatformSyncHandler handles platform sync stub HTTP requests.
type PlatformSyncHandler struct {
	svc *service.PlatformService
}

// NewPlatformSyncHandler creates a new PlatformSyncHandler.
func NewPlatformSyncHandler(svc *service.PlatformService) *PlatformSyncHandler {
	return &PlatformSyncHandler{svc: svc}
}

// SyncAll syncs all departments and users (stub).
// POST /platform/sync/all
func (h *PlatformSyncHandler) SyncAll(c *gin.Context) {
	response.Success(c, gin.H{"message": "同步请求已提交"})
}

// SyncDepartments syncs departments only (stub).
// POST /platform/sync/departments
func (h *PlatformSyncHandler) SyncDepartments(c *gin.Context) {
	response.Success(c, gin.H{"message": "部门同步请求已提交"})
}

// SyncUsers syncs users only (stub).
// POST /platform/sync/users
func (h *PlatformSyncHandler) SyncUsers(c *gin.Context) {
	response.Success(c, gin.H{"message": "用户同步请求已提交"})
}
