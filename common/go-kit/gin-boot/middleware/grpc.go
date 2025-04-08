package middleware

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
)

var GrpcInterceptors = []grpc.UnaryServerInterceptor{
	//OpentracingServerInterceptorSingleton(),
	apmgrpc.NewUnaryServerInterceptor(),
	GRPCLogger,
	GRPCPrintReqResp,
	RecoveryServerInterceptor(),
	BeforeInterceptor,
	GRPCEnvServerInterceptor(),
	ValidatorUnaryServerInterceptor(),
	grpc_prometheus.UnaryServerInterceptor,
}

var GrpcStreamInterceptors []grpc.StreamServerInterceptor

func AddGrpcUnaryInterceptors(interceptor grpc.UnaryServerInterceptor) {
	GrpcInterceptors = append(GrpcInterceptors, interceptor)
}

func AddGrpcUnaryInterceptorsFirst(interceptors ...grpc.UnaryServerInterceptor) {
	GrpcInterceptors = append(interceptors, GrpcInterceptors...)
}

func AddGrpcStreamInterceptors(interceptor grpc.StreamServerInterceptor) {
	GrpcStreamInterceptors = append(GrpcStreamInterceptors, interceptor)
}

func AddGrpcStreamInterceptorsFirst(interceptors ...grpc.StreamServerInterceptor) {
	GrpcStreamInterceptors = append(interceptors, GrpcStreamInterceptors...)
}

// RecoveryServerInterceptor recovery when server's method panic
// when panic occurs, return code.Internal default
func RecoveryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			if env.Env.Mode <= env.ModeTest {
				stack := logging.Stack(6)
				return status.Errorf(codes.Internal, "%v: %s", p, stack)
			}
			return status.Errorf(codes.Internal, "%v", p)
		}), // if you want handle recovery, modify here
	)
}
