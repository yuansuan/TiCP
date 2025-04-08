package storageQuota

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func ListStorageQuotaInfo(ctx context.Context, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao, pageOffset, pageSize int) ([]*model.StorageQuota, error, int64, int64) {
	session := engine.Context(ctx)
	defer session.Close()

	res, err, next, total := storageQuotaDao.List(session, pageOffset, pageSize)
	if err != nil {
		return nil, err, next, total
	}
	return res, nil, next, total

}
