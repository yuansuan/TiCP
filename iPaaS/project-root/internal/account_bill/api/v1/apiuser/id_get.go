package apiuser

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func AccountUserIdGet(ctx *gin.Context) {
	userID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}
	resp, err := module.AccountGetByUserID(ctx, userID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
