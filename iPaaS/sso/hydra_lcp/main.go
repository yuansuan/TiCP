package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/handler_rpc"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/migration"
)

func main() {
	server := boot.Default() //使用默认http server

	logger := logging.Default()
	logger.Info(config.InitConfig())

	logger.Infof("%#v", config.Custom)

	server.Register( //注册路由策略
		handler_rpc.InitGRPCServer,
	).
		DBAutoMigrate(migration.Mysql).
		Run() //启动运行

}
