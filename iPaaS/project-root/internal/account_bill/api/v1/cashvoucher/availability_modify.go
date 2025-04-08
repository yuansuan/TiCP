package cashvoucher

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	availabilityModify "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/availabilitymodify"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AvailabilityModify(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	cashVoucherID := ctx.Param(consts.CASH_VOUCHER_ID)
	if !validator.ValidCashVoucherID(ctx, cashVoucherID) {
		return
	}

	req := &availabilityModify.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	b, _ := voucherstatus.ValidAvailabilityStatusString(req.AvailabilityStatus)
	if !b {
		msg := fmt.Sprintf("invalid availability status")
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	err = module.AvailabilityModify(ctx, req, optUserID)
	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, nil)
}
