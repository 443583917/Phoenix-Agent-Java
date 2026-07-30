package platform

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// GroupInfoHandler handles group-info CRUD and management HTTP requests.
type GroupInfoHandler struct {
	svc *service.PlatformService
}

// NewGroupInfoHandler creates a new GroupInfoHandler.
func NewGroupInfoHandler(svc *service.PlatformService) *GroupInfoHandler {
	return &GroupInfoHandler{svc: svc}
}

// Page returns a paginated list of groups.
// GET /platform/group-info/page
func (h *GroupInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.GroupInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageGroupInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single group by its primary key.
// GET /platform/group-info/:id
func (h *GroupInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetGroupInfoByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, entity)
}

// GetBySN returns a group by its SN number.
// GET /platform/group-info/sn/:sn
func (h *GroupInfoHandler) GetBySN(c *gin.Context) {
	sn := c.Param("sn")
	query := &model.GroupInfo{SN: sn}
	list, _, err := h.svc.PageGroupInfo(c.Request.Context(), 1, 1, query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	if len(list) == 0 {
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, list[0])
}

// Create creates a new group.
// POST /platform/group-info
func (h *GroupInfoHandler) Create(c *gin.Context) {
	var entity model.GroupInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateGroupInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing group.
// PUT /platform/group-info
func (h *GroupInfoHandler) Update(c *gin.Context) {
	var entity model.GroupInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateGroupInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a group by its ID.
// DELETE /platform/group-info/:id
func (h *GroupInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteGroupInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ToggleStatus toggles the status of a group (0 <-> 1).
// PUT /platform/group-info/:id/toggle-status
func (h *GroupInfoHandler) ToggleStatus(c *gin.Context) {
	if err := h.svc.ToggleGroupInfoStatus(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// RemoveAgent removes an agent from a group by deleting the association.
// DELETE /platform/group-info/:groupId/agent/:agentId
func (h *GroupInfoHandler) RemoveAgent(c *gin.Context) {
	groupID := c.Param("groupId")
	agentIDStr := c.Param("agentId")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	// Find the association by group + agent IDs.
	ga, err := h.svc.GetGroupAgentInfoByGroupIdAndAgentId(c.Request.Context(), groupID, agentID)
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	if ga == nil {
		response.Error(c, errcode.NotFound)
		return
	}

	// Delete by association ID.
	if err := h.svc.DeleteGroupAgentInfo(c.Request.Context(), ga.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
