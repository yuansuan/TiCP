package api

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	storageQuotaService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageQuota"
	"xorm.io/xorm"
)

const (
	DefaultStorageUsageUpdateInterval          = 300
	DefaultDefaultUserStorageLimit             = 500
	DefaultMaxUserStorageLimit                 = 10 * 1024
	DefaultMaxSystemStorageLimit               = 1024 * 1024
	DefaultSystemStorageLimitWarningBufferSize = 100
)

type Quota struct {
	StorageUsageUpdateInterval          int64
	DefaultUserStorageLimit             int64
	MaxUserStorageLimit                 int64
	MaxSystemStorageLimit               int64
	SystemStorageLimitWarningBufferSize int64
	WarningFlag                         bool
	Engine                              *xorm.Engine
	StorageQuotaDao                     dao.StorageQuotaDao
	mu                                  sync.Mutex

	pathchecker.PathAccessCheckerImpl
}

func NewQuota(StorageQuotaDao dao.StorageQuotaDao, engine *xorm.Engine, pathchecker pathchecker.PathAccessCheckerImpl) *Quota {
	if StorageQuotaDao == nil {
		return nil
	}
	storageUsageUpdateInterval := config.GetConfig().Quota.StorageUsageUpdateInterval
	if storageUsageUpdateInterval <= 0 {
		storageUsageUpdateInterval = DefaultStorageUsageUpdateInterval
	}
	defaultUserStorageLimit := config.GetConfig().Quota.DefaultUserStorageLimit
	if defaultUserStorageLimit <= 0 {
		defaultUserStorageLimit = DefaultDefaultUserStorageLimit
	}
	maxUserStorageLimit := config.GetConfig().Quota.MaxUserStorageLimit
	if maxUserStorageLimit <= 0 {
		maxUserStorageLimit = DefaultMaxUserStorageLimit
	}

	maxSystemStorageLimit := config.GetConfig().Quota.MaxSystemStorageLimit
	if maxSystemStorageLimit <= 0 {
		maxSystemStorageLimit = DefaultMaxSystemStorageLimit
	}
	systemStorageLimitWarningBufferSize := config.GetConfig().Quota.SystemStorageLimitWarningBufferSize
	if systemStorageLimitWarningBufferSize <= 0 {
		systemStorageLimitWarningBufferSize = DefaultSystemStorageLimitWarningBufferSize
	}

	return &Quota{
		StorageUsageUpdateInterval:          storageUsageUpdateInterval,
		DefaultUserStorageLimit:             defaultUserStorageLimit,
		MaxUserStorageLimit:                 maxUserStorageLimit,
		MaxSystemStorageLimit:               maxSystemStorageLimit,
		SystemStorageLimitWarningBufferSize: systemStorageLimitWarningBufferSize,
		Engine:                              engine,
		StorageQuotaDao:                     StorageQuotaDao,
		mu:                                  sync.Mutex{},

		PathAccessCheckerImpl: pathchecker,
	}
}

func (q *Quota) CheckStorageUsageAndHandleError(userID string, uploadSize float64, logger *logging.Logger, ctx *gin.Context) bool {

	if q.WarningFlag {
		msg := fmt.Sprintf("system storage usage is over %d G, please contact administrator", q.MaxSystemStorageLimit-q.SystemStorageLimitWarningBufferSize)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.SystemQuotaExhausted, msg)
		return false
	}
	id, err := snowflake.ParseString(userID)
	if err != nil {
		logger.Infof("Error parsing userID: %v,got: %v", userID, err)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return false
	}
	exists, storageQuota, err := storageQuotaService.GetStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, &model.StorageQuota{UserId: id})
	if err != nil {
		logger.Errorf("Error getting storage quota: %v, id: %v", err, id)
		common.InternalServerError(ctx, "Error getting storage quota")
		return false
	}
	//first upload
	if !exists {
		if bytesToGB(uploadSize) > float64(q.DefaultUserStorageLimit) {
			msg := fmt.Sprintf("upload size %v G is over user storage limit %v G", uploadSize, q.DefaultUserStorageLimit)
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusForbidden, commoncode.QuotaExhausted, msg)
			return false

		}
		//insert storage quota
		storageQuota := &model.StorageQuota{
			UserId:       id,
			StorageUsage: bytesToGB(uploadSize),
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = storageQuotaService.InsertStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, storageQuota)
		if err != nil {
			logger.Errorf("Error inserting storage quota: %v", err)
			common.InternalServerError(ctx, "Error inserting storage quota")
			return false
		}
		return true
	}

	storageQuotaLimit := storageQuota.StorageLimit
	if storageQuotaLimit == 0 {
		storageQuotaLimit = float64(q.DefaultUserStorageLimit)
	}
	if bytesToGB(uploadSize)+storageQuota.StorageUsage > storageQuotaLimit {
		msg := fmt.Sprintf("upload size plus storage usage %v G is over user storage limit %v G", float64(math.Round((bytesToGB(uploadSize)+storageQuota.StorageUsage)*100))/100, storageQuotaLimit)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.QuotaExhausted, msg)
		return false
	}

	return true
}

