package dao

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

// StorageQuotaDaoImpl ...
type StorageQuotaDaoImpl struct{}

// NewStorageQuotaDaoImpl ...
func NewStorageQuotaDaoImpl() *StorageQuotaDaoImpl {
	return &StorageQuotaDaoImpl{}
}

// Insert ...
func (dao *StorageQuotaDaoImpl) Insert(session *xorm.Session, model *model.StorageQuota) error {
	_, err := session.Insert(model)
	if err != nil {
		return errors.Wrap(err, "dao:")
	}
	return nil
}

// Get ...
func (dao *StorageQuotaDaoImpl) Get(session *xorm.Session, model *model.StorageQuota) (bool, *model.StorageQuota, error) {
	exist, err := session.Where("user_id = ?", model.UserId).Get(model)
	if err != nil {
		return exist, model, errors.Wrap(err, "dao:")
	}
	return exist, model, nil
}

// Update ...
func (dao *StorageQuotaDaoImpl) Update(session *xorm.Session, userID snowflake.ID, model *model.StorageQuota) error {
	_, err := session.Where("user_id = ?", userID).Update(model)
	if err != nil {
		return errors.Wrap(err, "dao:")
	}
	return nil
}

// Total ...
func (dao *StorageQuotaDaoImpl) Total(session *xorm.Session) (float64, error) {
	var totalUsage float64
	sql := "SELECT SUM(storage_usage) FROM storage_quota"
	_, err := session.SQL(sql).Get(&totalUsage)
	if err != nil {
		return 0, err
	}
	return totalUsage, nil

}

// List ...
func (dao *StorageQuotaDaoImpl) List(session *xorm.Session, pageOffset, pageSize int) (res []*model.StorageQuota, err error, next, total int64) {
	res = make([]*model.StorageQuota, 0)
	total, err = session.Table("storage_quota").Limit(pageSize, pageOffset).FindAndCount(&res)
	if err != nil {
		return nil, err, -1, -1
	}

	next = int64(pageOffset + pageSize)
	if next > total-1 {
		next = -1
	}
	return res, nil, next, total
}
