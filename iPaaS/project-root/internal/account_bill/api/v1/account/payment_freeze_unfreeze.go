package account

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentfreezeunfreeze"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// PaymentFreezeUnfreeze ...
func PaymentFreezeUnfreeze(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	accountID := ctx.Param(consts.ACCONT_ID_KEY)
	if !validator.ValidAccountID(ctx, accountID) {
		return
	}

	req := &paymentfreezeunfreeze.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if !validator.ValidTradeID(ctx, req.TradeID) ||
		!validator.ValidComment(ctx, req.Comment) {
		return
	}

	// 冻结验证金额
	if req.IsFreezed && !validator.ValidAmount(ctx, req.Amount) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.PaymentFreezeUnfreeze(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
