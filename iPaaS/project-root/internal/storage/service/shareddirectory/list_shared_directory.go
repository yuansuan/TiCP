package shareddirectory

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// ListSharedDirectoryInfoByUserID 获取共享目录信息
func ListSharedDirectoryInfoByUserID(ctx context.Context, engine *xorm.Engine, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, userID string) ([]*model.SharedDirectory, error) {
	session := engine.Context(ctx)
	defer session.Close()

	return storageSharedDirectoryDao.ListByUserID(session, userID)
}

// ListSharedDirectoryInfoByPathPrefix 获取共享目录信息
func ListSharedDirectoryInfoByPathPrefix(ctx context.Context, engine *xorm.Engine, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, pathPrefix string) ([]*model.SharedDirectory, error) {
	session := engine.Context(ctx)
	defer session.Close()

	return storageSharedDirectoryDao.ListByPathPrefix(session, pathPrefix)
}
