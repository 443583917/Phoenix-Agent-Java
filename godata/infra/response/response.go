package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"traceId,omitempty"`
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	TraceID string      `json:"traceId,omitempty"`
}

func getTraceID(c *gin.Context) string {
	if tid := c.GetString("trace_id"); tid != "" {
		return tid
	}
	// 从 OTel span 上下文中获取 trace ID
	span := trace.SpanFromContext(c.Request.Context())
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success.Code,
		Message: errcode.Success.Msg,
		Data:    data,
		TraceID: getTraceID(c),
	})
}

func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    errcode.Success.Code,
		Message: errcode.Success.Msg,
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
		TraceID: getTraceID(c),
	})
}

func Error(c *gin.Context, ec errcode.ErrCode) {
	c.JSON(http.StatusOK, Response{
		Code:    ec.Code,
		Message: ec.Msg,
		TraceID: getTraceID(c),
	})
}

func ErrorWithMsg(c *gin.Context, ec errcode.ErrCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    ec.Code,
		Message: msg,
		TraceID: getTraceID(c),
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, ec errcode.ErrCode) {
	c.JSON(httpStatus, Response{
		Code:    ec.Code,
		Message: ec.Msg,
		TraceID: getTraceID(c),
	})
}
