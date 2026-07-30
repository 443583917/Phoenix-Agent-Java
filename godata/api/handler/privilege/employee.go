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

// EmployeeHandler handles employee binding CRUD and sync HTTP requests.
type EmployeeHandler struct {
	svc *service.PrivilegeService
}

// NewEmployeeHandler creates a new EmployeeHandler.
func NewEmployeeHandler(svc *service.PrivilegeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

// Page returns a paginated list of employee bindings.
// GET /api/privilege/employee/page
func (h *EmployeeHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.svc.PageEmployees(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

// GetByID returns a single employee binding by its primary key.
// GET /api/privilege/employee/:id
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	emp, err := h.svc.GetEmployeeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, emp)
}

// GetByEmpCode returns a single employee binding by its employee code.
// GET /api/privilege/employee/emp-code/:empCode
func (h *EmployeeHandler) GetByEmpCode(c *gin.Context) {
	emp, err := h.svc.GetEmployeeByEmpCode(c.Request.Context(), c.Param("empCode"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, emp)
}

// Create creates a new employee binding.
// POST /api/privilege/employee
func (h *EmployeeHandler) Create(c *gin.Context) {
	var dto model.PrivilegeEmployeeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateEmployee(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing employee binding. The ID comes from the JSON body.
// PUT /api/privilege/employee
func (h *EmployeeHandler) Update(c *gin.Context) {
	var dto model.PrivilegeEmployeeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateEmployee(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes an employee binding by its ID.
// DELETE /api/privilege/employee/:id
func (h *EmployeeHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteEmployee(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Sync batch-syncs employee binding data.
// POST /api/privilege/employee/sync
func (h *EmployeeHandler) Sync(c *gin.Context) {
	var dtos []model.PrivilegeEmployeeDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.SyncEmployees(c.Request.Context(), dtos); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// SyncByDept batch-syncs employee bindings for a given department.
// POST /api/privilege/employee/sync-by-dept/:deptId
func (h *EmployeeHandler) SyncByDept(c *gin.Context) {
	deptID := c.Param("deptId")
	var dtos []model.PrivilegeEmployeeDTO
	if err := c.ShouldBindJSON(&dtos); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	for i := range dtos {
		dtos[i].DeptID = deptID
	}

	if err := h.svc.SyncEmployees(c.Request.Context(), dtos); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
