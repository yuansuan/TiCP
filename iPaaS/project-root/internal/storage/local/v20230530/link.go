package v20230530

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"time"
)

// Link 链接一个文件或目录
func (s *Storage) Link(ctx *gin.Context) {

	userID, accessKey, systemFlag, err := s.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "Link", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &link.Request{}
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
	// generate absolute path
	srcAbsPath := filepath.Join(s.rootPath, fsutil.TrimPrefix(request.SrcPath, "/"))

	//check if user has access to the path
	if !systemFlag {
		if !s.CheckPathAccessAndHandleError(accessKey, userID, request.SrcPath, logger, ctx) ||
			!s.CheckPathAccessAndHandleError(accessKey, userID, request.DestPath, logger, ctx) {
			return
		}
	}

	// stat file or directory
	fileType := commoncode.FILE
	fileInfo, err := os.Lstat(srcAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := "file or directory not found, path: " + request.SrcPath
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.SrcPathNotFound, msg)
			return
		}
		msg := fmt.Sprintf("stat file or directory error, err: %v", err)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "stat file or directory error")
		return
	}
	if fileInfo != nil && fileInfo.IsDir() {
		fileType = commoncode.FOLDER
	}

	//check storage usage
	if !s.Quota.CheckStorageUsageAndHandleError(pathUserID, s.GetTotalSize(logger, srcAbsPath), logger, ctx) {
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

	// check if is hard link and src path is directory
	// use soft link instead
	if config.GetConfig().Local.LinkType == HardLink && fileInfo.IsDir() {
		sl := &linker.SoftLink{}
		err = sl.Link(srcAbsPath, destAbsPath)
	} else {
		err = s.linker.Link(srcAbsPath, destAbsPath)
	}

	if err != nil {
		msg := fmt.Sprintf("link file or directory error, err: %v", err)
		// 应该只是link exist的时候才返回400
		if strings.Contains(err.Error(), "file exists") {
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidParamErrorCode, "file exists")
		} else {
			//否则返回500
			logger.Errorf(msg)
			common.InternalServerError(ctx, "link file or directory error")
		}
		return
	}

	s.OperationLog.InsertOperationLog(logger, ctx, &model.StorageOperationLog{
		UserId:        userID,
		FileName:      filepath.Base(request.DestPath),
		SrcPath:       request.SrcPath,
		DestPath:      request.DestPath,
		FileType:      fileType,
		OperationType: commoncode.LINK,
		Size:          "",
		CreateTime:    time.Now(),
	})

	common.SuccessResp(ctx, nil)
}
