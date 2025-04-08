package compressInfo

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func InsertCompressInfo(ctx context.Context, engine *xorm.Engine, compressInfoDao dao.CompressInfoDao, model *model.CompressInfo) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := compressInfoDao.Insert(session, model)
	if err != nil {
		return err
	}
	return nil
}
