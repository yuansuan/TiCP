package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	logmiddleware "github.com/yuansuan/ticp/common/go-kit/logging/middleware"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/account"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/accountcashvoucher"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/apiuser"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/api/v1/cashvoucher"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/middleware"
)

// Init ...
func Init(drv *http.Driver) {
	logging.Default().Info("setup router")
	// add request id to context and response header
	drv.Use(middleware.RequestIDMiddleware).Use(
		logmiddleware.IngressLogger(logmiddleware.IngressLoggerConfig{
			IsLogRequestHeader:  true,
			IsLogRequestBody:    true,
			IsLogResponseHeader: true,
			IsLogResponseBody:   true,
		}))

	group := drv.Group("/internal")
	{
		accountGroup := group.Group("/accounts", AccountIDMiddleware())
		accountGroup.POST("", account.Create)
		accountGroup.GET("", account.AccountList)
		accountGroup.PATCH("/:AccountID/credit", account.CreditAdd)
		accountGroup.PATCH("/:AccountID/freezedAmount", account.PaymentReduce)
		accountGroup.GET("/:AccountID/billing", account.BillList)
		accountGroup.PATCH("/:AccountID/normalbalance", account.AccountIDReduce)
		accountGroup.GET("/:AccountID", account.AccountIDGet)
		accountGroup.PATCH("/:AccountID", account.AccountFreezeModify)
		accountGroup.PATCH("/:AccountID/creditquota", account.CreditQuotaModify)
		accountGroup.PATCH("/:AccountID/frozenamount", account.PaymentFreezeUnfreeze)
		accountGroup.PATCH("/:AccountID/refundamount", account.AmountRefund)
	}

	{
		userRouterGroup := group.Group("/users", UserIDMiddleware())
		userRouterGroup.PATCH("/:UserID/normalbalance", account.AccountYsiDReduce)
		userRouterGroup.GET("/:UserID/account", account.AccountYsIDGet)
	}

	{
		cashVoucherGroup := group.Group("/cashvouchers", CashVoucherIDMiddleware())
		cashVoucherGroup.GET("", cashvoucher.List)
		cashVoucherGroup.POST("", cashvoucher.Add)
		cashVoucherGroup.PATCH("/:CashVoucherID/availabilitystatus", cashvoucher.AvailabilityModify)
		cashVoucherGroup.GET("/:CashVoucherID", cashvoucher.Get)
	}

	{
		accountCashVoucherGroup := group.Group("/accountcashvouchers", AccountCashVoucherIDMiddleware())
		accountCashVoucherGroup.GET("", accountcashvoucher.List)
		accountCashVoucherGroup.POST("", accountcashvoucher.Add)
		accountCashVoucherGroup.GET("/:AccountCashVoucherID", accountcashvoucher.Get)
		accountCashVoucherGroup.PATCH("/:AccountCashVoucherID/status", accountcashvoucher.StatusModify)
	}

	apiGroup := drv.Group("/api")
	{
		userAccounts := apiGroup.Group("/accounts")
		userAccounts.GET("/account", apiuser.AccountUserIdGet)
		userAccounts.GET("/billing", apiuser.AccountUserBillList)
		userAccounts.GET("/resourcebilling", apiuser.ResourceBillList)
	}
}

func ParamValueMiddleware(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramValue := c.Param(paramKey)
		if strings.HasPrefix(paramValue, ":") {
			paramValue = paramValue[1:]
			c.Params = gin.Params{gin.Param{Key: paramKey, Value: paramValue}}
		}
		c.Next()
	}
}

func AccountIDMiddleware() gin.HandlerFunc {
	return ParamValueMiddleware(consts.ACCONT_ID_KEY)
}

func UserIDMiddleware() gin.HandlerFunc {
	return ParamValueMiddleware(consts.ACCOUNT_USER_ID_KEY)
}

func CashVoucherIDMiddleware() gin.HandlerFunc {
	return ParamValueMiddleware(consts.CASH_VOUCHER_ID)
}

func AccountCashVoucherIDMiddleware() gin.HandlerFunc {
	return ParamValueMiddleware(consts.Account_CashVoucher_ID_KEY)
}
