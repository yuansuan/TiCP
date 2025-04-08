package accountcashvoucher

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/add"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func Add(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &add.Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if !validator.ValidCashVoucherID(ctx, req.CashVoucherID) {
		return
	}

	if !validator.ValidAccountIDs(ctx, req.AccountIDs) {
		return
	}
	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	err := module.Add(ctx, req, optUserID)
	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, consts.Success)
}
