package router

import (
	"github.com/gin-gonic/gin"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/router/validation"
)

// InitHTTPHandlers ...
func InitHTTPHandlers(drv *http.Driver) {
	// 初始化自定义请求参数验证器
	err := validation.Initialize()
	if err != nil {
		logging.Default().Panicf("init request validation error, err: %v", err)
	}

	logging.Default().Info("setup router")
	drv.Use(overrideLoggerInGinCtx())
	drv.Use(middleware.IngressLogger(middleware.IngressLoggerConfig{
		IsLogRequestHeader:  true,
		IsLogRequestBody:    true,
		IsLogResponseHeader: true,
		IsLogResponseBody:   true,
	}))

	group := drv.Group("/admin")
	{
		licenseImpl := dao.NewLicenseImpl(boot.MW.DefaultORMEngine())
		licenseHandler := api.NewLicenseHandler(licenseImpl)
		licenseManageGroup := group.Group("/licenseManagers")
		licenseManageGroup.POST("", licenseHandler.AddLicenseManage)
		licenseManageGroup.DELETE("/:id", licenseHandler.DeleteLicenseManager)
		licenseManageGroup.PUT("/:id", licenseHandler.PutLicenseManager)
		licenseManageGroup.GET("", licenseHandler.ListLicenseManage)
		licenseManageGroup.GET("/:id", licenseHandler.GetLicenseManage)

		licenseGroup := group.Group("/licenses")
		licenseGroup.POST("", licenseHandler.AddLicense)
		licenseGroup.DELETE("/:id", licenseHandler.DeleteLicense)
		licenseGroup.PUT("/:id", licenseHandler.PutLicense)

		moduleCfgGroup := group.Group("/moduleConfigs")
		moduleCfgGroup.GET("", licenseHandler.ListModules)
		moduleCfgGroup.POST("", licenseHandler.AddModule)
		moduleCfgGroup.POST("/batch", licenseHandler.BatchAddModules)
		moduleCfgGroup.PUT("/:id", licenseHandler.PutModule)
		moduleCfgGroup.DELETE("/:id", licenseHandler.DeleteModule)

		//licenseManageGroup.GET("/usedStatic", api.UsedGroupByJobs)
	}
}

func overrideLoggerInGinCtx() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID := trace.GetRequestId(context)
		context.Set(trace.RequestIdKey, requestID)
		context.Set(logging.LoggerName, logging.Default().With(trace.RequestIdKey, requestID))
	}
}
