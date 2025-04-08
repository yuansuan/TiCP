package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/status"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/directoryUsage"
	"net/http"
	"strings"
)

const (
	directoryUsageCalculatingStatus = "calculating"
	directoryUsageSuccessStatus     = "success"
	directoryUsageFailedStatus      = "failed"
	directoryUsageCanceledStatus    = "canceled"
	directoryUsageUnknownStatus     = "unknown"
)

func (t *DirectoryUsage) DirectoryUsageStatus(ctx *gin.Context) {

	userID, _, systemFlag, err := t.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "DirectoryUsageStatus", "RequestId ", ctx.GetHeader(common.RequestIDKey), "UserId", userID)

	request := &status.Request{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
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
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.DirectoryUsageTaskNotFound, fmt.Sprintf("directory usage not exist, taskID: %s", request.DirectoryUsageTaskID))
		return
	}

	if !systemFlag && !strings.Contains(directoryUsage.Path, userID) {
		msg := fmt.Sprintf("user %s has no access to directory usage, path:%s", userID, directoryUsage.Path)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.DirectoryUsageTaskNoAccess, msg)
		return
	}
	s := ""
	switch directoryUsage.Status {
	case model.DirectoryUsageTaskFailed:
		s = directoryUsageFailedStatus
		logger.Infof("DirectoryUsageStatus, calculating err:%s,taskID: %s", directoryUsage.ErrMsg, request.DirectoryUsageTaskID)
	case model.DirectoryUsageTaskCalculating:
		s = directoryUsageCalculatingStatus
	case model.DirectoryUsageTaskFinished:
		s = directoryUsageSuccessStatus
	case model.DirectoryUsageTaskCanceled:
		s = directoryUsageCanceledStatus
	default:
		s = directoryUsageUnknownStatus
	}
	common.SuccessResp(ctx, &status.Data{
		Status:     s,
		Size:       directoryUsage.Size,
		LogicSize:  directoryUsage.LogicSize,
		ErrMessage: directoryUsage.ErrMsg,
	})
}
