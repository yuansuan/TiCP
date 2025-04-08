package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/impl"
)

type GRPCService struct {
	JobService service.JobService
}

func NewGRPCService() (*GRPCService, error) {
	jobService, err := impl.NewJobService()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		JobService: jobService,
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	jobServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterJobServer(s.Driver(), jobServer)
}
