package directoryUsage

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func GetDirectory(ctx context.Context, engine *xorm.Engine, directoryDao dao.DirectoryUsageDao, id string) (bool, *model.DirectoryUsage, error) {
	session := engine.Context(ctx)
	defer session.Close()

	exist, data, err := directoryDao.Get(session, id)
	if err != nil {
		return exist, nil, err
	}
	return exist, data, nil
}
