package account

import (
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/amountrefund"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// AmountRefund ...
func AmountRefund(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	accountID := ctx.Param(consts.ACCONT_ID_KEY)

	if !validator.ValidAccountID(ctx, accountID) {
		return
	}

	req := &amountrefund.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if !validator.ValidAmount(ctx, req.Amount) ||
		!validator.ValidRefundID(ctx, req.RefundID) ||
		(!util.IsBlank(req.ResourceID) && !validator.ValidResourceID(ctx, req.ResourceID)) {
		return
	}

	if len(req.Comment) > 256 {
		msg := fmt.Sprintf("invalid params, commnet length can not be exceed 256 characters")
		logging.GetLogger(ctx).Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	accountResp, err := module.AmountRefund(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}
	common.SuccessResp(ctx, accountResp)
}
