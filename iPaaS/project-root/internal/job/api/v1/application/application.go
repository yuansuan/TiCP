package application

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/utils"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/add"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/list"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

// Get ...
func (appController *Controller) Get(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.GetQuota")
	value := c.Param("AppID")
	logger = logger.With("AppID", value)
	if AppIDErr(c, value) {
		return
	}
	appID := snowflake.MustParseString(value)
	getApp, err := appController.srv.Apps().GetApp(c, appID)
	if NotExistErr(c, err) {
		return
	}
	if err != nil {
		logger.Warnf("get app failed, err:%v", err)
		common.InternalServerError(c, "internal server error")
		return
	}
	application := convertApp2Application(getApp)
	application.NeedLimitCore = getApp.NeedLimitCore
	common.SuccessResp(c, application)
}

// List ...
func (appController *Controller) List(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.List")

	userID, ok := checkListParam(c)
	if !ok {
		return
	}
	// 只返回已发布的应用列表
	publish := consts.PublishStatusPublished

	apps, _, err := appController.srv.Apps().ListApps(c, userID, publish)
	if err != nil {
		logger.Warnf("get app failed, err:%v", err)
		common.InternalServerError(c, "internal server error")
		return
	}
	var resApps []*schema.Application
	for _, app := range apps {
		application := convertApp2Application(app)
		resApps = append(resApps, application)
	}
	common.SuccessResp(c, resApps)
}

func checkListParam(c *gin.Context) (snowflake.ID, bool) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, "Invalid user id, "+err.Error())
		return 0, false
	}

	return userID, true
}

// AdminList ...
func (appController *Controller) AdminList(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.List")

	req := &list.Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		logger.Infof("get app failed, err:%v", err)
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}

	userIDStr := req.AllowUserID // 查询该用户下的有配额的应用列表
	userID, err := snowflake.ParseString(userIDStr)
	if err != nil {
		// 参数错误，用户ID不合法
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidUserID, "invalid params, userID")
		return
	}
	// 返回已发布和未发布的应用列表
	publish := consts.PublishStatusAll

	apps, _, err := appController.srv.Apps().ListApps(c, userID, publish)
	if err != nil {
		logger.Warnf("get app failed, err:%v", err)
		common.InternalServerError(c, "internal server error")
		return
	}
	var resApps []*schema.Application
	for _, app := range apps {
		application := convertApp2Application(app)
		application.NeedLimitCore = app.NeedLimitCore
		// only for admin
		resApps = append(resApps, application)
	}
	common.SuccessResp(c, resApps)
}

// Create ...
func (appController *Controller) Create(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.Create")
	req := &add.Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Infof("create app failed, err:%v", err)
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}
	if req.Command == "" {
		req.Command = `#YS_COMMAND_PREPARED`
	}
	if !checkAddParam(c, req) {
		return
	}

	logger = logger.With("Name", req.Name, "Type", req.Type, "Version", req.Version, "Command", req.Command)
	id, err := appController.srv.Apps().AddApp(c, req)
	if err != nil {
		if DuplicateErr(c, err) {
			return
		}
		logger.Warnf("create app failed, err:%v", err)
		common.InternalServerError(c, err.Error())
		return
	}
	res := new(add.Data)
	res.AppID = id.String()
	common.SuccessResp(c, res)
}

func checkAddParam(c *gin.Context, req *add.Request) bool {
	err := validateName(req.Name)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentName, "invalid params, "+err.Error())
		return false
	}

	err = validateType(req.Type)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentType, "invalid params, "+err.Error())
		return false
	}

	err = validateVersion(req.Version)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentVersion, "invalid params, "+err.Error())
		return false
	}

	err = validateDescription(req.Description)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDescription, "invalid params, "+err.Error())
		return false
	}

	err = validateCommand(req.Command)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCommand, "invalid params, "+err.Error())
		return false
	}

	zones := config.GetConfig().Zones
	err = validateBinPathAndImage(req.BinPath, req.Image, zones)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentBinPath, "invalid params, "+err.Error())
		return false
	}

	err = validateExtentionParams(req.ExtentionParams)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentExtentionParams, "invalid params, "+err.Error())
		return false
	}

	err = validateLicManagerID(req.LicManagerId)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentLicManagerId, "invalid params, "+err.Error())
		return false
	}

	err = validateResidualLogParser(req.ResidualLogParser)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentResidualLogParser, "invalid params, "+err.Error())
		return false
	}

	err = validateMonitorChartParser(req.MonitorChartParser)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentMonitorChartParser, "invalid params, "+err.Error())
		return false
	}

	err = validateSpecifyQueue(req.SpecifyQueue, zones)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentSpecifyQueue, "invalid params, "+err.Error())
		return false
	}

	return true
}

