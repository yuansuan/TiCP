package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/storage/operationLog/admin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	storageOperationLogService "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/service/storageOperationLog"
	"net/http"
)

func (o *OperationLog) ListOperationLogAdmin(ctx *gin.Context) {
	userID, _, err := o.GetUserIDAndAKAndHandleError(ctx)
	if err != nil {
		return
	}
	logger := logging.GetLogger(ctx).With("func", "ListOperationLog", "RequestId", ctx.GetHeader(common.RequestIDKey), "UserId", userID)
	request := &admin.ListOperationLogRequest{}
	if err := ctx.BindQuery(request); err != nil {
		msg := fmt.Sprintf("invalid params, err: %v", err)
		logger.Info(msg)
		common.InvalidParams(ctx, err.Error())
		return
	}

	var id snowflake.ID
	if request.UserID != "" {
		id, err = snowflake.ParseString(request.UserID)
		if err != nil {
			msg := fmt.Sprintf("parse userID failed, userID: %v, err: %v", request.UserID, err)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidUserID, msg)
			return
		}
	}

	if request.FileTypes != "" {
		if !CheckFileType(request.FileTypes) {
			msg := fmt.Sprintf("invalid file type, got: %v", request.FileTypes)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidFileType, msg)
			return
		}
	}

	if request.OperationTypes != "" {
		if !CheckOperationType(request.OperationTypes) {
			msg := fmt.Sprintf("invalid operation type, got: %v", request.OperationTypes)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidOperationType, msg)
			return
		}
	}

	if request.BeginTime != 0 {
		if !CheckTime(request.BeginTime) {
			msg := fmt.Sprintf("invalid begin time, got: %v", request.BeginTime)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidBeginTime, msg)
			return
		}
	}

	if request.EndTime != 0 {
		if !CheckTime(request.EndTime) || request.EndTime < request.BeginTime {
			msg := fmt.Sprintf("invalid end time, should be unix timestamp and less than begin time, got: %v", request.EndTime)
			logger.Info(msg)
			common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidEndTime, msg)
			return
		}
	}

	if request.PageOffset < 0 {
		msg := fmt.Sprintf("invalid page offset, page offset should be greater than or equal to 0, got: %v", request.PageOffset)
		logger.Error(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageOffset, msg)
		return
	}

	if request.PageSize < 1 || request.PageSize > 1000 {
		msg := fmt.Sprintf("invalid page size, page size should be in [1, 1000], got: %v", request.PageSize)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusBadRequest, commoncode.InvalidPageSize, msg)
		return
	}

	userIDs, err := storageOperationLogService.GetUserIDs(ctx, o.Engine, o.StorageOperationLogDao)
	flag := false
	for _, userID := range userIDs {
		if userID == fmt.Sprint(id) {
			flag = true
		}
	}

	if !flag {
		msg := fmt.Sprintf("this user has no operation log, userID: %v", fmt.Sprint(id))
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusNotFound, commoncode.UserNotExistsErrorCode, msg)
		return
	}

	param := &dao.StorageOperationLogQueryParam{
		UserIDs:        fmt.Sprint(id),
		FileName:       request.FileName,
		FileTypes:      request.FileTypes,
		OperationTypes: request.OperationTypes,
		BeginTime:      UnixToTime(request.BeginTime),
		EndTime:        UnixToTime(request.EndTime),
		PageOffset:     request.PageOffset,
		PageSize:       request.PageSize,
	}

	operationLog, err, next, total := storageOperationLogService.ListStorageOperationLog(ctx, o.Engine, o.StorageOperationLogDao, param)
	if err != nil {
		logger.Error(err.Error())
		common.InternalServerError(ctx, err.Error())
		return
	}

	common.SuccessResp(ctx, &admin.ListOperationLogResponseData{
		OperationLog: ToResponseOperationLogs(operationLog),
		NextMarker:   next,
		Total:        total,
	})
}
