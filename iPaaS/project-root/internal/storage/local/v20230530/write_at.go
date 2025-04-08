package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/writeAt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil/compress"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MaxLength = 10 * 1024 * 1024

// WriteAt 随机写文件
func (s *Storage) WriteAt(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "WriteAt", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &writeAt.Request{}
	if err := ctx.ShouldBind(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	flag, pathUserID, msg := fsutil.ValidateUserIDPath(request.Path)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		return
	}

	if request.Length < 0 || request.Length > MaxLength {
		msg := fmt.Sprintf("invalid length, length should be in [0, %v], got: %v", MaxLength, request.Length)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
		return
	}

	// generate absolute path
	filePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	//check if path exists and is a file
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("path not found, path: %s", request.Path)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("Lstat error, path: %s, err: %v", request.Path, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "Lstat error")
		return
	}

	if fileInfo.IsDir() {
		msg := fmt.Sprintf("path should be a file, got: %s", request.Path)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	// check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, float64(fileInfo.Size()), logger, ctx) {
		return
	}

	if request.Offset < 0 || request.Offset > fileInfo.Size() || (request.Offset == fileInfo.Size() && fileInfo.Size() != 0) {
		msg := fmt.Sprintf("invalid offset, offset should be greater than or equal to 0 and less than fileSize: %v, got: %d", fileInfo.Size(), request.Offset)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidOffset, msg)
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

	if err := s.writeAt(filePath, request.Offset, request.Length, decompressReader, logger); err != nil {
		msg := fmt.Sprintf("write at error, err: %v", err)
		if strings.Contains(err.Error(), pathNotSafeErrorMsg) {
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
			return
		}
		if strings.Contains(err.Error(), writeAtLengthNotMatchErrorMsg) {
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
			return
		}
		logger.Errorf(msg)
		common.InternalServerError(ctx, "write at failed")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		DestPath:      request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.WRITE_AT,
		Size:          s.OperationLog.FormatBytes(request.Length),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
