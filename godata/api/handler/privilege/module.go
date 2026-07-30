package privilege

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// ModuleHandler handles module CRUD HTTP requests.
type ModuleHandler struct {
	svc *service.PrivilegeService
}

// NewModuleHandler creates a new ModuleHandler.
func NewModuleHandler(svc *service.PrivilegeService) *ModuleHandler {
	return &ModuleHandler{svc: svc}
}

// Page returns a paginated list of modules.
// GET /api/privilege/module/page
func (h *ModuleHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")
	code := c.Query("code")
	systemID := c.Query("systemId")

	list, total, err := h.svc.PageModules(c.Request.Context(), page, size, name, code, systemID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// Tree returns the full module tree.
// GET /api/privilege/module/tree
func (h *ModuleHandler) Tree(c *gin.Context) {
	tree, err := h.svc.GetModuleTree(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tree)
}

// TreeWithACL returns the module tree with ACL info.
// GET /api/privilege/module/tree/acl
func (h *ModuleHandler) TreeWithACL(c *gin.Context) {
	tree, err := h.svc.GetModuleTreeWithACL(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tree)
}

// GetByID returns a single module by its primary key.
// GET /api/privilege/module/:id
func (h *ModuleHandler) GetByID(c *gin.Context) {
	module, err := h.svc.GetModuleByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, module)
}

// GetBySystem returns modules for the specified system.
// GET /api/privilege/module/system/:systemId
func (h *ModuleHandler) GetBySystem(c *gin.Context) {
	modules, err := h.svc.GetModulesBySystemID(c.Request.Context(), c.Param("systemId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, modules)
}

// GetByPID returns modules with the given parent ID.
// GET /api/privilege/module/pid/:pid
func (h *ModuleHandler) GetByPID(c *gin.Context) {
	modules, err := h.svc.GetModulesByPID(c.Request.Context(), c.Param("pid"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, modules)
}

// Create creates a new module.
// POST /api/privilege/module
func (h *ModuleHandler) Create(c *gin.Context) {
	var dto model.PrivilegeModuleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateModule(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing module. The module ID comes from the JSON body.
// PUT /api/privilege/module
func (h *ModuleHandler) Update(c *gin.Context) {
	var dto model.PrivilegeModuleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateModule(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
