package apiuser

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/billlist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AccountUserBillList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &billlist.Request{}
	err := ctx.BindQuery(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if req.PageIndex < 1 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidPageIndexErrorCode, "index can not be less than one")
		return
	}

	if req.PageSize < 1 || req.PageSize > 1000 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidPageSizeErrorCode, "page size should be in [1, 1000]")
		return
	}

	var startTime, endTime *time.Time
	var toTimeResult bool
	if req.StartTime != "" {
		startTime, toTimeResult = util.StringToTime(req.StartTime)
		if !toTimeResult {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "start time date format is incorrect，for example:2006-01-02 15:04:05")
			return
		}
	}

	if req.EndTime != "" {
		endTime, toTimeResult = util.StringToTime(req.EndTime)
		if !toTimeResult {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "end time date format is incorrect，for example:2006-01-02 15:04:05")
			return
		}
	}

	// 校验 end > start
	if req.StartTime != "" && req.EndTime != "" {
		if endTime.Before(*startTime) {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidEndTime, "end time should not be after start time")
			return
		}
	}

	userID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.UserBillList(ctx, req, userID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
