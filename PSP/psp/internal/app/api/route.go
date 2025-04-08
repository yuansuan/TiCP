package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/impl"
)

type RouteService struct {
	appService service.AppService
}

func NewRouteService() (*RouteService, error) {
	logger := logging.Default()
	appService, err := impl.NewAppService()
	if err != nil {
		logger.Errorf("init app server service err: %v", err)
		return nil, err
	}
	return &RouteService{
		appService: appService,
	}, nil
}

func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewRouteService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	appGroup := drv.Group("/api/v1/app")
	{
		appGroup.GET("/list", s.ListApp)
		appGroup.GET("/template/list", s.ListTemplate)
		appGroup.GET("/template", s.GetAppInfo)
		appGroup.POST("/template", s.AddApp)
		appGroup.PUT("/template", s.UpdateApp)
		appGroup.DELETE("/template", s.DeleteApp)
		appGroup.PUT("/template/publish", s.PublishApp)
		appGroup.POST("/template/syncAppContent", s.SyncAppContent)
		appGroup.GET("/zone", s.ListZone)
		appGroup.GET("/queue", s.ListQueue)
		appGroup.GET("/license", s.ListLicense)
		appGroup.GET("/schedulerResourceKey", s.GetSchedulerResourceKey)
		appGroup.GET("/schedulerResourceValue", s.GetSchedulerResourceValue)
	}
}
