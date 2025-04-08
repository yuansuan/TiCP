package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	licensetype "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info/type"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager/os"
	"github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager/publish"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/validation"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao"
	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LicenseHandler struct {
	licenseManagerDao dao.LicenseManagerDao
}

func NewLicenseHandler(licenseManagerDao dao.LicenseManagerDao) *LicenseHandler {
	return &LicenseHandler{
		licenseManagerDao: licenseManagerDao,
	}
}

func (l *LicenseHandler) ListLicenseManage(ctx *gin.Context) {
	// 不分页，全拉
	entities, _, err := l.licenseManagerDao.ListLicenseManagers(ctx, nil)
	if err != nil {
		logging.GetLogger(ctx).Warnf("ListLicenseManagersFail, Error: %s", err.Error())
		common.InternalServerError(ctx, "List Error")
		return
	}
	res := licmanager.ListLicManagerResponseData{
		Items: []*licmanager.GetLicManagerResponseData{},
	}
	lmList := dbModels.ToLicenseManagerExt(entities)

	for _, v := range lmList {
		licenseManagerRespData, err := toRespLicenseManager(v)
		if err != nil {
			errMessage := fmt.Sprintf("failed to convert license manager to proto, Error: %v", err)
			logging.GetLogger(ctx).Warnf(errMessage)
			common.InternalServerError(ctx, errMessage)
			return
		}
		res.Items = append(res.Items, licenseManagerRespData)
	}
	res.Total = len(res.Items)

	common.SuccessResp(ctx, &res)
}

func (l *LicenseHandler) GetLicenseManage(ctx *gin.Context) {
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}
	lMgr, err := l.licenseManagerDao.GetLicenseManager(ctx, id)
	if err != nil {
		logging.GetLogger(ctx).Warnf("GetLicenseManagerFail, Error: %s", err.Error())
		common.InternalServerError(ctx, "GetQuota Fail")
		return
	}
	if lMgr == nil {
		common.ErrorResp(ctx, 404, "LicenseManagerIdNotFound", fmt.Sprintf("%s not found", ctx.Param("id")))
		return
	}
	licenseManagerRespData, err := toRespLicenseManager(lMgr)
	if err != nil {
		errMessage := fmt.Sprintf("failed to convert license manager to proto, Error: %v", err)
		logging.GetLogger(ctx).Warnf(errMessage)
		common.InternalServerError(ctx, errMessage)
		return
	}

	common.SuccessResp(ctx, &licenseManagerRespData)
}

func (l *LicenseHandler) AddLicenseManage(ctx *gin.Context) {
	var req = licmanager.AddLicManagerRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}
	id, ok := genSnowFlakId(ctx)
	if !ok {
		return
	}

	os, err := os.ToOS(req.Os)
	if err != nil {
		errMessage := fmt.Sprintf("invalid os params, err: %v", err)
		logging.GetLogger(ctx).Error(errMessage)
		common.InvalidParams(ctx, errMessage)
		return
	}

	lm := &dbModels.LicenseManager{
		Id:          id,
		AppType:     req.AppType,
		Os:          os,
		Description: req.Desc,
		ComputeRule: req.ComputeRule,
		Status:      2, // 未发布
	}
	err = l.licenseManagerDao.AddLicenseManager(ctx, lm)
	if err != nil {
		logging.GetLogger(ctx).Warnf("add to db fail, error: %v", err)
		common.InternalServerError(ctx, "add license manger fail")
		return
	}
	res := licmanager.AddLicManagerResponseData{Id: id.String()}
	common.SuccessResp(ctx, &res)
}

