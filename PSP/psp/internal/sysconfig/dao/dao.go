package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type SysConfigDao interface {
	Get(ctx context.Context, key string) (*model.SysConfig, bool, error)
	Set(ctx context.Context, id snowflake.ID, key, value string) error
}

type AlertManagerConfigDao interface {
	Set(ctx context.Context, in []*model.AlertNotification) error
	Get(ctx context.Context, configType string) ([]*model.AlertNotification, error)
}
