package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "phoenix.langfuse"

type TracingService struct {
	tracer trace.Tracer
}

func NewTracingService() *TracingService {
	return &TracingService{
		tracer: otel.Tracer(tracerName),
	}
}

func (s *TracingService) StartGraphSpan(ctx context.Context, threadID, query string) (context.Context, trace.Span) {
	ctx, span := s.tracer.Start(ctx, "graph-stream",
		trace.WithAttributes(
			attribute.String("thread.id", threadID),
			attribute.String("graph.input", query),
		),
	)
	return ctx, span
}

func (s *TracingService) StartNodeSpan(ctx context.Context, nodeName string) (context.Context, trace.Span) {
	ctx, span := s.tracer.Start(ctx, "node."+nodeName,
		trace.WithAttributes(
			attribute.String("node.name", nodeName),
		),
	)
	return ctx, span
}

func (s *TracingService) RecordTokens(span trace.Span, prompt, completion int) {
	span.SetAttributes(
		attribute.Int("tokens.prompt", prompt),
		attribute.Int("tokens.completion", completion),
		attribute.Int("tokens.total", prompt+completion),
	)
}

func (s *TracingService) EndSpan(span trace.Span, output string) {
	span.SetAttributes(attribute.String("graph.output", truncate(output, 500)))
	span.End()
}

func (s *TracingService) RecordError(span trace.Span, err error) {
	span.RecordError(err)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
