package storageQuota

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func UpdateStorageQuotaInfo(ctx context.Context, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao, userID snowflake.ID, model *model.StorageQuota) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageQuotaDao.Update(session, userID, model)
	if err != nil {
		return err
	}
	return nil
}
