package mysql

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	app "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
	"xorm.io/xorm"
)

type Application struct {
	db *xorm.Engine
}

func newApplication(ds *datastore) *Application {
	return &Application{ds.db}
}

func (a *Application) GetApp(ctx context.Context, id snowflake.ID) (*app.Application, error) {
	app := &app.Application{}
	get, err := a.db.Where("id = ?", id).Get(app)
	if err != nil {
		return nil, errors.Wrap(err, "get application error")
	}
	if !get {
		return nil, xorm.ErrNotExist
	}
	return app, nil
}

// ListApps 列出所有应用
func (a *Application) ListApps(ctx context.Context, userID snowflake.ID, publishStatus consts.PublishStatus) ([]*app.Application, int64, error) {
	session := a.db.Context(ctx)
	if publishStatus != consts.PublishStatusAll {
		session = session.Where("application.publish_status = ?", publishStatus)
	}

	if userID != 0 {
		session = session.Join("INNER", "application_quota", "application_quota.application_id = application.id").
			Where("application_quota.ys_id = ?", userID)
	}

	apps := make([]*app.Application, 0)
	count, err := session.Table("application").FindAndCount(&apps)
	if err != nil {
		return nil, 0, errors.Wrap(err, "list application error")
	}
	return apps, count, nil
}

type mysqlErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

func (a *Application) AddApp(ctx context.Context, appInfo *app.Application) error {
	insert, err := a.db.Insert(appInfo)
	if err != nil {
		if db.IsDuplicatedError(err) {
			return common.ErrDuplicateEntry
		}
		return errors.Wrap(err, "add application error")
	}
	if insert == 0 {
		return errors.New("add application error")
	}
	return nil
}

func (a *Application) UpdateApp(ctx context.Context, appInfo *app.Application) error {
	update, err := a.db.ID(appInfo.ID).AllCols().Update(appInfo)
	if err != nil {
		if db.IsDuplicatedError(err) {
			return common.ErrDuplicateEntry
		}
		return errors.Wrap(err, "update application error")
	}
	if update == 0 {
		return common.ErrNoEffect
	}
	return nil
}
