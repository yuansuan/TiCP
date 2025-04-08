package dataloader

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	apiconf "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
	openapijob "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

const (
	FirstNum        = 1
	Bursting        = 1
	BurstSuccess    = 2
	BurstFailed     = 3
	DownloadSuccess = 4
	DownloadFailed  = 5
)

var (
	burstNotBind = errors.New("local app not bind burst cloud app")
)

// JobLoader ...
type JobLoader struct {
	rpc            *client.GRPC
	localCfg       *apiconf.Local
	localAPI       *openapi.OpenAPI
	jobDao         dao.JobDao
	jobAttrDao     dao.JobAttrDao
	jobTimelineDao dao.JobTimelineDao
}

// UpdateJob ...
type UpdateJob struct {
	Job      *model.Job
	Cols     []string
	PreState string
}

// NewJobLoader 作业数据同步
func NewJobLoader() (*JobLoader, error) {
	jobDao, err := dao.NewJobDao()
	if err != nil {
		return nil, err
	}

	jobAttrDao, err := dao.NewJobAttrDao()
	if err != nil {
		return nil, err
	}

	jobTimelineDao, err := dao.NewJobTimelineDao()
	if err != nil {
		return nil, err
	}

	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	rpc := client.GetInstance()
	localCfg := apiconf.GetConfig().Local
	if localCfg == nil {
		return nil, errors.New("openapi configuration is invalid")
	}

	impl := &JobLoader{
		rpc:            rpc,
		localCfg:       localCfg,
		localAPI:       localAPI,
		jobDao:         jobDao,
		jobAttrDao:     jobAttrDao,
		jobTimelineDao: jobTimelineDao,
	}

	return impl, nil
}

// JobLoaderStart 作业信息同步
func (loader *JobLoader) JobLoaderStart() {
	go loader.jobSyncTicker()
}

// jobSyncTicker 作业同步定时器
func (loader *JobLoader) jobSyncTicker() {
	logger := logging.Default()

	// 根据配置启用
	syncData := config.GetConfig().SyncData
	if !syncData.Enable {
		logger.Infof("sync job data routine has disabled")
		return
	}

	timerDuration := time.Second * time.Duration(syncData.Interval)
	timer := time.NewTimer(timerDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := loader.jobSync(); err != nil {
				logging.GetLogger(context.Background()).Errorf("job sync err: %v", err)
			}
			timer.Reset(timerDuration)
		}
	}
}

// jobSync 作业数据同步
func (loader *JobLoader) jobSync() error {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	unfinishedJobList, total, err := loader.jobDao.GetUnfinishedJobList(ctx, &xtype.Page{
		Index: xtype.MinPageIndex,
		Size:  xtype.MaxPageSize,
	})
	if err != nil {
		return err
	}

	if xtype.MaxPageSize < total {
		nextPageIndex := xtype.MinPageIndex + 1
		pageIndexes := util.GetTotalPageIndexes(nextPageIndex, xtype.MaxPageSize, total)
		for _, index := range pageIndexes {
			jobList, _, err := loader.jobDao.GetUnfinishedJobList(ctx, &xtype.Page{
				Index: index,
				Size:  xtype.MaxPageSize,
			})
			if err != nil {
				return err
			}

			unfinishedJobList = append(unfinishedJobList, jobList...)
		}
	}

	updateJobMap := make(map[string]*UpdateJob, len(unfinishedJobList))
	for _, job := range unfinishedJobList {
		var openapiAdminJob *schema.AdminJobInfo
		if job.Type == common.Local {
			openapiAdminJob, err = openapijob.AdminGetJob(loader.localAPI, job.OutJobId)
			if err != nil {
				logger.Errorf("get openapi local job [%v] err: %v", job.OutJobId, err)
				continue
			}
		} else {
			continue
		}

		if err = loader.updateJob(ctx, job, openapiAdminJob); err != nil {
			continue
		}
	}

	logger.Debugf("sync [%v] job info from openapi", len(updateJobMap))

	return nil
}

func (loader *JobLoader) updateJob(ctx context.Context, job *model.Job, openapiJob *schema.AdminJobInfo) error {
	logger := logging.GetLogger(ctx)

	updateJob, cols, preJobState := loader.generateJobUpdateInfo(ctx, openapiJob, job)
	if len(cols) > 0 {
		if err := loader.jobDao.UpdateJobWithCols(ctx, updateJob, cols); err != nil {
			logger.Errorf("sync update job [outJobId=%v] err: %v", updateJob.OutJobId, err)
			return err
		}

		if !strutil.IsEmpty(preJobState) && preJobState != updateJob.State {
			go loader.sendStateSyncMessage(updateJob)
		}
	}
	return nil
}

func (loader *JobLoader) sendStateSyncMessage(job *model.Job) {
	ctx := context.Background()
	SaveJobTimeLine(ctx, job, loader.jobTimelineDao)

	stateMsg := util.ConvertJobStateMsg(job.State)
	if strutil.IsEmpty(stateMsg) {
		stateMsg = job.State
	}

	loader.sendMessage(ctx, job.UserId.String(), job.Id.String(), job.Name, stateMsg)
}

func (loader *JobLoader) sendMessage(ctx context.Context, userID, jobID, jobName, contentSuffix string) {
	logger := logging.GetLogger(ctx)

	content := fmt.Sprintf("作业[编号:%v 名称:%v]%v", jobID, jobName, contentSuffix)
	msg := &pbnotice.WebsocketMessage{
		UserId:  userID,
		Type:    common.JobEventType,
		Content: content,
	}

	if _, err := loader.rpc.Notice.SendWebsocketMessage(ctx, msg); err != nil {
		logger.Errorf("job submit send ws message err: %v", err)
	}
}

// SaveJobTimeLine 保存时间线信息
func SaveJobTimeLine(ctx context.Context, job *model.Job, jobTimelineDao dao.JobTimelineDao) {
	logger := logging.GetLogger(ctx)

	if job.State == consts.JobStatePending {
		return
	}

	eventName, eventTime := util.ConvertJobTimelineEvent(job)
	if eventName == "" {
		return
	}

	timeline := &model.JobTimeline{
		JobId:     job.Id,
		EventName: eventName,
		EventTime: eventTime,
	}
	if err := jobTimelineDao.InsertJobTimeline(ctx, timeline); err != nil {
		logger.Errorf("save job [%v] timeline event name [%v] err: %v", job.Id, eventName, err)
	}
}
