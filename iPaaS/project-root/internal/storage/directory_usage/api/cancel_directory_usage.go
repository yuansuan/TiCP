package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/cancel"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/directoryUsage"
	"net/http"
	"strings"
)

func (t *DirectoryUsage) DirectoryUsageCancel(ctx *gin.Context) {
	userID, _, systemFlag, err := t.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "DirectoryUsageCancel", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &cancel.Request{}
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
	if len(strings.TrimSpace(request.DirectoryUsageTaskID)) == 0 {
		msg := fmt.Sprintf("DirectoryUsageTaskID is empty")
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidDirectoryUsageTaskID, msg)
		return
	}

	exist, directoryUsage, err := directoryUsage.GetDirectory(ctx, t.Engine, t.DirectoryUsageDao, request.DirectoryUsageTaskID)
	if err != nil {
		msg := fmt.Sprintf("get directory usage failed, taskID: %s, userID: %v,err: %v", request.DirectoryUsageTaskID, userID, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "get directory usage failed")
		return
	}
	if !exist {
		msg := fmt.Sprintf("directory usage not existed, taskID: %s, userID: %v,err: %v", request.DirectoryUsageTaskID, userID, err)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.DirectoryUsageTaskNotFound, "directory usage not exist")
		return
	}

	if !systemFlag && !strings.Contains(directoryUsage.Path, userID) {
		msg := fmt.Sprintf("user %s has no access to directory usage, path:%s", userID, directoryUsage.Path)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.DirectoryUsageTaskNoAccess, msg)
		return
	}

	err = t.CancelCalculateDirectoryUsage(ctx, logger, request.DirectoryUsageTaskID)
	if err != nil {
		msg := fmt.Sprintf("cancel directory usage failed, taskID: %s, userID: %v,err: %v", request.DirectoryUsageTaskID, userID, err)
		logger.Error(msg)
		common.InternalServerError(ctx, "cancel directory usage failed")
		return
	}

	common.SuccessResp(ctx, nil)
}
