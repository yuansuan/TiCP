package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"strconv"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

// Migration 用于对数据库进行自动升级和降级
type Migration struct {
	db  database.Driver
	src source.Driver
}

// AutoMigrate 根据 desc 自动决定调用哪个升级方法
func (m *Migration) AutoMigrate(desc string) error {
	if desc == "up" {
		return m.UpToDate()
	}

	if n, err := strconv.Atoi(desc); err == nil && n > 0 {
		return m.MigrateTo(uint(n))
	}

	return errors.New("migration: invalid migration description")
}

// Close closing the source and keeping the database
func (m *Migration) Close() error {
	return m.src.Close()
}

// MigrateTo 升级到一个指定的版本
func (m *Migration) MigrateTo(version uint) error {
	mm, err := migrate.NewWithInstance("iofs", m.src, "standard-compute", m.db)
	if err != nil {
		return err
	}

	log.Infof("migrate to version %d", version)
	if err = mm.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == migrate.ErrNoChange {
		log.Info("database has not changes")
	}
	return nil
}

// UpToDate 升级到最新的一个版本
func (m *Migration) UpToDate() error {
	mm, err := migrate.NewWithInstance("iofs", m.src, "standard-compute", m.db)
	if err != nil {
		return err
	}

	log.Infof("migrate up to date")
	if err = mm.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == migrate.ErrNoChange {
		log.Info("database has not changes")
	}
	return nil
}

// NewMigration 创建数据库自动升级工具
func NewMigration(engine *xorm.Engine, cfg *config.Config, src source.Driver) (_ *Migration, err error) {
	var db database.Driver
	switch cfg.Database.Type {
	case "mysql":
		db, err = mysql.WithInstance(engine.DB().DB, &mysql.Config{})
	case "sqlite":
		db, err = sqlite.WithInstance(engine.DB().DB, &sqlite.Config{})
	default:
		return nil, errors.New("unsupported database type")
	}

	if err != nil {
		return nil, err
	}

	return &Migration{src: src, db: db}, nil
}
