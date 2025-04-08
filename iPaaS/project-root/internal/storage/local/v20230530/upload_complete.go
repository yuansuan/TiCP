package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	uploadcomplete "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	uploadInfoService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/uploadInfo"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadComplete 用于完成一个文件的上传。
func (s *Storage) UploadComplete(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "UploadComplete", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &uploadcomplete.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if request.UploadID == "" {
		msg := fmt.Sprintf("uploadID required")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.UploadIDNotFound, msg)
		return
	}

	uploadCacheKey := fmt.Sprintf("%s_%s", uploadCachePrefix, request.UploadID)

	flag, _, msg := fsutil.ValidateUserIDPath(request.Path)
	if !flag {
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		uploadCache.Delete(uploadCacheKey)
		return
	}

	uploadFilePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	content, exist := uploadCache.Get(uploadCacheKey)
	var uploadCacheContent *UploadCacheContent
	var size int64
	if !exist {
		//check if upload id exists in db
		exist, uploadInfo, err := uploadInfoService.GetUploadInfo(ctx, s.Engine, s.UploadInfoDao, &model.UploadInfo{Id: request.UploadID})
		if err != nil {
			msg := fmt.Sprintf("get upload info failed, uploadID: %s, userID: %v,err: %v", request.UploadID, userID, err)
			logger.Error(msg)
			common.InternalServerError(ctx, msg)
			return
		}
		if !exist {
			msg := fmt.Sprintf("upload info not existed, uploadID: %s, userID: %v,err: %v", request.UploadID, userID, err)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.UploadIDNotFound, "uploadID not exist")
			return
		}
		size = uploadInfo.Size
		uploadCacheContent = &UploadCacheContent{
			UploadID:  uploadInfo.Id,
			Path:      uploadInfo.Path,
			Overwrite: uploadInfo.Overwrite,
			Size:      size,
		}
	} else {
		var ok bool
		if uploadCacheContent, ok = content.(*UploadCacheContent); !ok {
			msg := fmt.Sprintf("uploadCache type err, uploadCacheKey:%s", uploadCacheKey)
			logger.Error(msg)
			uploadCache.Delete(uploadCacheKey)
			common.InternalServerError(ctx, "uploadCache type err")
			return
		} else {
			size = uploadCacheContent.Size
		}
	}

	//check if uploadFilePath matches the path of the upload init request (uploadCacheContent.Path)
	if uploadFilePath != uploadCacheContent.Path {
		msg := fmt.Sprintf("path not matches with upload init request, path: %s, uploadInfo.Path: %s", request.Path, uploadCacheContent.Path)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathNotMatchUploadInit, msg)
		return
	}

	var fileInfo os.FileInfo
	if flag, fileInfo, err = s.HandlePathContainsFileError(ctx, logger, uploadFilePath, fmt.Sprintf("path contains file, path: %s", request.Path)); !flag {
		return
	}
	//check if file exists
	if !uploadCacheContent.Overwrite {
		if err != nil && !os.IsNotExist(err) {
			msg := fmt.Sprintf("get file info failed, path: %s, err: %v", request.Path, err)
			logger.Errorf(msg)
			uploadCache.Delete(uploadCacheKey)
			common.InternalServerError(ctx, msg)
			return
		}

		if fileInfo != nil {
			msg := fmt.Sprintf("file already exists, path: %s", request.Path)
			logger.Infof(msg)
			uploadCache.Delete(uploadCacheKey)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathExists, msg)
			return
		}
	}

	uploadTmpFilePath := filepath.Join(s.rootPath, fsutil.TmpBucketUploadingJoin(uploadCacheContent.UploadID))
	//move file from tmp path to uploadFilePath
	uploadCache.Delete(uploadCacheKey)
	if err := fsutil.Move(uploadTmpFilePath, uploadFilePath, logger); err != nil {
		msg := fmt.Sprintf("move upload file failed, file path: %v, userID: %v, err: %v", request.Path, userID, err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
		return
	}

	logger.Infof("upload file succeed, file path: %v, userID: %v", uploadFilePath, userID)
	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(uploadFilePath),
		DestPath:      request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.UPLOAD,
		Size:          s.OperationLog.FormatBytes(size),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
