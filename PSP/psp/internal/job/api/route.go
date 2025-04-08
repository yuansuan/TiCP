package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/impl"
)

type apiRoute struct {
	jobService service.JobService
}

func NewAPIRoute() (*apiRoute, error) {
	jobService, err := impl.NewJobService()
	if err != nil {
		return nil, err
	}

	return &apiRoute{
		jobService: jobService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewAPIRoute()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	group := drv.Group("/api/v1")
	{
		jobGroup := group.Group("/job")

		jobGroup.POST("/list", api.JobList)
		jobGroup.GET("/detail", api.JobDetail)
		jobGroup.GET("/jobSetDetail", api.JobSetDetail)
		jobGroup.GET("/residual", api.JobResidual)
		jobGroup.GET("/snapshots", api.JobSnapshotList)
		jobGroup.GET("/snapshot", api.JobSnapshot)

		jobGroup.GET("/computeTypes", api.JobComputeTypeList)
		jobGroup.GET("/jobSetNames", api.JobSetNameList)
		jobGroup.GET("/appNames", api.JobAppNameList)
		jobGroup.GET("/userNames", api.JobUserNameList)
		jobGroup.GET("/queueNames", api.JobQueueNameList)

		jobGroup.POST("/submit", api.Submit)
		jobGroup.POST("/resubmit", api.Resubmit)
		jobGroup.POST("/terminate", api.JobTerminate)

		jobGroup.GET("/workspace", api.GetWorkSpace)
		jobGroup.POST("/createTempDir", api.CreateTempDir)

		jobGroup.GET("/statistics/totalCPUTime", api.GetJobStatisticsTotalCPUTime)
		jobGroup.GET("/statistics/overview", api.GetJobStatisticsOverview)
		jobGroup.GET("/statistics/detail", api.GetJobStatisticsDetail)
		jobGroup.GET("/statistics/export", api.GetJobStatisticsExport)

		// 大屏展示
		jobGroup.GET("/statistics/top5ProjectInfo", api.GetTop5ProjectInfo)

		//应用作业数
		jobGroup.GET("/appJobNum", api.AppJobNum)
		//用户作业数
		jobGroup.GET("/userJobNum", api.UserJobNum)
	}

}
