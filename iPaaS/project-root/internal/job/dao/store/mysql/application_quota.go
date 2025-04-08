package mysql

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"xorm.io/xorm"
)

// ApplicationQuota 应用配额
type ApplicationQuota struct {
	db *xorm.Engine
}

func newApplicationQuota(ds *datastore) *ApplicationQuota {
	return &ApplicationQuota{ds.db}
}

// GetByUser 根据用户id和应用id查询配额
func (a *ApplicationQuota) GetByUser(ctx context.Context, session *xorm.Session, applicationID, userID snowflake.ID, forUpdate bool) (*models.ApplicationQuota, error) {
	if session == nil {
		session = a.db.Context(ctx)
	}

	quota := &models.ApplicationQuota{YsID: userID, ApplicationID: applicationID}
	if forUpdate {
		session = session.ForUpdate()
	}
	has, err := session.Get(quota)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, xorm.ErrNotExist
	}
	return quota, nil
}
