package application

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

func validateAllowParams(c *gin.Context, logger *logging.Logger, appIDStr string) (snowflake.ID, bool) {
	if AppIDErr(c, appIDStr) {
		return 0, false
	}

	appID, err := snowflake.ParseString(appIDStr)
	if err != nil {
		logger.Infof("invalid params, %s", err.Error())
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, "invalid params, invalid appid")
		return 0, false
	}
	return appID, true
}

func handleAllowError(c *gin.Context, logger *logging.Logger, err error) {
	if errors.Is(err, common.ErrAppIDNotFound) {
		logger.Infof("app not found, %s", err.Error())
		common.ErrorResp(c, http.StatusNotFound, api.AppIDNotFoundErrorCode, "app not found")
		return
	}
	if errors.Is(err, common.ErrAppAllowNotFound) {
		logger.Infof("allow not found, %s", err.Error())
		common.ErrorResp(c, http.StatusNotFound, api.AppAllowNotFound, "allow not found")
		return
	}
	if errors.Is(err, common.ErrAppAllowAlreadyExist) {
		logger.Infof("allow already exist, %s", err.Error())
		common.ErrorResp(c, http.StatusConflict, api.AppAllowAlreadyExist, "allow already exist")
		return
	}

	if err != nil {
		logger.Warnf("internal server error, %s", err.Error())
		common.InternalServerError(c, "internal server error")
		return
	}
}

// AllowGet 获取app的白名单信息
func (appController *Controller) AllowGet(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.AllowGet")
	appIDStr := c.Param("AppID")

	logger = logger.With("AppID", appIDStr)
	appID, valid := validateAllowParams(c, logger, appIDStr)
	if !valid {
		return
	}

	allow, err := appController.srv.AppsAllow().GetAllow(c, appID)
	if err != nil {
		handleAllowError(c, logger, err)
		return
	}

	common.SuccessResp(c, allow)
}

// AllowAdd 添加一个app到白名单
func (appController *Controller) AllowAdd(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.AllowAdd")
	appIDStr := c.Param("AppID")

	logger = logger.With("AppID", appIDStr)
	appID, valid := validateAllowParams(c, logger, appIDStr)
	if !valid {
		return
	}

	aq, err := appController.srv.AppsAllow().AddAllow(c, appID)
	if err != nil {
		handleAllowError(c, logger, err)
		return
	}

	common.SuccessResp(c, aq)
}

// AllowDelete 从白名单中删除一个app
func (appController *Controller) AllowDelete(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.AllowDelete")
	appIDStr := c.Param("AppID")

	logger = logger.With("AppID", appIDStr)
	appID, valid := validateAllowParams(c, logger, appIDStr)
	if !valid {
		return
	}

	err := appController.srv.AppsAllow().DeleteAllow(c, appID)
	if err != nil {
		handleAllowError(c, logger, err)
		return
	}

	common.SuccessResp(c, nil)
}
