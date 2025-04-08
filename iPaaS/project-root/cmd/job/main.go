package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/handler_rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/migration"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/mongo"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/router"
)

func main() {
	server := boot.Default() //使用默认http server

	if err := config.InitConfig(); err != nil {
		logging.Default().DPanicf("init config error. err: %v", err)
	}

	server.Register(
		func(_ *http.Driver) {
			// before handler_rpc.InitGRPCServer and router.Init
			if config.GetConfig().Mongo != nil && config.GetConfig().Mongo.Enable {
				if err := mongo.Init(config.GetConfig().Mongo.URI()); err != nil {
					logging.Default().DPanicf("init mongo error. err: %v", err)
				}
			}
		},
		handler_rpc.InitGRPCServer,
		router.Init).
		OnShutdown(func(_ *http.Driver) {
			if config.GetConfig().Mongo != nil && config.GetConfig().Mongo.Enable {
				mongo.Shutdown()
			}
		}).
		DBAutoMigrate(migration.Mysql).Run()
}
