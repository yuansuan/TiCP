package job

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// judgeCommonError 判断公共错误
func judgeCommonError(c *gin.Context, err error) (isHandled bool) {
	if errors.Is(err, common.ErrInvalidUserID) {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, err.Error())
		return true
	}
	if errors.Is(err, common.ErrUserNotExists) {
		common.ErrorResp(c, http.StatusBadRequest, api.UserNotExistsErrorCode, err.Error())
		return true
	}
	if errors.Is(err, common.ErrJobIDNotFound) {
		common.ErrorResp(c, http.StatusNotFound, api.JobIDNotFound, err.Error())
		return true
	}
	if errors.Is(err, common.ErrInvalidJobID) {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidJobID, err.Error())
		return true
	}
	if errors.Is(err, common.ErrJobAccessDenied) {
		common.ErrorResp(c, http.StatusForbidden, api.JobAccessDenied, err.Error())
		return true
	}
	if errors.Is(err, common.ErrInvalidArgumentAllocType) {
		common.ErrorResp(c, http.StatusForbidden, api.InvalidArgumentAllocType, err.Error())
		return true
	}
	return false
}

// JudgeError 判断错误
func JudgeError(c *gin.Context, err error, errorHandlers ...func(*gin.Context, error) bool) (isHandled bool) {
	if err == nil {
		return false
	}
	if judgeCommonError(c, err) {
		return true
	}
	for _, errorHandler := range errorHandlers {
		if errorHandler(c, err) {
			return true
		}
	}
	common.InternalServerError(c, err.Error())
	return true
}

// JudgeGetError 判断获取错误
func JudgeGetError(c *gin.Context, err error) (isHandled bool) {
	return JudgeError(c, err)
}

// JudgeBatchGetError 判断批量获取错误
func JudgeBatchGetError(c *gin.Context, err error) bool {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentJobIDs) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentJobIDs, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeListError 判断获取列表错误
func JudgeListError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentPageSize) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPageSize, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentPageOffset) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPageOffset, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentZone) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentZone, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentJobState) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentJobState, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidAppID) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidAccountId) {
			common.ErrorResp(c, http.StatusBadRequest, api.InternalErrorInvalidAccountId, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeListSyncNeedFileError 需要同步文件作业列表错误
func JudgeListSyncNeedFileError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentPageSize) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPageSize, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentPageOffset) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPageOffset, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentZone) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentZone, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeSubmitError 判断提交错误
func JudgeSubmitError(c *gin.Context, err error) (isHandled bool) {
	preScheduleErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidPreScheduleID) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPreScheduleID, err.Error())
			return true
		}
		if errors.Is(err, common.ErrPreScheduleNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.PreScheduleNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrPreScheduleUsed) {
			common.ErrorResp(c, http.StatusForbidden, api.PreScheduleUsed, err.Error())
			return true
		}
		return false
	}

	appErrorHandler := func(c *gin.Context, err error) bool {
		// ErrInvalidAppID
		if errors.Is(err, common.ErrInvalidAppID) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, err.Error())
			return true
		}
		// ErrAppIDNotFound
		if errors.Is(err, common.ErrAppIDNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, err.Error())
			return true
		}
		// ErrUserNoAppQuota
		if errors.Is(err, common.ErrUserNoAppQuota) {
			common.ErrorResp(c, http.StatusForbidden, api.UserNoAppQuota, err.Error())
			return true
		}
		// ErrAppNotPublished
		if errors.Is(err, common.ErrAppNotPublished) {
			common.ErrorResp(c, http.StatusForbidden, api.AppNotPublished, err.Error())
			return true
		}
		// ErrInvalidArgumentCommand
		if errors.Is(err, common.ErrInvalidArgumentCommand) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCommand, err.Error())
			return true
		}
		// ErrInvalidArgumentEnv
		if errors.Is(err, common.ErrInvalidArgumentEnv) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentEnv, err.Error())
			return true
		}
		return false
	}

	resourceErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentResource) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentResource, err.Error())
			return true
		}
		if errors.Is(err, common.ErrQuotaExhaustedResource) {
			common.ErrorResp(c, http.StatusBadRequest, api.QuotaExhaustedResource, err.Error())
			return true
		}
		return false
	}

	fileErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentInput) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentInput, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentOutput) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentOutput, err.Error())
			return true
		}
		if errors.Is(err, common.ErrJobPathUnauthorized) {
			common.ErrorResp(c, http.StatusForbidden, api.JobPathUnauthorized, err.Error())
			return true
		}
		return false
	}

	customRuleErrorHandler := func(c *gin.Context, err error) bool {
		// ErrInvalidArgumentCustomStateRuleKeyStatement
		if errors.Is(err, common.ErrInvalidArgumentCustomStateRuleKeyStatement) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCustomStateRuleKeyStatement, err.Error())
			return true
		}
		// ErrInvalidArgumentCustomStateRuleResultState
		if errors.Is(err, common.ErrInvalidArgumentCustomStateRuleResultState) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCustomStateRuleResultState, err.Error())
			return true
		}
		return false
	}

	chargeErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidChargeParams) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentChargeParams, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidAccountId) {
			common.ErrorResp(c, http.StatusInternalServerError, api.InternalErrorInvalidAccountId, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidAccountStatusNotEnoughBalance) {
			common.ErrorResp(c, http.StatusForbidden, api.InvalidAccountStatusNotEnoughBalance, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidAccountStatusFrozen) {
			common.ErrorResp(c, http.StatusForbidden, api.InvalidAccountStatusFrozen, err.Error())
			return true
		}
		return false
	}

	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidArgumentName) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentName, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentComment) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentComment, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentZone) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentZone, err.Error())
			return true
		}
		return false
	}

	payByErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidPayBy) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentPayBy, err.Error())
			return true
		}
		if errors.Is(err, common.ErrPayByTokenExpire) {
			common.ErrorResp(c, http.StatusBadRequest, api.PayByTokenExpire, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentPayBySignature) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentPayBySignature, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler, preScheduleErrorHandler,
		appErrorHandler, resourceErrorHandler, fileErrorHandler,
		customRuleErrorHandler, chargeErrorHandler, payByErrorHandler)
}

