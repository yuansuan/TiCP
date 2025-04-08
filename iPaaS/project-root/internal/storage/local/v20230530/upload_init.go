package v20230530

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	uploadInit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	uploadInfoService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/uploadInfo"
)

var (
	uploadCache = cache.New(10*time.Hour, 24*time.Hour)
)

const (
	maxUploadSize         = 1024 * 1024 * 1024 * 1024 // 1TB
	maxSliceSize          = 1024 * 1024 * 20          // 20MB
	uploadCachePrefix     = "upload"
	uploadCacheExpireTime = 12 * time.Hour
)

type UploadCacheContent struct {
	UploadID  string   `json:"upload_id"`
	Path      string   `json:"path"`
	Size      int64    `json:"size"`
	Overwrite bool     `json:"overwrite"`
	Fd        *os.File `json:"fd"`
}

// UploadInit 用于初始化一个文件上传，返回一个uploadID，后续的上传slice都需要带上这个uploadID。
func (s *Storage) UploadInit(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "UploadInit", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &uploadInit.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	if request.Size < 0 {
		msg := fmt.Sprintf("file size cannot be negative, got: %d", request.Size)
		logger.Error(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidSize, msg)
		return

	}
	if request.Size > maxUploadSize {
		msg := fmt.Sprintf("file size should less than 1TB, got: %d", request.Size)
		logger.Error(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.SizeTooLarge, msg)
		return
	}

	flag, pathUserID, msg := fsutil.ValidateUserIDPath(request.Path)
	if !flag {
		logger.Error(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path and check storage usage
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) ||
		!s.Quota.CheckStorageUsageAndHandleError(pathUserID, float64(request.Size), logger, ctx) {
		return
	}

	// generate absolute path
	filePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	// check if the file name is too long
	if err := fsutil.ValidateFileNameLength(request.Path); err != nil {
		msg := err.Error()
		logger.Infof(msg)
		if errors.Is(err, fsutil.ErrFileNameTooLong) {
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidFileName, msg)
			return
		} else if errors.Is(err, fsutil.ErrPathTooLong) {
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
	}

	var fileInfo os.FileInfo
	if flag, fileInfo, err = s.HandlePathContainsFileError(ctx, logger, filePath, fmt.Sprintf("path contains file, path: %s", request.Path)); !flag {
		return
	}
	//check if file exists
	if !request.Overwrite {
		if err != nil && !os.IsNotExist(err) {
			msg := fmt.Sprintf("get file info failed, path: %s, err: %v", request.Path, err)
			logger.Error(msg)
			common.InternalServerError(ctx, "get file info failed")
			return
		}
		if fileInfo != nil {
			msg := fmt.Sprintf("file already exists, path: %s", request.Path)
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathExists, msg)
			return
		}
	}

	//create and truncate tmp file
	uploadID := uuid.New().String()
	uploadTmpFilePath := filepath.Join(s.rootPath, fsutil.TmpBucketUploadingJoin(uploadID))

	file, err := fsutil.CreateAndReturnFile(uploadTmpFilePath, request.Size, logger)
	if err != nil {
		msg := fmt.Sprintf("create and truncate tmp file failed, file path: %s, err: %v", uploadTmpFilePath, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "create and truncate tmp file failed")
		return
	}

	uploadCacheKey := fmt.Sprintf("%s_%s", uploadCachePrefix, uploadID)
	content, exist := uploadCache.Get(uploadCacheKey)
	if exist {
		if uploadCacheStruct, ok := content.(*UploadCacheContent); ok {
			logger.Infof("uploadCache key:%s exist, content:%+v", uploadCacheKey, uploadCacheStruct)
		} else {
			logger.Infof("uploadCache key:%s exist, get content is nil", uploadCacheKey)
		}
		common.ErrorResp(ctx, http.StatusConflict, commoncode.UploadTaskExists, "upload task exists")
		return
	}

	uploadCacheContent := &UploadCacheContent{UploadID: uploadID, Path: filePath, Size: request.Size, Overwrite: request.Overwrite, Fd: file}
	uploadCache.Set(uploadCacheKey, uploadCacheContent, uploadCacheExpireTime)
	uploadCache.OnEvicted(func(key string, value interface{}) {
		if uploadContent, ok := value.(*UploadCacheContent); ok {
			err := uploadContent.Fd.Close()
			if err != nil {
				logger.Errorf("close upload file failed, err: %v", err)
			}
		}
	})

	//insert upload info
	uploadInfo := &model.UploadInfo{
		Id:         uploadID,
		UserId:     userID,
		Path:       filePath,
		TmpPath:    uploadTmpFilePath,
		Size:       request.Size,
		Overwrite:  request.Overwrite,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = uploadInfoService.InsertUploadInfo(ctx, s.Engine, s.UploadInfoDao, uploadInfo)
	if err != nil {
		msg := fmt.Sprintf("insert upload info failed, uploadInfo: %+v, err: %v", uploadInfo, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "insert upload info failed")
		return
	}

	common.SuccessResp(ctx, &uploadInit.Data{
		UploadID: uploadID,
	})
}
