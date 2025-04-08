package dao

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"strings"
	"time"
	"xorm.io/xorm"
)

// StorageOperationLogDaoImpl ...
type StorageOperationLogDaoImpl struct{}

const TableNameStorageOperationLog = "storage_operation_log"

// NewStorageOperationLogDaoImpl ...
func NewStorageOperationLogDaoImpl() *StorageOperationLogDaoImpl {
	return &StorageOperationLogDaoImpl{}
}

func (dao *StorageOperationLogDaoImpl) Insert(session *xorm.Session, model *model.StorageOperationLog) error {
	_, err := session.Insert(model)
	if err != nil {
		return errors.Wrap(err, "dao:")
	}
	return nil
}

func (dao *StorageOperationLogDaoImpl) List(session *xorm.Session, param *StorageOperationLogQueryParam) (res []*model.StorageOperationLog, err error, next, total int64) {
	res = make([]*model.StorageOperationLog, 0)

	total, err = buildQuery(session, param).Count()
	if err != nil {
		return nil, err, 0, 0
	}

	query := buildQuery(session, param)
	if param.PageSize != 0 {
		query = query.Limit(int(param.PageSize), int(param.PageOffset))
	}

	err = query.Find(&res)
	if err != nil {
		return nil, err, 0, 0
	}

	next = param.PageOffset + int64(len(res))
	if next >= total {
		next = -1
	}

	return res, nil, next, total
}

func buildQuery(session *xorm.Session, param *StorageOperationLogQueryParam) *xorm.Session {
	query := session.Table(TableNameStorageOperationLog).Asc("create_time").Where("is_deleted = ?", 0)

	if param.UserIDs != "" {
		userIDs := strings.Split(param.UserIDs, ",")
		query = query.In("user_id", userIDs)
	}
	if param.FileName != "" {
		query = query.Where("file_name LIKE ?", "%"+param.FileName+"%")
	}
	if param.FileTypes != "" {
		fileTypes := strings.Split(param.FileTypes, ",")
		query = query.In("file_type", fileTypes)
	}
	if param.OperationTypes != "" {
		operationTypes := strings.Split(param.OperationTypes, ",")
		query = query.In("operation_type", operationTypes)
	}
	if param.BeginTime.Unix() > 0 {
		query = query.Where("create_time >= ?", param.BeginTime)
	}
	if param.EndTime.Unix() > 0 {
		query = query.Where("create_time <= ?", param.EndTime)
	}

	return query
}

func (dao *StorageOperationLogDaoImpl) GetUserIDs(session *xorm.Session) (res []string, err error) {
	res = make([]string, 0)
	err = session.Table(TableNameStorageOperationLog).Distinct("user_id").Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *StorageOperationLogDaoImpl) DeleteExpiredLog(session *xorm.Session, retentionPeriod int64) error {
	period := time.Now().AddDate(0, 0, int(-retentionPeriod))
	sql := "delete from `storage_operation_log` where create_time <= ?"
	_, err := session.Exec(sql, period)
	return err
}
