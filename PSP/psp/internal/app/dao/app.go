package dao

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type appDaoImpl struct{}

func NewAppDao() AppDao {
	return &appDaoImpl{}
}

func (d *appDaoImpl) GetApp(ctx context.Context, app *model.App) (bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	exist, err := session.Get(app)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (d *appDaoImpl) GetAppTotalNum(ctx context.Context) (int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	session.Where("state = ?", "published")
	count, err := session.Count(&model.App{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *appDaoImpl) ListApp(ctx context.Context, ids []snowflake.ID, app *model.App) ([]*model.App, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var apps []*model.App
	if len(ids) > 0 {
		session.In("id", ids)
	}
	if app.ComputeType != "" {
		session.Where("compute_type = ?", app.ComputeType)
	}
	if app.State != "" {
		session.Where("state = ?", app.State)
	}
	if app.LicenseManagerId != "" {
		session.Where("license_manager_id = ?", app.LicenseManagerId)
	}

	session.Desc("compute_type").Asc("name", "version_num")

	err := session.Find(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (d *appDaoImpl) AddApps(ctx context.Context, apps []*model.App) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Insert(apps)
	if err != nil {
		return err
	}
	return nil
}

func (d *appDaoImpl) UpdateApp(ctx context.Context, app *model.App) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.ID(app.ID).Omit("id").AllCols().Update(app)
	if err != nil {
		return err
	}
	return nil
}

func (d *appDaoImpl) UpdateAppsState(ctx context.Context, names []string, computeType, state string) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	app := &model.App{State: state}
	if len(names) > 0 {
		session.In("name", names)
	}
	if computeType != "" {
		session.Where("compute_type = ?", computeType)
	}

	_, err := session.Cols("state").Update(app)
	if err != nil {
		return err
	}
	return nil
}

func (d *appDaoImpl) DeleteApp(ctx context.Context, app *model.App) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Delete(app)
	if err != nil {
		return err
	}
	return nil
}
