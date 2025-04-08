package handler_rpc

import (
	"context"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store/mysql"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/mongo"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/leader"
)

var NSPrefix string

// InitGRPCServer grpc server init
func InitGRPCServer(drv *http.Driver) {
	cfg := config.GetConfig()

	appService := srvv1.NewService(mysql.GetMysqlFactoryWithEngine())
	engine := boot.MW.DefaultORMEngine()
	jobDao := dao.NewJobDaoImpl(engine)
	sender := wx.NewHTTPSender(cfg.WebhookURL)
	var residualDao dao.ResidualDao
	if config.GetConfig().Mongo != nil && config.GetConfig().Mongo.Enable {
		residualDao = dao.NewResidualDaoImpl(mongo.Client().Database(cfg.Mongo.Database()).Collection(models.Residual{}.TableName()))
	} else {
		residualDao = jobDao
	}
	residualHandler := residual.NewHandler(config.GetConfig().Zones)
	logger := logging.Default()
	ctx := context.TODO()
	if NSPrefix != "" {
		NSPrefix = NSPrefix + "."
	}
	leader.Runner(ctx, NSPrefix+"paas_job_scheduler",
		scheduler.NewScheduler(logger, appService, jobDao).Run, leader.SetInterval(1*time.Second))
	leader.Runner(ctx, NSPrefix+"paas_job_updater",
		scheduler.NewJobTimer(jobDao, appService, sender).Run, leader.SetInterval(1*time.Second))
	leader.Runner(ctx, NSPrefix+"paas_job_residual_updater",
		residual.NewUpdateTimer(jobDao, residualDao, appService, residualHandler).Run, leader.SetInterval(1*time.Second))
	leader.Runner(ctx, NSPrefix+"paas_job_monitor_chart_updater",
		monitorchart.NewUpdateTimer(jobDao, appService).Run, leader.SetInterval(1*time.Second))
}
