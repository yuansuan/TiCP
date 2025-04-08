package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/impl"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/app"
)

type GRPCService struct {
	appService service.AppService
}

func NewGRPCService() (*GRPCService, error) {
	appService, err := impl.NewAppService()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		appService: appService,
	}, nil
}

func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	appServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}
	app.RegisterAppServiceServer(s.Driver(), appServer)
}
