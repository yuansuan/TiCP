package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service/impl"
)

type RouteService struct {
	SysConfigService service.SysConfigService
}

func NewUserService() (*RouteService, error) {
	sysConfigService, err := impl.NewSysConfigService()
	if err != nil {
		logging.Default().Errorf("init sys config server service err: %v", err)
		return nil, err
	}

	return &RouteService{
		SysConfigService: sysConfigService,
	}, nil
}

func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewUserService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	sysConfigGroup := drv.Group("/api/v1/sysconfig")
	{
		sysConfigGroup.GET("/global", s.GetGlobalSysConfig)
		sysConfigGroup.GET("/getJobConfig", s.GetJobConfig)
		sysConfigGroup.POST("/setJobConfig", s.SetJobConfig)
		sysConfigGroup.GET("/getJobBurstConfig", s.GetJobBurstConfig)
		sysConfigGroup.POST("/setJobBurstConfig", s.SetJobBurstConfig)

		sysConfigGroup.POST("/setEmailConfig", s.SetEmailConfig)
		sysConfigGroup.GET("/getEmailConfig", s.GetEmailConfig)

		// email
		sysConfigGroup.GET("/globalEmail", s.GetGlobalEmail)
		sysConfigGroup.POST("/globalEmail", s.SetGlobalEmail)
		sysConfigGroup.POST("/email/testSend", s.TestSendEmail)

		sysConfigGroup.GET("/getThreePersonManagementConfig", s.GetThreePersonManagementConfig)
		sysConfigGroup.POST("/setThreePersonManagementConfig", s.SetThreePersonManagementConfig)

	}
}
