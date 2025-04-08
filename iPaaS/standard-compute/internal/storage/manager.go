package storage

import (
	"context"
	"fmt"
	"path/filepath"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage/client"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

// Manager 存储管理器
type Manager struct {
	cfg *config.Config
	cli *client.Client
}

type inputStorageExt struct {
	*v20230530.JobInHPCInputStorage
	SrcSize int64
}

func (m *Manager) Download(ctx context.Context, jobID int64, inputs []*v20230530.JobInHPCInputStorage, workspace string) error {
	var err error
	if !filepath.IsAbs(workspace) {
		err = fmt.Errorf("workspace [%s] is not absolute", workspace)
		log.Error(err)
		return err
	}

	// get all file stat, and update to db
	downloaderMap := make(map[*inputStorageExt]Downloader)
	var totalSize int64
	for _, input := range inputs {
		if err = m.fulfillDownloadMap(input, downloaderMap, &totalSize); err != nil {
			err = fmt.Errorf("fulfill download map failed, %w", err)
			log.Error(err)
			return err
		}
	}

	if err = dao.Default.UpdateJobFileProgress(ctx, jobID, dao.UpdateJobFileProgressArgs{
		DownloadTotalSize: &totalSize,
	}); err != nil {
		err = fmt.Errorf("update job file progress failed, %w", err)
		log.Error(err)
		return err
	}

	var currentSize int64
	for input, downloader := range downloaderMap {
		log.Infof("Starting download: jobID=%v, src=%v, dest=%v", jobID, input.Src, filepath.Join(workspace, input.Dst))
		if err = downloader.Download(ctx, input.Src, filepath.Join(workspace, input.Dst)); err != nil {
			err = fmt.Errorf("download file failed, %w", err)
			log.Error(err)
			return err
		}
		currentSize += input.SrcSize
		log.Infof("Download success: jobID=%v, src=%v, dest=%v", jobID, input.Src, filepath.Join(workspace, input.Dst))
		if err = dao.Default.UpdateJobFileProgress(ctx, jobID, dao.UpdateJobFileProgressArgs{
			DownloadCurrentSize: &currentSize,
		}); err != nil {
			err = fmt.Errorf("update job file progress failed, %w", err)
			log.Error(err)
			return err
		}
	}

	return nil
}

// example
func (m *Manager) fulfillDownloadMap(input *v20230530.JobInHPCInputStorage, downloadMap map[*inputStorageExt]Downloader, totalSize *int64) error {
	srcEndpoint, path, err := util.ParseRawStorageUrl(input.Src)
	if err != nil {
		return fmt.Errorf("parse raw storage url [%s] failed, %w", input.Src, err)
	}

	fileInfo, err := m.cli.Stat(srcEndpoint, path)
	if err != nil {
		return fmt.Errorf("call storage client file stat failed, endpoint: %s, path: %s, err : %w", srcEndpoint, path, err)
	}

	if fileInfo.IsDir() {
		subFileInfos, err := m.cli.Ls(srcEndpoint, path)
		if err != nil {
			return fmt.Errorf("call storage client ls failed, endpoint: %s, path: %s, err: %w", srcEndpoint, path, err)
		}

		for _, subFileInfo := range subFileInfos {
			if err = m.fulfillDownloadMap(&v20230530.JobInHPCInputStorage{
				Type: input.Type,
				// start with shema:// cannot use filepath.Join
				Src: fmt.Sprintf("%s/%s", input.Src, subFileInfo.Name()),
				Dst: filepath.Join(input.Dst, subFileInfo.Name()),
			}, downloadMap, totalSize); err != nil {
				return err
			}
		}
	} else {
		*totalSize += fileInfo.Size()

		downloader, err := m.getDownloader(input.Type)
		if err != nil {
			return fmt.Errorf("get downloader failed, %w", err)
		}

		downloadMap[&inputStorageExt{
			JobInHPCInputStorage: input,
			SrcSize:              fileInfo.Size(),
		}] = downloader
	}

	return nil
}

func (m *Manager) getDownloader(storageType v20230530.StorageType) (Downloader, error) {
	switch storageType {
	case v20230530.HPCStorageType:
		return &localDownloader{cli: m.cli}, nil
	case v20230530.CloudStorageType:
		return &remoteDownloader{cli: m.cli}, nil
	default:
		err := fmt.Errorf("unsupported downloader")
		log.Error(err)
		return nil, err
	}
}

// NewManager 创建一个存储管理器
func NewManager(cfg *config.Config) (*Manager, error) {
	return &Manager{
		cfg: cfg,
		cli: client.New(cfg),
	}, nil
}
