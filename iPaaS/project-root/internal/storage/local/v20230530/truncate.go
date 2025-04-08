package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/truncate"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Truncate 改变文件大小
func (s *Storage) Truncate(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Truncate", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &truncate.Request{}
	if err := ctx.BindJSON(request); err != nil {
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

	//check if user has access to the path and check storage usage
	if !systemFlag && !s.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) ||
		!s.Quota.CheckStorageUsageAndHandleError(pathUserID, float64(request.Size), logger, ctx) {
		return
	}

	if request.Size < 0 {
		msg := fmt.Sprintf("invalid size, size: %d", request.Size)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidSize, msg)
		return
	}

	// generate absolute path
	fileAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	// stat file
	if flag, _, err = s.HandlePathContainsFileError(ctx, logger, fileAbsPath, fmt.Sprintf("path contains file, path: %s", request.Path)); !flag {
		return
	}

	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("stat file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, msg)
		return
	}
	if !request.CreateIfNotExists && os.IsNotExist(err) {
		msg := "file not exists and no need to create, path: " + request.Path
		logger.Info(msg)
		common.SuccessResp(ctx, nil)
		return
	}

	// truncate file
	if err = fsutil.Truncate(fileAbsPath, request.Size, logger); err != nil {
		msg := fmt.Sprintf("truncate file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "truncate file error")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		DestPath:      request.Path,
		FileType:      commoncode.FILE,
		OperationType: commoncode.TRUNCATE,
		Size:          s.OperationLog.FormatBytes(request.Size),
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
