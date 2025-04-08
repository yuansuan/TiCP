package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copyRange"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// CopyRange 复制一个文件的一部分到另一个文件
func (s *Storage) CopyRange(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "CopyRange", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &copyRange.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	flag, _, msg := fsutil.ValidateUserIDPath(request.SrcPath)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	flag, pathUserID, msg := fsutil.ValidateUserIDPath(request.DestPath)
	if !flag {
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag {
		if !s.CheckPathAccessAndHandleError(accessKey, userID, request.SrcPath, logger, ctx) ||
			!s.CheckPathAccessAndHandleError(accessKey, userID, request.DestPath, logger, ctx) {
			return
		}
	}

	// generate absolute path
	srcAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.SrcPath, "/"))

	// check src path exist and is a file
	srcFileInfo, err := os.Stat(srcAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("src file not found, path: %s", request.SrcPath)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.SrcPathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
		return
	}
	if srcFileInfo.IsDir() {
		msg := fmt.Sprintf("src path should be a file, got: %s", request.SrcPath)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, float64(srcFileInfo.Size()), logger, ctx) {
		return
	}

	// generate absolute path
	destAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.DestPath, "/"))

	// check dest path exist and is a file
	var destFileInfo os.FileInfo
	if flag, destFileInfo, err = s.HandlePathContainsFileError(ctx, logger, destAbsPath, fmt.Sprintf("dest path contains file, path: %s", request.DestPath)); !flag {
		return
	}

	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("dest path not found, path: %s", request.DestPath)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.DestPathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("Lstat error, path: %s, err: %v", request.DestPath, err)
		logger.Error(msg)
		common.InternalServerError(ctx, msg)
		return
	}
	if destFileInfo.IsDir() {
		msg := fmt.Sprintf("dest path should be a file, got: %s", request.DestPath)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	// open src
	src, err := os.Open(srcAbsPath)
	if err != nil {
		msg := fmt.Sprintf("open file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "open file error")
		return
	}
	defer func(src *os.File) {
		if src == nil {
			return
		}
		err = src.Close()
		if err != nil {
			logger.Errorf("close file error, err: %v", err)
		}
	}(src)

	if request.SrcOffset < 0 || request.SrcOffset > srcFileInfo.Size() {
		msg := fmt.Sprintf("invalid src offset, srcOffset should be in [0, srcFileSize], srcFileSize: %d", srcFileInfo.Size())
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidSrcOffset, msg)
		return
	}

	if request.DestOffset < 0 || request.DestOffset > destFileInfo.Size() {
		msg := fmt.Sprintf("invalid dest offset, destOffset should be in [0, destFileSize], destFileSize: %d", destFileInfo.Size())
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidDestOffset, msg)
		return
	}

	if request.Length < 0 || request.SrcOffset+request.Length > srcFileInfo.Size() {
		msg := fmt.Sprintf("invalid length of data, length should be in [0, srcFileSize - srcOffset], srcFileSize: %d, srcOffset: %d", srcFileInfo.Size(), request.SrcOffset)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidLength, msg)
		return
	}

	//read data
	data := make([]byte, request.Length)
	_, err = src.ReadAt(data, request.SrcOffset)
	if err != nil && err != io.EOF {
		msg := fmt.Sprintf("read file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, msg)
		return
	}

	// open dest
	dest, err := os.OpenFile(destAbsPath, os.O_RDWR, 0666)
	if err != nil {
		msg := fmt.Sprintf("open file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "open file error")
		return
	}
	defer func(dest *os.File) {
		if dest == nil {
			return
		}
		err = dest.Close()
		if err != nil {
			logger.Errorf("close file error, err: %v", err)
		}
	}(dest)

	// write data
	_, err = dest.WriteAt(data, request.DestOffset)
	if err != nil {
		msg := fmt.Sprintf("write file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "write file error")
		return
	}

	// sync it
	if err = dest.Sync(); err != nil {
		msg := fmt.Sprintf("sync file error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "sync file error")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.DestPath),
		SrcPath:       request.SrcPath,
		DestPath:      request.DestPath,
		FileType:      commoncode.FILE,
		OperationType: commoncode.COPY_RANGE,
		Size:          s.OperationLog.FormatBytes(request.Length),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
