package dao

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// CompressInfoDaoImpl ...
type CompressInfoDaoImpl struct{}

// NewCompressInfoDaoImpl ...
func NewCompressInfoDaoImpl() *CompressInfoDaoImpl {
	return &CompressInfoDaoImpl{}
}

// Insert ...
func (dao *CompressInfoDaoImpl) Insert(session *xorm.Session, model *model.CompressInfo) error {
	_, err := session.Insert(model)
	if err != nil {
		return err
	}
	return nil
}

// Update ...
func (dao *CompressInfoDaoImpl) Update(session *xorm.Session, model *model.CompressInfo) error {
	_, err := session.ID(model.Id).Update(model)
	if err != nil {
		return err
	}
	return nil
}

// ListUnfinishedTask ...
func (dao *CompressInfoDaoImpl) ListUnfinishedTask(session *xorm.Session) ([]*model.CompressInfo, error) {
	models := make([]*model.CompressInfo, 0)
	//查询所有未完成的压缩任务
	err := session.Where("status = ?", model.CompressTaskRunning).Find(&models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

// Get ...
func (dao *CompressInfoDaoImpl) Get(session *xorm.Session, id string) (bool, *model.CompressInfo, error) {
	compressInfo := &model.CompressInfo{}
	has, err := session.ID(id).Get(compressInfo)
	if err != nil {
		return false, nil, err
	}
	return has, compressInfo, nil
}
