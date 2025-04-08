package application

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/quotaadd"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/quotadelete"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/validation"
)

func handleValidationError(c *gin.Context, err error, req any, handlers ...validation.FieldErrorHandler) {
	err = validation.HandleError(err, req)
	if err != nil {
		var fe validation.Error
		if errors.As(err, &fe) && fe.HandleFieldError(handlers...) {
			return
		}

		common.InvalidParams(c, "invalid params, "+err.Error())
	}
}

func handleUserIDRequired(c *gin.Context, fe validation.Error) bool {
	if fe.Tag() == "required" && fe.Field() == "UserID" {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, "invalid params, "+fe.Error())
		return true
	}
	return false
}

func validateQuotaParams(c *gin.Context, logger *logging.Logger, appIDStr, userIDStr string) (snowflake.ID, snowflake.ID, bool) {
	if AppIDErr(c, appIDStr) {
		return 0, 0, false
	}

	if userIDStr == "" {
		logger.Infof("invalid params, %s", "userid is empty")
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, "invalid params, empty userid")
		return 0, 0, false
	}

	appID, err := snowflake.ParseString(appIDStr)
	if err != nil {
		logger.Infof("invalid params, %s", err.Error())
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, "invalid params, invalid appid")
		return 0, 0, false
	}

	userID, err := snowflake.ParseString(userIDStr)
	if err != nil {
		logger.Infof("invalid params, %s", err.Error())
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, "invalid params, invalid userid")
		return 0, 0, false
	}

	return appID, userID, true
}

func handleQuotaError(c *gin.Context, logger *logging.Logger, err error) {
	if errors.Is(err, common.ErrAppIDNotFound) {
		logger.Infof("app not found, %s", err.Error())
		common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, "app not found")
		return
	}
	if errors.Is(err, common.ErrAppQuotaNotFound) {
		logger.Infof("quota not found, %s", err.Error())
		common.ErrorResp(c, http.StatusNotFound, api.AppQuotaNotFound, "quota not found")
		return
	}

	if errors.Is(err, common.ErrAppQuotaAlreadyExist) {
		logger.Infof("quota already exist, %s", err.Error())
		common.ErrorResp(c, http.StatusConflict, api.AppQuotaAlreadyExist, "quota already exist")
		return
	}
	if errors.Is(err, common.ErrUserNotExists) {
		logger.Infof("user not exists, %s", err.Error())
		common.ErrorResp(c, http.StatusNotFound, api.UserNotExistsErrorCode, "user not exists")
		return
	}

	if err != nil {
		logger.Warnf("internal server error, %s", err.Error())
		common.InternalServerError(c, "internal server error")
		return
	}
}

// QuotaGet 获取app的某个user配额
func (appController *Controller) QuotaGet(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.QuotaGet")
	appIDStr := c.Param("AppID")
	userIDStr := c.Query("UserID")

	logger = logger.With("AppID", appIDStr, "UserID", userIDStr)

	appID, userID, valid := validateQuotaParams(c, logger, appIDStr, userIDStr)
	if !valid {
		return
	}

	aq, err := appController.srv.AppsQuota().GetQuota(c, appID, userID)
	if err != nil {
		handleQuotaError(c, logger, err)
		return
	}

	common.SuccessResp(c, aq)
}

// QuotaAdd 添加一个app的user配额
func (appController *Controller) QuotaAdd(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.QuotaAdd")
	appIDStr := c.Param("AppID")

	req := quotaadd.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Infof("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleUserIDRequired(c, fe)
		})
		return
	}

	userIDStr := req.UserID

	logger = logger.With("AppID", appIDStr, "UserID", userIDStr)
	appID, userID, valid := validateQuotaParams(c, logger, appIDStr, userIDStr)
	if !valid {
		return
	}

	aq, err := appController.srv.AppsQuota().AddQuota(c, appID, userID)
	if err != nil {
		handleQuotaError(c, logger, err)
		return
	}

	common.SuccessResp(c, aq)
}

// QuotaDelete 删除一个app的user配额
func (appController *Controller) QuotaDelete(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.QuotaDelete")
	appIDStr := c.Param("AppID")

	req := quotadelete.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Infof("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleUserIDRequired(c, fe)
		})
		return
	}

	userIDStr := req.UserID

	logger = logger.With("AppID", appIDStr, "UserID", userIDStr)
	appID, userID, valid := validateQuotaParams(c, logger, appIDStr, userIDStr)
	if !valid {
		return
	}

	err := appController.srv.AppsQuota().DeleteQuota(c, appID, userID)
	if err != nil {
		handleQuotaError(c, logger, err)
		return
	}

	common.SuccessResp(c, nil)
}
