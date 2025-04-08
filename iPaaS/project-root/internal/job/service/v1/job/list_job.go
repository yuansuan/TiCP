package job

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// List 获取作业列表
func (srv *jobService) List(ctx context.Context, in *joblist.Request, userID, appID snowflake.ID, withDelete, isSystemFailed bool) (int64, []*models.Job, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job list start")
	defer logger.Info("job list end")

	offsetPage := int(*in.PageOffset)
	size := int(*in.PageSize)
	offset := offsetPage * size

	return srv.jobdao.ListJobs(ctx, offset, size, userID, appID, in.Zone, in.JobState, withDelete, isSystemFailed)
}
