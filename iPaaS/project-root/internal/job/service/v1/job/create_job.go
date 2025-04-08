package job

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Create 创建作业
func (srv *jobService) Create(ctx context.Context, req *jobcreate.Request, userID snowflake.ID, chargeParams schema.ChargeParams, appInfo *models.Application, scheduleInfo *models.PreSchedule, convert createConvertFunc) (string, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job create start")
	defer logger.Info("job create end")

	if scheduleInfo != nil {
		logger = logger.With("scheduleID", scheduleInfo.ID)
		logger.Infof("job create with schedule")
	}

	// job id
	jobID, jobIDStr, err := GenJobID(ctx, srv.IDGen)
	if err != nil {
		logger.Warnf("genJobID error: %v", err)
		return "", err // internal error
	}

	// TODO: 创建作业订阅Event，paas反向通知

	job, err := convert(ctx, jobID)
	if err != nil {
		return "", err // internal error
	}
	if config.GetConfig().BillEnabled {
		accountId, err := util.CheckAccountAndMerchandiseInCreateJob(logger, userID, job.AppID)
		if err != nil {
			logger.Warnf("check account in create job failed, %v", err)
			return "", err
		}

		job.AccountID = accountId
		job.ChargeType = *chargeParams.ChargeType

		// 代支付
		if req.PayBy != "" {
			payByAccountID, err := util.CheckPayByAccount(logger, req.PayBy)
			if err != nil {
				logger.Infof("check payBy account failed, %v", err)
				return "", err
			}

			job.PayByAccountID = payByAccountID
		}
	}

	logger = logger.With("jobID", jobIDStr)

	// insert job to db
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err = db.Insert(job)
		return err
	})
	if err != nil {
		// Failed to create the job
		logger.Warnf("session.Insert error: %v", err)
		return "", err // internal error
	}

	logger.Info("job create success")

	if scheduleInfo != nil {
		// 更新预调度状态
		err = with.DefaultSession(ctx, func(db *xorm.Session) error {
			scheduleInfo.Used = true
			_, err = db.ID(scheduleInfo.ID).Cols("used").UseBool("used").Update(scheduleInfo)
			return err
		})
		if err != nil {
			logger.Warnf("update schedule status error: %v", err)
			return "", err
		}
	}

	for _, plugin := range srv.jobPlugins {
		logger.Info("insert plugin: " + plugin.Name())
		plugin.Insert(ctx, appInfo, job)
	}

	return jobIDStr, nil
}

// GenJobID 生成作业ID
func GenJobID(ctx context.Context, idgen snowflake.IDGen) (snowflake.ID, string, error) {
	logger := logging.GetLogger(ctx).With("func", "jobcreate.genJobID")

	jobID, err := idgen.GenID(ctx)
	if err != nil {
		logger.Errorf("generate a snowflake id fail: %v", err)
		return 0, "", err
	}
	return jobID, jobID.String(), nil
}
