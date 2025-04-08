package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/impl"
)

type GRPCService struct {
	RoleService service.RoleService
	PermService service.PermService
}

func NewGRPCService() (*GRPCService, error) {
	roleService, err := impl.NewRoleService()
	permService, err := impl.NewPermService()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		RoleService: roleService,
		PermService: permService,
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

	pb.RegisterRoleManagerServer(s.Driver(), grpcServer)
	pb.RegisterPermissionManagerServer(s.Driver(), grpcServer)
}
