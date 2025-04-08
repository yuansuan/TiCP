package migration

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/yuansuan/ticp/common/go-kit/migration"
)

func Migrate(dsn, microServiceName, version string, migrationSourceFS fs.FS, forceMigrate bool) error {
	migrateToLatest := false
	var v int
	var err error
	if version == "up" {
		migrateToLatest = true
	} else {
		v, err = strconv.Atoi(version)
		if err != nil {
			return fmt.Errorf("version not set 'up' but convert version %s to int failed, %w", version, err)
		}
	}

	migrator, err := migration.NewMigrator(dsn, microServiceName, migrationSourceFS, migration.Mysql, forceMigrate)
	if err != nil {
		return fmt.Errorf("new migrator failed, %w", err)
	}

	if migrateToLatest {
		return migrator.MigrateToLatest()
	}

	return migrator.MigrateTo(v)
}
