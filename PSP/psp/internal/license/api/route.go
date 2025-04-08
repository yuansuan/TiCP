package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/license/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/service/impl"
)

type apiRoute struct {
	licManagerService   service.LicenseManagerService
	licInfoService      service.LicenseInfoService
	moduleConfigService service.ModuleConfigService
}

func NewAPIRoute() (*apiRoute, error) {
	licManagerService, err := impl.NewLicenseManagerService()
	if err != nil {
		return nil, err
	}

	licInfoService, err := impl.NewLicenseInfoService()
	if err != nil {
		return nil, err
	}

	moduleConfigService, err := impl.NewModuleConfigService()
	if err != nil {
		return nil, err
	}

	return &apiRoute{
		licManagerService:   licManagerService,
		licInfoService:      licInfoService,
		moduleConfigService: moduleConfigService,
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
		licenseManagersGroup := group.Group("/licenseManagers")
		licenseManagersGroup.GET("", api.ListLicenseManager)
		licenseManagersGroup.GET("/:id", api.LicenseManageInfo)
		licenseManagersGroup.GET("/typeList", api.LicenseTypeList)
		licenseManagersGroup.POST("", api.AddLicenseManager)
		licenseManagersGroup.PUT("/:id", api.EditLicenseManager)
		licenseManagersGroup.DELETE("/:id", api.DeleteLicenseManager)
	}

	{
		licenseInfoGroup := group.Group("/licenseInfos")
		licenseInfoGroup.POST("", api.AddLicenseInfo)
		licenseInfoGroup.PUT("/:id", api.EditLicenseInfo)
		licenseInfoGroup.DELETE("/:id", api.DeleteLicenseInfo)
		licenseInfoGroup.GET("/:id/moduleConfigs", api.ListModuleConfig)
	}

	{
		moduleConfigsGroup := group.Group("/moduleConfigs")
		moduleConfigsGroup.POST("", api.AddModuleConfig)
		moduleConfigsGroup.PUT("/:id", api.EditModuleConfig)
		moduleConfigsGroup.DELETE("/:id", api.DeleteModuleConfig)
	}

}
