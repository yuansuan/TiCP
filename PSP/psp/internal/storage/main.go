package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/cmd/docs"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/api/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/rpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/service/daemon"
)

func main() {
	// 创建默认服务
	server := boot.Default()
	defer boot.Recovery()

	// 初始化配置
	config.InitConfig()
	daemon.InitDaemon()

	// 启动服务
	server.Register(
		InitSwaggerDoc,
		api.InitAPI,
		rpc.InitGRPC,
		openapi.InitOpenapiAPI,
	).Run()
}

func InitSwaggerDoc(drv *http.Driver) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	drv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
