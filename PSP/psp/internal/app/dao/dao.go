package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type AppDao interface {
	GetApp(ctx context.Context, app *model.App) (bool, error)
	ListApp(ctx context.Context, ids []snowflake.ID, app *model.App) ([]*model.App, error)
	AddApps(ctx context.Context, apps []*model.App) error
	UpdateApp(ctx context.Context, app *model.App) error
	UpdateAppsState(ctx context.Context, names []string, computeType, state string) error
	DeleteApp(ctx context.Context, app *model.App) error
	GetAppTotalNum(ctx context.Context) (int64, error)
}
