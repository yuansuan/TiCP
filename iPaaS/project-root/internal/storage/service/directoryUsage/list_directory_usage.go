package directoryUsage

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func ListDirectoryUsage(ctx context.Context, engine *xorm.Engine, directoryUsageDao dao.DirectoryUsageDao) ([]*model.DirectoryUsage, error) {
	session := engine.Context(ctx)
	defer session.Close()

	data, err := directoryUsageDao.ListCalculatingTask(session)
	if err != nil {
		return nil, err
	}
	return data, nil
}
