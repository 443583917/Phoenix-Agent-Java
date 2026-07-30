package kg

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

// ──────────────────────────── KGEntity ────────────────────────────

type KGEntityHandler struct {
	svc *service.KgService
}

func NewKGEntityHandler(svc *service.KgService) *KGEntityHandler {
	return &KGEntityHandler{svc: svc}
}

func (h *KGEntityHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.KGEntity
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageKGEntity(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *KGEntityHandler) List(c *gin.Context) {
	list, err := h.svc.ListKGEntity(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *KGEntityHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetKGEntityByID(c.Request.Context(), c.Param("id"))
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

func (h *KGEntityHandler) Create(c *gin.Context) {
	var entity model.KGEntity
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if _, err := h.svc.CreateKGEntity(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, entity)
}

func (h *KGEntityHandler) Update(c *gin.Context) {
	var entity model.KGEntity
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdateKGEntity(c.Request.Context(), &entity); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *KGEntityHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteKGEntity(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

// ──────────────────────────── KGRelation ────────────────────────────

type KGRelationHandler struct {
	svc *service.KgService
}

func NewKGRelationHandler(svc *service.KgService) *KGRelationHandler {
	return &KGRelationHandler{svc: svc}
}

func (h *KGRelationHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.KGRelation
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageKGRelation(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *KGRelationHandler) GetByID(c *gin.Context) {
	relation, err := h.svc.GetKGRelationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, relation)
}

func (h *KGRelationHandler) Create(c *gin.Context) {
	var relation model.KGRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if _, err := h.svc.CreateKGRelation(c.Request.Context(), &relation); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, relation)
}

func (h *KGRelationHandler) Update(c *gin.Context) {
	var relation model.KGRelation
	if err := c.ShouldBindJSON(&relation); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdateKGRelation(c.Request.Context(), &relation); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *KGRelationHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteKGRelation(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *KGRelationHandler) FindByEntity(c *gin.Context) {
	list, err := h.svc.FindRelationsByEntityID(c.Request.Context(), c.Param("entityId"))
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

// ──────────────────────────── KGDomain ────────────────────────────

type KGDomainHandler struct {
	svc *service.KgService
}

func NewKGDomainHandler(svc *service.KgService) *KGDomainHandler {
	return &KGDomainHandler{svc: svc}
}

func (h *KGDomainHandler) Page(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.KGDomain
	_ = c.ShouldBindQuery(&query)

	list, total, err := h.svc.PageKGDomain(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *KGDomainHandler) List(c *gin.Context) {
	list, err := h.svc.ListKGDomain(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, list)
}

func (h *KGDomainHandler) GetByID(c *gin.Context) {
	domain, err := h.svc.GetKGDomainByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.NotFound)
		return
	}
	response.Success(c, domain)
}

func (h *KGDomainHandler) Create(c *gin.Context) {
	var domain model.KGDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if _, err := h.svc.CreateKGDomain(c.Request.Context(), &domain); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, domain)
}

func (h *KGDomainHandler) Update(c *gin.Context) {
	var domain model.KGDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := h.svc.UpdateKGDomain(c.Request.Context(), &domain); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}

func (h *KGDomainHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteKGDomain(c.Request.Context(), c.Param("id")); err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, true)
}
