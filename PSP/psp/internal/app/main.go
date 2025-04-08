package main

import (
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/yuansuan/ticp/PSP/psp/cmd/docs"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/api/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/rpc"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/service"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

func main() {
	server := boot.Default()
	defer boot.Recovery()

	config.InitConfig()
	service.InitDaemon()

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
