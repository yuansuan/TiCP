package store

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"xorm.io/xorm"
)

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store -destination mock_store.go -package store github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store Factory,Quota,Allow,FactoryNew,ApplicationStore,ApplicationQuotaStore,ApplicationAllowStore

var client Factory

// ApplicationStore align dao.ApplicationDao
type ApplicationStore = dao.ApplicationDao

// Factory is a store factory.
type Factory interface {
	Applications() ApplicationStore
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	client = factory
}

// ApplicationQuotaStore align dao.ApplicationQuotaDao
type ApplicationQuotaStore = dao.ApplicationQuotaDao

// ApplicationAllowStore align dao.ApplicationAllowDao
type ApplicationAllowStore = dao.ApplicationAllowDao

// Quota is a quota store.
type Quota interface {
	ApplicationQuota() ApplicationQuotaStore
}

// Allow is a allow store.
type Allow interface {
	ApplicationAllow() ApplicationAllowStore
}

// FactoryNew is a factory with engine.
type FactoryNew interface {
	Factory
	Quota
	Allow
	Engine() *xorm.Engine
}
