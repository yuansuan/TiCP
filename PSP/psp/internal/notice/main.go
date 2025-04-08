package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/notice/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/rpc"
)

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
