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

// DepartmentHandler handles department CRUD and sync HTTP requests.
type DepartmentHandler struct {
	svc *service.PrivilegeService
}

// NewDepartmentHandler creates a new DepartmentHandler.
func NewDepartmentHandler(svc *service.PrivilegeService) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

// OrgTree returns the full organization tree (companies with nested departments).
// GET /api/privilege/department/orgTree
func (h *DepartmentHandler) OrgTree(c *gin.Context) {
	tree, err := h.svc.GetOrgTree(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tree)
}

// Page returns a paginated list of departments.
// GET /api/privilege/department/page
func (h *DepartmentHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")
	code := c.Query("code")
	companyID := c.Query("companyId")

	list, total, err := h.svc.PageDepartments(c.Request.Context(), page, size, name, code, companyID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// Tree returns the department tree for a specific company.
// GET /api/privilege/department/tree
func (h *DepartmentHandler) Tree(c *gin.Context) {
	companyID := c.Query("companyId")
	tree, err := h.svc.GetDeptTree(c.Request.Context(), companyID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, tree)
}

// GetByID returns a single department by its primary key.
// GET /api/privilege/department/:id
func (h *DepartmentHandler) GetByID(c *gin.Context) {
	dept, err := h.svc.GetDepartmentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, dept)
}

// GetByPID returns departments with the given parent ID.
// GET /api/privilege/department/pid/:pid
func (h *DepartmentHandler) GetByPID(c *gin.Context) {
	list, err := h.svc.GetDepartmentsByPID(c.Request.Context(), c.Param("pid"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByCompany returns all departments for a company.
// GET /api/privilege/department/company/:companyId
func (h *DepartmentHandler) GetByCompany(c *gin.Context) {
	list, err := h.svc.GetDepartmentsByCompanyID(c.Request.Context(), c.Param("companyId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// GetByCode returns a single department by its code.
// GET /api/privilege/department/code/:code
func (h *DepartmentHandler) GetByCode(c *gin.Context) {
	dept, err := h.svc.GetDepartmentByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, dept)
}

// Create creates a new department.
// POST /api/privilege/department
func (h *DepartmentHandler) Create(c *gin.Context) {
	var dto model.PrivilegeDepartmentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateDepartment(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing department. The department ID comes from the JSON body.
// PUT /api/privilege/department
func (h *DepartmentHandler) Update(c *gin.Context) {
	var dto model.PrivilegeDepartmentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateDepartment(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a department by its ID.
// DELETE /api/privilege/department/:id
func (h *DepartmentHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteDepartment(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Sync batch-syncs department data.
// POST /api/privilege/department/sync
func (h *DepartmentHandler) Sync(c *gin.Context) {
	var dtos []model.PrivilegeDepartmentDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SyncDepartments(c.Request.Context(), dtos); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// SyncChildren batch-syncs child departments for a given parent department.
// POST /api/privilege/department/sync-children/:deptId
func (h *DepartmentHandler) SyncChildren(c *gin.Context) {
	deptID := c.Param("deptId")
	var dtos []model.PrivilegeDepartmentDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	// Set the parent ID for each DTO before syncing.
	for i := range dtos {
		dtos[i].PID = deptID
	}

	if err := h.svc.SyncDepartments(c.Request.Context(), dtos); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
