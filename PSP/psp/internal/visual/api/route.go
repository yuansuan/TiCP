package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/cmd/config"
	visualconfig "github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/impl"
)

type RouteService struct {
	visualService service.VisualService
}

func NewRouteService() (*RouteService, error) {
	logger := logging.Default()
	visualService, err := impl.NewVisualService()
	if err != nil {
		logger.Errorf("init visual server service err: %v", err)
		return nil, err
	}
	return &RouteService{
		visualService: visualService,
	}, nil
}

func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	// 根据全局配置启用
	if !config.Custom.Main.EnableVisual && !visualconfig.GetConfig().Local {
		logger.Infof("visual route init has disabled by global config")
		return
	}

	s, err := NewRouteService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	appGroup := drv.Group("/api/v1/vis")
	{
		appGroup.GET("/session", s.ListSession)
		appGroup.POST("/session", s.StartSession)
		appGroup.GET("/session/getMountInfo", s.GetMountInfo)
		appGroup.POST("/session/reboot", s.RebootSession)
		appGroup.POST("/session/powerOff", s.PowerOffSession)
		appGroup.POST("/session/powerOn", s.PowerOnSession)
		appGroup.POST("/session/close", s.CloseSession)
		appGroup.GET("/session/ready", s.ReadySession)
		appGroup.GET("/session/remoteAppUrl", s.GetRemoteAppURL)
		appGroup.GET("/session/projectNames", s.ListUsedProjectNames)
		appGroup.GET("/session/export", s.ExportSessionInfo)
	}

	{
		appGroup.GET("/hardware", s.ListHardware)
		appGroup.POST("/hardware", s.AddHardware)
		appGroup.PUT("/hardware", s.UpdateHardware)
		appGroup.DELETE("/hardware", s.DeleteHardware)
	}

	{
		appGroup.GET("/software", s.ListSoftware)
		appGroup.POST("/software", s.AddSoftware)
		appGroup.PUT("/software", s.UpdateSoftware)
		appGroup.DELETE("/software", s.DeleteSoftware)
		appGroup.PUT("/software/publish", s.PublishSoftware)
		appGroup.GET("/software/usingStatuses", s.ListSoftwareUseStatuses)

		appGroup.GET("/software/preset", s.GetSoftwarePresets)
		appGroup.POST("/software/preset", s.SetSoftwarePresets)

		appGroup.POST("/software/remote/app", s.AddRemoteApp)
		appGroup.PUT("/software/remote/app", s.UpdateRemoteApp)
		appGroup.DELETE("/software/remote/app", s.DeleteRemoteApp)
	}

	{
		appGroup.GET("/statistic/duration", s.DurationStatistic)
		appGroup.GET("/statistic/duration/list", s.ListHistoryDuration)
		appGroup.GET("/statistic/report/duration", s.SessionUsageDurationStatistic)
		appGroup.GET("/statistic/report/duration/export", s.ExportUsageDurationStatistic)
		appGroup.GET("/statistic/report/createNumber", s.SessionCreateNumberStatistic)
		appGroup.GET("/statistic/report/numberStatus", s.SessionNumberStatusStatistic)
	}
}