// JudgeDeleteError 判断删除错误
func JudgeDeleteError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobStateNotAllowDelete) {
			common.ErrorResp(c, http.StatusForbidden, api.JobStateNotAllowDelete, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

func JudgeRetransmitError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobNotAllowedRetransmit) {
			common.ErrorResp(c, http.StatusForbidden, api.JobNotAllowedRetransmit, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeTerminateError 判断终止错误
func JudgeTerminateError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		// ErrJobStateNotAllowTerminate
		if errors.Is(err, common.ErrJobStateNotAllowTerminate) {
			common.ErrorResp(c, http.StatusForbidden, api.JobStateNotAllowTerminate, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeResumeError 判断恢复错误
func JudgeResumeError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		// ErrJobStateNotAllowTerminate
		if errors.Is(err, common.ErrJobStateNotAllowResume) {
			common.ErrorResp(c, http.StatusForbidden, api.JobStateNotAllowResume, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeTransmitSuspendError 判断传输暂停错误
func JudgeTransmitSuspendError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobStateNotAllowTransmitSuspend) {
			common.ErrorResp(c, http.StatusForbidden, api.JobStateNotAllowTransmitSuspend, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeTransmitResumeError 判断传输恢复错误
func JudgeTransmitResumeError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobStateNotAllowTransmitResume) {
			common.ErrorResp(c, http.StatusForbidden, api.JobStateNotAllowTransmitResume, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeGetResidualError 判断获取残差图错误
func JudgeGetResidualError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobResidualNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.JobResidualNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrHpcResidual) {
			common.ErrorResp(c, http.StatusNotFound, api.JobGetHpcResidualFailed, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeListSnapshotError 判断获取云图列表错误
func JudgeListSnapshotError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobSnapshotNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.JobSnapshotNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidPath) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPath, err.Error())
			return true
		}
		if errors.Is(err, common.ErrPathNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.PathNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrAppIDNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeGetSnapshotError 判断获取云图错误
func JudgeGetSnapshotError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobSnapshotNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.JobSnapshotNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidPath) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidPath, err.Error())
			return true
		}
		if errors.Is(err, common.ErrPathNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.PathNotFound, err.Error())
			return true
		}
		if errors.Is(err, common.ErrAppIDNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeGetMonitorChartError 判断获取监控图表错误
func JudgeGetMonitorChartError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		// ErrJobMonitorChartNotFound
		if errors.Is(err, common.ErrJobMonitorChartNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.JobMonitorChartNotFound, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgePreScheduleError 判断预调度错误
func JudgePreScheduleError(c *gin.Context, err error) (isHandled bool) {
	appErrorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrInvalidAppID) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, err.Error())
			return true
		}
		if errors.Is(err, common.ErrAppIDNotFound) {
			common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, err.Error())
			return true
		}
		if errors.Is(err, common.ErrUserNoAppQuota) {
			common.ErrorResp(c, http.StatusForbidden, api.UserNoAppQuota, err.Error())
			return true
		}
		if errors.Is(err, common.ErrAppNotPublished) {
			common.ErrorResp(c, http.StatusForbidden, api.AppNotPublished, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentCommand) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCommand, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentEnv) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentEnv, err.Error())
			return true
		}
		return false
	}

	resourceErrorHandler := func(c *gin.Context, err error) bool {
		// ErrQuotaExhaustedResource
		if errors.Is(err, common.ErrQuotaExhaustedResource) {
			common.ErrorResp(c, http.StatusBadRequest, api.QuotaExhaustedResource, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, appErrorHandler, resourceErrorHandler)
}

// JudgeSyncFileStateError 判断更新文件传输状态错误
func JudgeSyncFileStateError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobFileSyncStateUpdateFailed) {
			common.ErrorResp(c, http.StatusBadRequest, api.JobFileSyncStateUpdateFailed, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentDownloadFinishedTime) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDownloadFinishedTime, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentDownloadFinishedTime) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDownloadFinishedTime, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentDownloadFileSizeCurrent) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDownloadFileSizeCurrent, err.Error())
			return true
		}
		if errors.Is(err, common.ErrInvalidArgumentDownloadFileSizeTotal) {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDownloadFileSizeTotal, err.Error())
			return true
		}
		return false
	}
	return JudgeError(c, err, errorHandler)
}

// JudgeCpuUsageError 判断获取CPU使用率错误
func JudgeCpuUsageError(c *gin.Context, err error) (isHandled bool) {
	errorHandler := func(c *gin.Context, err error) bool {
		if errors.Is(err, common.ErrJobStateNotAllowQuery) {
			common.ErrorResp(c, http.StatusBadRequest, api.JobStateNotAllowQuery, err.Error())
			return true
		} else if errors.Is(err, common.ErrWrongCPUUsage) {
			common.ErrorResp(c, http.StatusServiceUnavailable, api.WrongCPUUsage, err.Error()) // 返回503
			return true
		} else {
			common.ErrorResp(c, http.StatusInternalServerError, "unknown", err.Error())
			return true
		}
	}
	return JudgeError(c, err, errorHandler)
}

func JudgeUpdateError(c *gin.Context, err error) (isHandled bool) {
	return JudgeError(c, err)
}
