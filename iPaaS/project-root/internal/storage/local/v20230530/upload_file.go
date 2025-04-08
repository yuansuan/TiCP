package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	uploadFile "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/file"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const maxFileSize = 1024 * 1024 * 100 // 100MB

func (s *Storage) UploadFile(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func ", "UploadFile ", "RequestId ", ctx.GetHeader(common.RequestIDKey), " UserId ", userID)

	request := &uploadFile.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, "invalid params")
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

	// generate absolute path
	filePath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

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
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathExists, msg)
			return
		}
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		msg := fmt.Sprintf("read body failed, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "read body failed")
		return
	}
	defer func() { ctx.Request.Body.Close() }()

	if len(body) > maxFileSize {
		msg := fmt.Sprintf("file size too large, size: %d, max: %d", len(body), maxFileSize)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.SizeTooLarge, msg)
		return
	}

	//check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, float64(len(body)), logger, ctx) {
		return
	}

	file, err := fsutil.CreateAndReturnFile(filePath, int64(len(body)), logger)
	if err != nil {
		msg := fmt.Sprintf("create file failed, file path: %s, err: %v", request.Path, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "create file failed")
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			msg := fmt.Sprintf("close file failed, err: %v", err)
			logger.Error(msg)
			common.InternalServerError(ctx, "close file failed")
			return
		}
	}()

	if err = s.uploadSlice(filePath, 0, body, file, logger); err != nil {
		msg := fmt.Sprintf("write file failed, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "write file failed")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(filePath),
		DestPath:      request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.UPLOAD,
		Size:          s.OperationLog.FormatBytes(int64(len(body))),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)

}
