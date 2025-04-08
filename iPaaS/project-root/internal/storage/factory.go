package storage

import (
	"fmt"

	"github.com/marmotedu/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	directoryUsageApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/directory_usage/api"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	operationLogApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/operationlog/api"
	quota2 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/quota/api"
	sharedDirectoryApi "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/shared_directory/api"
)

const (
	Version20230530 = "2023-05-30"
)

var (
	storages        = make(map[string]Storage)
	loaded          = false
	quota           *quota2.Quota
	operationLog    *operationLogApi.OperationLog
	sharedDirectory *sharedDirectoryApi.SharedDirectory
	directoryUsage  *directoryUsageApi.DirectoryUsage
)

var keyFormatter = func(storageType, version string) string {
	return fmt.Sprintf("%s-%s", storageType, version)
}

func lazyInit() {
	key := keyFormatter("local", Version20230530)
	rootPath := config.GetConfig().Local.RootPath
	storage, err := v20230530.NewStorage(rootPath, boot.MW.DefaultORMEngine(), dao.StorageQuota, dao.StorageOperationLog, dao.StorageSharedDirectory, dao.CompressInfo, dao.DirectoryUsage)
	if err != nil {
		panic(err)
	}
	storages[key] = storage
	loaded = true
	quota = storage.Quota
	operationLog = storage.OperationLog
	sharedDirectory = storage.SharedDirectory
	directoryUsage = storage.DirectoryUsage
}

func New() {
	if !loaded {
		lazyInit()
	}
}

func GetStorage(version string) (Storage, *quota2.Quota, *operationLogApi.OperationLog, *sharedDirectoryApi.SharedDirectory, *directoryUsageApi.DirectoryUsage, error) {
	storageType := config.GetConfig().StorageType
	key := keyFormatter(storageType, version)
	if storages[key] == nil {
		msg := fmt.Sprintf("storage not found, version: %s", version)
		logging.Default().Infof(msg)
		return nil, nil, nil, nil, nil, errors.New(msg)
	}
	return storages[key], quota, operationLog, sharedDirectory, directoryUsage, nil
}
