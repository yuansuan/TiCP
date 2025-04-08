package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Mv 移动一个文件或文件夹。
func (s *Storage) Mv(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Mv", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &mv.Request{}
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

	//check if user has access to the path and check storage usage
	if !systemFlag {
		if !s.CheckPathAccessAndHandleError(accessKey, userID, request.SrcPath, logger, ctx) ||
			!s.CheckPathAccessAndHandleError(accessKey, userID, request.DestPath, logger, ctx) {
			return
		}
	}

	// generate absolute path
	srcAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.SrcPath, "/"))

	// stat file or directory
	srcFileInfo, err := os.Stat(srcAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := "file not found, path: " + request.SrcPath
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.SrcPathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file error")
		return
	}
	filetype := commoncode.FILE
	if srcFileInfo != nil && srcFileInfo.IsDir() {
		filetype = commoncode.FOLDER
	}

	//check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, 0, logger, ctx) {
		return
	}

	// generate absolute path
	destAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.DestPath, "/"))

	// stat file or directory
	var destFileInfo os.FileInfo
	if flag, destFileInfo, err = s.HandlePathContainsFileError(ctx, logger, destAbsPath, fmt.Sprintf("dest path contains file, path: %s", request.DestPath)); !flag {
		return
	}

	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("stat file or directory error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file or directory error")
		return
	}

	if destFileInfo != nil {
		msg := "file or directory exist, path: " + request.DestPath
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.DestPathExists, msg)
		return
	}

	if err := os.MkdirAll(filepath.Dir(destAbsPath), filemode.Directory); err != nil {
		msg := fmt.Sprintf("mkdir error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "mkdir error")
		return
	}

	if err := os.Rename(srcAbsPath, destAbsPath); err != nil {
		msg := fmt.Sprintf("move error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "move error")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.DestPath),
		SrcPath:       request.SrcPath,
		DestPath:      request.DestPath,
		FileType:      filetype,
		OperationType: commoncode.MOVE,
		Size:          "",
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
