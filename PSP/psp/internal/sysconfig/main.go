package main

import (
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/yuansuan/ticp/PSP/psp/cmd/docs"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/rpc"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

// main
//
//	@title			PSP
//	@version		1.1
//	@description	PSP server
//	@host			127.0.0.1:8889
//	@BasePath		/api/v1
func main() {
	server := boot.Default()
	defer boot.Recovery()

	config.InitConfig()

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