func (appController *Controller) Update(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.Update")
	value := c.Param("AppID")
	if AppIDErr(c, value) {
		return
	}
	req := &update.Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Infof("update app failed, err:%v", err)
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}
	if !checkUpdateParam(c, req) {
		return
	}
	logger = logger.With("AppID", req.AppID, "Name", req.Name, "Type", req.Type, "Version", req.Version)
	req.AppID = value
	err := appController.srv.Apps().UpdateApp(c, req)
	if err != nil {
		if DuplicateErr(c, err) {
			return
		}
		if NoEffectErr(c, err) {
			return
		}
		logger.Warnf("update app failed, err:%v", err)
		common.InternalServerError(c, "internal server error")
		return
	}
	common.SuccessResp(c, nil)
}

func checkUpdateParam(c *gin.Context, req *update.Request) bool {
	err := validateName(req.Name)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentName, "invalid params, "+err.Error())
		return false
	}

	err = validateType(req.Type)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentType, "invalid params, "+err.Error())
		return false
	}

	err = validateVersion(req.Version)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentVersion, "invalid params, "+err.Error())
		return false
	}

	err = validateDescription(req.Description)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentDescription, "invalid params, "+err.Error())
		return false
	}

	err = validateCommand(req.Command)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentCommand, "invalid params, "+err.Error())
		return false
	}

	zones := config.GetConfig().Zones
	err = validateBinPathAndImage(req.BinPath, req.Image, zones)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentBinPath, "invalid params, "+err.Error())
		return false
	}

	err = validateExtentionParams(req.ExtentionParams)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentExtentionParams, "invalid params, "+err.Error())
		return false
	}

	err = validatePublishStatus(req.PublishStatus)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentPublishStatus, "invalid params, "+err.Error())
		return false
	}

	err = validateLicManagerID(req.LicManagerId)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentLicManagerId, "invalid params, "+err.Error())
		return false
	}

	err = validateResidualLogParser(req.ResidualLogParser)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentResidualLogParser, "invalid params, "+err.Error())
		return false
	}

	err = validateMonitorChartParser(req.MonitorChartParser)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentMonitorChartParser, "invalid params, "+err.Error())
		return false
	}

	return true
}

// Delete ...
func (appController *Controller) Delete(c *gin.Context) {
	logger := logging.GetLogger(c).With("func", "application.Delete")
	value := c.Param("AppID")
	if AppIDErr(c, value) {
		return
	}
	appID := snowflake.MustParseString(value)
	err := appController.srv.Apps().DeleteApp(c, appID)
	if err != nil {
		if NotExistErr(c, err) {
			return
		}
		logger.Warnf("delete app failed, err:%v", err)
		common.InternalServerError(c, "internal server error")
		return
	}
	common.SuccessResp(c, nil)
}

func (appController *Controller) InvalidAppID(c *gin.Context) {
	common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, "invalid params, appID")
}

func NotExistErr(c *gin.Context, err error) bool {
	if errors.Is(err, xorm.ErrNotExist) {
		common.ErrorResp(c, http.StatusNotFound, api.AppIdNotFound, "app not found")
		return true
	}
	return false
}

func NoEffectErr(c *gin.Context, err error) bool {
	if errors.Is(err, common.ErrNoEffect) {
		common.ErrorResp(c, http.StatusOK, api.NoEffect, "no effect")
		return true
	}
	return false
}

func DuplicateErr(c *gin.Context, err error) bool {
	if errors.Is(err, common.ErrDuplicateEntry) {
		common.InvalidParams(c, "invalid params, type and version duplicate")
		return true
	}
	return false
}

func AppIDErr(c *gin.Context, value string) bool {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, "invalid params, empty appID")
		return true
	}
	if _, err := snowflake.ParseString(value); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidAppID, "invalid params, appID")
		return true
	}
	return false
}

// inside to outside
func convertApp2Application(app *models.Application) *schema.Application {
	return &schema.Application{
		AppID:              app.ID.String(),
		Name:               app.Name,
		Type:               app.Type,
		Version:            app.Version,
		AppParamsVersion:   app.AppParamsVersion,
		Image:              app.Image,
		Endpoint:           app.Endpoint,
		Command:            app.Command,
		PublishStatus:      app.PublishStatus,
		Description:        app.Description,
		IconUrl:            app.IconUrl,
		CoresMaxLimit:      app.CoresMaxLimit,
		CoresPlaceholder:   app.CoresPlaceholder,
		FileFilterRule:     app.FileFilterRule,
		ResidualEnable:     app.ResidualEnable,
		ResidualLogRegexp:  app.ResidualLogRegexp,
		ResidualLogParser:  app.ResidualLogParser,
		MonitorChartEnable: app.MonitorChartEnable,
		MonitorChartRegexp: app.MonitorChartRegexp,
		MonitorChartParser: app.MonitorChartParser,
		SnapshotEnable:     app.SnapshotEnable,
		BinPath:            app.BinPath,
		ExtentionParams:    app.ExtentionParams,
		LicManagerId:       app.LicManagerId.String(),
		SpecifyQueue:       util.ToStringMap(app.SpecifyQueue),
		CreateTime:         util.ModelTimeToString(app.CreateTime),
		UpdateTime:         util.ModelTimeToString(app.UpdateTime),
	}
}
