package v20230530

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/compress_task"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/bytespool"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	directoryUsage "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/directory_usage/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	operationlog "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/operationlog/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	quota "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/quota/api"
	sharedDirectory "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/shared_directory/api"
	"xorm.io/xorm"
)

const (
	bufferSize                           = 1024 * 1024 //1m
	defaultUploadingFileExpireDuration   = 24 * 10
	defaultUploadingFileCleanInterval    = 1
	defaultCompressingFileExpireDuration = 24 * 10
	defaultCompressingFileCleanInterval  = 1
	pageSize                             = 500
	userIDKey                            = "x-ys-user-id"
	userAppKeyInQuery                    = "AccessKeyId"
	SoftLink                             = "soft-link"
	HardLink                             = "hard-link"
	Indirect                             = "indirect"
	Direct                               = "direct"
)

// Storage 存储api
type Storage struct {
	rootPath        string
	pool            *bytespool.BytesPool
	linker          linker.Linker
	fileLocks       sync.Map
	UploadInfoDao   dao.UploadInfoDao
	Engine          *xorm.Engine
	Quota           *quota.Quota
	OperationLog    *operationlog.OperationLog
	CompressTask    *compress_task.CompressTask
	SharedDirectory *sharedDirectory.SharedDirectory
	DirectoryUsage  *directoryUsage.DirectoryUsage

	pathchecker.PathAccessCheckerImpl
}

// NewStorage 新建存储api
func NewStorage(rootPath string, engine *xorm.Engine, storageQuotaDao dao.StorageQuotaDao, storageOperationLogDao dao.StorageOperationLogDao, storageSharedDirectoryDao dao.StorageSharedDirectoryDao, compressInfoDao dao.CompressInfoDao, directoryUsageDao dao.DirectoryUsageDao) (*Storage, error) {
	var l linker.Linker
	switch config.GetConfig().Local.LinkType {
	case SoftLink:
		l = &linker.SoftLink{}
	case HardLink:
		l = &linker.HardLink{}
	default:
		return nil, errors.New("unsupported local storage type")
	}

	startCleanupDaemon(rootPath)

	iamClient := iam_client.NewClient(
		config.GetConfig().IamServerUrl,
		config.GetConfig().AccessKeyId,
		config.GetConfig().AccessKeySecret,
	)

	authEnabled := true
	if config.GetConfig().AuthEnabled != nil {
		authEnabled = *config.GetConfig().AuthEnabled
	}

	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       authEnabled,
		IamClient:         iamClient,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}

	storageQuota := quota.NewQuota(storageQuotaDao, engine, pathchecker)
	go storageQuota.CheckAndUpdateStorageUsage("checkAndUpdateStorageUsage", rootPath, time.Duration(storageQuota.StorageUsageUpdateInterval)*time.Second)

	operationLog := operationlog.NewOperationLog(storageOperationLogDao, engine, pathchecker)
	go operationLog.CleanExpiredOperationLog("cleanExpiredOperationLog", time.Duration(operationLog.OperationLogCleanInterval)*time.Second)

	sharedDirectory := sharedDirectory.NewSharedDirectory(storageSharedDirectoryDao, engine, sharedDirectory.GetHC(), pathchecker)

	directoryUsage := directoryUsage.NewDirectoryUsage(directoryUsageDao, engine, pathchecker, rootPath)
	go directoryUsage.Recover()

	compressTask := compress_task.NewCompressTask(compressInfoDao, engine, rootPath)
	go compressTask.Recover()

	return &Storage{
		rootPath:              rootPath,
		pool:                  bytespool.NewBytesPool(bufferSize),
		linker:                l,
		UploadInfoDao:         dao.UploadInfo,
		Engine:                engine,
		Quota:                 storageQuota,
		OperationLog:          operationLog,
		SharedDirectory:       sharedDirectory,
		CompressTask:          compressTask,
		DirectoryUsage:        directoryUsage,
		PathAccessCheckerImpl: pathchecker,
	}, nil
}

func cleanTmpFiles(rootPath string, folderName string, expire time.Duration, tick time.Duration) {
	tmpFileFolder := filepath.Join(rootPath, folderName)

	cleanExpireFiles(
		fmt.Sprintf("file-manager-%s-clean", folderName),
		tmpFileFolder,
		expire,
		tick,
	)
}

func timeoutFiles(rootPath, mode, tmpPath string, paths []string, expire, interval time.Duration) {
	ctx := logging.AppendWith(context.Background(), "daemon", "timeoutFiles")
	logger := logging.GetLogger(ctx)

	if mode != Direct && mode != Indirect {
		logger.Fatal("Currently, only \"direct\" and \"indirect\" modes are supported\n")
		return
	}

	if paths == nil || len(paths) == 0 {
		logger.Warnf("You should specify a cleanup directory!")
		return
	}

	if mode == Indirect && len(tmpPath) == 0 {
		logger.Fatal("tmpPath must be set in \"indirect\" mode\n")
		return
	}

	for range time.Tick(interval) {
		for _, path := range paths {

			cleanupPath := filepath.Join(rootPath, path)
			if mode == Direct {
				err := cleanTimeoutFiles(cleanupPath, expire, nil)
				if err != nil {
					logger.Errorf("remove %s failed. error:%s", path, err)
				}
			} else {
				tmpPath := filepath.Join(rootPath, tmpPath)
				err := moveTimeoutFiles(cleanupPath, tmpPath, expire)
				if err != nil {
					logger.Errorf("move %s failed. error:%s", path, err)
				}
			}
		}
	}
}

