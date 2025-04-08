package dao

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	dblib "github.com/yuansuan/ticp/iPaaS/project-root/pkg/db"
)

// GetRemoteAppByName 按会话id，remoteapp名称获取remote配置
func GetRemoteAppByName(ctx context.Context, userID, sessionID snowflake.ID, appName string) (*models.RemoteApp, bool, error) {
	var (
		remoteApp models.RemoteApp
		exists    bool
	)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		detail, exist, err := GetSessionDetailsBySessionID(ctx, userID, sessionID)
		if err != nil {
			return errors.Wrap(err, "remoteapp.GetRemoteAppByName")
		}
		if !exist {
			return errors.Wrap(err, "remoteapp.GetRemoteAppByName.GetSessionDetailsBySessionID not exist")
		}

		if exists, err = db.Where("software_id = ? AND name = ?", detail.SoftwareId, appName).Get(&remoteApp); err != nil {
			return errors.Wrap(err, "dao")
		}
		return nil
	})
	return &remoteApp, exists, err
}

// ListRemoteAppBySoftwareID list remote app by software id
func ListRemoteAppBySoftwareID(ctx context.Context, softwareID snowflake.ID) ([]*models.RemoteApp, error) {
	remoteApps := make([]*models.RemoteApp, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		err := db.Where("software_id = ?", softwareID).Find(&remoteApps)
		if err != nil {
			return errors.Wrap(err, "dao")
		}

		return nil
	})
	return remoteApps, err
}

// AddRemoteApp add remote app
func AddRemoteApp(ctx context.Context, data *models.RemoteApp) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(data)
		if dblib.IsDuplicatedError(err) {
			return fmt.Errorf("duplicated error, %w, msg: %s", dblib.ErrDuplicatedEntry, err.Error())
		}

		return errors.Wrap(err, "remoteapp.AddRemoteApp")
	})
	return err
}

// UpdateRemoteApp update remote app
func UpdateRemoteApp(ctx context.Context, data *models.RemoteApp, mustCol ...string) (bool, error) {
	var count int64
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		count, e = db.ID(data.Id).Cols(mustCol...).Update(data)
		return errors.Wrap(e, "remoteapp.UpdateRemoteApp")
	})
	return count > 0, err
}

// DeleteRemoteApp delete remote app
func DeleteRemoteApp(ctx context.Context, id snowflake.ID) (bool, error) {
	var count int64
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		count, e = db.ID(id).Delete(new(models.RemoteApp))
		return errors.Wrap(e, "remoteapp.DeleteRemoteApp")
	})
	return count > 0, err
}

func BatchInsertRemoteAppUserPass(ctx context.Context, remoteAppsUserPass []*models.RemoteAppUserPass) error {
	if len(remoteAppsUserPass) == 0 {
		return nil
	}

	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.InsertMulti(&remoteAppsUserPass)
		return err
	})
}

func GetRemoteAppUserPass(ctx context.Context, sessionId snowflake.ID, remoteAppName string) (*models.RemoteAppUserPass, bool, error) {
	remoteAppUserPass := new(models.RemoteAppUserPass)
	var exist bool
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Where("session_id = ?", sessionId).
			Where("remote_app_name = ?", remoteAppName).
			Get(remoteAppUserPass)
		return e
	})
	if err != nil {
		return nil, false, fmt.Errorf("dao, %w", err)
	}
	return remoteAppUserPass, exist, nil
}
