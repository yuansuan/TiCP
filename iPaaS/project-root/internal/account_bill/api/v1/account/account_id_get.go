package account

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AccountIDGet(ctx *gin.Context) {
	accountID := ctx.Param(consts.ACCONT_ID_KEY)
	if !validator.ValidAccountID(ctx, accountID) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.AccountGetByID(ctx, &idget.Request{
		AccountID: accountID,
	}, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
