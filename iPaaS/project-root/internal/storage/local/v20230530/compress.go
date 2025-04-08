package v20230530

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/cancel"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/start"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/status"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/compress_task"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"golang.org/x/sync/semaphore"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	userTaskNum            = 20
	compressingStatus      = "compressing"
	compressSuccessStatus  = "success"
	compressFailedStatus   = "failed"
	compressCanceledStatus = "canceled"
	compressUnknownStatus  = "unknown"
)

func (s *Storage) CompressStart(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "CompressStart", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &start.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if len(request.Paths) == 0 {
		msg := "paths is empty"
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	for _, path := range request.Paths {
		flag, _, msg := fsutil.ValidateUserIDPath(path)
		if !flag {
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
		//check if user has access to the path
		if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, path, logger, ctx) {
			return
		}
	}

	flag, pathUserID, msg := fsutil.ValidateUserIDPath(request.TargetPath)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidTargetPath, msg)
		return
	}

	basePath := filepath.Join(s.rootPath, userID)
	if request.BasePath != "" {
		flag, _, msg := fsutil.ValidateUserIDPath(request.BasePath)
		if !flag {
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidBasePath, msg)
			return
		}
		basePath = filepath.Join(s.rootPath, fsutil.TrimPrefix(request.BasePath, "/"))
	}

	//check if user has access to the path
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.TargetPath, logger, ctx) {
		return
	}

	//check if target path exists
	targetFilePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.TargetPath, "/"))

	if filepath.Ext(targetFilePath) != ".zip" {
		msg := fmt.Sprintf("unsupported compress file type, path: %s", filepath.Ext(request.TargetPath))
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.UnsupportedCompressFileType, msg)
		return
	}

	var fileInfo os.FileInfo
	if flag, fileInfo, err = s.HandlePathContainsFileError(ctx, logger, targetFilePath, fmt.Sprintf("target path contains file, path: %s", request.TargetPath)); !flag {
		return
	}
	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("Lstat error, path: %s, err: %v", request.TargetPath, err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "Lstat error")
		return
	}
	if fileInfo != nil {
		msg := fmt.Sprintf("target file already exists,path: %v", request.TargetPath)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.TargetPathExists, msg)
		return
	}

	if err := os.MkdirAll(filepath.Dir(targetFilePath), filemode.Directory); err != nil {
		msg := fmt.Sprintf("create target dir error, path: %s, err: %v", request.TargetPath, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "create target dir error")
		return
	}
	pathsMap := make(map[string]struct{})
	emptyDirs := make([]string, 0)
	if s.CompressTask.GetFilePaths(ctx, request.Paths, pathsMap, emptyDirs, logger, true) {
		return
	}

	var paths []string
	for path := range pathsMap {
		paths = append(paths, path)
	}

	// check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, s.GetTotalSize(logger, paths...), logger, ctx) {
		return
	}

	sem, _ := compress_task.Semaphores.LoadOrStore(userID, semaphore.NewWeighted(int64(userTaskNum)))
	if !sem.(*semaphore.Weighted).TryAcquire(1) {
		msg := fmt.Sprintf("too many compress tasks, userID: %s,limit: %d", userID, userTaskNum)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusTooManyRequests, commoncode.TooManyCompressTask, msg)
		return
	}

	fileName := filepath.Base(targetFilePath)
	tmpCompressFilePath := filepath.Join(s.rootPath, fsutil.TmpFileCompressJoin(fileName))
	targetAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.TargetPath, "/"))
	compressCacheKey := fmt.Sprintf("%s_%s", compress_task.ZipCachePrefix, userID)
	logger.Infof("compressCacheKey:%s, tmpCompressFilePath:%s, targetAbsPath:%s", compressCacheKey, tmpCompressFilePath, targetAbsPath)

	compressID := uuid.New().String()
	compressInfo := &model.CompressInfo{
		Id:         compressID,
		UserId:     userID,
		TmpPath:    fsutil.TmpFileCompressJoin(fileName),
		Paths:      strings.Join(request.Paths, ","),
		TargetPath: request.TargetPath,
		BasePath:   request.BasePath,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err = s.CompressTask.InsertCompressInfo(logger, ctx, compressInfo); err != nil {
		common.InternalServerError(ctx, fmt.Sprintf("insert compress info to database error:%v,compressID:%v", err, compressInfo.Id))
		return
	}

	zipCache := s.CompressTask.GetOrCreateZipCache(userID, compressCacheKey)
	zipCache.ContentMap.Store(compressID, &compress_task.ZipCacheContent{IsFinished: false, TargetPath: targetAbsPath})
	compress_task.CompressCache.Set(compressCacheKey, zipCache, compress_task.ZipCacheExpireTime)

	compressParams := &compress_task.CompressParams{
		TargetPath:          targetAbsPath,
		TmpCompressFilePath: tmpCompressFilePath,
		BasePath:            basePath,
		Paths:               paths,
		EmptyDirs:           emptyDirs,
		CompressCacheKey:    compressCacheKey,
		Logger:              logger,
		ZipCache:            zipCache,
		CompressID:          compressID,
		Sem:                 sem.(*semaphore.Weighted),
		CompressInfo:        compressInfo,
	}
	go s.CompressTask.StartCompressTask(context.Background(), compressParams)

	resp := new(start.Data)
	resp.FileName = fileName
	resp.CompressID = compressID

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      fileName,
		SrcPath:       strings.Join(request.Paths, ","),
		DestPath:      request.TargetPath,
		FileType:      commoncode.Batch,
		OperationType: commoncode.COMPRESS,
		Size:          "",
		CreateTime:    time.Now(),
	})
	common.SuccessResp(ctx, resp)
}

