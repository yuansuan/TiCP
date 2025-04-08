package main

import (
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/yuansuan/ticp/PSP/psp/cmd/docs"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/rpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/dataloader"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

func main() {
	// 创建默认服务
	server := boot.Default()
	defer boot.Recovery()

	// 初始化
	config.InitConfig()

	dataloader.InitDaemon()

	// 启动服务
	server.Register(
		InitSwaggerDoc,
		api.InitAPI,
		rpc.InitGRPC,
	).Run()
}

func InitSwaggerDoc(drv *http.Driver) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	drv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
