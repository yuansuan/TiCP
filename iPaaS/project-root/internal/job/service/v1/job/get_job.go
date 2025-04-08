package job

import (
	"context"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// Get 获取作业详情
func (srv *jobService) Get(ctx context.Context, req *jobget.Request, userID snowflake.ID, allow allowFunc, withDelete bool) (*models.Job, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job get start")
	defer logger.Info("job get end")

	jobID := snowflake.MustParseString(req.JobID)

	job, err := srv.jobdao.Get(ctx, jobID, false, withDelete)
	if err != nil {
		if !errors.Is(err, common.ErrJobIDNotFound) {
			logger.Warnf("get Job error! err: %v", err)
		}

		return nil, err
	}

	// 用户权限验证
	if !allow(userID.String(), job.UserID.String()) {
		logger.Warnf("no permission to operate other's job")
		return nil, errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
	}

	return job, err
}
