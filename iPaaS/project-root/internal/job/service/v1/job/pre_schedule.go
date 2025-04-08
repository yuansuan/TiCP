package job

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/hpc/openapi"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// PreSchedule 创建作业预调度
func (srv *jobService) PreSchedule(ctx context.Context, req *jobpreschedule.Request,
	zones schema.Zones, userID snowflake.ID, appInfo *models.Application,
) (jobpreschedule.Data, error) {
	logger := logging.GetLogger(ctx).With("func", "jobpreschedule.PreSchedule", "userID", userID)
	logger.Info("job pre schedule start")
	defer logger.Info("job pre schedule end")

	// job pre schedule id
	preScheduleID, preScheduleStr, err := GenJobID(ctx, srv.IDGen)
	if err != nil {
		logger.Warnf("genJobID error: %v", err)
		return jobpreschedule.Data{}, err // internal error
	}

	// 预调度，分区
	// TODO:这样调用会导致单测依赖实际的hpc client以及selector的实现，后续尝试mock
	zone, err := preScheduleZone(ctx, req, zones, appInfo, preScheduleID, userID)
	if err != nil {
		logger.Warnf("preScheduleZone error: %v", err)
		return jobpreschedule.Data{}, err // internal error
	}

	preSchedule, err := util.ConvertPreScheduleModel(ctx, logger, req, userID, preScheduleID, appInfo, zone)
	if err != nil {
		return jobpreschedule.Data{}, err // internal error
	}

	logger = logger.With("preScheduleID", preScheduleStr)

	// insert preSchedule to db
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err = db.Insert(preSchedule)
		return err
	})
	if err != nil {
		// Failed to create the job
		logger.Warnf("session.Insert error: %v", err)
		return jobpreschedule.Data{}, err // internal error
	}

	logger.Info("job pre schedule create success")

	return jobpreschedule.Data{
		ScheduleID: preScheduleStr,
		Workdir:    preSchedule.WorkDir,
	}, nil
}

// 预调度返回zone
func preScheduleZone(ctx context.Context,
	req *jobpreschedule.Request, zones schema.Zones, appInfo *models.Application,
	preScheduleID snowflake.ID, userID snowflake.ID,
) (consts.Zone, error) {
	logger := logging.GetLogger(ctx).With("func", "jobpreschedule.preScheduleZone", "preScheduleID", preScheduleID, "userID", userID)
	logger.Infof("preScheduleZone start, req: %s", spew.Sdump(req))

	// fixed为true时，只从request.Zones的分区范围中选取，request.Zones按照传入顺序排优先级(此时request.Zones不能为空)
	// fixed为false时，从所有分区中选取，但是优先级为request.Zones中的分区更高(此时request.Zones可为空)
	// 这里传入的zones已校验空，且过滤了config里没有的分区

	appSpecifyQueueMap := util.ToStringMap(appInfo.SpecifyQueue)

	// 区域选择器
	// zones,hpcClient,appSpecifyQueueMap,queue
	zs := scheduler.NewZoneSelector(zones, openapi.Client(), appSpecifyQueueMap, "", "")

	// 资源聚合
	resources, err := zs.ResourceAggregation(ctx)
	if err != nil {
		logger.Warnf("ResourceAggregation error: %v", err)
		return consts.ZoneUnknown, err
	}

	// 过滤器
	cpuRange := util.CoresRange{
		MinExpectedCores: int64(*req.Params.Resource.MinCores),
		MaxExpectedCores: int64(*req.Params.Resource.MaxCores),
	}

	licenseFilterParams := &scheduler.LicenseFilterParams{
		IdentifierID: preScheduleID,
		CPURange:     cpuRange,
		Shared:       req.Shared,
	}
	zs.RegisterFilter(
		scheduler.NewBaseFilter(cpuRange, int64(*req.Params.Resource.Memory), "", "", "", req.Shared),
		scheduler.NewLicenseFilter(zones, licenseFilterParams, appInfo, rpc.GetInstance().License.LicenseServer),
	)

	// 优选器
	selectorWeights := config.GetConfig().SelectorWeights
	rc := scheduler.NewResourceCount()
	rc.SetWeight(selectorWeights[rc.Name()])
	// sf:=scheduler.NewStorageFirst(job.FileInputStorageZone, 10000), // !因为没有FileInputStorageZone，所以这里不需要
	zps := scheduler.NewZonePrioritySelector(req.Zones) // 按照传入顺序排优先级的选择器，传入了config里没有的分区也没关系
	zps.SetWeight(selectorWeights[zps.Name()])
	qps := scheduler.NewQueuePrioritySelector() // 队列优先级选择器
	qps.SetWeight(selectorWeights[qps.Name()])
	zs.RegisterOptimalSelector(rc, zps, qps)

	params := &scheduler.SelectParams{
		Resources: resources,
	}
	zone, _, err := zs.Select(ctx, params, scheduler.WithPreScheduleID(preScheduleID), scheduler.WithUserID(userID))
	if err != nil {
		logger.Warnf("Select error: %v", err)
		return consts.ZoneUnknown, err
	}

	logger.Infof("preScheduleZone end, zone: %s", zone)

	return consts.Zone(zone), nil
}