func cleanTimeoutFiles(cleanupPath string, expire time.Duration, modTimePointer *time.Time) error {
	now := time.Now()

	var postCleanup func(string) error
	postCleanup = func(filePath string) error {
		dir, err := os.ReadDir(filePath)

		if err != nil {
			return err
		}

		for _, file := range dir {
			currentPath := filepath.Join(filePath, file.Name())
			info, err := file.Info()
			if err != nil {
				return err
			}

			modTime := info.ModTime()
			if modTimePointer != nil {
				modTime = *modTimePointer
			}

			if file.IsDir() {
				err := postCleanup(currentPath)
				if err != nil {
					return err
				}
			}

			if file.IsDir() {
				cp, err := os.ReadDir(currentPath)
				if err != nil {
					return err
				}
				if len(cp) == 0 && now.Sub(modTime) >= expire {
					err := os.Remove(currentPath)
					if err != nil {
						return err
					}
				}
			} else {
				if now.Sub(modTime) >= expire {
					err := os.Remove(currentPath)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	err := postCleanup(cleanupPath)
	if err != nil {
		return err
	}

	return nil
}

func moveTimeoutFiles(src, dst string, expire time.Duration) error {
	now := time.Now()
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	modTime := stat.ModTime()
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == src {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if now.Sub(modTime) > expire {
			relativePath := strings.TrimPrefix(path, src+string(filepath.Separator))
			destinationPath := filepath.Join(dst, relativePath)

			if d.IsDir() {
				if err := os.MkdirAll(destinationPath, info.Mode()); err != nil {
					return err
				}
			} else {
				if _, err := os.Stat(destinationPath); err == nil || os.IsExist(err) {
					if err := os.Remove(destinationPath); err != nil {
						return err
					}
				}

				err = os.Rename(path, destinationPath)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	modTimePointer := &modTime
	err = cleanTimeoutFiles(src, expire, modTimePointer)
	if err != nil {
		return err
	}

	return nil
}

func startCleanupDaemon(rootPath string) {
	uploadingFileExpireDuration := config.GetConfig().TmpFileCleanup.UploadingFileExpireDuration
	if uploadingFileExpireDuration == 0 {
		uploadingFileExpireDuration = defaultUploadingFileExpireDuration
	}
	uploadingFileCleanInterval := config.GetConfig().TmpFileCleanup.UploadingFileCleanInterval
	if uploadingFileCleanInterval == 0 {
		uploadingFileCleanInterval = defaultUploadingFileCleanInterval
	}
	compressingFileExpireDuration := config.GetConfig().TmpFileCleanup.CompressingFileExpireDuration
	if compressingFileExpireDuration == 0 {
		compressingFileExpireDuration = defaultCompressingFileExpireDuration
	}
	compressingFileCleanInterval := config.GetConfig().TmpFileCleanup.CompressingFileCleanInterval
	if compressingFileCleanInterval == 0 {
		compressingFileCleanInterval = defaultCompressingFileCleanInterval
	}

	cleanupEnabled := config.GetConfig().TimeoutCleanup.Enabled
	if cleanupEnabled {
		cleanupMode := config.GetConfig().TimeoutCleanup.Mode
		paths := config.GetConfig().TimeoutCleanup.CleanupPaths
		duration := config.GetConfig().TimeoutCleanup.TimeoutDuration
		interval := config.GetConfig().TimeoutCleanup.CleanupInterval
		tmpPath := config.GetConfig().TimeoutCleanup.TmpPath
		go timeoutFiles(rootPath, cleanupMode, tmpPath, paths, time.Duration(duration)*time.Minute, time.Duration(interval)*time.Minute)
	}

	go cleanTmpFiles(rootPath, fsutil.TmpBucketUploading, time.Duration(uploadingFileExpireDuration)*time.Hour, time.Duration(uploadingFileCleanInterval)*time.Hour)
	go cleanTmpFiles(rootPath, fsutil.TmpFileCompress, time.Duration(compressingFileExpireDuration)*time.Hour, time.Duration(compressingFileCleanInterval)*time.Hour)
}

func cleanExpireFiles(name string, uploadTmpFileFolder string, expire time.Duration, interval time.Duration) {
	ctx := logging.AppendWith(context.Background(), "daemon", name)
	logger := logging.GetLogger(ctx)

	for range time.Tick(interval) {
		logger.Infof("start cleaning upload tmp files")
		if _, err := os.Stat(uploadTmpFileFolder); os.IsNotExist(err) {
			logger.Infof("uploadTmpFileFolder %v not exist", uploadTmpFileFolder)
			continue
		}
		infos, err := fsutil.Ls(uploadTmpFileFolder, pageSize, logger)
		if err != nil {
			logger.Warnf("error when ls uploadTmpFileFolder=%v, err=%v", uploadTmpFileFolder, err)
			continue
		}

		for _, info := range infos {
			if time.Now().Sub(info.ModTime()) < expire {
				continue
			}
			logger.Infof("file removing: %v %v %v", info.Name(), info.ModTime(), info.Size())
			if err := os.RemoveAll(filepath.Join(uploadTmpFileFolder, info.Name())); err != nil {
				logger.Warnf("remove file failed: %v", err)
			} else {
				logger.Infof("file remove success: %v", info.Name())
			}
		}
	}
}
