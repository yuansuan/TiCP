package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	creditadd "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditadd"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// CreditAdd ...
func CreditAdd(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	accountID := ctx.Param(consts.ACCONT_ID_KEY)

	if !validator.ValidAccountID(ctx, accountID) {
		return
	}

	req := &creditadd.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if req.DeltaAwardBalance < 0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "DeltaAwardBalance can not be less than 0")
		return
	}

	if req.DeltaNormalBalance < 0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "DeltaNormalBalance can not be less than 0")
		return
	}

	if !validator.ValidComment(ctx, req.Comment) ||
		!validator.ValidTradeID(ctx, req.TradeId) ||
		!validator.ValidIdempotentID(ctx, req.IdempotentID) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	accountResp, err := module.CreditAdd(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}
	common.SuccessResp(ctx, accountResp)
}
