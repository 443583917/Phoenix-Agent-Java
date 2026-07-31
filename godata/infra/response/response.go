package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"go.opentelemetry.io/otel/trace"
)

// 兼容前端契约：前端 successCode 用严格比较 code === '100'（字符串）。
// Java 版 ReturnVo.code 为 String 类型，Go 版序列化时需输出字符串。
func codeString(c int) string {
	return strconv.Itoa(c)
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"traceId,omitempty"`
}

type PageResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
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
		Code:    codeString(errcode.Success.Code),
		Message: errcode.Success.Msg,
		Success: true,
		Data:    data,
		TraceID: getTraceID(c),
	})
}

func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    codeString(errcode.Success.Code),
		Message: errcode.Success.Msg,
		Success: true,
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
		TraceID: getTraceID(c),
	})
}

func Error(c *gin.Context, ec errcode.ErrCode) {
	c.JSON(http.StatusOK, Response{
		Code:    codeString(ec.Code),
		Message: ec.Msg,
		Success: false,
		Error:   ec.Msg,
		Msg:     ec.Msg,
		TraceID: getTraceID(c),
	})
}

func ErrorWithMsg(c *gin.Context, ec errcode.ErrCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    codeString(ec.Code),
		Message: msg,
		Success: false,
		Error:   msg,
		Msg:     msg,
		TraceID: getTraceID(c),
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, ec errcode.ErrCode) {
	c.JSON(httpStatus, Response{
		Code:    codeString(ec.Code),
		Message: ec.Msg,
		Success: false,
		Error:   ec.Msg,
		Msg:     ec.Msg,
		TraceID: getTraceID(c),
	})
}
