package dao

import (
	"context"
	"errors"

	"github.com/yuansuan/ticp/common/go-kit/example/sqlitedb/internal/dao/models"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
)

func Get(ctx context.Context, key string) (*models.KeyMap, error) {
	sess := boot.MW.DefaultSession(ctx)
	defer func() { _ = sess.Close() }()

	m := new(models.KeyMap)
	exists, err := sess.Where("key = ?", key).Get(m)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New("keyMap not found")
	}

	return m, nil
}

func Set(ctx context.Context, key, val string) error {
	sess := boot.MW.DefaultSession(ctx)
	defer func() { _ = sess.Close() }()

	exists, err := sess.Where("key = ?", key).Exist(new(models.KeyMap))
	if err != nil {
		return err
	} else if !exists {
		_, err = sess.Insert(&models.KeyMap{Key: key, Value: val})
	} else {
		_, err = sess.Where("key = ?", key).Update(&models.KeyMap{Value: val})
	}

	return err
}