func (q *Quota) CheckAndUpdateStorageUsage(name string, rootPath string, interval time.Duration) {

	ctx := logging.AppendWith(context.Background(), "daemon", name)
	logger := logging.GetLogger(ctx)

	for range time.Tick(interval) {
		logger.Infof("start check and update storage usage")
		var totalSystemSize = new(float64)

		files, err := os.ReadDir(rootPath)
		if err != nil {
			logger.Warnf("read root path error：%v", err)
			return
		}
		var wg sync.WaitGroup
		for _, file := range files {
			if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
				absPath := filepath.Join(rootPath, file.Name())
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					q.updateUserAndSystemStorageUsage(absPath, logger, q.DefaultUserStorageLimit, totalSystemSize)
				}(absPath)
			}
		}
		wg.Wait()

		if bytesToGB(*totalSystemSize) > float64(q.MaxSystemStorageLimit-q.SystemStorageLimitWarningBufferSize) {
			logger.Warnf("system storage usage exceeded threshold. Total size: %v GB, threshold size: %v GB ", bytesToGB(*totalSystemSize), q.MaxSystemStorageLimit-q.SystemStorageLimitWarningBufferSize)
			// 触发警报逻辑
			q.WarningFlag = true
			logger.Errorf("system storage usage is over %d, please contact administrator!", q.MaxSystemStorageLimit-q.SystemStorageLimitWarningBufferSize)
		} else {
			logger.Infof("system storage usage is within the threshold. Total size: %v GB, threshold size: %v GB ", bytesToGB(*totalSystemSize), q.MaxSystemStorageLimit-q.SystemStorageLimitWarningBufferSize)
			q.WarningFlag = false
		}
		logger.Infof("end check and update storage usage")
	}

}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (q *Quota) GetUserIDAndAKAndHandleError(ctx *gin.Context) (string, string, error) {
	userID, accessKey, _, err := q.PathAccessCheckerImpl.GetUserIDAndAKAndHandleError(ctx, pathchecker.AdminURLPrefix)
	return userID, accessKey, err
}

func (q *Quota) updateUserAndSystemStorageUsage(userFolderPath string, logger *logging.Logger, defaultUserStorageLimit int64, systemTotalSize *float64) {
	totalSize := int64(0)
	inodes := make(map[uint64]struct{})
	ctx := context.Background()
	logger.Infof("start update user storage usage for user: %v", userFolderPath)
	err := filepath.Walk(userFolderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warnf("Error accessing path: %v", err)
			return err
		}

		// 获取文件的 inode 号
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("unable to get inode for file %s", info.Name())
		}

		// 硬链接只统计一次
		inode := stat.Ino
		if _, found := inodes[inode]; !found {
			inodes[inode] = struct{}{}
			totalSize += info.Size()
		}

		return nil
	})
	if err != nil {
		logger.Warnf("Error walking path: %v", err)
		return
	}

	userID := filepath.Base(userFolderPath)
	//check if user storage quota exists
	id, err := snowflake.ParseString(userID)
	if err != nil || id <= 0 {
		logger.Warnf("invalid userID: err %v,userID %v", err, userID)
		return
	}
	exists, _, err := storageQuotaService.GetStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, &model.StorageQuota{UserId: id})
	if err != nil {
		logger.Warnf("Error getting storage quota: %v", err)
		return
	}
	if !exists {
		//insert storage quota
		storageQuota := &model.StorageQuota{
			UserId:       id,
			StorageUsage: bytesToGB(float64(totalSize)),
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = storageQuotaService.InsertStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, storageQuota)
		if err != nil {
			logger.Warnf("Error inserting storage quota: %v", err)
		}
	} else {
		//update storage quota
		storageQuota := &model.StorageQuota{
			UserId:       id,
			StorageUsage: bytesToGB(float64(totalSize)),
			UpdateTime:   time.Now(),
		}
		err = storageQuotaService.UpdateStorageQuotaInfo(ctx, q.Engine, q.StorageQuotaDao, id, storageQuota)
		if err != nil {
			logger.Warnf("Error updating storage quota: %v", err)
		}
	}
	logger.Infof("end update user storage usage for user: %v, totalSize: %v", userFolderPath, bytesToGB(float64(totalSize)))

	q.mu.Lock()
	defer q.mu.Unlock()
	*systemTotalSize += float64(totalSize)
}

func bytesToGB(bytes float64) float64 {
	gb := bytes / (1024 * 1024 * 1024)
	return float64(math.Round(gb*100)) / 100
}
