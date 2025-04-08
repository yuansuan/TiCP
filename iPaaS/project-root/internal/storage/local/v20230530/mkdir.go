package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (s *Storage) Mkdir(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Mkdir", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &mkdir.Request{}
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
		!s.Quota.CheckStorageUsageAndHandleError(pathUserID, 0, logger, ctx) {
		return
	}

	// generate absolute path
	folderAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	var fileInfo os.FileInfo
	if flag, fileInfo, err = s.HandlePathContainsFileError(ctx, logger, folderAbsPath, fmt.Sprintf("path contains file, path: %s", request.Path)); !flag {
		return
	}
	// if not ignore exist, check if dir exists
	if !request.IgnoreExist {
		if err != nil && !os.IsNotExist(err) {
			msg := fmt.Sprintf("stat file error, err: %v", err)
			logger.Errorf(msg)
			common.InternalServerError(ctx, msg)
			return
		}
		if fileInfo != nil {
			msg := "directory already exists, path: " + request.Path
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.PathExists, msg)
			return
		}
	}

	// make dir
	if err := os.MkdirAll(folderAbsPath, filemode.Directory); err != nil {
		msg := fmt.Sprintf("mkdir error, err: %v", err)
		logger.Error(msg)
		common.InternalServerError(ctx, "mkdir error")
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.Path),
		DestPath:      request.Path,
		FileType:      commoncode.FOLDER,
		OperationType: commoncode.MKDIR,
		Size:          "",
		CreateTime:    time.Now(),
	})

	logger.Infof("mkdir succeed, directory path: %v, userID: %v", folderAbsPath, userID)

	common.SuccessResp(ctx, nil)
}
