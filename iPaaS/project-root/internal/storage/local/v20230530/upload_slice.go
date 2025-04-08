package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	uploadSlice "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil/compress"
	uploadInfoService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/uploadInfo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadSlice 用于上传一个文件的一个分片。
func (s *Storage) UploadSlice(ctx *gin.Context) {

	userID, _, _, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "UploadSlice", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &uploadSlice.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	//get upload info
	if request.UploadID == "" {
		msg := fmt.Sprintf("uploadID required")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.UploadIDNotFound, msg)
		return
	}

	uploadCacheKey := fmt.Sprintf("%s_%s", uploadCachePrefix, request.UploadID)
	var uploadCacheContent *UploadCacheContent
	content, exist := uploadCache.Get(uploadCacheKey)
	if !exist {
		//check if upload id exists in db
		exist, uploadInfo, err := uploadInfoService.GetUploadInfo(ctx, s.Engine, s.UploadInfoDao, &model.UploadInfo{Id: request.UploadID})
		if err != nil {
			msg := fmt.Sprintf("get upload info failed, uploadID: %s, userID: %v,err: %v", request.UploadID, userID, err)
			logger.Error(msg)
			common.InternalServerError(ctx, "get upload info failed")
			return
		}
		if !exist {
			msg := fmt.Sprintf("upload info not existed, uploadID: %s, userID: %v,err: %v", request.UploadID, userID, err)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.UploadIDNotFound, "uploadID not exist")
			return
		}
		// set upload cache
		uploadCacheContent = &UploadCacheContent{
			UploadID:  uploadInfo.Id,
			Path:      uploadInfo.Path,
			Size:      uploadInfo.Size,
			Overwrite: uploadInfo.Overwrite,
		}

		file, err := os.OpenFile(uploadInfo.TmpPath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			msg := fmt.Sprintf("open file error, path: %s, err: %v", uploadInfo.TmpPath, err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, "open file error")
			return
		}
		uploadCacheContent.Fd = file

		uploadCache.Set(uploadCacheKey, uploadCacheContent, uploadCacheExpireTime)
		uploadCache.OnEvicted(func(key string, value interface{}) {
			if uploadContent, ok := value.(*UploadCacheContent); ok {
				err := uploadContent.Fd.Close()
				if err != nil {
					logger.Errorf("close upload file failed, err: %v", err)
				}
			}
		})

	} else {
		var ok bool
		if uploadCacheContent, ok = content.(*UploadCacheContent); !ok {
			msg := fmt.Sprintf("uploadCache type err, uploadCacheKey:%s", uploadCacheKey)
			uploadCache.Delete(uploadCacheKey)
			logger.Error(msg)
			common.InternalServerError(ctx, "uploadCache type err")
			return
		}
	}

	//check if offset and length are valid
	if request.Offset < 0 {
		msg := fmt.Sprintf("offset must be greater than or equal to 0, offset: %v", request.Offset)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidOffset, msg)
		return
	}

	if request.Length < 0 {
		msg := fmt.Sprintf("length must be greater than or equal to 0, length: %v", request.Length)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
		return
	}

	if request.Offset > uploadCacheContent.Size {
		msg := fmt.Sprintf("offset exceeds file size, offset: %v, uploadInfo.Size: %v", request.Offset, uploadCacheContent.Size)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidOffset, msg)
		return
	}

	if request.Length > maxSliceSize {
		msg := fmt.Sprintf("length must be less than 20MB , length: %v", request.Length)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.LengthTooLarge, msg)
		return
	}
	if request.Offset+request.Length > uploadCacheContent.Size {
		msg := fmt.Sprintf("offset plus length must be less than file size,  offset: %v, length: %v", request.Offset, request.Length)
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
		return
	}

	if ctx.Request.Body == nil || ctx.Request.Body == http.NoBody {
		msg := fmt.Sprintf("slice invalid")
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.InvalidParams(ctx, "slice invalid")
		return
	}

	compressor, err := compress.GetCompressor(request.Compressor)
	if err != nil {
		msg := fmt.Sprintf("get compressor error, err: %v", err)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidCompressor, "get compressor error")
		return
	}

	decompressReader, err := compressor.Decompress(ctx.Request.Body)
	if err != nil {
		msg := fmt.Sprintf("uncompress error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "uncompress error")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("close request body error, err: %v", err)
		}
	}(ctx.Request.Body)

	body, err := io.ReadAll(decompressReader)
	if err != nil {
		msg := fmt.Sprintf("read body failed, err: %v", err)
		logger.Error(msg)
		uploadCache.Delete(uploadCacheKey)
		common.InternalServerError(ctx, "read body failed")
		return
	}
	defer func() { ctx.Request.Body.Close() }()

	// check storage usage
	parts := strings.Split(fsutil.TrimPrefix(uploadCacheContent.Path, s.rootPath), "/")
	if len(parts) < 2 {
		msg := fmt.Sprintf("invalid path, path: %s", fsutil.TrimPrefix(uploadCacheContent.Path, s.rootPath))
		logger.Info(msg)
		uploadCache.Delete(uploadCacheKey)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}
	if !s.Quota.CheckStorageUsageAndHandleError(parts[1], float64(len(body)), logger, ctx) {
		uploadCache.Delete(uploadCacheKey)
		return
	}

	//upload file slice
	uploadTmpFilePath := filepath.Join(s.rootPath, fsutil.TmpBucketUploadingJoin(request.UploadID))

	if err := s.uploadSlice(uploadTmpFilePath, request.Offset, body, uploadCacheContent.Fd, logger); err != nil {
		msg := fmt.Sprintf("write slice failed, err: %v", err)
		logger.Error(msg)
		uploadCache.Delete(uploadCacheKey)
		common.InternalServerError(ctx, "write slice failed")
		return
	}

	//last upload, sync file & delete uploadCache & lock
	if request.Offset+request.Length == uploadCacheContent.Size {
		s.fileLocks.Delete(uploadTmpFilePath)
		if err := uploadCacheContent.Fd.Sync(); err != nil {
			msg := fmt.Sprintf("sync file error, err: %v", err)
			logger.Error(msg)
			uploadCache.Delete(uploadCacheKey)
			common.InternalServerError(ctx, "internal server error")
			return
		}
	}

	common.SuccessResp(ctx, nil)

}
