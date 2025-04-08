package shareddirectory

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"xorm.io/xorm"
)

// DeleteSharedDirectoryInfo 删除共享目录信息
func DeleteSharedDirectoryInfo(ctx context.Context, engine *xorm.Engine, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, path string) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := storageSharedDirectoryDao.DeleteByPath(session, path)
	if err != nil {
		return err
	}

	return nil
}