func (l *LicenseHandler) PutLicenseManager(ctx *gin.Context) {
	var req = licmanager.PutLicManagerRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	lm, err := convertLicenseManagerModelFromPutReq(ctx, &req)
	if err != nil {
		errorMessage := fmt.Sprintf("convert to license manager fail, err: %v", err)
		logging.GetLogger(ctx).Errorf(errorMessage)
		common.InternalServerError(ctx, errorMessage)
		return
	}

	suc, err := l.licenseManagerDao.UpdateLicenseManager(ctx, lm)
	if err != nil {
		logging.GetLogger(ctx).Warnf("UpdateFail, Error: %s, RequestId: %s", err.Error(), common.GetRequestID(ctx))
		common.InternalServerError(ctx, "update fail")
		return
	}
	if !suc {
		common.ErrorResp(ctx, 404, "LicenseManagerIdNotFound", fmt.Sprintf("%s not found", ctx.Param("id")))
		return
	}
	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) DeleteLicenseManager(ctx *gin.Context) {
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}
	suc, canDelete, err := l.licenseManagerDao.DeleteLicenseManager(ctx, id)
	if err != nil {
		logging.GetLogger(ctx).Warnf("DeleteLicenseManagerFail, Error: %s, Request: %s", err.Error(), common.GetRequestID(ctx))
		common.InternalServerError(ctx, "Delete license manager fail")
		return
	}
	if !canDelete {
		common.ErrorResp(ctx, 403, "Forbidden", "This license manager still has license providers.")
		return
	}
	if !suc {
		common.ErrorResp(ctx, 404, "LicenseManagerIdNotFound", fmt.Sprintf("%s not found", ctx.Param("id")))
		return
	}
	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) AddLicense(ctx *gin.Context) {
	req := licenseinfo.AddLicenseInfoRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	if err := l.validateAddLicenseRequest(req); err != nil {
		common.ErrorResp(ctx, http.StatusBadRequest, err.ErrorCode, err.ErrorMsg)
		return
	}

	licInfo, err := convertLicenseModelFromAddReq(ctx, &req)
	if err != nil {
		errMessage := fmt.Sprintf("To license info dao fail, Error: %s, RequestId: %s", err.Error(), common.GetRequestID(ctx))
		logging.GetLogger(ctx).Warnf(errMessage)
		common.InternalServerError(ctx, errMessage)
		return
	}

	// 检查manager是否存在
	// TODO 最好做成事务，不然可能产生垃圾数据
	licM, err := l.licenseManagerDao.GetLicenseManager(ctx, licInfo.ManagerId)
	if err != nil {
		common.InternalServerError(ctx, "GetQuota license manager fail")
		return
	}
	if licM == nil {
		common.ErrorResp(ctx, 403, "Forbidden", fmt.Sprintf("license manager %s not found", licInfo.ManagerId))
		return
	}

	id, ok := genSnowFlakId(ctx)
	if !ok {
		return
	}
	licInfo.Id = id
	err = l.licenseManagerDao.AddLicenseInfo(ctx, licInfo)
	if err != nil {
		logging.GetLogger(ctx).Warnf("AddLicenseInfoFail, Error: %s, RequestId: %s", err.Error(), common.GetRequestID(ctx))
		common.InternalServerError(ctx, "Add license fail")
		return
	}
	res := licenseinfo.AddLicenseInfoResponseData{
		Id: id.String(),
	}
	common.SuccessResp(ctx, &res)
}

func (l *LicenseHandler) DeleteLicense(ctx *gin.Context) {
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}
	suc, canDelete, err := l.licenseManagerDao.DeleteLicenseInfo(ctx, id)
	if err != nil {
		common.InternalServerError(ctx, "Delete provider fail")
		return
	}
	if !canDelete {
		common.ErrorResp(ctx, 403, "Forbidden", "License is using")
		return
	}
	if !suc {
		common.ErrorResp(ctx, 404, "LicenseInfoIdNotFound", fmt.Sprintf("%s not found", ctx.Param("id")))
		return
	}
	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) PutLicense(ctx *gin.Context) {
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}
	req := licenseinfo.PutLicenseInfoRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	if err := l.validatePutLicenseRequest(req); err != nil {
		common.ErrorResp(ctx, http.StatusBadRequest, err.ErrorCode, err.ErrorMsg)
		return
	}

	licInfo, err := convertLicenseModelFromPutReq(ctx, &req)
	if err != nil {
		errMessage := fmt.Sprintf("To license info dao fail, Error: %s, RequestId: %s", err.Error(), common.GetRequestID(ctx))
		logging.GetLogger(ctx).Warn(errMessage)
		common.InternalServerError(ctx, errMessage)
		return
	}

	suc, err := l.licenseManagerDao.UpdateLicenseInfo(ctx, licInfo)
	if err != nil {
		common.InternalServerError(ctx, "Update license info fail")
		return
	}

	if !suc {
		common.ErrorResp(ctx, 404, "LicenseInfoIdNotFound", fmt.Sprintf("%s not found", id.String()))
		return
	}

	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) ListModules(ctx *gin.Context) {
	req := moduleconfig.ListModuleConfigRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	licenseID, err := snowflake.ParseString(req.LicenseId)
	if err != nil {
		logging.GetLogger(ctx).Errorf("parse license id err, err: %v", err)
		common.InvalidParams(ctx, fmt.Sprintf("Invalid license id: %s", req.LicenseId))
	}

	moduleConfigs, err := l.licenseManagerDao.ListModuleConfig(ctx, licenseID)
	if err != nil {
		logging.GetLogger(ctx).Errorf("list module config err, err: %v", err)
		common.InternalServerError(ctx, "List module config fail")
		return
	}

	res := moduleconfig.ListModuleConfigResponseData{}
	for _, moduleConfig := range moduleConfigs {
		moduleConfigRespData := toRespModuleConfig(moduleConfig)
		if err != nil {
			errMessage := fmt.Sprintf("failed to convert module config to proto, error: %v", err)
			logging.GetLogger(ctx).Errorf(errMessage)
			common.InternalServerError(ctx, errMessage)
			return
		}
		res.ModuleConfigs = append(res.ModuleConfigs, moduleConfigRespData)
	}

	common.SuccessResp(ctx, res)
}

