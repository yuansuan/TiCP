package shareddirectory

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// GetSharedDirectoryInfoByPath 获取共享目录信息
func GetSharedDirectoryInfoByPath(ctx context.Context, engine *xorm.Engine, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, path string) (bool, *model.SharedDirectory, error) {
	session := engine.Context(ctx)
	defer session.Close()

	return storageSharedDirectoryDao.GetByPath(session, path)
}
