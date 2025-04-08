package dao

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// StorageSharedDirectoryDaoImpl 共享目录数据实现
type StorageSharedDirectoryDaoImpl struct{}

// NewStorageSharedDirectoryDaoImpl 共享目录数据
func NewStorageSharedDirectoryDaoImpl() *StorageSharedDirectoryDaoImpl {
	return &StorageSharedDirectoryDaoImpl{}
}

// GetByPath 获取
func (s *StorageSharedDirectoryDaoImpl) GetByPath(session *xorm.Session, path string) (bool, *model.SharedDirectory, error) {
	mp := &model.SharedDirectory{}
	has, err := session.Where("path = ?", path).And("is_deleted = 0").Get(mp)
	if err != nil {
		return false, nil, errors.Wrap(err, "dao:")
	}
	return has, mp, nil
}

// ListByUserID 根据用户ID获取
func (s *StorageSharedDirectoryDaoImpl) ListByUserID(session *xorm.Session, userID string) ([]*model.SharedDirectory, error) {
	mps := make([]*model.SharedDirectory, 0)
	err := session.Where("user_id = ?", userID).And("is_deleted = 0").Find(&mps)
	if err != nil {
		return nil, errors.Wrap(err, "dao:")
	}
	return mps, nil
}

// ListByPathPrefix 根据用户ID获取
func (s *StorageSharedDirectoryDaoImpl) ListByPathPrefix(session *xorm.Session, pathPrefix string) ([]*model.SharedDirectory, error) {
	mps := make([]*model.SharedDirectory, 0)
	// where path like 'pathPrefix%'
	err := session.Where("path like ?", pathPrefix+"%").And("is_deleted = 0").Find(&mps)
	if err != nil {
		return nil, errors.Wrap(err, "dao:")
	}
	return mps, nil
}

// Insert 插入
func (s *StorageSharedDirectoryDaoImpl) Insert(session *xorm.Session, mp *model.SharedDirectory) error {
	_, err := session.Insert(mp)
	if err != nil {
		return errors.Wrap(err, "dao:")
	}
	return nil
}

// DeleteByPath 删除
func (s *StorageSharedDirectoryDaoImpl) DeleteByPath(session *xorm.Session, path string) error {
	deleteRow, err := session.Where("path = ?", path).Update(&model.SharedDirectory{IsDeleted: 1})
	if err != nil {
		return errors.Wrap(err, "delete application error")
	}
	if deleteRow == 0 {
		return xorm.ErrNotExist
	}
	return nil
}
