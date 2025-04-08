package job

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// BatchGet 批量获取作业
func (srv *jobService) BatchGet(ctx context.Context, ids []string, userID snowflake.ID) ([]*models.Job, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job batch get start")
	defer logger.Info("job batch get end")

	jobIDs := []snowflake.ID{}

	for _, id := range ids {
		jobID := snowflake.MustParseString(id)
		jobIDs = append(jobIDs, jobID)
	}

	forUpdate := false
	jobs, err := srv.jobdao.BatchGet(ctx, jobIDs, userID, forUpdate, false) // 数据库查询带上了userID所以只会返回用户有权限的数据
	if err != nil {
		logger.Warnf("get Job error! err: %v", err)
		return nil, err
	}

	return jobs, err
}
