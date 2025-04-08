package accountcashvoucher

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/get"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func Get(ctx *gin.Context) {
	accountCashVoucherID := ctx.Param(consts.Account_CashVoucher_ID_KEY)
	if !validator.ValidAccountCashVoucherID(ctx, accountCashVoucherID) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.AccountCashVoucherGetByID(ctx, &get.Request{
		AccountCashVoucherID: accountCashVoucherID,
	}, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
