package privilege

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// CompanyHandler handles company CRUD HTTP requests.
type CompanyHandler struct {
	svc *service.PrivilegeService
}

// NewCompanyHandler creates a new CompanyHandler.
func NewCompanyHandler(svc *service.PrivilegeService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// Page returns a paginated list of companies.
// POST /api/privilege/company/page
func (h *CompanyHandler) Page(c *gin.Context) {
	var query model.PrivilegeCompanyQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	list, total, err := h.svc.PageCompanies(c.Request.Context(), query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, query.Page, query.Size)
}

// GetByID returns a single company by its primary key.
// GET /api/privilege/company/:id
func (h *CompanyHandler) GetByID(c *gin.Context) {
	company, err := h.svc.GetCompanyByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, company)
}

// GetByCode returns a single company by its code.
// GET /api/privilege/company/code/:code
func (h *CompanyHandler) GetByCode(c *gin.Context) {
	company, err := h.svc.GetCompanyByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, company)
}

// Create creates a new company.
// POST /api/privilege/company
func (h *CompanyHandler) Create(c *gin.Context) {
	var dto model.PrivilegeCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if _, err := h.svc.CreateCompany(c.Request.Context(), dto); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Update updates an existing company. The company ID comes from the JSON body.
// PUT /api/privilege/company
func (h *CompanyHandler) Update(c *gin.Context) {
	var dto model.PrivilegeCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if err := h.svc.UpdateCompany(c.Request.Context(), dto, dto.ID); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// Delete soft-deletes a company by its ID.
// DELETE /api/privilege/company/:id
func (h *CompanyHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteCompany(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