func (l *LicenseHandler) AddModule(ctx *gin.Context) {
	req := moduleconfig.AddModuleConfigRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	licId, err := snowflake.ParseString(req.LicenseId)
	if err != nil {
		common.InvalidParams(ctx, fmt.Sprintf("Invalid license id: %s", req.LicenseId))
		return
	}
	// TODO 查询license和add module需要做成事务，不做事务最话情况是产生一条垃圾module数据
	exist, _, err := l.licenseManagerDao.GetLicenseInfoByID(ctx, licId)
	if err != nil {
		common.InternalServerError(ctx, "GetQuota license info fail")
		return
	}
	if !exist {
		common.ErrorResp(ctx, 404, "LicenseIdNotFound", fmt.Sprintf("License id %s not exist", licId))
		return
	}
	id, ok := genSnowFlakId(ctx)
	if !ok {
		return
	}
	cfg := dbModels.ModuleConfig{
		Id:         id,
		LicenseId:  licId,
		Total:      req.Total,
		ModuleName: req.ModuleName,
	}
	suc, err := l.licenseManagerDao.AddModuleConfig(ctx, &cfg)
	if err != nil {
		common.InternalServerError(ctx, "Add module fail")
		return
	}
	if !suc {
		common.ErrorResp(ctx, 403, "Forbidden", fmt.Sprintf("module %s existed", cfg.ModuleName))
		return
	}
	res := moduleconfig.AddModuleConfigResponseData{
		Id: id.String(),
	}
	common.SuccessResp(ctx, &res)
}

func (l *LicenseHandler) BatchAddModules(ctx *gin.Context) {
	req := moduleconfig.BatchAddModuleConfigRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	if err := l.validateBachAddModuleConfigsRequest(req); err != nil {
		common.ErrorResp(ctx, http.StatusBadRequest, err.ErrorCode, err.ErrorMsg)
		return
	}

	licenseID, err := snowflake.ParseString(req.LicenseId)
	if err != nil {
		message := fmt.Sprintf("invalid license id: %s", req.LicenseId)
		logging.Default().Errorf(message)
		common.InvalidParams(ctx, message)
		return
	}

	moduleConfigModels, err := toModuleConfigModels(licenseID, req.ModuleConfigs)
	if err != nil {
		message := fmt.Sprintf("convert to config module config models err, err: %v", err)
		logging.Default().Error(message)
		common.InternalServerError(ctx, message)
		return
	}

	// 放在事务里面，防止脏数据产生
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		exist, _, err := l.licenseManagerDao.GetLicenseInfoByID(ctx, licenseID)
		if err != nil {
			message := fmt.Sprintf("get license info err, err: %v", err)
			logging.Default().Errorf(message)
			return errors.New(message)
		}
		if !exist {
			message := fmt.Sprintf("license id not found, license id: %s", req.LicenseId)
			logging.Default().Errorf(message)
			return common.ErrLicenseIDNotFound
		}

		err = l.licenseManagerDao.BatchAddModuleConfigs(ctx, moduleConfigModels)
		if err != nil {
			message := fmt.Sprintf("batch add module configs err, err: %v", err)
			logging.Default().Error(message)
			return errors.New(message)
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, common.ErrLicenseIDNotFound) {
			common.ErrorResp(ctx, http.StatusNotFound, api.LicenseIdNotFound, err.Error())
		} else {
			common.InternalServerError(ctx, err.Error())
		}
		return
	}

	var ids []string
	for _, model := range moduleConfigModels {
		ids = append(ids, model.Id.String())
	}

	resp := moduleconfig.BatchAddModuleConfigResponseData{Ids: ids}
	common.SuccessResp(ctx, resp)
}

