package compressInfo

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func UpdateCompressInfo(ctx context.Context, engine *xorm.Engine, compressInfoDao dao.CompressInfoDao, model *model.CompressInfo) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := compressInfoDao.Update(session, model)
	if err != nil {
		return err
	}
	return nil
}
