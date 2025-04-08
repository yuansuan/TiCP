//go:generate mv generate.text generate.go
//go:generate go build -o make-license generate.go
//go:generate ./make-license
//go:generate mv generate.go generate.text

package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/user/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/rpc"
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
		rpc.InitGRPC,
		api.InitAPI,
	).Run()
}
