package storageOperationLog

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"xorm.io/xorm"
)

func DeleteExpiredStorageOperationLog(ctx context.Context, engine *xorm.Engine, storageOperationLogDao dao.StorageOperationLogDao, retentionPeriod int64) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageOperationLogDao.DeleteExpiredLog(session, retentionPeriod)
	return err
}
