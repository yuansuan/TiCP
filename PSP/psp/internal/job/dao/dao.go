package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// JobDao 作业表数据访问
type JobDao interface {
	// InsertJob 保存作业信息
	InsertJob(ctx context.Context, job *model.Job) (snowflake.ID, error)

	// UpdateJob 更新作业信息
	UpdateJob(ctx context.Context, job *model.Job) error

	// UpdateDataState 更新作业数据状态
	UpdateDataState(ctx context.Context, outJobID, dataState string) error

	// UpdateBurstJob 更新爆发作业信息
	UpdateBurstJob(ctx context.Context, job *model.Job) error

	// UpdateJobWithCols 指定列更新作业信息
	UpdateJobWithCols(ctx context.Context, job *model.Job, cols []string) error

	// GetJobUserNameList 获取作业用户名称列表
	GetJobUserNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error)

	// GetJobComputeTypeList 获取作业计算类型列表
	GetJobComputeTypeList(ctx context.Context, isAdmin bool, loginUserID snowflake.ID) ([]string, error)

	// GetJobSetNameList 获取作业集名称列表
	GetJobSetNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error)

	// GetJobAppNameList 获取作业应用名称列表
	GetJobAppNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error)

	// GetJobQueueNameList 获取作业队列名称列表
	GetJobQueueNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error)

	// GetAppJobNum 获取应用作业数
	GetAppJobNum(ctx context.Context, start, end int64) ([]*dto.AppJobInfo, error)

	// GetUserJobNum 获取用户作业数
	GetUserJobNum(ctx context.Context, start, end int64) ([]*dto.UserJobInfo, error)

	// GetJobByOutID 获取作业详细信息
	GetJobByOutID(ctx context.Context, outJobID, jobType string) (bool, *model.Job, error)

	// GetJobListByOutJobID 获取作业信息列表
	GetJobListByOutJobID(ctx context.Context, outJobIDList []string) ([]*model.Job, error)

	// GetJobDetail 获取作业详细信息
	GetJobDetail(ctx context.Context, jobID snowflake.ID) (bool, *model.Job, error)

	// GetUnfinishedJobList 分页获取未结束作业信息列表
	GetUnfinishedJobList(ctx context.Context, page *xtype.Page) ([]*model.Job, int64, error)

	// GetJobList 分页获取作业列表
	GetJobList(ctx context.Context, filter *dto.JobFilter, page *xtype.Page, orderSort *xtype.OrderSort, isAdmin bool, loginUserID snowflake.ID) ([]*model.Job, int64, error)

	// GetJobCPUTimeTotal 获取作业核时总数据
	GetJobCPUTimeTotal(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64) (float64, error)

	// GetJobStatisticsOverview 获取作业统计总览
	GetJobStatisticsOverview(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64, pageIndex, pageSize int) ([]*model.StatisticsJob, int64, error)

	// GetJobStatisticsDetail 获取作业统计详情
	GetJobStatisticsDetail(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64, pageIndex, pageSize int) ([]*model.StatisticsJob, int64, error)

	// GetJobCPUTimeMetric 作业核时运行指标统计
	GetJobCPUTimeMetric(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string, states []string) ([]*dto.JobCPUTimeQueryMetric, error)

	// GetJobCountMetric 获取应用和用户数量统计指标
	GetJobCountMetric(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string, states []string) ([]*dto.JobQueryResultMetric, error)

	// GetJobDeliverCount 作业提交数量指标
	GetJobDeliverCount(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string) ([]*dto.JobQueryResultMetric, error)

	// GetJobWaitStatistic 作业等待指标
	GetJobWaitStatistic(ctx context.Context, filter *dto.JobMetricFiler, statisticType string, states []string) ([]*dto.JobQueryResultMetric, error)

	// GetJobStatusNum 获取作业状态统计
	GetJobStatusNum(ctx context.Context, start int64, states []string) ([]*dto.JobStatus, error)

	//GetTop5ProjectByCPUTime 获取CPU时间最多的5个项目
	GetTop5ProjectByCPUTime(ctx context.Context, projectIds []snowflake.ID, start, end int64) ([]*dto.ProjectCPUTime, error)

	// GetJobCountByProjectIds 获取项目下作业数
	GetJobCountByProjectIds(ctx context.Context, projectIds []snowflake.ID, start, end int64) ([]*dto.ProjectJobCount, error)
}

// JobAttrDao 作业属性表数据访问
type JobAttrDao interface {
	// GetJobAttrList 获取作业属性信息
	GetJobAttrList(ctx context.Context, jobID snowflake.ID) ([]*model.JobAttr, error)

	// GetJobAttrByKey 获取指定key的作业属性信息
	GetJobAttrByKey(ctx context.Context, jobID snowflake.ID, key string) (bool, *model.JobAttr, error)

	// UpdateJobAttr 更新作业属性信息
	UpdateJobAttr(ctx context.Context, attrs *model.JobAttr) error

	// InsertJobAttr 保存作业属性信息
	InsertJobAttr(ctx context.Context, attrs *model.JobAttr) error
}

// JobTimelineDao 作业时间线表数据访问
type JobTimelineDao interface {
	// GetJobTimeline 获取指定作业的时间线信息
	GetJobTimeline(ctx context.Context, jobID snowflake.ID) ([]*model.JobTimeline, error)

	// InsertJobTimeline 保存作业时间线信息
	InsertJobTimeline(ctx context.Context, timeline *model.JobTimeline) error

	// GetJobTimelineByName 获取指定名称的作业时间线信息
	GetJobTimelineByName(ctx context.Context, jobID snowflake.ID, eventName string) (bool, *model.JobTimeline, error)
}
