package main

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	migration "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/docs/migrations"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

func main() {
	server := boot.Default()
	if err := config.InitConfig(); err != nil {
		logging.Default().Fatalf("init config failed: %s", err)
	}

	logger := logging.Default()
	logger.Infow("dump config", "config", config.GetConfig())

	server.Register(
		api.Init,
	).RegisterRoutine().OnShutdown().DBAutoMigrate(migration.Mysql).Run()
}
