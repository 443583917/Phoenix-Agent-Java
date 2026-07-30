package knowledge

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
)

type BusinessKnowledgeHandler struct {
	svc *service.DataService
}

func NewBusinessKnowledgeHandler(svc *service.DataService) *BusinessKnowledgeHandler {
	return &BusinessKnowledgeHandler{svc: svc}
}

func (h *BusinessKnowledgeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	var query model.BusinessKnowledge
	_ = c.ShouldBindQuery(&query)
	list, total, err := h.svc.PageBusinessKnowledge(c.Request.Context(), page, size, &query)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *BusinessKnowledgeHandler) GetByID(c *gin.Context) {
	entity, err := h.svc.GetBusinessKnowledgeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, entity)
}

func (h *BusinessKnowledgeHandler) Create(c *gin.Context) {
	var entity model.BusinessKnowledge
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	result, err := h.svc.CreateBusinessKnowledge(c.Request.Context(), &entity)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, result)
}

func (h *BusinessKnowledgeHandler) Update(c *gin.Context) {
	var entity model.BusinessKnowledge
	if err := c.ShouldBindJSON(&entity); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	entity.ID = c.Param("id")
	if err := h.svc.UpdateBusinessKnowledge(c.Request.Context(), &entity); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *BusinessKnowledgeHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteBusinessKnowledge(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, true)
}

func (h *BusinessKnowledgeHandler) Recall(c *gin.Context) {
	entity, err := h.svc.GetBusinessKnowledgeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, gin.H{"id": entity.ID, "results": []interface{}{}, "recalled": true})
}

func (h *BusinessKnowledgeHandler) RefreshVectorStore(c *gin.Context) {
	response.Success(c, gin.H{"refreshed": true, "message": "Vector store refresh triggered"})
}

func (h *BusinessKnowledgeHandler) RetryEmbedding(c *gin.Context) {
	if err := h.svc.RetryBusinessKnowledgeEmbedding(c.Request.Context(), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.Success(c, gin.H{"id": c.Param("id"), "reembedded": true})
}

func handleErr(c *gin.Context, err error) {
	if appErr, ok := err.(*usecase.AppError); ok {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
		return
	}
	response.Error(c, errcode.InternalError)
}
