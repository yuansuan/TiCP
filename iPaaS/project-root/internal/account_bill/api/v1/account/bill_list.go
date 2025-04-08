package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/billlist"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func BillList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &billlist.Request{}
	err := ctx.BindQuery(req)
	if err != nil {
		msg := fmt.Sprintf("invalid params, %v", err)
		logger.Error(msg)
		common.InvalidParams(ctx, msg)
		return
	}

	if req.PageIndex < 1 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidPageIndexErrorCode, "index can not be less than one")
		return
	}

	if req.PageSize < 1 || req.PageSize > 1000 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidPageSizeErrorCode, "page size should be in [1, 1000]")
		return
	}

	if !validator.ValidStartTimeAndEndTime(ctx, req.StartTime, req.EndTime) {
		return
	}

	optUserID, b := validator.ValidAuthUserID(ctx)
	if !b {
		return
	}

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	if req.AccountID != "" {
		if !validator.ValidAccountID(ctx, req.AccountID) {
			return
		}

		accountIDSno := snowflake.MustParseString(req.AccountID)
		_, err = module.IsAccountExists(ctx, accountIDSno, sess)
		if errors.Is(err, common.ErrAccountNotExists) {
			// 资金账户不存在则查询ysid
			resp, err := module.AccountGetByYsID(ctx, &ysidget.Request{
				UserID: req.AccountID,
			}, optUserID)

			if err == nil {
				req.AccountID = resp.AccountID
			} else if !validator.ErrJudge(ctx, err) {
				return
			}
		}
	}

	resp, err := module.BillList(ctx, req, optUserID)

	if !validator.ErrJudge(ctx, err) {
		return
	}

	common.SuccessResp(ctx, resp)
}
