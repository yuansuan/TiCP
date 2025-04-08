package daemon

import (
	"context"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/impl"
)

type DaemonService struct {
	LocalFileService service.FileService
}

func (s *DaemonService) InitCommonPath() error {
	ctx := context.Background()
	// 项目启动检查项目路径与企业路径是否存在，不存在则创建
	if !s.LocalFileService.Exist(ctx, "", true, common.ProjectFolderPath) {
		err := s.LocalFileService.CreateDir(ctx, "", common.ProjectFolderPath, true)
		if err != nil {
			return err
		}
	}
	if !s.LocalFileService.Exist(ctx, "", true, common.PublicFolderPath) {
		err := s.LocalFileService.CreateDir(ctx, "", common.PublicFolderPath, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DaemonService) CheckUserPublicPath() error {
	ctx := context.Background()
	// 如果未开启，检查用户目录是否已存在公共目录的软连接，如果有则删除
	if !config.GetConfig().PublicFolderEnable {
		rsp, err := client.GetInstance().Users.GetAllUserName(ctx, &user.GetAllUserRequest{})
		if err != nil {
			return err
		}

		names := rsp.GetNames()
		for _, userName := range names {
			if s.LocalFileService.Exist(ctx, userName, false, common.PublicFolderPath) {
				err := s.LocalFileService.Remove(ctx, userName, common.PublicFolderPath, false)
				return err
			}
		}
	}

	return nil
}

func NewDaemonService() (*DaemonService, error) {
	localFileService, err := impl.NewLocalFileService()
	if err != nil {
		return nil, err
	}

	return &DaemonService{
		LocalFileService: localFileService,
	}, nil
}

// InitDaemon ...
func InitDaemon() {
	daemonService, err := NewDaemonService()
	if err != nil {
		panic(errors.Wrap(err, "failed to init daemon"))
	}

	err = daemonService.InitCommonPath()
	if err != nil {
		panic(errors.Wrap(err, "failed to init common path"))
	}

	err = daemonService.CheckUserPublicPath()
	if err != nil {
		panic(errors.Wrap(err, "failed to init user public path"))
	}

}
