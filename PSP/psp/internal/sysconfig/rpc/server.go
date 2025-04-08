package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service/impl"
)

type GRPCService struct {
	SysConfigService service.SysConfigService
}

func NewGRPCService() (*GRPCService, error) {
	sysConfigService, err := impl.NewSysConfigService()
	if err != nil {
		logging.Default().Errorf("init sys config server service err: %v", err)
		return nil, err
	}

	return &GRPCService{
		SysConfigService: sysConfigService,
	}, nil
}

func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	grpcServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	sysconfig.RegisterSysConfigServer(s.Driver(), grpcServer)
}
