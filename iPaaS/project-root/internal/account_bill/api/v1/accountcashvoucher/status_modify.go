package accountcashvoucher

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/statusmodify"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func StatusModify(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	accountCashVoucherID := ctx.Param(consts.Account_CashVoucher_ID_KEY)
	if !validator.ValidAccountCashVoucherID(ctx, accountCashVoucherID) {
		return
	}

	req := &statusmodify.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	b, _ := voucherstatus.ValidAccountCashVoucherStatusString(req.Status)
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

	err = module.StatusModify(ctx, req, optUserID)
	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, consts.Success)
}