func (s *Storage) CompressStatus(ctx *gin.Context) {

	userID, _, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "CompressStatus", "RequestId ", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &status.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if len(strings.TrimSpace(request.CompressID)) == 0 {
		msg := fmt.Sprintf("CompressID is empty")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidCompressID, msg)
		return
	}

	compressCacheKey := fmt.Sprintf("%s_%s", compress_task.ZipCachePrefix, userID)
	logger.Infof("CompressStatus, userID:%s,compressCacheKey:%s", userID, compressCacheKey)

	var cacheContent *compress_task.ZipCacheContent

	getCompressInfo := func() (bool, *model.CompressInfo, error) {
		exist, compressInfo, err := s.CompressTask.GetCompressInfo(logger, ctx, request.CompressID)
		if err != nil {
			msg := fmt.Sprintf("getCompressInfo err, compressID:%s", request.CompressID)
			logger.Error(msg)
			common.InternalServerError(ctx, msg)
			return false, nil, err
		}
		if !exist {
			msg := fmt.Sprintf("can not find compress task, compressID:%s", request.CompressID)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.CompressTaskNotFound, msg)
			return false, nil, nil
		}
		return true, compressInfo, nil
	}

	var compressInfo *model.CompressInfo
	content, exist := compress_task.CompressCache.Get(compressCacheKey)
	if !exist {
		if exist, compressInfo, err = getCompressInfo(); err != nil || !exist {
			return
		}

		cacheContent = &compress_task.ZipCacheContent{
			IsFinished: compressInfo.Status != model.CompressTaskRunning,
			TargetPath: compressInfo.TargetPath,
			Status:     int(compressInfo.Status),
			Err:        compressInfo.ErrorMsg,
		}
	} else {
		switch content := content.(type) {
		case *compress_task.ZipCache:
			value, ok := content.ContentMap.Load(request.CompressID)
			if !ok {
				if exist, compressInfo, err = getCompressInfo(); err != nil || !exist {
					return
				}

				cacheContent = &compress_task.ZipCacheContent{
					IsFinished: compressInfo.Status != model.CompressTaskRunning,
					TargetPath: compressInfo.TargetPath,
					Status:     int(compressInfo.Status),
					Err:        compressInfo.ErrorMsg,
				}
			} else {
				var ok bool
				cacheContent, ok = value.(*compress_task.ZipCacheContent)
				if !ok || cacheContent == nil {
					msg := fmt.Sprintf("compressCache type err, compressCacheKey:%s", compressCacheKey)
					logger.Error(msg)
					common.InternalServerError(ctx, "compressCache type err")
					return
				}
			}
		default:
			msg := fmt.Sprintf("compressCache type err, compressCacheKey:%s", compressCacheKey)
			logger.Error(msg)
			common.InternalServerError(ctx, "compressCache type err")
			return
		}

		if cacheContent == nil {
			msg := fmt.Sprintf("can not find compress task, compressCacheKey:%s", compressCacheKey)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.CompressTaskNotFound, fmt.Sprintf("can not find compress task, compressID:%s", request.CompressID))
			return
		}
	}

	logger.Infof("CompressStatus, compressCacheKey:%s, request compress status info:%+v", compressCacheKey, cacheContent)

	if !systemFlag && !strings.Contains(cacheContent.TargetPath, userID) {
		msg := fmt.Sprintf("user %s has no access to compress task, compressID:%s", userID, request.CompressID)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.CompressTaskNoAccess, msg)
		return
	}

	resp := new(status.Data)
	resp.IsFinished = cacheContent.IsFinished
	switch cacheContent.Status {
	case int(model.CompressTaskFailed):
		resp.Status = compressFailedStatus
		logger.Infof("CompressStatus, compressing err:%s", cacheContent.Err)
	case int(model.CompressTaskRunning):
		resp.Status = compressingStatus
	case int(model.CompressTaskFinished):
		resp.Status = compressSuccessStatus
	case int(model.CompressTaskCanceled):
		resp.Status = compressCanceledStatus
	default:
		resp.Status = compressUnknownStatus
	}
	common.SuccessResp(ctx, resp)

}

func (s *Storage) CompressCancel(ctx *gin.Context) {
	userID, _, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "CompressCancel", "RequestId ", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &cancel.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if len(strings.TrimSpace(request.CompressID)) == 0 {
		msg := fmt.Sprintf("CompressID is empty")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidCompressID, msg)
		return
	}
	compressCacheKey := fmt.Sprintf("%s_%s", compress_task.ZipCachePrefix, userID)
	zipCache := s.CompressTask.GetOrCreateZipCache(userID, compressCacheKey)
	exist, compressInfo, err := s.CompressTask.GetCompressInfo(logger, ctx, request.CompressID)
	if err != nil {
		msg := fmt.Sprintf("getCompressInfo err, compressID:%s", request.CompressID)
		logger.Error(msg)
		common.InternalServerError(ctx, msg)
		return
	}
	if !exist {
		msg := fmt.Sprintf("can not find compress task, compressID:%s", request.CompressID)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.CompressTaskNotFound, msg)
		return
	}
	if !systemFlag && !strings.Contains(compressInfo.TargetPath, userID) {
		msg := fmt.Sprintf("user %s has no access to compress task, compressID:%s", userID, request.CompressID)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.CompressTaskNoAccess, msg)
		return
	}
	if err = s.CompressTask.CancelCompressTask(ctx, logger, request.CompressID, zipCache); err != nil {
		return
	}
	common.SuccessResp(ctx, nil)
}
