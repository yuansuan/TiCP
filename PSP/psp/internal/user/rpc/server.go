package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/impl"
)

type GRPCService struct {
	AuthService service.AuthService
	UserService service.UserService
}

func NewGRPCService() (*GRPCService, error) {
	authService := impl.NewAuthService()
	userService := impl.NewUserService()

	return &GRPCService{
		AuthService: authService,
		UserService: userService,
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

	user.RegisterUsersServer(s.Driver(), grpcServer)
}
