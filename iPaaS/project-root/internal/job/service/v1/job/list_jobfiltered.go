package job

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// ListFiltered 获取作业列表
func (srv *jobService) ListFiltered(ctx context.Context, in *joblistfiltered.Request,
	userID, appID snowflake.ID) (int64, []*models.Job, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job listfiltered start")
	defer logger.Info("job listfiltered end")

	offset := int(consts.DefaultPageOffset)
	if in.PageOffset != nil {
		offset = int(*in.PageOffset)
	}
	size := int(consts.DefaultPageSize)
	if in.PageSize != nil {
		size = int(*in.PageSize)
	}

	return srv.jobdao.ListJobsFiltered(ctx, offset, size, in, userID, appID)
}
