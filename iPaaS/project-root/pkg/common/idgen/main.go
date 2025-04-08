/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/idgen/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/idgen/handler_rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/idgen/router"
)

func main() {
	server := boot.Default() //使用默认http server
	config.InitConfig()

	logger := logging.Default()
	logger.Infof("%#v", config.GetConfig())

	server.Register( //注册路由策略
		router.UseRoutersGenerated,
		handler_rpc.InitGRPCServer,
	).
		RegisterRoutine( //注册go-routine在后台运行
		).
		OnShutdown( //注册退出事件
			handler_rpc.OnShutdown,
		).
		Run() //启动运行

}
