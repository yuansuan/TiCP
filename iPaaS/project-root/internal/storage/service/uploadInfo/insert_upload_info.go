package uploadInfo

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func InsertUploadInfo(ctx context.Context, engine *xorm.Engine, uploadInfoDao dao.UploadInfoDao, model *model.UploadInfo) error {
	session := engine.Context(ctx)
	defer session.Close()

	err := uploadInfoDao.Insert(session, model)
	if err != nil {
		return err
	}
	return nil
}
