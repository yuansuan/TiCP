package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/migration"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/router"
)

func main() {
	server := boot.Default() //使用默认http server

	if err := config.InitConfig(); err != nil {
		logging.Default().DPanicf("init config error. err: %v", err)
	}

	server.Register(
		router.Init,
	).DBAutoMigrate(migration.Mysql).Run()
}
