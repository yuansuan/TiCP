package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	storageOperationLogService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageOperationLog"
	"xorm.io/xorm"
)

const (
	DefaultRetentionPeriod = 30
)

type OperationLog struct {
	OperationLogCleanInterval int64
	MaxUserOperationLogCount  int64
	RetentionPeriod           int64
	Engine                    *xorm.Engine
	StorageOperationLogDao    dao.StorageOperationLogDao

	pathchecker.PathAccessCheckerImpl
}

func NewOperationLog(storageOperationLogDao dao.StorageOperationLogDao, engine *xorm.Engine, pathchecker pathchecker.PathAccessCheckerImpl) *OperationLog {
	if storageOperationLogDao == nil {
		return nil
	}
	storageOperationLogInterval := config.GetConfig().OperationLog.OperationLogCleanInterval
	retentionPeriod := config.GetConfig().OperationLog.RetentionPeriod
	if retentionPeriod <= 0 {
		retentionPeriod = DefaultRetentionPeriod
	}

	return &OperationLog{
		OperationLogCleanInterval: storageOperationLogInterval,
		RetentionPeriod:           retentionPeriod,
		Engine:                    engine,
		StorageOperationLogDao:    storageOperationLogDao,
		PathAccessCheckerImpl:     pathchecker,
	}
}

func (o *OperationLog) CleanExpiredOperationLog(name string, interval time.Duration) {
	ctx := logging.AppendWith(context.Background(), "daemon", name, "retentionPeriod", o.RetentionPeriod)
	logger := logging.GetLogger(ctx)

	for range time.Tick(interval) {
		expiredDate := time.Now().AddDate(0, 0, int(-o.RetentionPeriod))
		logger.Infof("start cleaning expired operation log, expired date:%v ...", expiredDate)
		err := storageOperationLogService.DeleteExpiredStorageOperationLog(ctx, o.Engine, o.StorageOperationLogDao, o.RetentionPeriod)
		if err != nil {
			logger.Errorf("failed to delete expired operation log, err: %v", err)
			continue
		}
		logger.Infof("finish cleaning expired operation log ...")
	}
}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (o *OperationLog) GetUserIDAndAKAndHandleError(ctx *gin.Context) (string, string, error) {
	userID, accessKey, _, err := o.PathAccessCheckerImpl.GetUserIDAndAKAndHandleError(ctx, pathchecker.AdminURLPrefix)
	return userID, accessKey, err
}

func (o *OperationLog) InsertOperationLog(logger *logging.Logger, ctx *gin.Context, operationLog *model.StorageOperationLog) {
	err := storageOperationLogService.InsertStorageOperationLog(ctx, o.Engine, o.StorageOperationLogDao, operationLog)
	if err != nil {
		logger.Warnf("insert storage operation log failed, err: %v", err)
	}
	return
}

func (o *OperationLog) FormatBytes(bits int64) string {
	const (
		KB = 1 << (10 * (iota + 1))
		MB
		GB
		TB
	)

	var value float64
	var unit string

	switch {
	case bits >= TB:
		value = float64(bits) / TB
		unit = "TB"
	case bits >= GB:
		value = float64(bits) / GB
		unit = "GB"
	case bits >= MB:
		value = float64(bits) / MB
		unit = "MB"
	case bits >= KB:
		value = float64(bits) / KB
		unit = "KB"
	default:
		value = float64(bits)
		unit = "B"
	}

	return fmt.Sprintf("%.2f%s", value, unit)
}
