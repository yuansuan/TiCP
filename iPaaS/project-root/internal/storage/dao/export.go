package dao

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"xorm.io/xorm"
)

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_upload_info_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao UploadInfoDao
//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_storage_quota_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao StorageQuotaDao
//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_storage_operation_log_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao StorageOperationLogDao
//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_storage_shared_directory_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao StorageSharedDirectoryDao
//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_compress_info_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao CompressInfoDao
//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao -destination mock_directory_usage_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao DirectoryUsageDao

var (
	UploadInfo             = NewUploadInfoDaoImpl()
	StorageQuota           = NewStorageQuotaDaoImpl()
	StorageOperationLog    = NewStorageOperationLogDaoImpl()
	StorageSharedDirectory = NewStorageSharedDirectoryDaoImpl()
	CompressInfo           = NewCompressInfoDaoImpl()
	DirectoryUsage         = NewDirectoryUsageDaoImpl()
)

// UploadInfoDao upload info interface
type UploadInfoDao interface {
	Insert(session *xorm.Session, model *model.UploadInfo) error
	Get(session *xorm.Session, model *model.UploadInfo) (bool, *model.UploadInfo, error)
}

type StorageQuotaDao interface {
	Insert(session *xorm.Session, model *model.StorageQuota) error
	Get(session *xorm.Session, model *model.StorageQuota) (bool, *model.StorageQuota, error)
	Update(session *xorm.Session, userID snowflake.ID, model *model.StorageQuota) error
	Total(session *xorm.Session) (float64, error)
	List(session *xorm.Session, pageOffset, pageSize int) (res []*model.StorageQuota, err error, next, total int64)
}

type StorageOperationLogQueryParam struct {
	// 用户 ID
	UserIDs string
	// 文件名
	FileName string
	// 文件类型, 可选值: file-普通文件, folder-文件夹, batch-批量操作
	FileTypes string
	// 操作类型, 可选值: upload-上传, download-下载, delete-删除, move-移动, mkdir-添加文件夹, copy-拷贝, copy_range-指定范围拷贝,compress-压缩, create-创建, link-链接, read_at-读, write_at-写
	OperationTypes string
	// 开始时间
	BeginTime time.Time
	// 结束时间
	EndTime time.Time
	//分页偏移量
	PageOffset int64
	//分页大小
	PageSize int64
}

type StorageOperationLogDao interface {
	Insert(session *xorm.Session, model *model.StorageOperationLog) error
	List(session *xorm.Session, param *StorageOperationLogQueryParam) (res []*model.StorageOperationLog, err error, next, total int64)
	GetUserIDs(session *xorm.Session) (res []string, err error)
	DeleteExpiredLog(session *xorm.Session, retentionPeriod int64) error
}

// StorageSharedDirectoryDao 共享目录数据
type StorageSharedDirectoryDao interface {
	Insert(session *xorm.Session, model *model.SharedDirectory) error
	GetByPath(session *xorm.Session, path string) (bool, *model.SharedDirectory, error)
	ListByUserID(session *xorm.Session, userID string) ([]*model.SharedDirectory, error)
	ListByPathPrefix(session *xorm.Session, pathPrefix string) ([]*model.SharedDirectory, error)
	DeleteByPath(session *xorm.Session, path string) error
}

// CompressInfoDao 压缩信息数据
type CompressInfoDao interface {
	Insert(session *xorm.Session, model *model.CompressInfo) error
	Update(session *xorm.Session, model *model.CompressInfo) error
	ListUnfinishedTask(session *xorm.Session) ([]*model.CompressInfo, error)
	Get(session *xorm.Session, id string) (bool, *model.CompressInfo, error)
}

// DirectoryUsageDao 目录使用情况数据
type DirectoryUsageDao interface {
	Insert(session *xorm.Session, model *model.DirectoryUsage) error
	Update(session *xorm.Session, model *model.DirectoryUsage) error
	ListCalculatingTask(session *xorm.Session) ([]*model.DirectoryUsage, error)
	Get(session *xorm.Session, id string) (bool, *model.DirectoryUsage, error)
}
