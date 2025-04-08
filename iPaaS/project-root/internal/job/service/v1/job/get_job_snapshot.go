package job

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	jg "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"xorm.io/xorm"
)

// GetJobSnapshot 获取作业云图数据
func (srv *jobService) GetJobSnapshot(ctx context.Context, appSrv application.AppSrv, req *jg.Request, userID snowflake.ID, allow allowFunc) (string, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job get snapshot start")
	defer logger.Info("job get snapshot end")

	jobID := snowflake.MustParseString(req.JobID)
	fileName := req.Path

	// JobID存在性验证
	job, err := srv.jobdao.Get(ctx, jobID, false, false)
	if err != nil {
		if !errors.Is(err, common.ErrJobIDNotFound) {
			logger.Warnf("get Job error! err: %v", err)
		}
		return "", err // internal error OR job not exist
	}

	// 用户权限验证
	if !allow(userID.String(), job.UserID.String()) {
		logger.Warnf("no permission to operate other's job")
		return "", errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
	}

	// 获取app信息
	app, err := appSrv.GetApp(ctx, job.AppID)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return "", errors.WithMessagef(common.ErrAppIDNotFound, "job app not found, appID: '%s'", job.AppID)
		}
		logger.Warnf("get app error: %v", err)
		return "", err
	}
	// 判断app是否开启云图功能
	if !app.SnapshotEnable {
		logger.Warnf("app not enable snapshot")
		return "", errors.WithMessage(common.ErrJobSnapshotNotFound, "app not enable snapshot")
	}

	// 获取数据
	zones, path, storageZone, storageType, state, err := GetSnapshotNeededJobInfo(ctx, job)
	if err != nil {
		return "", err
	}
	if state.IsInitiated() || state.IsInitiallySuspended() || state.IsPending() {
		logger.Warnf("job not started")
		return "", errors.WithMessage(common.ErrJobSnapshotNotFound, "job not started")
	}

	if path == "" {
		logger.Info("snapshot path is empty")
		return "", errors.WithMessage(common.ErrJobSnapshotNotFound, "snapshot path is empty")
	}

	zone, ok := zones[storageZone]
	if !ok {
		logger.Warnf("residual: zone not found: %v", job.Zone)
		return "", errors.New("zone not found")
	}

	clientParams := storage.ClientParams{}
	// 获取存储client
	if storageType == consts.HpcStorage {
		clientParams.Endpoint = zone.HPCEndpoint
		if clientParams.Endpoint == "" {
			logger.Warnf("hpc domain is empty")
			return "", errors.New("hpc domain is empty")
		}
		clientParams.AdminAPI = true // hpc存储使用adminAPI
	} else {
		clientParams.Endpoint = zone.StorageEndpoint
		if clientParams.Endpoint == "" {
			logger.Warnf("cloud storage domain is empty")
			return "", errors.New("cloud storage domain is empty")
		}
		clientParams.UserID = util.ParseYsID(path) // 云存储非adminLs需要assumerole
	}

	path = strings.TrimPrefix(path, clientParams.Endpoint)
	path = util.AddPrefixSlash(path)
	path = util.AddSuffixSlash(path)
	path = filepath.Join(path, fileName)

	ctx = logging.AppendWith(ctx, "func", "job.ReadSnapshotFile", "endpoint", clientParams.Endpoint, "path", path)
	image, err := ReadSnapshotFile(ctx, clientParams, path)
	if err != nil {
		return "", err
	}

	// 格式统一转换为png
	reader := bytes.NewReader(image.Data)
	pngImage, err := util.ConvertImageTo(reader, fileName, util.PNG)
	if err != nil {
		logger.Warnf("convert image error: %v", err)
		return "", err
	}

	// 转换为base64
	base64Data := base64.StdEncoding.EncodeToString(pngImage)
	data := "data:image/png;base64," + base64Data

	return data, nil
}

func ReadSnapshotFile(ctx context.Context, clientParams storage.ClientParams, path string) (*readAt.Response, error) {
	logger := logging.GetLogger(ctx)
	fileInfo, err := storage.Client().Stat(clientParams, path)
	if err != nil {
		if fileInfo.ErrorCode == api.InvalidPath {
			logger.Infof("invalid path: %v", err)
			return nil, errors.WithMessage(common.ErrInvalidPath, err.Error())
		}
		if fileInfo.ErrorCode == api.PathNotFound {
			logger.Infof("file not found: %v", err)
			return nil, errors.WithMessage(common.ErrPathNotFound, err.Error())
		}
		logger.Warnf("stat error: %v", err)
		return nil, err
	}

	// size 大于 50M 不返回
	maxSize := int64(50 * 1024 * 1024)
	if fileInfo.Data.File.Size > maxSize {
		logger.Warnf("file size too large: %v", fileInfo.Data.File.Size)
		return nil, errors.New("file size too large")
	}

	readAtParams := storage.ReadAtParams{
		Readpath: path,
		Length:   fileInfo.Data.File.Size,
		Offset:   0,
		Resolver: nil,
	}

	image, err := storage.Client().ReadAt(clientParams, readAtParams)
	if err != nil {
		if fileInfo.ErrorCode == api.PathNotFound {
			logger.Infof("file not found: %v", err)
			return nil, errors.WithMessage(common.ErrPathNotFound, err.Error())
		}
		logger.Warnf("readAt error: %v", err)
		return nil, err
	}
	return image, nil
}

// GetSnapshotNeededJobInfo return zones, path, state, storageZone, storageType
func GetSnapshotNeededJobInfo(ctx context.Context, job *models.Job) (schema.Zones, string, string, consts.FileType, consts.State, error) {
	zones := config.GetConfig().Zones
	path := job.OutputDir
	state := consts.NewState(job.State, job.SubState)
	storageZone := job.FileOutputStorageZone
	storageType := consts.FileType(job.OutputType)

	params := &models.AdminParams{}
	err := json.Unmarshal([]byte(job.Params), params)
	if err != nil {
		logging.GetLogger(ctx).Warnf("unmarshal job params error: %v", err)
		return zones, path, storageZone, storageType, state, fmt.Errorf("unmarshal job params error: %v", err)
	}

	noTrans := params.Input.Type == consts.HpcStorage.String() && params.Input.Destination == ""

	if !state.IsFinal() || path == "" || noTrans { // 作业运行中直接从workdir获取文件
		path = job.WorkDir
		storageZone = job.Zone
		storageType = consts.HpcStorage
	}

	return zones, path, storageZone, storageType, state, nil
}
