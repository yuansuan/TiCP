package job

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	jl "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"xorm.io/xorm"
)

// RegxpList 过滤正则列表
// 存储正则参数只能过滤匹配正则的文件，不能过滤不匹配正则的文件
// 这里先过滤掉无'_'和无'.'的文件，再自己过滤筛出图片
var RegxpList = []string{`^[^_]*$`, `^[^.]*$`}

// ListJobSnapshot 获取作业云图集
func (srv *jobService) ListJobSnapshot(ctx context.Context, appSrv application.AppSrv, req *jl.Request, userID snowflake.ID, allow allowFunc) (map[string][]string, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job list snapshot start")
	defer logger.Info("job list snapshot end")

	jobID := snowflake.MustParseString(req.JobID)

	// JobID存在性验证
	job, err := srv.jobdao.Get(ctx, jobID, false, false)
	if err != nil {
		if !errors.Is(err, common.ErrJobIDNotFound) {
			logger.Warnf("get Job error! err: %v", err)
		}
		return nil, err // internal error OR job not exist
	}

	// 用户权限验证
	if !allow(userID.String(), job.UserID.String()) {
		logger.Info("no permission to operate other's job")
		return nil, errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
	}

	// 获取app信息
	app, err := appSrv.GetApp(ctx, job.AppID)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			return nil, errors.WithMessagef(common.ErrAppIDNotFound, "job app not found, appID: '%s'", job.AppID)
		}
		logger.Warnf("get app error: %v", err)
		return nil, err
	}
	// 判断app是否开启云图功能
	if !app.SnapshotEnable {
		logger.Info("app not enable snapshot")
		return nil, errors.WithMessage(common.ErrJobSnapshotNotFound, "app not enable snapshot")
	}

	// 获取数据
	zones, path, storageZone, storageType, state, err := GetSnapshotNeededJobInfo(ctx, job)
	if err != nil {
		return nil, err
	}
	if state.IsInitiated() || state.IsInitiallySuspended() || state.IsPending() {
		logger.Info("job not started")
		return nil, errors.WithMessage(common.ErrJobSnapshotNotFound, "job not started")
	}

	if path == "" {
		logger.Info("snapshot path is empty")
		return nil, errors.WithMessage(common.ErrJobSnapshotNotFound, "snapshot path is empty")
	}

	zone, ok := zones[storageZone]
	if !ok {
		logger.Warnf("residual: zone not found: %v", job.Zone)
		return nil, errors.New("zone not found")
	}

	clientParams := storage.ClientParams{}
	// 获取存储client
	if storageType == consts.HpcStorage {
		clientParams.Endpoint = zone.HPCEndpoint
		if clientParams.Endpoint == "" {
			logger.Warnf("hpc domain is empty")
			return nil, errors.New("hpc domain is empty")
		}
		clientParams.AdminAPI = true // hpc存储使用adminLs
	} else {
		clientParams.Endpoint = zone.StorageEndpoint
		if clientParams.Endpoint == "" {
			logger.Warnf("cloud storage domain is empty")
			return nil, errors.New("cloud storage domain is empty")
		}
		clientParams.UserID = util.ParseYsID(path) // 云存储非adminLs需要assumerole
	}

	path = strings.TrimPrefix(path, clientParams.Endpoint)
	path = util.AddPrefixSlash(path)
	path = util.AddSuffixSlash(path)

	lsParams := storage.LsParams{
		Offset:    0, // offset 从0开始
		Lspath:    path,
		RegxpList: RegxpList,
	}

	// 获取文件列表
	ctx = logging.AppendWith(ctx, "func", "job.ProcesAllImageFiles", "endpoint", clientParams.Endpoint, "path", path)
	return ProcesAllImageFiles(ctx, clientParams, lsParams)
}

// ProcesAllImageFiles 处理所有图片文件
func ProcesAllImageFiles(ctx context.Context, clientParams storage.ClientParams, lsParams storage.LsParams) (map[string][]string, error) {
	logger := logging.GetLogger(ctx)
	imageNameMap := make(map[string][]string)
	for {
		if lsParams.Offset == -1 { // -1表示已经最后一页
			break
		}

		resp, err := storage.Client().LsWithPage(clientParams, lsParams)
		if err != nil {
			if resp.ErrorCode == api.InvalidPath {
				logger.Infof("invalid path: %v", err)
				return nil, errors.WithMessage(common.ErrInvalidPath, err.Error())
			}
			if resp.ErrorCode == api.PathNotFound {
				logger.Infof("path not found: %v", err)
				return nil, errors.WithMessage(common.ErrPathNotFound, err.Error())
			}
			logger.Warnf("ls error: %v", err)
			return nil, errors.Wrap(err, "ls error")
		}

		files := resp.Data.Files
		ProcessImageFiles(ctx, files, imageNameMap)
		lsParams.Offset = resp.Data.NextMarker
	}
	return imageNameMap, nil
}

// ProcessImageFiles 处理图片文件
func ProcessImageFiles(ctx context.Context, files []*schema.FileInfo, imageNameMap map[string][]string) {
	for _, file := range files {
		isImage, err := util.IsImageFile(file.IsDir, file.Size, file.Name)
		if err != nil {
			logging.GetLogger(ctx).Infof("is image file error: %v", err)
			continue
		}
		if !isImage {
			continue
		}

		name := util.ExtractNameFromFileName(file.Name)
		if name != "" {
			regName := name + "_*.*"
			UpdateImageNameMap(imageNameMap, regName, file.Name)
		}
	}
}

// UpdateImageNameMap 更新imageNameMap
func UpdateImageNameMap(imageNameMap map[string][]string, regName, fileName string) {
	if _, ok := imageNameMap[regName]; !ok {
		imageNameMap[regName] = []string{}
	}
	imageNameMap[regName] = append(imageNameMap[regName], fileName)
}
