package account

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidreduce"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AccountYsiDReduce(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	userID := ctx.Param(consts.ACCOUNT_USER_ID_KEY)
	if !validator.ValidUserID(ctx, userID) {
		return
	}

	req := &ysidreduce.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if !validator.ValidAmount(ctx, req.Amount) ||
		!validator.ValidComment(ctx, req.Comment) ||
		!validator.ValidTradeID(ctx, req.TradeID) ||
		!validator.ValidIdempotentID(ctx, req.IdempotentID) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.AccountYsIDReduce(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
