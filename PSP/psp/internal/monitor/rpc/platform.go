package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
)

// GetPlatformCores 获取Platform及核数
func (s *GRPCService) GetPlatformCores(ctx context.Context, req *monitor.GetPlatformCoresRequest) (*monitor.GetPlatformCoresResponse, error) {
	logger := logging.GetLogger(ctx)
	platformInfos, err := s.nodeDao.GetPlatformCores(ctx)
	if err != nil {
		logger.Errorf("get platform cores error, err: %v", err)
		return nil, err
	}

	return &monitor.GetPlatformCoresResponse{
		PlatformCores: platformInfos,
	}, nil
}
