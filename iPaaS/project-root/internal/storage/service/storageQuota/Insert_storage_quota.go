package storageQuota

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func InsertStorageQuotaInfo(ctx context.Context, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao, model *model.StorageQuota) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageQuotaDao.Insert(session, model)
	if err != nil {
		return err
	}
	return nil
}
