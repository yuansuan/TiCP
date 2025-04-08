package router

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	logmiddleware "github.com/yuansuan/ticp/common/go-kit/logging/middleware"

	app "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/api/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/api/v1/job"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/api/v1/zone"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store/mysql"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/middleware"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/iam"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/mongo"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual"
	jobservice "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/job"
)

// Init ...
func Init(drv *http.Driver) {
	logging.Default().Info("setup router")
	cfg := config.GetConfig()

	// add request id to context and response header
	drv.Use(middleware.RequestIDMiddleware)
	drv.Use(logmiddleware.IngressLogger(logmiddleware.IngressLoggerConfig{
		IsLogRequestHeader:  true,
		IsLogRequestBody:    true,
		IsLogResponseHeader: true,
		IsLogResponseBody:   false,
	}))

	appController := app.NewApplicationController(mysql.GetMysqlFactoryWithEngine())
	engine := boot.MW.DefaultORMEngine()
	jobDao := dao.NewJobDaoImpl(engine)
	var residualDao dao.ResidualDao
	if config.GetConfig().Mongo != nil && config.GetConfig().Mongo.Enable {
		residualDao = dao.NewResidualDaoImpl(mongo.Client().Database(cfg.Mongo.Database()).Collection(models.Residual{}.TableName()))
	} else {
		residualDao = jobDao
	}
	residualHandler := residual.NewHandler(cfg.Zones)
	jobService := jobservice.NewJobService(rpc.GetInstance(), jobDao, residualDao, residualHandler)

	iamClient := iam.Client()
	jh := job.NewJobHandler(appController, jobService, iamClient)
	// /apiGroup 路由组
	apiGroup := drv.Group("/api")
	{
		// jobGroup 路由组
		jobGroup := apiGroup.Group("/jobs")
		{
			// GET请求 - 获取指定作业
			jobGroup.GET(":JobID", jh.Get)

			// POST请求 - 批量获取指定作业
			jobGroup.POST("batch", jh.BatchGet)

			// GET请求 - 获取所有作业
			jobGroup.GET("", jh.List)

			// POST请求 - 创建作业
			jobGroup.POST("", jh.Create)

			// DELETE请求 - 删除指定作业
			jobGroup.DELETE(":JobID", jh.Delete)
			jobGroup.DELETE("", jh.InvalidJobID)

			// PATCH请求 - 终止指定作业
			jobGroup.PATCH(":JobID/terminate", jh.Terminate)
			jobGroup.PATCH("terminate", jh.InvalidJobID)

			// PATCH请求 - 恢复指定作业
			jobGroup.PATCH(":JobID/resume", jh.Resume)
			jobGroup.PATCH("resume", jh.InvalidJobID)

			// PATCH请求 - 指定作业传输暂停
			jobGroup.PATCH(":JobID/transmit/suspend", jh.TransmitSuspend)
			jobGroup.PATCH("transmit/suspend", jh.InvalidJobID)

			// PATCH请求 - 指定作业传输恢复
			jobGroup.PATCH(":JobID/transmit/resume", jh.TransmitResume)
			jobGroup.PATCH("transmit/resume", jh.InvalidJobID)

			// GET请求 - 获取指定作业的残差图
			jobGroup.GET(":JobID/residual", jh.Residual)
			jobGroup.GET("residual", jh.InvalidJobID)

			// GET请求 - 获取指定作业的云图集
			jobGroup.GET(":JobID/snapshots", jh.Snapshots)
			jobGroup.GET("snapshots", jh.InvalidJobID)

			// GET请求 - 获取指定作业的云图数据
			jobGroup.GET(":JobID/snapshots/img", jh.SnapshotImg)
			jobGroup.GET("snapshots/img", jh.InvalidJobID)

			// GET请求 - 获取指定作业的监控图表
			jobGroup.GET(":JobID/monitorchart", jh.MonitorChart)
			jobGroup.GET("monitorchart", jh.InvalidJobID)

			// POST请求 - 创建作业预调度
			jobGroup.POST("preschedule", jh.PreSchedule)

			// GET请求 - 获取指定作业的CPU使用率
			jobGroup.GET(":JobID/cpuusage", jh.CpuUsage)
		}

		// appGroup 路由组
		appGroup := apiGroup.Group("/apps")
		{
			// GET请求 - 获取所有应用程序
			appGroup.GET("", appController.List)
		}

		// zoneGroup 路由组
		zoneGroup := apiGroup.Group("/zones")
		{
			zh := zone.NewZoneHandler(cfg)
			// GET请求 - 获取所有区域
			zoneGroup.GET("", zh.List)
		}

	}

	// /adminGroup 路由组
	adminGroup := drv.Group("/admin", middleware.AdminAccessCheck())
	{
		// adminJobGroup 路由组
		adminJobGroup := adminGroup.Group("/jobs")
		{
			// GET请求 - 管理员获取指定作业
			adminJobGroup.GET(":JobID", jh.AdminGet)

			// GET请求 - 管理员获取所有作业
			adminJobGroup.GET("", jh.AdminList)

			// GET请求 - 管理员获取所有作业
			adminJobGroup.GET("filtered", jh.AdminJobListFiltered)

			// POST请求 - 管理员创建作业
			adminJobGroup.POST("", jh.AdminCreate)

			// DELETE请求 - 删除指定作业
			adminJobGroup.DELETE(":JobID", jh.AdminDelete)

			adminJobGroup.DELETE("", jh.InvalidJobID)

			// PATCH请求 - 管理员终止指定作业
			adminJobGroup.PATCH(":JobID/terminate", jh.AdminTerminate)

			// GET请求 - 管理员获取指定作业的残差图
			adminJobGroup.GET(":JobID/residual", jh.AdminResidual)
			adminJobGroup.GET("/residual", jh.InvalidJobID)

			// GET请求 - 获取指定作业的云图集
			adminJobGroup.GET(":JobID/snapshots", jh.AdminSnapshots)
			adminJobGroup.GET("/snapshots", jh.InvalidJobID)

			// GET请求 - 获取指定作业的云图数据
			adminJobGroup.GET(":JobID/snapshots/img", jh.AdminSnapshotImg)
			adminJobGroup.GET("/snapshots/img", jh.InvalidJobID)

			// GET请求 - 获取指定作业的监控图表
			adminJobGroup.GET(":JobID/monitorchart", jh.AdminMonitorChart)
			adminJobGroup.GET("/monitorchart", jh.InvalidJobID)

			// GET请求 - 获取指定作业的CPU使用率
			adminJobGroup.GET(":JobID/cpuusage", jh.AdminCpuUsage)
		}

		// adminAppGroup 路由组
		adminAppGroup := adminGroup.Group("/apps")
		{
			// GET请求 - 管理员获取所有应用程序
			adminAppGroup.GET("", appController.AdminList)

			// GET请求 - 管理员获取指定应用程序
			adminAppGroup.GET(":AppID", appController.Get)

			// POST请求 - 管理员创建应用程序
			adminAppGroup.POST("", appController.Create)

			// PUT请求 - 管理员修改应用程序
			adminAppGroup.PUT(":AppID", appController.Update)

			adminAppGroup.PUT("", appController.InvalidAppID)

			// DELETE请求 - 管理员删除指定应用程序
			adminAppGroup.DELETE(":AppID", appController.Delete)

			adminAppGroup.DELETE("", appController.InvalidAppID)

			// GET请求 - 查看某个应用对应用户的配额
			adminAppGroup.GET(":AppID/quota", appController.QuotaGet)

			// PATCH请求 - 修改某个应用对应用户的配额(即允许用户使用该应用)
			adminAppGroup.PATCH(":AppID/quota", appController.QuotaAdd)

			// DELETE请求 - 删除某个应用对应用户的配额
			adminAppGroup.DELETE(":AppID/quota", appController.QuotaDelete)

			// GET请求 - 查看某个应用是否在白名单中
			adminAppGroup.GET(":AppID/allow", appController.AllowGet)

			// POST请求 - 添加某个应用到白名单（任何人都可以提交作业）
			adminAppGroup.POST(":AppID/allow", appController.AllowAdd)

			// DELETE请求 - 从白名单中删除某个应用
			adminAppGroup.DELETE(":AppID/allow", appController.AllowDelete)
		}

		// /systemGroup 路由组
		systemGroup := drv.Group("/system")
		{
			jobGroup := systemGroup.Group("/jobs")
			// GET请求 - 请求需要同步的作业列表
			jobGroup.GET("/syncfile", jh.ListNeedSyncFileJobs)

			// PATCH请求 - 更新同步文件状态进度
			jobGroup.PATCH(":JobID/syncfile", jh.SyncFileState)
		}

	}
}
