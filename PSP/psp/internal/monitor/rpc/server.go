package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
)

type GRPCService struct {
	nodeDao dao.NodeDao
}

func NewGRPCService() (*GRPCService, error) {
	return &GRPCService{
		nodeDao: dao.NewNodeDao(),
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	nodeServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterMonitorServer(s.Driver(), nodeServer)
}
