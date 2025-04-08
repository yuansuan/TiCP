package job

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

func (srv *jobService) ListNeedSyncFileJobs(ctx context.Context, req *jobneedsyncfile.Request) (*jobneedsyncfile.Data, error) {
	offset := *req.PageOffset * *req.PageSize
	jobs, total, err := srv.jobdao.ListNeedFileSyncJobs(ctx, req.Zone, offset, *req.PageSize)
	if err != nil {
		logging.GetLogger(ctx).Warnf("get Job error! err: %v", err)
		return nil, err
	}

	return &jobneedsyncfile.Data{
		Jobs:  util.JobModelToNeedSyncFileJobs(jobs),
		Total: total,
	}, err
}
