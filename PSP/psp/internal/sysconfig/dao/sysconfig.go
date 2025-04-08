package dao

import (
	"context"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type sysConfigImpl struct{}

func NewAppDao() SysConfigDao {
	return &sysConfigImpl{}
}

func (d *sysConfigImpl) Get(ctx context.Context, key string) (*model.SysConfig, bool, error) {
	session := boot.MW.DefaultSession(ctx)

	sysConfig := &model.SysConfig{}
	exist, err := session.Where("`key` = ?", key).Get(sysConfig)
	if err != nil {
		return nil, false, err
	}

	return sysConfig, exist, err
}

func (d *sysConfigImpl) Set(ctx context.Context, id snowflake.ID, key, value string) error {
	session := boot.MW.DefaultSession(ctx)

	sysConfig := &model.SysConfig{}
	_, err := session.Where("`key` = ?", key).Delete(sysConfig)
	if err != nil {
		return err
	}

	_, err = session.Insert(&model.SysConfig{
		Id:    id,
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}

	return nil
}
