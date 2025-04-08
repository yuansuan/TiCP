package database

import (
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

func NewOrm(cfg *config.Config) (*xorm.Engine, error) {
	return xorm.NewEngine(cfg.Database.Type, cfg.Database.DSN)
}
