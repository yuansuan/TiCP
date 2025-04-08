package mysql

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"

	"github.com/marmotedu/errors"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/db"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/logger"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) MigrateDatabase() error {
	return MigrateDatabase(ds.db)
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds.db)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

func (ds *datastore) PolicyAudits() store.PolicyAuditStore {
	return newPolicyAudit(ds)
}

func (ds *datastore) Roles() store.RoleStore {
	return newRoles(ds)
}

func (ds *datastore) RolePolicyRelations() store.RolePolicyRelationStore {
	return newRolePolicyRelation(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func GetDB(factory store.Factory) (*gorm.DB, error) {
	d, ok := factory.(*datastore)
	if !ok {
		return nil, errors.New("invalid factory")
	}
	return d.db, nil
}

// GetMySQLFactoryOr create dao factory with the given configs.
func GetMySQLFactoryOr() (store.Factory, error) {
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:               config.GetConfig().Database.Host,
			Username:           config.GetConfig().Database.UserName,
			Password:           config.GetConfig().Database.Password,
			Database:           config.GetConfig().Database.DBName,
			MaxIdleConnections: 10,
			MaxOpenConnections: 50,
			LogLevel:           2,
			Logger:             logger.New(2),
		}
		dbIns, err = db.New(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		mysqlFactory = &datastore{dbIns}

	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get dao store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	go cleanExpireSecretInDB(mysqlFactory)

	return mysqlFactory, nil
}

// only for test transaction
func SetMySQLFactory(dbIns *gorm.DB) (store.Factory, error) {

	mysqlFactory = &datastore{dbIns}
	return mysqlFactory, nil
}

func cleanExpireSecretInDB(factory store.Factory) {
	logging.Default().Info("start clean expire secret in db")
	for range time.Tick(1 * time.Minute) {
		// a hour ago
		expireTime := time.Now().Add(-1 * time.Hour)
		if err := factory.Secrets().CleanExpireSecret(context.Background(), expireTime); err != nil {
			logging.Default().Warnf("clean expire secret failed, error: %+v", err)
		}
	}
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&dao.Role{}); err != nil {
		return errors.Wrap(err, "migrate user model failed")
	}
	if err := db.AutoMigrate(&dao.Policy{}); err != nil {
		return errors.Wrap(err, "migrate policy model failed")
	}
	if err := db.AutoMigrate(&dao.Secret{}); err != nil {
		return errors.Wrap(err, "migrate secret model failed")
	}
	if err := db.AutoMigrate(&dao.RolePolicyRelation{}); err != nil {
		return errors.Wrap(err, "migrate role policy relation model failed")
	}
	return nil
}
