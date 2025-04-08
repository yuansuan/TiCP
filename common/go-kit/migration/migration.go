package migration

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Migrator struct {
	microServiceName string
	forceMigrate     bool

	db                *gorm.DB
	dbType            DatabaseType
	migrationSourceFS fs.FS
	source            *source
}

func NewMigrator(dsn, microServiceName string, migrationSourceFS fs.FS, dbType DatabaseType, forceMigrate bool) (*Migrator, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}

	if microServiceName == "" {
		return nil, fmt.Errorf("microServiceName cannot be empty")
	}

	if migrationSourceFS == nil {
		return nil, fmt.Errorf("migration source fs cannot be nil")
	}

	// TODO just support mysql for now
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("connect mysql failed, %w", err)
	}

	m := &Migrator{
		microServiceName:  microServiceName,
		forceMigrate:      forceMigrate,
		db:                db,
		dbType:            dbType,
		migrationSourceFS: migrationSourceFS,
	}

	if err = m.init(); err != nil {
		return nil, fmt.Errorf("init migrator failed, %w", err)
	}

	return m, nil
}

type versionInfo struct {
	Current int `gorm:"primaryKey;autoIncrement:false"`
}

func (m *Migrator) init() error {
	// get migrate content from migrationSourceFS
	source, err := newSource(m.migrationSourceFS, m.dbType)
	if err != nil {
		return fmt.Errorf("new source failed, %w", err)
	}
	m.source = source

	// auto migrate version table
	if err = m.db.Table(m.versionTableName()).AutoMigrate(&versionInfo{}); err != nil {
		return fmt.Errorf("auto migrate table %s failed, %w", m.versionTableName(), err)
	}

	// check if not version marked in version table, mark it to the latest version
	_, exist, err := m.getVersionFromDB()
	if err != nil {
		return fmt.Errorf("get version from db failed, %w", err)
	}
	if !exist {
		version := 0
		if !m.forceMigrate {
			version = m.source.getLatestVersion()
		}

		res := m.db.Table(m.versionTableName()).Create(&versionInfo{
			Current: version,
		})
		if res.Error != nil {
			return fmt.Errorf("update current version from %d to %d failed, %w", 0, version, res.Error)
		}
	}

	return nil
}

func (m *Migrator) getVersionFromDB() (int, bool, error) {
	verInfo := &versionInfo{}
	res := m.db.Table(m.versionTableName()).First(verInfo)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("get version info from %s failed, %w", m.versionTableName(), res.Error)
	}

	return verInfo.Current, true, nil
}

func (m *Migrator) MigrateTo(version int) error {
	if version <= 0 {
		return fmt.Errorf("version should larger than 0")
	}

	oldVersion, exist, err := m.getVersionFromDB()
	if err != nil {
		return fmt.Errorf("get version from db failed, %w", err)
	}
	if !exist {
		// should not happen here
		return fmt.Errorf("version record not found in db")
	}

	if oldVersion == version { // no need to migrate database
		return nil
	}

	return m.doMigrate(oldVersion, version)
}

func (m *Migrator) MigrateToLatest() error {
	oldVersion, exist, err := m.getVersionFromDB()
	if err != nil {
		return fmt.Errorf("get version from db failed, %w", err)
	}
	if !exist {
		// should not happen here
		return fmt.Errorf("version record not found in db")
	}

	latestVersion := m.source.getLatestVersion()
	if oldVersion == latestVersion { // no need to migrate database
		return nil
	}

	return m.doMigrate(oldVersion, latestVersion)
}

func (m *Migrator) doMigrate(oldVersion, expectVersion int) error {
	var low, high int
	var mType migrationType
	if oldVersion < expectVersion {
		low = oldVersion
		high = expectVersion
		mType = up
	} else {
		low = expectVersion
		high = oldVersion
		mType = down
	}

	head, tail, err := m.source.subLinkList(low, high)
	if err != nil {
		return fmt.Errorf("get sub link list failed, %w", err)
	}

	var curr *sourceNode
	switch mType {
	case up:
		// start from head to exec
		curr = head
		for curr != nil {
			// the script in head should not be executed
			// for example upgrade from version 1 to version 3, got the list are 1,2,3 three source nodes. should only exec 2,3 upgrade script
			if curr.prev != nil {
				if err = m.execScriptAndUpdateVersion(curr, up); err != nil {
					return err
				}
			}

			curr = curr.next
		}
	case down:
		// start from tail to exec
		curr = tail
		for curr != nil {
			// the script in head should not be executed
			// for example downgrade from version 3 to version 1, got the list are 1,2,3 three source nodes. should only exec 3,2 downgrade script
			if curr.prev != nil {
				if err = m.execScriptAndUpdateVersion(curr, down); err != nil {
					return err
				}
			}
			curr = curr.prev
		}
	}

	return nil
}

func (m *Migrator) versionTableName() string {
	return fmt.Sprintf("%s_version", m.microServiceName)
}

func (m *Migrator) execScriptAndUpdateVersion(curr *sourceNode, mType migrationType) error {
	var sqlScript, filename string
	switch mType {
	case up:
		sqlScript = string(curr.upgrade.sqlScript)
		filename = curr.upgrade.filename
	case down:
		sqlScript = string(curr.downgrade.sqlScript)
		filename = curr.downgrade.filename
	}

	sqlStats := strings.Split(sqlScript, ";")
	err := m.db.Transaction(func(tx *gorm.DB) error {
		var e error
		for _, sqlStat := range sqlStats {
			if strings.TrimSpace(sqlStat) == "" {
				continue
			}
			e = tx.Exec(sqlStat).Error
			if e != nil {
				return fmt.Errorf("exec sql script [%s] failed, %w", filename, e)
			}
		}

		switch mType {
		case up:
			if e = tx.Table(m.versionTableName()).Where("current = ?", curr.version-1).Update("current", curr.version).Error; e != nil {
				return fmt.Errorf("upgrade version info table failed, %w", e)
			}
		case down:
			if e = tx.Table(m.versionTableName()).Where("current = ?", curr.version).Update("current", curr.version-1).Error; e != nil {
				return fmt.Errorf("downgrade version info table failed, %w", e)
			}
		}

		return nil
	})
	if err != nil {
		switch mType {
		case up:
			return fmt.Errorf("upgrade db version from %d to %d failed, %w", curr.version-1, curr.version, err)
		case down:
			return fmt.Errorf("downgrade db version from %d to %d failed, %w", curr.version, curr.version-1, err)
		}
	}

	return nil
}
