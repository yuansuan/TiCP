package mysql

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"xorm.io/xorm"
)

// ApplicationAllow 应用白名单
type ApplicationAllow struct {
	db *xorm.Engine
}

func newApplicationAllow(ds *datastore) *ApplicationAllow {
	return &ApplicationAllow{ds.db}
}

// GetByAppId 根据应用id查询白名单
func (a *ApplicationAllow) GetByAppId(ctx context.Context, session *xorm.Session,
	applicationID snowflake.ID) (*models.ApplicationAllow, error) {
	if session == nil {
		session = a.db.Context(ctx)
	}

	allow := &models.ApplicationAllow{ApplicationID: applicationID}
	has, err := session.Get(allow)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, xorm.ErrNotExist
	}
	return allow, nil
}
