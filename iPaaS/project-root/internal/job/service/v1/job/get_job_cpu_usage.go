package job

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcpuusage"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/hpc/openapi"
)

// GetCpuUsage 获取作业CPU使用率
func (srv *jobService) GetCpuUsage(ctx context.Context, req *jobcpuusage.Request,
	userID snowflake.ID, allow allowFunc) (*schema.JobCpuUsage, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("GetCpuUsage() start")
	defer logger.Debug("GetCpuUsage() end")

	jobID := snowflake.MustParseString(req.JobID)
	// JobID存在性验证
	job, err := srv.jobdao.Get(ctx, jobID, false, false)
	if err != nil {
		if !errors.Is(err, common.ErrJobIDNotFound) {
			logger.Warnf("get Job error! err: %v", err)
		}
		return nil, err // internal error OR job not exist
	}

	// 用户权限验证
	if !allow(userID.String(), job.UserID.String()) {
		logger.Info("no permission to operate other's job")
		return nil, errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
	}

	// 作业状态验证
	state := consts.NewState(job.State, job.SubState)
	if !state.IsRunning() { //
		logger.Info("job is not running")
		return nil, errors.WithMessage(common.ErrJobStateNotAllowQuery, "job is not running")
	}

	// 查询作业CPU使用率
	zones := config.GetConfig().Zones
	zone, ok := zones[job.Zone]
	if !ok {
		logger.Warnf("get zone error! zone: %s", job.Zone)
		return nil, errors.WithMessage(common.ErrInvalidArgumentZone, "get zone error")
	}

	zoneDomain := zone.HPCEndpoint
	if zoneDomain == "" {
		logger.Warnf("zone hpc endpoint is empty")
		return nil, errors.WithMessage(common.ErrInvalidArgumentZone, "zone hpc endpoint is empty")
	}

	hpcResp, err := openapi.Client().GetJobCpuUsage(zoneDomain, module.DefaultTimeout, job.HPCJobID)
	if err != nil {
		if hpcResp.ErrorCode != "" {
			switch hpcResp.ErrorCode {
			case "JobNotRunning": // where defined?
				return nil, errors.WithMessage(common.ErrJobStateNotAllowQuery, "job is not running")
			case api.WrongCPUUsage:
				if strings.Contains(err.Error(), "statusCode: 503") {
					logger.Warnf("GetCpuUsage() warning! err: %v", err)
				} else {
					logger.Errorf("GetCpuUsage() error! err: %v", err)
				}
				return nil, errors.WithMessage(common.ErrWrongCPUUsage,
					"some node cpu usage is wrong, please try again later")
			default:
				logger.Errorf("GetCpuUsage() error! err: [%v](%v)", hpcResp.ErrorCode, hpcResp.ErrorMsg)
			}
		}
		return nil, errors.WithMessage(err, "GetCpuUsage() error")
	}

	cpuUsage := &schema.JobCpuUsage{
		JobID:           jobID.String(),
		AverageCpuUsage: hpcResp.Data.AverageCpuUsage,
		NodeUsages:      hpcResp.Data.NodeUsages,
	}
	return cpuUsage, nil
}
