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

// GroupAgentInfoHandler handles group-agent-association CRUD HTTP requests.
type GroupAgentInfoHandler struct {
	svc *service.PlatformService
}

// NewGroupAgentInfoHandler creates a new GroupAgentInfoHandler.
func NewGroupAgentInfoHandler(svc *service.PlatformService) *GroupAgentInfoHandler {
	return &GroupAgentInfoHandler{svc: svc}
}

// Page returns a paginated list of group-agent associations.
// GET /platform/group-agent-info/page
func (h *GroupAgentInfoHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.GroupAgentInfo
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageGroupAgentInfo(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single group-agent association by its primary key.
// GET /platform/group-agent-info/:id
func (h *GroupAgentInfoHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetGroupAgentInfoByID(c.Request.Context(), c.Param("id"))
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

// GetByGroupID returns all agents associated with a group.
// GET /platform/group-agent-info/group/:groupId
func (h *GroupAgentInfoHandler) GetByGroupID(c *gin.Context) {
	list, err := h.svc.FindGroupAgentInfoByGroupID(c.Request.Context(), c.Param("groupId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByAgentID returns all groups associated with an agent.
// GET /platform/group-agent-info/agent/:agentId
func (h *GroupAgentInfoHandler) GetByAgentID(c *gin.Context) {
	agentIDStr := c.Param("agentId")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	list, err := h.svc.FindGroupAgentInfoByAgentID(c.Request.Context(), agentID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// Create creates a new group-agent association.
// POST /platform/group-agent-info
func (h *GroupAgentInfoHandler) Create(c *gin.Context) {
	var entity model.GroupAgentInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateGroupAgentInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

// Update updates an existing group-agent association.
// PUT /platform/group-agent-info
func (h *GroupAgentInfoHandler) Update(c *gin.Context) {
	var entity model.GroupAgentInfo
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateGroupAgentInfo(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a group-agent association by its ID.
// DELETE /platform/group-agent-info/:id
func (h *GroupAgentInfoHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteGroupAgentInfo(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// DeleteByGroupAndAgent removes an agent from a group by group + agent IDs.
// DELETE /platform/group-agent-info/group/:groupId/agent/:agentId
func (h *GroupAgentInfoHandler) DeleteByGroupAndAgent(c *gin.Context) {
	groupID := c.Param("groupId")
	agentIDStr := c.Param("agentId")
	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

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

// List returns all group-agent associations.
// GET /platform/group-agent-info/list
func (h *GroupAgentInfoHandler) List(c *gin.Context) {
	list, err := h.svc.ListGroupAgentInfo(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}
