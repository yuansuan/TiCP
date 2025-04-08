package storageQuota

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func GetStorageQuotaInfo(ctx context.Context, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao, model *model.StorageQuota) (bool, *model.StorageQuota, error) {
	session := engine.Context(ctx)
	defer session.Close()

	exist, model, err := storageQuotaDao.Get(session, model)
	if err != nil {
		return exist, model, err
	}
	return exist, model, nil

}
