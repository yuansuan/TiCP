package migrations

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

//go:embed mysql/*.sql
var mysql embed.FS

//go:embed sqlite/*.sql
var sqlite embed.FS

// NewSource 创建
func NewSource(cfg *config.Config) (source.Driver, error) {
	switch cfg.Database.Type {
	case "mysql":
		return iofs.New(mysql, cfg.Database.Type)
	case "sqlite":
		return iofs.New(sqlite, cfg.Database.Type)
	default:
		return nil, fmt.Errorf("unsupport database source")
	}
}
