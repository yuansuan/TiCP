package dao

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"

	"github.com/pkg/errors"

	"xorm.io/xorm"
)

// UploadInfoDaoImpl ...
type UploadInfoDaoImpl struct{}

// NewUploadInfoDaoImpl ...
func NewUploadInfoDaoImpl() *UploadInfoDaoImpl {
	return &UploadInfoDaoImpl{}
}

// Insert ...
func (dao *UploadInfoDaoImpl) Insert(session *xorm.Session, model *model.UploadInfo) error {
	_, err := session.Insert(model)
	if err != nil {
		return errors.Wrap(err, "dao:")
	}
	return nil
}

// Get ...
func (dao *UploadInfoDaoImpl) Get(session *xorm.Session, model *model.UploadInfo) (bool, *model.UploadInfo, error) {
	exist, err := session.ID(model.Id).Get(model)
	if err != nil {
		return exist, model, errors.Wrap(err, "dao:")
	}
	return exist, model, nil
}
