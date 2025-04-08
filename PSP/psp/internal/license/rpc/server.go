package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/license"
)

type GRPCService struct {
	localAPI *openapi.OpenAPI
}

func NewGRPCService() (*GRPCService, error) {
	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		localAPI: localAPI,
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	licenseServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterLicenseServer(s.Driver(), licenseServer)
}