func (l *LicenseHandler) PutModule(ctx *gin.Context) {
	req := moduleconfig.PutModuleConfigRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err, req)
		return
	}

	moduleConfig, err := convertModelFromPutReq(ctx, &req)
	if err != nil {
		errMessage := fmt.Sprintf("convert module fail, err: %v", err)
		logging.GetLogger(ctx).Error(errMessage)
		common.InternalServerError(ctx, errMessage)
		return
	}
	suc, err := l.licenseManagerDao.UpdateModuleConfigTotal(ctx, moduleConfig)
	if err != nil {
		errorMsg := fmt.Sprintf("Update module fail, err: %v", err)
		common.InternalServerError(ctx, errorMsg)
		return
	}
	if !suc {
		common.ErrorResp(ctx, 404, "ModuleConfigIdNotFound", fmt.Sprintf("%s not found", ctx.Param("id")))
		return
	}
	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) DeleteModule(ctx *gin.Context) {
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}
	suc, canDelete, err := l.licenseManagerDao.DeleteModuleConfig(ctx, id)
	if err != nil {
		common.InternalServerError(ctx, "Delete module fail")
		return
	}
	if !suc {
		common.ErrorResp(ctx, 404, "ModuleConfigIdNotFound", fmt.Sprintf("%s not found", id.String()))
		return
	}
	if !canDelete {
		common.ErrorResp(ctx, 403, "Forbidden", "Module is using")
		return
	}
	common.SuccessResp(ctx, nil)
}

func (l *LicenseHandler) validateAddLicenseRequest(req licenseinfo.AddLicenseInfoRequest) *common.Error {
	_, _, err := validateTimeForAdd(req.BeginTime, req.EndTime)
	if err != nil {
		errMessage := fmt.Sprintf("begin: %s, end: %s, err: %v", req.BeginTime, req.EndTime, err)
		logging.Default().Error(errMessage)
		return err
	}

	return nil
}

func (l *LicenseHandler) validatePutLicenseRequest(req licenseinfo.PutLicenseInfoRequest) *common.Error {
	_, _, err := validateTimeForPut(req.BeginTime, req.EndTime)
	if err != nil {
		errMessage := fmt.Sprintf("begin: %s, end: %s, err: %v", req.BeginTime, req.EndTime, err)
		logging.Default().Error(errMessage)
		return err
	}

	return nil
}

func (l *LicenseHandler) validateBachAddModuleConfigsRequest(req moduleconfig.BatchAddModuleConfigRequest) *common.Error {
	licenseID := req.LicenseId
	for _, moduleConfig := range req.ModuleConfigs {
		if moduleConfig.LicenseId != licenseID {
			errMessage := fmt.Sprintf("license id not match, license id: %s, module config license id: %s", licenseID, moduleConfig.LicenseId)
			logging.Default().Error(errMessage)
			return &common.Error{
				ErrorCode: api.InvalidParams,
				ErrorMsg:  errMessage,
			}
		}
	}

	return nil
}

func convertLicenseManagerModelFromPutReq(ctx context.Context, req *licmanager.PutLicManagerRequest) (*dbModels.LicenseManager, error) {
	ID, err := snowflake.ParseString(req.Id)
	if err != nil {
		logging.GetLogger(ctx).Errorf("parse string fail, err: %v", err)
		return nil, err
	}

	osType, err := os.ToOS(req.Os)
	if err != nil {
		logging.GetLogger(ctx).Errorf("convert to os fail, err: %v", err)
		return nil, err
	}

	stat, err := publish.ToStatus(req.Status)
	if err != nil {
		logging.GetLogger(ctx).Errorf("convert to status fail, err: %v", err)
		return nil, err
	}

	var publishTime time.Time
	if stat.Published() {
		publishTime = time.Now()
	}

	return &dbModels.LicenseManager{
		Id:          ID,
		AppType:     req.AppType,
		Os:          osType,
		Status:      stat,
		Description: req.Desc,
		ComputeRule: req.ComputeRule,
		PublishTime: publishTime,
	}, nil
}

