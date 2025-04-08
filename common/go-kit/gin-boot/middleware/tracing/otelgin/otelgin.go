package otelgin

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing/otellog"
)

const (
	ginTraceKey         = "ys-otel-go-tracer"
	instrumentationName = "yuansuan.cn/tracing/otelgin"
)

// Middleware returns middleware that will trace incoming requests.
// The service parameter should describe the name of the (virtual)
// server handling the request.
//
// copy from go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
func Middleware(service string, options ...Option) gin.HandlerFunc {
	cfg := newConfig(options...)
	tracer := cfg.tp.Tracer(instrumentationName, oteltrace.WithInstrumentationVersion(tracing.SemVersion()))

	return func(c *gin.Context) {
		if !cfg.enabled || (cfg.ignoreMatcher != nil && cfg.ignoreMatcher(c)) {
			c.Next()
			return
		}
		if cfg.excludes.Match(c.Request.URL.Path) {
			c.Next()
			return
		}

		c.Set(ginTraceKey, tracer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx := cfg.pg.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		ctx, span := tracer.Start(ctx, cfg.sn(c), opts...)
		defer span.End()

		// injecting the logger with tracing
		ctx = otellog.WithTracingLogger(ctx)
		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
