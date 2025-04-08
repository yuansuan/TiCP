package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/start"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/fsutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (t *DirectoryUsage) DirectoryUsageStart(ctx *gin.Context) {
	userID, accessKey, systemFlag, err := t.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "DirectoryUsageStart", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &start.Request{}
	if err := ctx.BindJSON(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if userID == "" {
		msg := fmt.Sprintf("Error parsing userID: %v,got: %v", err, userID)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}
	_, err = snowflake.ParseString(userID)
	if err != nil {
		msg := fmt.Sprintf("Error parsing userID: %v,got: %v", err, userID)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, "invalid user id")
		return
	}

	if len(strings.TrimSpace(request.Path)) == 0 {
		msg := fmt.Sprintf("Path is empty")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, msg)
		return
	}

	//check if user has access to the path
	if !systemFlag && !t.CheckPathAccessAndHandleError(accessKey, userID, request.Path, logger, ctx) {
		return
	}

	taskID := uuid.New().String()
	targetAbsPath := filepath.Join(t.rootPath, fsutil.TrimPrefix(request.Path, "/"))

	fileInfo, err := os.Stat(targetAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("Path does not exist: %s", request.Path)
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, msg)
			return
		}
		logger.Errorf("Error checking path: %s", err)
		common.InternalServerError(ctx, "internal server error")
		return
	}
	realPath, err := filepath.EvalSymlinks(targetAbsPath)
	if err != nil {
		msg := fmt.Sprintf("Error evaluating symlink: %s, path: %v", err, targetAbsPath)
		logger.Errorf(msg)
		common.InternalServerError(ctx, "internal server error")
		return
	}
	fileInfo, err = os.Stat(realPath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("realPath does not exist: %s", realPath)
			logger.Infof(msg)
			common.ErrorResp(ctx, http.StatusNotFound, commoncode.PathNotFound, fmt.Sprintf("realPath does not exist: %s", fsutil.TrimPrefix(request.Path, t.rootPath)))
			return
		}
		logger.Errorf("Error checking path: %s", err)
		common.InternalServerError(ctx, "internal server error")
		return
	}
	if !fileInfo.IsDir() {
		msg := fmt.Sprintf("Path is not a directory: %s", realPath)
		logger.Infof(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPath, fmt.Sprintf("Path is not a directory: %s", fsutil.TrimPrefix(request.Path, t.rootPath)))
		return
	}
	go t.StartCalculateDirectoryUsage(context.Background(), logger, taskID, targetAbsPath, userID, false)

	common.SuccessResp(ctx, &start.Data{
		DirectoryUsageTaskID: taskID,
	})
}
