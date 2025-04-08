package account

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	accountCreate "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/create"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accounttype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// Create ...
func Create(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &accountCreate.Request{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		msg := fmt.Sprintf("create account failed, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if util.IsBlank(req.AccountName) {
		msg := fmt.Sprintf("create account failed, account name can not be empty!")
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	accountType := accounttype.AccountType(req.AccountType)
	if !accounttype.ValidAccountType(accountType) {
		msg := fmt.Sprintf("create account failed, account type invalid!")
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if accountType == accounttype.PERSONAL &&
		util.IsBlank(req.UserID) {
		msg := fmt.Sprintf("create account failed, create user account need a user id!")
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.Create(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)

}
