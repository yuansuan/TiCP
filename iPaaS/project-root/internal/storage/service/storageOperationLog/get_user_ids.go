package storageOperationLog

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"xorm.io/xorm"
)

func GetUserIDs(ctx context.Context, engine *xorm.Engine, storageOperationLogDao dao.StorageOperationLogDao) ([]string, error) {
	session := engine.Context(ctx)
	defer session.Close()

	res, err := storageOperationLogDao.GetUserIDs(session)
	if err != nil {
		return nil, err
	}
	return res, nil
}
