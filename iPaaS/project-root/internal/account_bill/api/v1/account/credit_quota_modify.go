package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditquotamodify"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// CreditQuotaModify ...
func CreditQuotaModify(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	accountID := ctx.Param(consts.ACCONT_ID_KEY)
	if !validator.ValidAccountID(ctx, accountID) {
		return
	}

	req := &creditquotamodify.Request{}
	err := ctx.BindJSON(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if req.CreditQuotaAmount < 0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "CreditQuotaAmount should not be less than 0")
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	resp, err := module.CreditQuotaModify(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
