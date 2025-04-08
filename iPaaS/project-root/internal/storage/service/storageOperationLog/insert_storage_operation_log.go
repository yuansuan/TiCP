package storageOperationLog

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func InsertStorageOperationLog(ctx context.Context, engine *xorm.Engine, storageOperationLogDao dao.StorageOperationLogDao, model *model.StorageOperationLog) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageOperationLogDao.Insert(session, model)
	if err != nil {
		return err
	}
	return nil
}
