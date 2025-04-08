package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil/compress"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ReadAt 随机读文件
func (s *Storage) ReadAt(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "ReadAt", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &readAt.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	flag, _, msg := fsutil.ValidateUserIDPath(request.Path)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		return
	}

	// generate absolute path
	filePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	//check if path exists and is a file
	fileInfo, err := os.Stat(filePath)
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

	//deal with symlink
	if flag, filePath, fileInfo = s.HandleSymbolicLink(filePath, logger, fileInfo, ctx); !flag {
		return
	}

	if request.Offset < 0 || request.Offset > fileInfo.Size() || (request.Offset == fileInfo.Size() && fileInfo.Size() != 0) {
		msg := fmt.Sprintf("invalid offset, offset should be greater than or equal to 0 and less than fileSize: %v, got: %d", fileInfo.Size(), request.Offset)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidOffset, msg)
		return
	}

	if request.Length <= 0 {
		msg := fmt.Sprintf("invalid length, length should be greater than 0, got: %d", request.Length)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
		return
	}

	compressor, err := compress.GetCompressor(request.Compressor)
	if err != nil {
		msg := fmt.Sprintf("get compressor error, err: %v", err)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidCompressor, "get compressor error")
		return
	}

	reader, err := fsutil.ReadAt(filePath, request.Offset, request.Length)
	if err != nil {
		msg := fmt.Sprintf("read file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "read file error")
		return
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("close reader error, err: %v", err))
		}
	}(reader)

	compressWriter, err := compressor.Compress(ctx.Writer)
	if err != nil {
		msg := fmt.Sprintf("compress error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "compress error")
		return
	}
	if _, err := io.Copy(compressWriter, reader); err != nil {
		msg := fmt.Sprintf("copy error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "copy error")
		return
	}
	if err := compressWriter.Close(); err != nil {
		msg := fmt.Sprintf("close error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "close error")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		SrcPath:       request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.READ_AT,
		Size:          s.OperationLog.FormatBytes(request.Length),
		CreateTime:    time.Now(),
	})

	return
}
