package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/handler_rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/migration"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/router"
)

func main() {
	server := boot.Default() //使用默认http server
	logger := logging.Default()

	err := config.InitConfig()
	if err != nil {
		logger.Fatal("err_init_config_failed", err)
	}

	server.
		Register( //注册路由策略
			handler_rpc.InitGRPCServer,
			router.Init,
		).
		RegisterRoutine( //注册go-routine在后台运行
			handler_rpc.InitGRPCClient,
		).
		OnShutdown( //注册退出事件
			handler_rpc.OnShutdown,
		).
		DBAutoMigrate(migration.Mysql).
		Run() //启动运行

}
