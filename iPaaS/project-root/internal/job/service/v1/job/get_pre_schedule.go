package job

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

func (srv *jobService) GetPreSchedule(ctx context.Context, preScheduleID string) (*models.PreSchedule, bool, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("get pre schedule start")
	defer logger.Info("get pre schedule end")

	preScheduleIDInt := snowflake.MustParseString(preScheduleID)

	return srv.jobdao.GetPreSchedule(ctx, preScheduleIDInt)
}
