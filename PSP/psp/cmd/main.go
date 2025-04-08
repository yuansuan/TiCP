package main

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yuansuan/ticp/PSP/psp/cmd/consts"
	"os"

	"github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/cmd/docs"
	appapi "github.com/yuansuan/ticp/PSP/psp/internal/app/api"
	appopenapi "github.com/yuansuan/ticp/PSP/psp/internal/app/api/openapi"
	apprpc "github.com/yuansuan/ticp/PSP/psp/internal/app/rpc"
	appsrv "github.com/yuansuan/ticp/PSP/psp/internal/app/service"
	approveapi "github.com/yuansuan/ticp/PSP/psp/internal/approve/api"
	approverpc "github.com/yuansuan/ticp/PSP/psp/internal/approve/rpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway"
	jobapi "github.com/yuansuan/ticp/PSP/psp/internal/job/api"
	jobopenapi "github.com/yuansuan/ticp/PSP/psp/internal/job/api/openapi"
	jobrpc "github.com/yuansuan/ticp/PSP/psp/internal/job/rpc"
	jobsrv "github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	licenseapi "github.com/yuansuan/ticp/PSP/psp/internal/license/api"
	licenserpc "github.com/yuansuan/ticp/PSP/psp/internal/license/rpc"
	monitorapi "github.com/yuansuan/ticp/PSP/psp/internal/monitor/api"
	monitorrpc "github.com/yuansuan/ticp/PSP/psp/internal/monitor/rpc"
	monitor "github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/dataloader"
	noticeapi "github.com/yuansuan/ticp/PSP/psp/internal/notice/api"
	noticerpc "github.com/yuansuan/ticp/PSP/psp/internal/notice/rpc"
	projectapi "github.com/yuansuan/ticp/PSP/psp/internal/project/api"
	projectopenapi "github.com/yuansuan/ticp/PSP/psp/internal/project/api/openapi"
	projectrpc "github.com/yuansuan/ticp/PSP/psp/internal/project/rpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/daemon"
	rbacapi "github.com/yuansuan/ticp/PSP/psp/internal/rbac/api"
	rbacrpc "github.com/yuansuan/ticp/PSP/psp/internal/rbac/rpc"
	storageapi "github.com/yuansuan/ticp/PSP/psp/internal/storage/api"
	storageopenapi "github.com/yuansuan/ticp/PSP/psp/internal/storage/api/openapi"
	storagerpc "github.com/yuansuan/ticp/PSP/psp/internal/storage/rpc"
	sysconfigapi "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/api"
	sysconfigrpc "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/rpc"
	userapi "github.com/yuansuan/ticp/PSP/psp/internal/user/api"
	userrpc "github.com/yuansuan/ticp/PSP/psp/internal/user/rpc"
	visualapi "github.com/yuansuan/ticp/PSP/psp/internal/visual/api"
	visualopenapi "github.com/yuansuan/ticp/PSP/psp/internal/visual/api/openapi"
	visualsrv "github.com/yuansuan/ticp/PSP/psp/internal/visual/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

// main
//
//	@title			PSP
//	@version		4.2.0
//	@description	PSP server
//	@BasePath		/api/v1
//	@host			127.0.0.1:8889
func main() {
	logLevel := config.GetLoaderLogLevel()
	configPath := os.Getenv(consts.ConfigPath)
	server := boot.DefaultServer(configPath, logLevel) //使用默认http server

	config.InitConfig()

	//初始化日志
	defer tracelog.CloseLogger()
	tracelog.InitLogger()

	// service daemon
	{
		appsrv.InitDaemon()
		jobsrv.InitDaemon()
		visualsrv.InitDaemon()
		monitor.InitDaemon()
		//storagesrv.InitDaemon()
		daemon.InitDaemon()
	}

	server.
		Register(
			InitSwaggerDoc,
			gateway.InitGateway,
			jobapi.InitAPI,
			jobrpc.InitGRPC,
			jobopenapi.InitOpenapiAPI,
			rbacrpc.InitGRPC,
			rbacapi.InitAPI,
			storagerpc.InitGRPC,
			storageapi.InitAPI,
			storageopenapi.InitOpenapiAPI,
			userrpc.InitGRPC,
			userapi.InitAPI,
			noticerpc.InitGRPC,
			noticeapi.InitAPI,
			appapi.InitAPI,
			apprpc.InitGRPC,
			appopenapi.InitOpenapiAPI,
			monitorapi.InitAPI,
			visualapi.InitAPI,
			visualopenapi.InitOpenapiAPI,
			sysconfigapi.InitAPI,
			sysconfigrpc.InitGRPC,
			monitorrpc.InitGRPC,
			licenseapi.InitAPI,
			licenserpc.InitGRPC,
			approverpc.InitGRPC,
			approveapi.InitAPI,
			projectapi.InitAPI,
			projectrpc.InitGRPC,
			projectopenapi.InitOpenapiAPI,
		).
		RegisterRoutine().
		OnShutdown().
		Run()
}

func InitSwaggerDoc(drv *http.Driver) {
	swagger := config.Custom.Main.Swagger
	if swagger.Enable {
		if swagger.Host != "" && swagger.Port != "" {
			docs.SwaggerInfo.Host = fmt.Sprintf("%v:%v", swagger.Host, swagger.Port)
		}
		docs.SwaggerInfo.BasePath = "/api/v1"
		drv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
