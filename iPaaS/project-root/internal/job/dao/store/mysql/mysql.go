package mysql

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	"xorm.io/xorm"
)

type datastore struct {
	db *xorm.Engine
}

func (ds *datastore) Applications() dao.ApplicationDao {
	return newApplication(ds)
}

func GetMysqlFactory() store.Factory {
	engine := boot.MW.DefaultORMEngine()
	return &datastore{engine}
}

type datastoreAdapter struct {
	store.Factory
}

func (ds *datastoreAdapter) Engine() *xorm.Engine {
	return ds.Factory.(*datastore).db
}

func (ds *datastoreAdapter) ApplicationQuota() dao.ApplicationQuotaDao {
	return newApplicationQuota(ds.Factory.(*datastore))
}

func (ds *datastoreAdapter) ApplicationAllow() dao.ApplicationAllowDao {
	return newApplicationAllow(ds.Factory.(*datastore))
}

// GetMysqlFactoryWithEngine return a factory with engine.
func GetMysqlFactoryWithEngine() store.FactoryNew {
	return &datastoreAdapter{GetMysqlFactory()}
}