// stringToTimeStamp ...
func stringToTimeStamp(timeString string) (*timestamppb.Timestamp, error) {
	timeParse, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return nil, status.Errorf(consts.InvalidArgument, "The date format is incorrect, for example:2006-01-02 15:04:05")
	}
	return timestamppb.New(timeParse), nil
}

func convertLicenseModelFromAddReq(ctx *gin.Context, req *licenseinfo.AddLicenseInfoRequest) (*dbModels.LicenseInfo, error) {
	mgId, err := snowflake.ParseString(req.ManagerId)
	if err != nil {
		errMessage := fmt.Sprintf("Bad manager id: %s", req.ManagerId)
		logging.GetLogger(ctx).Error(errMessage)
		return nil, errors.New(errMessage)
	}
	auth := parseAuth(req.Auth)

	licenseType, err := licensetype.ToType(req.LicenseType)
	if err != nil {
		logging.Default().Errorf("invalid license type: %d, err: %v", req.LicenseType, err)
		return nil, err
	}

	beginTime, err := stringToTimeStamp(req.BeginTime)
	if err != nil {
		logging.Default().Errorf("invalid begin time: %s, err: %v", req.BeginTime, err)
		return nil, err
	}
	endTime, err := stringToTimeStamp(req.EndTime)
	if err != nil {
		logging.Default().Errorf("invalid end time: %s, err: %v", req.EndTime, err)
		return nil, err
	}

	var licenseAddresses string
	if req.LicenseProxies != nil && len(req.LicenseProxies) > 0 {
		content, err := json.Marshal(req.LicenseProxies)
		if err != nil {
			logging.Default().Errorf("marshal license addresses fail, err: %v", err)
			return nil, err
		}
		licenseAddresses = string(content)
	}

	licInfo := &dbModels.LicenseInfo{
		ManagerId:             mgId,
		MacAddr:               req.MacAddr,
		ToolPath:              req.ToolPath,
		Provider:              req.Provider,
		LicensePort:           req.Port,
		LicenseNum:            req.LicenseNum,
		LicenseServer:         req.LicenseEnvVar,
		LicenseUrl:            req.LicenseUrl,
		LicenseProxies:        licenseAddresses,
		LicenseType:           licenseType,
		Weight:                req.Weight,
		Auth:                  auth,
		HpcEndpoint:           req.HpcEndpoint,
		AllowableHpcEndpoints: req.AllowableHpcEndpoints,
		CollectorType:         req.CollectorType,
		BeginTime:             beginTime.AsTime(),
		EndTime:               endTime.AsTime(),
		LicenseServerStatus:   dbModels.LicenseServerStatusAbnormal,
	}
	return licInfo, nil
}

func convertLicenseModelFromPutReq(ctx context.Context, req *licenseinfo.PutLicenseInfoRequest) (*dbModels.LicenseInfo, error) {
	license := &dbModels.LicenseInfo{}

	ID, err := snowflake.ParseString(req.Id)
	if err != nil {
		logging.GetLogger(ctx).Errorf("parse id fail, id: %s, err: %v", req.Id, err)
		return nil, err
	}
	license.Id = ID

	license.Weight = req.Weight

	license.Auth = parseAuth(req.Auth)

	license.LicensePort = req.Port

	licenseType, err := licensetype.ToType(req.LicenseType)
	if err != nil {
		logging.GetLogger(ctx).Errorf("convert license type fail, license type: %d, err: %v", req.LicenseType, err)
		return nil, err
	}
	license.LicenseType = licenseType

	// 以下是非必填项校验
	if req.BeginTime != "" {
		beginTime, err := stringToTimeStamp(req.BeginTime)
		if err != nil {
			logging.Default().Errorf("invalid begin time: %s, err: %v", req.BeginTime, err)
			return nil, err
		}
		license.BeginTime = beginTime.AsTime()
	}

	if req.EndTime != "" {
		endTime, err := stringToTimeStamp(req.EndTime)
		if err != nil {
			logging.Default().Errorf("invalid end time: %s, err: %v", req.EndTime, err)
			return nil, err
		}
		license.EndTime = endTime.AsTime()
	}

	if req.ManagerId != "" {
		managerID, err := snowflake.ParseString(req.ManagerId)
		if err != nil {
			logging.GetLogger(ctx).Errorf("parse manager id fail, id: %s, err: %v", req.ManagerId, err)
			return nil, err
		}
		license.ManagerId = managerID
	}

	if req.LicenseProxies != nil && len(req.LicenseProxies) > 0 {
		content, err := json.Marshal(req.LicenseProxies)
		if err != nil {
			logging.Default().Errorf("marshal license addresses fail, err: %v", err)
			return nil, err
		}
		license.LicenseProxies = string(content)
	}

	license.Provider = req.Provider
	license.MacAddr = req.MacAddr
	license.ToolPath = req.ToolPath
	license.LicenseUrl = req.LicenseUrl
	license.LicenseNum = req.LicenseNum
	license.LicenseServer = req.LicenseEnvVar
	license.HpcEndpoint = req.HpcEndpoint
	license.CollectorType = req.CollectorType

	if req.AllowableHpcEndpoints != nil && len(req.AllowableHpcEndpoints) > 0 {
		license.AllowableHpcEndpoints = req.AllowableHpcEndpoints
	}

	return license, nil
}

