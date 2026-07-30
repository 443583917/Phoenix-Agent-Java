package langfuse

import (
	"context"

	"go.uber.org/zap"
)

type TraceService struct {
	logger *zap.Logger
}

func NewTraceService() *TraceService {
	return &TraceService{logger: zap.L().Named("langfuse.trace")}
}

type Span struct {
	TraceID   string `json:"traceId"`
	SpanID    string `json:"spanId"`
	Name      string `json:"name"`
	Input     string `json:"input,omitempty"`
	Output    string `json:"output,omitempty"`
	Model     string `json:"model,omitempty"`
	PromptTokens     int `json:"promptTokens,omitempty"`
	CompletionTokens int `json:"completionTokens,omitempty"`
	TotalTokens      int `json:"totalTokens,omitempty"`
}

func (s *TraceService) StartSpan(ctx context.Context, name string) *Span {
	s.logger.Debug("start span", zap.String("name", name))
	return &Span{Name: name}
}

func (s *TraceService) EndSpan(ctx context.Context, span *Span, output string) {
	s.logger.Debug("end span",
		zap.String("name", span.Name),
		zap.Int("promptTokens", span.PromptTokens),
		zap.Int("completionTokens", span.CompletionTokens),
	)
	span.Output = output
}

func (s *TraceService) RecordError(ctx context.Context, span *Span, err error) {
	s.logger.Error("span error",
		zap.String("name", span.Name),
		zap.Error(err),
	)
}

func (s *TraceService) Flush(ctx context.Context) {
	s.logger.Debug("flushing traces")
}
