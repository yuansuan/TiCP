package cashvoucher

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/add"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func Add(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &add.Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		msg := fmt.Sprintf("create cash voucher failed, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if util.IsBlank(req.CashVoucherName) {
		msg := fmt.Sprintf("create cash voucher failed, voucher name can not be empty!")
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if !validator.ValidAmount(ctx, req.Amount) ||
		!validator.ValidComment(ctx, req.Comment) {
		return
	}

	expiredType := voucherexpiredtype.ExpiredType(req.ExpiredType)
	if !voucherexpiredtype.Valid(expiredType) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "voucher expired type invalid")
		return
	}

	if expiredType == voucherexpiredtype.RelativeExpired {
		if req.RelExpiredTime <= 0 {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "voucher relative expired time can not less than or equal to 0")
			return
		}
	}

	if expiredType == voucherexpiredtype.AbsExpired {
		if util.IsBlank(req.AbsExpiredTime) {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "voucher absolute expired time can not be blank")
			return
		}

		_, toTimeResult := util.StringToTime(req.AbsExpiredTime)
		if !toTimeResult {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "voucher absolute expired time date format is incorrect，for example:2006-01-02 15:04:05")
			return
		}
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.AddCashVoucher(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
