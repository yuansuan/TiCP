package shareddirectory

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// InsertSharedDirectoryInfo 插入共享目录信息
func InsertSharedDirectoryInfo(ctx context.Context, engine *xorm.Engine, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, model *model.SharedDirectory) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageSharedDirectoryDao.Insert(session, model)
	if err != nil {
		return err
	}
	return nil
}
