package uploadInfo

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

func GetUploadInfo(ctx context.Context, engine *xorm.Engine, uploadInfoDao dao.UploadInfoDao, model *model.UploadInfo) (bool, *model.UploadInfo, error) {
	session := engine.Context(ctx)
	defer session.Close()

	exist, model, err := uploadInfoDao.Get(session, model)
	if err != nil {
		return exist, model, err
	}
	return exist, model, nil

}
