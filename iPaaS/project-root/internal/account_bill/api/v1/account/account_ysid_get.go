package account

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AccountYsIDGet(ctx *gin.Context) {
	userID := ctx.Param(consts.ACCOUNT_USER_ID_KEY)
	if !validator.ValidUserID(ctx, userID) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.AccountGetByYsID(ctx, &ysidget.Request{
		UserID: userID,
	}, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
