package storageQuota

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"xorm.io/xorm"
)

func GetStorageQuotaTotal(ctx context.Context, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao) (float64, error) {
	session := engine.Context(ctx)
	defer session.Close()

	total, err := storageQuotaDao.Total(session)
	if err != nil {
		return total, err
	}
	return total, nil

}
