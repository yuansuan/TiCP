/*
 * Copyright (C) 2019 LambdaCal Inc.
 */

package casbin

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"sync"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"xorm.io/xorm"
)

// Dao Dao
type Dao struct {
	enforcer *casbin.SyncedEnforcer
}

var (
	once sync.Once
	dao  *Dao
)

// NewEnforcer NewEnforcer
func NewEnforcer(xormEngine *xorm.Engine, rbacConfPath string) *casbin.SyncedEnforcer {
	logger := logging.Default()
	a, err := xormadapter.NewAdapterByEngine(xormEngine)
	if err != nil {
		logger.Fatalf("create xorm adapter error: %v", err)
	}
	e, err := casbin.NewSyncedEnforcer(rbacConfPath, a)
	if err != nil {
		logger.Fatalf("create casbin enforcer error: %v", err)
	}
	e.LoadPolicy()
	return e
}

// New New
func New(xormEngine *xorm.Engine, rbacConfPath string) *Dao {
	once.Do(func() {
		dao = &Dao{
			NewEnforcer(xormEngine, rbacConfPath),
		}
	})

	return dao
}
