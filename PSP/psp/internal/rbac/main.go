package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/rpc"
)

// main
//
//	@title			PSP
//	@version		1.1
//	@description	PSP server
//	@host			127.0.0.1:8889
//	@BasePath		/api/v1
func main() {
	// 创建默认服务
	server := boot.Default()
	defer boot.Recovery()

	// 初始化
	config.InitConfig()

	// 启动服务
	server.Register(
		api.InitAPI,
		rpc.InitGRPC,
	).Run()
}
