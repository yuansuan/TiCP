package impl

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	openapijob "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *jobServiceImpl) GetJobResidual(ctx context.Context, jobID string) (*dto.JobResidualResponse, error) {
	logger := logging.GetLogger(ctx)

	exist, job, err := s.jobDao.GetJobDetail(ctx, snowflake.MustParseString(jobID))
	if err != nil {
		logger.Errorf("get the job: [%v] detail err: %v", jobID, err)
		return nil, err
	}
	if !exist {
		logger.Infof("the job: [%v] not found", jobID)
		return nil, JobNotFoundError
	}

	residualData := &schema.Residual{}
	switch job.Type {
	case common.Local:
		residual, err := openapijob.AdminGetResidual(s.localAPI, job.OutJobId)
		if err != nil {
			logger.Errorf("get the admin job: [%v] residual err: %v", jobID, err)
			return &dto.JobResidualResponse{}, nil
		}
		residualData = residual
	default:
		logger.Errorf("job compute type: [%v] not match", job.Type)
		return &dto.JobResidualResponse{}, nil
	}

	return util.ConvertJobResidual(residualData), nil
}
