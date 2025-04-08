package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao"
)

type GRPCService struct {
	projectMemberDao dao.ProjectMemberDao
	projectDao       dao.ProjectDao
}

func NewGRPCService() (*GRPCService, error) {
	projectDao, err := dao.NewProjectDao()
	if err != nil {
		return nil, err
	}

	return &GRPCService{
		projectMemberDao: dao.NewProjectMemberDao(),
		projectDao:       projectDao,
	}, nil
}

// InitGRPC 初始化GRPC服务
func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	projectServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	pb.RegisterProjectServer(s.Driver(), projectServer)
}
