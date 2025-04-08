package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/api"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/rpc"
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
		api.InitAPI,
		rpc.InitGRPC,
	).Run()
}
