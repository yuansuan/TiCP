package impl

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

// JobTerminate 作业终止
func (s *jobServiceImpl) JobTerminate(ctx context.Context, outJobID, computeType string) error {
	logger := logging.GetLogger(ctx)

	if _, err := job.AdminTerminate(s.localAPI, outJobID); err != nil {
		logger.Errorf("terminate the local job [%v] err: %v", outJobID, err)
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("terminate %v job success, params: outJobID=%v", computeType, outJobID))

	return nil
}
