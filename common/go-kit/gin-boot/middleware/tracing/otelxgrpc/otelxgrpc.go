package otelxgrpc

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing/otellog"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(otellog.WithTracingLogger(ctx), req)
	}
}

func UnaryDumper(request, response bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span := trace.SpanFromContext(ctx)
		if span.IsRecording() && request {
			span.AddEvent("request coming", trace.WithAttributes(
				attribute.String("body", util.MustParseJSON(req)),
			))
		}

		resp, err = handler(ctx, req)
		if span.IsRecording() && response {
			span.AddEvent("response sent", trace.WithAttributes(
				attribute.String("body", util.MustParseJSON(resp)),
			))
		}

		return
	}
}
