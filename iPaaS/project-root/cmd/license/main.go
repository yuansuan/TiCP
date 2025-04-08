package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/handler_rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/migration"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/router"
)

func main() {
	server := boot.Default() //使用默认http server
	logger := logging.Default()

	err := config.InitConfig()
	if err != nil {
		logger.Fatalf("init config error, err: %v", err)
	}

	server.
		Register( //注册路由策略
			handler_rpc.InitGRPCServer,
			router.InitHTTPHandlers,
		).RegisterRoutine().OnShutdown().
		DBAutoMigrate(migration.Mysql).Run() //启动运行
}
