package otellog

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.opentelemetry.io/otel/trace"
)

// WithTracingLogger 将链路数据注入到日志中
func WithTracingLogger(ctx context.Context) context.Context {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		logger := logging.GetLogger(ctx)
		if logger == logging.Default() { // inject tracing
			logger = logger.With(
				"__trace_id", span.SpanContext().TraceID().String(),
				"__span_id", span.SpanContext().SpanID().String(),
			)
			ctx = context.WithValue(ctx, logging.LoggerName, logger)
		}
	}
	return ctx
}
