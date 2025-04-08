package directoryUsage

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func UpdateDirectoryUsage(ctx context.Context, engine *xorm.Engine, directoryUsageDao dao.DirectoryUsageDao, model *model.DirectoryUsage) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := directoryUsageDao.Update(session, model)
	if err != nil {
		return err
	}
	return nil
}
