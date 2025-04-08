package fake

import (
	"context"
	"errors"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type application struct {
	ds *datastore
}

func newApplication(ds *datastore) *application {
	return &application{ds}
}

func (a *application) GetApp(ctx context.Context, id snowflake.ID) (*dbModels.Application, error) {
	a.ds.Lock()
	defer a.ds.Unlock()
	for _, app := range a.ds.apps {
		if app.ID == id {
			return app, nil
		}
	}
	return nil, errors.New("application not found")
}

func (a *application) ListApps(ctx context.Context, userID snowflake.ID, publishStatus consts.PublishStatus) ([]*dbModels.Application, int64, error) {
	a.ds.Lock()
	defer a.ds.Unlock()
	return a.ds.apps, int64(len(a.ds.apps)), nil
}

func (a *application) AddApp(ctx context.Context, appInfo *dbModels.Application) error {
	a.ds.Lock()
	defer a.ds.Unlock()

	for _, app := range a.ds.apps {
		if app.ID == appInfo.ID || app.Name == appInfo.Name {
			return errors.New("application already exists")
		}
	}
	a.ds.apps = append(a.ds.apps, appInfo)
	return nil
}

func (a *application) UpdateApp(ctx context.Context, appInfo *dbModels.Application) error {
	a.ds.Lock()
	defer a.ds.Unlock()

	for i, app := range a.ds.apps {
		if app.ID == appInfo.ID {
			a.ds.apps[i] = appInfo
			return nil
		}
	}
	return errors.New("application not found")
}
