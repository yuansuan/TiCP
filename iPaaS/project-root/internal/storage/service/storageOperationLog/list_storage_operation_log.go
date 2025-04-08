package storageOperationLog

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func ListStorageOperationLog(ctx context.Context, engine *xorm.Engine, storageOperationLogDao dao.StorageOperationLogDao, param *dao.StorageOperationLogQueryParam) ([]*model.StorageOperationLog, error, int64, int64) {
	session := engine.Context(ctx)
	defer session.Close()

	res, err, next, total := storageOperationLogDao.List(session, param)
	if err != nil {
		return nil, err, 0, 0
	}
	return res, nil, next, total
}
