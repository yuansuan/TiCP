package dao

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// DirectoryUsageDaoImpl ...
type DirectoryUsageDaoImpl struct{}

// NewDirectoryUsageDaoImpl ...
func NewDirectoryUsageDaoImpl() *DirectoryUsageDaoImpl {
	return &DirectoryUsageDaoImpl{}
}

func (dao *DirectoryUsageDaoImpl) Insert(session *xorm.Session, model *model.DirectoryUsage) error {
	_, err := session.Insert(model)
	if err != nil {
		return err
	}
	return nil
}

func (dao *DirectoryUsageDaoImpl) Update(session *xorm.Session, model *model.DirectoryUsage) error {
	_, err := session.ID(model.Id).Update(model)
	if err != nil {
		return err
	}
	return nil
}

func (dao *DirectoryUsageDaoImpl) Get(session *xorm.Session, id string) (bool, *model.DirectoryUsage, error) {
	directoryUsage := &model.DirectoryUsage{}
	exist, err := session.ID(id).Get(directoryUsage)
	if err != nil {
		return exist, directoryUsage, err
	}
	return exist, directoryUsage, nil
}

func (dao *DirectoryUsageDaoImpl) ListCalculatingTask(session *xorm.Session) ([]*model.DirectoryUsage, error) {
	models := make([]*model.DirectoryUsage, 0)
	//查询所有计算中的任务
	err := session.Where("status = ?", model.DirectoryUsageTaskCalculating).Find(&models)
	if err != nil {
		return nil, err
	}
	return models, nil
}