func convertModelFromPutReq(ctx context.Context, req *moduleconfig.PutModuleConfigRequest) (*dbModels.ModuleConfig, error) {
	var moduleID snowflake.ID
	moduleID, err := snowflake.ParseString(req.Id)
	if err != nil {
		logging.GetLogger(ctx).Errorf("parse module id fail, id: %s, err: %v", req.Id, err)
		return nil, err
	}

	moduleName := req.ModuleName
	return &dbModels.ModuleConfig{
		Id:         moduleID,
		ModuleName: moduleName,
		Total:      req.Total,
	}, nil
}

func validateTimeForAdd(beginTime, endTime string) (*timestamppb.Timestamp, *timestamppb.Timestamp, *common.Error) {
	begin, err := stringToTimeStamp(beginTime)
	if err != nil {
		return nil, nil, common.WrapError("InvalidArgument.BeginTime", err.Error())
	}
	end, err := stringToTimeStamp(endTime)
	if err != nil {
		return nil, nil, common.WrapError("InvalidArgument.EndTime", err.Error())
	}

	// 校验 end > start
	if end.AsTime().Before(begin.AsTime()) {
		return nil, nil, common.WrapError("InvalidArgument.Time", "the end date should precede the start date")
	}
	return begin, end, nil
}

func validateTimeForPut(beginTime, endTime string) (begin, end *timestamppb.Timestamp, err *common.Error) {
	if beginTime == "" && endTime == "" {
		return nil, nil, nil
	}

	if beginTime != "" && endTime == "" {
		return nil, nil, common.WrapError("InvalidArgument.Time", "begin time and end time must be existed together")
	}

	if beginTime == "" && endTime != "" {
		return nil, nil, common.WrapError("InvalidArgument.Time", "begin time and end time must be existed together")
	}

	return validateTimeForAdd(beginTime, endTime)
}

func getResourceId(ctx *gin.Context) (snowflake.ID, bool) {
	idOrg := ctx.Param("id")
	id, err := snowflake.ParseString(idOrg)
	if err != nil || id == 0 {
		common.InvalidParams(ctx, "invalid id")
		return 0, false
	}
	return id, true
}

