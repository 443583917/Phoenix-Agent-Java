package common

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

type PlatformSyncHandler struct {
	svc *service.PlatformService
}

func NewPlatformSyncHandler(svc *service.PlatformService) *PlatformSyncHandler {
	return &PlatformSyncHandler{svc: svc}
}

func (h *PlatformSyncHandler) SyncAll(c *gin.Context) {
	ctx := c.Request.Context()
	_ = h.svc.SyncDepartments(ctx)
	_ = h.svc.SyncUsers(ctx)
	response.Success(c, gin.H{"message": "全量同步请求已提交"})
}

func (h *PlatformSyncHandler) SyncDepartments(c *gin.Context) {
	if err := h.svc.SyncDepartments(c.Request.Context()); err != nil {
		response.Success(c, gin.H{"success": false, "message": err.Error()})
		return
	}
	response.Success(c, gin.H{"success": true, "message": "部门同步完成"})
}

func (h *PlatformSyncHandler) SyncUsers(c *gin.Context) {
	if err := h.svc.SyncUsers(c.Request.Context()); err != nil {
		response.Success(c, gin.H{"success": false, "message": err.Error()})
		return
	}
	response.Success(c, gin.H{"success": true, "message": "用户同步完成"})
}
