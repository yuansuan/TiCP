package main

import (
	"github.com/yuansuan/ticp/common/go-kit/example/dbautomigrate/migration"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
)

func main() {
	server := boot.Default()

	server.DBAutoMigrate(migration.Mysql).Run()
}
