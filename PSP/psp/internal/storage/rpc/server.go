package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/impl"
)

type GRPCService struct {
	LocalFileService service.FileService
}

func NewGRPCService() (*GRPCService, error) {
	localFileService, err := impl.NewLocalFileService()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		LocalFileService: localFileService,
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	fileServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterStorageServer(s.Driver(), fileServer)
}
