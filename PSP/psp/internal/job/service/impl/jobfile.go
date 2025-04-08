package impl

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/config"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

// CreateJobTempDir 创建作业临时目录
func (s *jobServiceImpl) CreateJobTempDir(ctx context.Context, userName, computeType string) (string, error) {
	logger := logging.GetLogger(ctx)

	randomStr := s.sid.Generate().String()
	tempUpload := config.GetConfig().TempUpload
	tempPath := fmt.Sprintf("%v/%v", tempUpload, randomStr)
	createDirReq := &pb.CreateDirReq{
		Path:     tempPath,
		UserName: userName,
		Cross:    true,
		IsCloud:  false,
	}

	if _, err := s.rpc.Storage.CreateDir(ctx, createDirReq); err != nil {
		logger.Errorf("[%v] create job temp directory [%v] err: %v", userName, tempPath, err)
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("%v create %v job submit template directory [%v] success", userName, computeType, tempPath))

	return tempPath, nil
}

// GetWorkSpace 获取工作空间
func (s *jobServiceImpl) GetWorkSpace(ctx context.Context) string {
	return config.GetConfig().WorkDir.Workspace
}
