package otelgrpcgw

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func WithTracing(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}
	if otelgrpc.Inject(ctx, &md); len(md) == 0 {
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	}

	return md
}
