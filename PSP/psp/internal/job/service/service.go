package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type JobService interface {
	// GetWorkSpace 获取工作空间
	GetWorkSpace(ctx context.Context) string

	// GetJobUserNameList 获取作业用户名称列表
	GetJobUserNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error)

	// GetJobComputeTypeList 获取作业应用名称列表
	GetJobComputeTypeList(ctx context.Context, loginUserID snowflake.ID) ([]*dto.ComputeTypeName, error)

	// GetJobSetNameList 获取作业集名称列表
	GetJobSetNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error)

	// GetJobAppNameList 获取作业应用名称列表
	GetJobAppNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error)

	// GetJobQueueNameList 获取作业队列名称列表
	GetJobQueueNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error)

	// GetJobDetail 获取作业详情
	GetJobDetail(ctx context.Context, jobID string) (*dto.JobDetailInfo, error)

	// GetJobSetDetail 获取作业集详情
	GetJobSetDetail(ctx context.Context, jobSetID string, loginUserID snowflake.ID) (*dto.JobSetInfo, []*dto.JobListInfo, error)

	// JobTerminate 作业终止
	JobTerminate(ctx context.Context, outJobID, computeType string) error

	// GetJobDetailByOutID 获取作业详情
	GetJobDetailByOutID(ctx context.Context, outJobID, jobType string) (*model.Job, error)

	// CreateJobTempDir 创建作业临时目录
	CreateJobTempDir(ctx context.Context, userName string, computeType string) (string, error)

	// JobSubmit 作业提交
	JobSubmit(ctx *gin.Context, param *dto.SubmitParam) ([]string, error)

	// JobResubmit 作业重提交
	JobResubmit(ctx context.Context, req *dto.ResubmitRequest, loginUserId snowflake.ID, username string) (*dto.ResubmitResponse, error)

	// GetJobList 获取作业列表
	GetJobList(ctx context.Context, filter *dto.JobFilter, page *xtype.Page, orderSort *xtype.OrderSort, loginUserID snowflake.ID) ([]*model.Job, int64, error)

	// GetJobCPUTimeTotal 获取作业核时总数据
	GetJobCPUTimeTotal(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64) (float64, error)

	// GetJobStatisticsOverview 获取作业统计总览
	GetJobStatisticsOverview(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64, pageIndex, pageSize int) ([]*dto.StatisticsOverview, int64, error)

	// GetJobStatisticsDetail 获取作业统计详情
	GetJobStatisticsDetail(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64, pageIndex, pageSize int) ([]*dto.JobDetailInfo, int64, error)

	// GetJobStatisticsExport 导出作业统计数据
	GetJobStatisticsExport(ctx *gin.Context, queryType, computeType, showType string, names, projectIds []string, startTime, endTime int64) error

	// AppJobNum 应用作业数
	AppJobNum(ctx context.Context, start, end int64) ([]*dto.AppJobInfo, int, error)

	// UserJobNum 用户作业数
	UserJobNum(ctx context.Context, start, end int64) ([]*dto.UserJobInfo, error)

	// GetJobCPUTimeMetric 作业核时运行指标统计
	GetJobCPUTimeMetric(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCPUTimeMetric, error)

	// GetJobCountMetric 应用和用户数量统计指标
	GetJobCountMetric(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCountMetric, error)

	// GetJobDeliverCount 作业提交数量指标
	GetJobDeliverCount(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCountMetric, error)

	// GetJobWaitStatistic 作业等待指标
	GetJobWaitStatistic(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobWaitStatistic, error)

	// GetJobStatus 作业状态统计
	GetJobStatus(ctx context.Context) (map[string]int64, error)

	// GetJobTimeline 获取作业时间线
	GetJobTimeline(ctx context.Context, jobID, uploadFileTaskID, jobState, dataState string) ([]*dto.JobTimeLine, error)

	// GetJobResidual 获取作业残差图
	GetJobResidual(ctx context.Context, jobID string) (*dto.JobResidualResponse, error)

	// GetJobSnapshotList 获取云图集
	GetJobSnapshotList(ctx context.Context, jobID string) (*dto.JobSnapshotListResponse, error)

	// GetJobSnapshot 获取云图资源
	GetJobSnapshot(ctx context.Context, jobID, path string) (*dto.JobSnapshotResponse, error)

	// GetOutIDByJobID 根据作业id获取paas作业id
	GetOutIDByJobID(ctx context.Context, jobID snowflake.ID) (string, error)

	// GetTop5ProjectInfo 获取top5项目信息
	GetTop5ProjectInfo(ctx context.Context, start, end int64) (*dto.GetTop5ProjectInfoResponse, error)
}
