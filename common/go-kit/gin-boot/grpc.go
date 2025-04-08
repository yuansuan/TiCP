package boot

import (
	"time"

	_grpc "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"google.golang.org/grpc"
)

type grpcType struct {
}

var (
	GRPC = &grpcType{}
)

// DefaultClient DefaultClient
//
// Deprecated: please use grpc_boot.InjectAllClient()
func (g *grpcType) DefaultClient() (*grpc.ClientConn, error) {
	return _grpc.DefaultClient()
}

// GetClient GetClient
//
// Deprecated: please use grpc_boot.InjectAllClient()
func (g *grpcType) GetClient(name string) (*grpc.ClientConn, error) {
	return _grpc.GetClient(name)
}

// NewClient NewClient
func (g *grpcType) NewClient(addr string, timeout time.Duration, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	return _grpc.NewClient(addr, timeout, dialOptions...)
}

// DefaultServer DefaultServer
func (g *grpcType) DefaultServer() (*_grpc.ServerType, error) {
	return _grpc.DefaultServer()
}

// DefaultServer DefaultServer
func (g *grpcType) DefaultGateway() (*_grpc.ServerType, error) {
	return _grpc.DefaultServer()
}
