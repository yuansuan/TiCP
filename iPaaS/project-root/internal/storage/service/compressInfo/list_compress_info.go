package compressInfo

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func ListCompressInfo(ctx context.Context, engine *xorm.Engine, compressInfoDao dao.CompressInfoDao) ([]*model.CompressInfo, error) {
	session := engine.Context(ctx)
	defer session.Close()

	data, err := compressInfoDao.ListUnfinishedTask(session)
	if err != nil {
		return nil, err
	}
	return data, nil
}
