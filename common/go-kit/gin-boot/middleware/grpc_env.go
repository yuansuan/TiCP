package middleware

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	"google.golang.org/grpc"
)

// ContextKey ...
type ContextKey string

func (c ContextKey) String() string {
	return "ys-" + string(c)
}

const (
	// EnvContextKey ...
	EnvContextKey = ContextKey("req-env")
)

// GRPCEnvServerInterceptor ...
func GRPCEnvServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		header := util.HeaderFromIncomingContext(ctx)
		reqEnv := header.Get(EnvContextKey.String())
		ctx = context.WithValue(ctx, EnvContextKey, reqEnv)
		resp, err = handler(ctx, req)
		return
	}
}

// GRPCEnvClientInterceptor ...
func GRPCEnvClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = util.AppendToOutgoingContext(ctx, EnvContextKey.String(), env.Env.Type)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
