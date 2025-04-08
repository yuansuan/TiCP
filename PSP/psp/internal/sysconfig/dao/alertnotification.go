package dao

import (
	"context"

	"github.com/pkg/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dao/model"
)

type alertNotification struct{}

func NewAlertNotificationDao() AlertManagerConfigDao {
	return &alertNotification{}
}

func (d *alertNotification) Set(ctx context.Context, alertNotifications []*model.AlertNotification) error {
	if alertNotifications == nil || len(alertNotifications) == 0 {
		return errors.New("request parameter alertNotifications is empty")
	}

	_, err := boot.MW.DefaultTransaction(ctx, func(session *xorm.Session) (interface{}, error) {
		// 先删除
		var keys []string
		for _, notification := range alertNotifications {
			keys = append(keys, notification.Key)
		}
		alertConfig := &model.AlertNotification{}
		_, err := session.In("key", keys).Where("type = ?", alertNotifications[0].Type).Delete(alertConfig)
		if err != nil {
			return nil, err
		}

		// 再插入
		_, err = session.Insert(alertNotifications)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *alertNotification) Get(ctx context.Context, configType string) ([]*model.AlertNotification, error) {
	session := boot.MW.DefaultSession(ctx)

	var alertConfig []*model.AlertNotification
	err := session.Table(model.AlertNotificationTableName).Where("type = ?", configType).Find(&alertConfig)
	if err != nil {
		return nil, err
	}

	return alertConfig, nil

}
