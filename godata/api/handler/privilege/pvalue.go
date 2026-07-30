package privilege

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// PvalueHandler handles permission value CRUD HTTP requests.
type PvalueHandler struct {
	svc *service.PrivilegeService
}

// NewPvalueHandler creates a new PvalueHandler.
func NewPvalueHandler(svc *service.PrivilegeService) *PvalueHandler {
	return &PvalueHandler{svc: svc}
}

// Page returns a paginated list of permission values.
// POST /api/privilege/pvalue/page
func (h *PvalueHandler) Page(c *gin.Context) {
	var query model.PrivilegePvalueQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	list, total, err := h.svc.PagePvalues(c.Request.Context(), query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, query.Page, query.Size)
}

// GetByID returns a single permission value by its primary key.
// GET /api/privilege/pvalue/:id
func (h *PvalueHandler) GetByID(c *gin.Context) {
	pv, err := h.svc.GetPvalueByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, pv)
}

// GetBySystem returns all permission values for the given system.
// GET /api/privilege/pvalue/system
func (h *PvalueHandler) GetBySystem(c *gin.Context) {
	systemID := c.Query("systemId")
	list, err := h.svc.GetPvaluesBySystemID(c.Request.Context(), systemID)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// Create creates a new permission value.
// POST /api/privilege/pvalue
func (h *PvalueHandler) Create(c *gin.Context) {
	var dto model.PrivilegePvalueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreatePvalue(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing permission value. The ID comes from the JSON body.
// PUT /api/privilege/pvalue
func (h *PvalueHandler) Update(c *gin.Context) {
	var dto model.PrivilegePvalueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdatePvalue(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a permission value by its ID.
// DELETE /api/privilege/pvalue/:id
func (h *PvalueHandler) Delete(c *gin.Context) {
	if err := h.svc.DeletePvalue(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
