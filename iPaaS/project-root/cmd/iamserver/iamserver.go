package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/migration"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

func main() {
	server := boot.Default()

	if err := config.InitConfig(); err != nil {
		logging.Default().DPanicf("init config error. err: %v", err)
	}

	server.Register(iamserver.InitRouter).DBAutoMigrate(migration.Mysql).Run()
}
