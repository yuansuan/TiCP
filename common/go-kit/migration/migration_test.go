package migration

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/common/go-kit/migration/example"
)

const (
	dbName  = "migration_unit_test_db"
	address = "0.0.0.0"
)

type TableA struct {
	Id             int64
	AlterAddColumn string
}

// database initialized and code got auto migrate scripts without version table
func TestMigrateFromNoVersionCaseWhenDatabaseAlreadyInitDone(t *testing.T) {
	testDB := func() *memory.Database {
		db := memory.NewDatabase(dbName)
		db.EnablePrimaryKeyIndexes()

		personTableName := "person"
		db.AddTable(personTableName, memory.NewTable(personTableName, sql.NewPrimaryKeySchema(sql.Schema{
			{Name: "id", Type: types.Int64, Nullable: false, Source: personTableName, PrimaryKey: true},
			{Name: "name", Type: types.Text, Nullable: false, Source: personTableName},
		}), db.GetForeignKeyCollection()))

		tableATableName := "tableA"
		db.AddTable(tableATableName, memory.NewTable(tableATableName, sql.NewPrimaryKeySchema(sql.Schema{
			{Name: "id", Type: types.Int64, Nullable: false, Source: tableATableName, PrimaryKey: true},
			{Name: "alter_add_column", Type: types.Text, Nullable: false, Source: tableATableName},
		}), db.GetForeignKeyCollection()))

		return db
	}()

	dsn, err := mockMysqlServer(t, testDB)
	assert.NoError(t, err)

	migrator, err := NewMigrator(dsn, "micro_service_a", example.Mysql, Mysql, false)
	assert.NoError(t, err)
	assert.NotNil(t, migrator)

	hasTableA := migrator.db.Migrator().HasTable("tableA")
	assert.True(t, hasTableA)
	hasColumn := migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column")
	assert.True(t, hasColumn)

	err = migrator.MigrateTo(2)
	assert.NoError(t, err)
	hasColumn = migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column")
	assert.False(t, hasColumn)

	err = migrator.MigrateTo(1)
	assert.NoError(t, err)
	hasTableA = migrator.db.Migrator().HasTable("tableA")
	assert.False(t, hasTableA)

	err = migrator.MigrateTo(3)
	assert.NoError(t, err)
	hasTableA = migrator.db.Migrator().HasTable("tableA")
	assert.True(t, hasTableA)
	hasColumn = migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column")
	assert.True(t, hasColumn)

	err = migrator.MigrateTo(4)
	assert.Error(t, err)
	t.Logf("%v", err)
}

// database not initialized but code contains auto migrate scripts
func TestMigrateFromDatabaseNotInitialized(t *testing.T) {
	dsn, err := mockMysqlServer(t, memory.NewDatabase(dbName))
	assert.NoError(t, err)

	migrator, err := NewMigrator(dsn, "micro_service_a", example.Mysql, Mysql, true)
	assert.NoError(t, err)
	assert.NotNil(t, migrator)

	assert.NoError(t, migrator.MigrateToLatest())

	assert.True(t, migrator.db.Migrator().HasTable("person"))
	assert.True(t, migrator.db.Migrator().HasTable("tableA"))
	assert.True(t, migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column"))

	assert.NoError(t, migrator.MigrateTo(1))
	assert.True(t, migrator.db.Migrator().HasTable("person"))
	assert.False(t, migrator.db.Migrator().HasTable("tableA"))
	assert.False(t, migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column"))
}

func TestMigrateToLatest(t *testing.T) {
	microServerName := "micro_service_test_name"

	testDB := func() *memory.Database {
		db := memory.NewDatabase(dbName)
		db.EnablePrimaryKeyIndexes()

		personTableName := "person"
		db.AddTable(personTableName, memory.NewTable(personTableName, sql.NewPrimaryKeySchema(sql.Schema{
			{Name: "id", Type: types.Int64, Nullable: false, Source: personTableName, PrimaryKey: true},
			{Name: "name", Type: types.Text, Nullable: false, Source: personTableName},
		}), db.GetForeignKeyCollection()))

		versionTableName := fmt.Sprintf("%s_version", microServerName)
		versionTable := memory.NewTable(versionTableName, sql.NewPrimaryKeySchema(sql.Schema{
			{Name: "current", Type: types.Int64, Nullable: false, Source: versionTableName},
		}), db.GetForeignKeyCollection())
		db.AddTable(versionTableName, versionTable)
		ctx := sql.NewEmptyContext()
		err := versionTable.Insert(ctx, sql.NewRow(int64(1)))
		assert.NoError(t, err)

		return db
	}()

	dsn, err := mockMysqlServer(t, testDB)
	assert.NoError(t, err)

	migrator, err := NewMigrator(dsn, microServerName, example.Mysql, Mysql, false)
	assert.NoError(t, err)
	assert.NotNil(t, migrator)

	version, exist, err := migrator.getVersionFromDB()
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, 1, version)

	err = migrator.MigrateToLatest()
	assert.NoError(t, err)

	version, exist, err = migrator.getVersionFromDB()
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, 3, version)

	hasTableA := migrator.db.Migrator().HasTable("tableA")
	assert.True(t, hasTableA)
	hasColumn := migrator.db.Table("tableA").Migrator().HasColumn(&TableA{}, "alter_add_column")
	assert.True(t, hasColumn)
}

func TestNilCheck(t *testing.T) {
	_, err := NewMigrator("", "micro_service_test", example.Mysql, Mysql, false)
	assert.Error(t, err)

	_, err = NewMigrator("fake_dsn", "", example.Mysql, Mysql, false)
	assert.Error(t, err)

	_, err = NewMigrator("fake_dsn", "fake_micro_service_name", nil, Mysql, false)
	assert.Error(t, err)
}

func getRandomPortUsable() int {
	p := randomPort()
	if !isPortUnused(p) {
		return getRandomPortUsable()
	}

	return p
}

// from 1000 to 65535
func randomPort() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(64536) + 1000
}

func isPortUnused(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// 如果监听失败，说明端口已经被占用
		return false
	}
	defer func() {
		_ = listener.Close()
	}()

	return true
}

func mockMysqlServer(t *testing.T, dbs ...sql.Database) (string, error) {
	engine := sqle.NewDefault(
		memory.NewDBProvider(dbs...),
	)

	port := getRandomPortUsable()
	t.Logf("listen on %s:%d", address, port)

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		return "", fmt.Errorf("new default server failed, %w", err)
	}

	startErrChan := make(chan error)
	go func() {
		if err = s.Start(); err != nil {
			startErrChan <- fmt.Errorf("start server failed, %w", err)
		}
	}()
	select {
	case err = <-startErrChan:
		return "", err
	case <-time.After(1 * time.Second):
		return fmt.Sprintf("tcp(%s:%d)/%s", "127.0.0.1", port, dbName), nil
	}
}