func toRespLicenseManager(lm *dbModels.LicenseManagerExt) (*licmanager.GetLicManagerResponseData, error) {
	res := &licmanager.GetLicManagerResponseData{
		Id:         lm.Id.String(),
		CreateTime: lm.CreateTime,
		Status:     lm.Status.GetValue(),
	}
	res.AppType = lm.AppType
	res.Os = lm.Os.GetValue()
	res.ComputeRule = lm.ComputeRule
	res.Desc = lm.Description
	res.LicenseInfos = make([]*licenseinfo.GetLicenseInfoResponseData, len(lm.Licenses))
	for i, licInfo := range lm.Licenses {
		resLicInfo := &licenseinfo.GetLicenseInfoResponseData{
			Id: licInfo.Id.String(),
		}

		if licInfo.LicenseProxies != "" {
			LicenseProxy := map[string]licenseinfo.LicenseProxy{}
			if err := json.Unmarshal([]byte(licInfo.LicenseProxies), &LicenseProxy); err != nil {
				logging.Default().Errorf("ummarshal license address fail, err: %v", err)
				return nil, err
			}
			resLicInfo.LicenseProxies = LicenseProxy
		}

		resLicInfo.ManagerId = licInfo.ManagerId.String()
		resLicInfo.CollectorType = licInfo.CollectorType
		resLicInfo.LicenseNum = licInfo.LicenseNum
		resLicInfo.LicenseUrl = licInfo.LicenseUrl
		resLicInfo.Provider = licInfo.Provider
		resLicInfo.HpcEndpoint = licInfo.HpcEndpoint
		resLicInfo.AllowableHpcEndpoints = licInfo.AllowableHpcEndpoints
		resLicInfo.ToolPath = licInfo.ToolPath
		resLicInfo.LicenseType = licInfo.LicenseType.GetValue()
		resLicInfo.Port = licInfo.LicensePort
		resLicInfo.Auth = licInfo.Auth == 1
		resLicInfo.BeginTime = licInfo.BeginTime.String()
		resLicInfo.EndTime = licInfo.EndTime.String()
		resLicInfo.MacAddr = licInfo.MacAddr
		resLicInfo.Weight = licInfo.Weight
		resLicInfo.LicenseEnvVar = licInfo.LicenseServer
		resLicInfo.LicenseServerStatus = licInfo.LicenseServerStatus
		resLicInfo.ModuleConfigs = make([]*moduleconfig.GetModuleConfigResponseData, len(licInfo.Modules))
		for j, mCfg := range licInfo.Modules {
			resCfg := toRespModuleConfig(mCfg)
			resLicInfo.ModuleConfigs[j] = resCfg
		}
		res.LicenseInfos[i] = resLicInfo
	}
	return res, nil
}

func toRespModuleConfig(moduleConfig *dbModels.ModuleConfig) *moduleconfig.GetModuleConfigResponseData {
	res := &moduleconfig.GetModuleConfigResponseData{
		Id:          moduleConfig.Id.String(),
		LicenseId:   moduleConfig.LicenseId.String(),
		ModuleName:  moduleConfig.ModuleName,
		Total:       moduleConfig.Total,
		UsedNum:     moduleConfig.Used,
		ActualTotal: moduleConfig.ActualTotal,
		ActualUsed:  moduleConfig.ActualUsed,
	}
	return res
}

func genSnowFlakId(ctx *gin.Context) (snowflake.ID, bool) {
	idRep, err := rpc.GetInstance().IDGen.GenerateID(ctx, &idgen.GenRequest{})
	if err != nil {
		logging.GetLogger(ctx).Warnf("GenIdFail, Error: %s", err.Error())
		common.InternalServerError(ctx, "Gen Id Fail")
		return 0, false
	}
	return snowflake.ID(idRep.Id), true
}

func handleValidationError(ctx *gin.Context, err error, req interface{}) {
	logging.GetLogger(ctx).Error(err)
	validationError := validation.HandleError(err, req)
	common.InvalidParams(ctx, fmt.Sprintf("invalid params, err: %v", validationError))
}

func parseAuth(auth bool) int {
	if auth {
		return 1
	} else {
		return 2
	}
}

func toModuleConfigModels(licenseID snowflake.ID, moduleConfigs []*moduleconfig.AddModuleConfigRequest) ([]*dbModels.ModuleConfig, error) {
	var moduleConfigModels []*dbModels.ModuleConfig
	for _, moduleConfig := range moduleConfigs {
		moduleConfigModel, err := toModuleConfigModel(licenseID, moduleConfig)
		if err != nil {
			logging.Default().Errorf("convert module config err, module config: %v, err: %v", moduleConfig, err)
			return nil, err
		}
		moduleConfigModels = append(moduleConfigModels, moduleConfigModel)
	}

	return moduleConfigModels, nil
}

func toModuleConfigModel(licenseID snowflake.ID, moduleConfig *moduleconfig.AddModuleConfigRequest) (*dbModels.ModuleConfig, error) {
	snowflakeID, err := rpc.GenID(context.TODO())
	if err != nil {
		return nil, err
	}

	return &dbModels.ModuleConfig{
		Id:         snowflakeID,
		LicenseId:  licenseID,
		ModuleName: moduleConfig.ModuleName,
		Total:      moduleConfig.Total,
	}, nil
}
