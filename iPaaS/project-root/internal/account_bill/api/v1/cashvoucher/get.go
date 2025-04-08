package cashvoucher

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/get"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func Get(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	cashVoucherID := ctx.Param(consts.CASH_VOUCHER_ID)
	if !validator.ValidCashVoucherID(ctx, cashVoucherID) {
		return
	}

	req := &get.Request{CashVoucherID: cashVoucherID}
	err := ctx.BindQuery(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.GetCashVoucherByID(ctx, req, optUserID)
	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
