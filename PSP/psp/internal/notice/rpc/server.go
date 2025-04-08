package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service/impl"
)

type GRPCService struct {
	emailService   service.EmailService
	messageService service.MessageService
}

func NewGRPCService() (*GRPCService, error) {
	emailService, err := impl.NewEmailService()
	if err != nil {
		return nil, err
	}

	messageService, err := impl.NewMessageService()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		emailService:   emailService,
		messageService: messageService,
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	grpcServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterNoticeServer(s.Driver(), grpcServer)
}
