package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/utils"
	"net/http"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/validation"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"xorm.io/xorm"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobbatchget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
	jobsnapshotget "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
)

const (
	demoInput  string = "http[s]://domain/ys_id/input_path/"
	demoOutput string = "http[s]://domain/ys_id/output_path/"
	demoDest   string = "ys_id/dest_path"
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

func handleJobIDRequired(c *gin.Context, fe validation.Error) bool {
	if fe.Tag() == "required" && fe.Field() == "JobID" {
		common.ErrorResp(c, http.StatusBadRequest, api.InvalidJobID, "invalid params, "+fe.Error())
		return true
	}
	return false
}

// AppQuotaChecker app配额检查器
type AppQuotaChecker func(ctx context.Context, appID snowflake.ID) bool

// SharedChecker 不取整检查器
type SharedChecker func(ctx context.Context, userID snowflake.ID) (bool, error)

// AllocTypeChecker 是否能手动指定分配方式检查器
type AllocTypeChecker func(ctx context.Context, userID snowflake.ID) (bool, error)

// ValidateCreate 验证创建作业参数 return appInfo, inputZone, outputZone, error
func (h *Handler) ValidateCreate(ctx context.Context, req *jobcreate.Request,
	userID snowflake.ID, appQuotaChecker AppQuotaChecker, sharedChecker SharedChecker,
	allocTypeChecker AllocTypeChecker) (*models.Application, consts.Zone, consts.Zone,
	schema.ChargeParams, error) {
	// 校验app信息
	appID := req.Params.Application.AppID
	appYSID, err := snowflake.ParseString(appID)
	if err != nil {
		logging.GetLogger(ctx).Infof("parse appID error: %v", err)
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams,
			errors.WithMessagef(common.ErrInvalidAppID, "parse appID error, your input is '%s'", appID)
	}
	app, err := h.appSrv.Apps().GetApp(ctx, appYSID)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			// [404 AppIDNotFound 未找到app]
			return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams,
				errors.WithMessagef(common.ErrAppIDNotFound, "app not found, your input is '%s'", appID)
		}
		logging.GetLogger(ctx).Infof("get app info error: %v", err)
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams,
			errors.WithMessagef(err, "get app info error, your input is '%s'", appID)
	}

	appInfo, err := validateApplication(ctx, appID, req.Params.Application.Command, req.Zone,
		req.Params.EnvVars, app, userID, appQuotaChecker)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// 校验资源
	err = validateResource(ctx, req.Params.Resource)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// 校验文件信息
	inputZone, outputZone, err := validateFile(ctx, req.Params.Input, req.Params.Output,
		config.GetConfig().Zones, userID)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// 校验自定义规则
	err = validateCustomRule(ctx, req.Params.CustomStateRule)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	chargeParams, err := validateChargeParams(ctx, req.ChargeParam)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// job req.NoRound means shared
	err = validateShared(ctx, req.NoRound, userID, sharedChecker)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// AllocType字段仅内部用户使用
	err = validateAllocType(ctx, req.AllocType, userID, allocTypeChecker)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	// 校验其他信息
	err = validateOther(ctx, req.Name, req.Comment, req.Zone, req.Params.Application.Command)
	if err != nil {
		return nil, consts.ZoneUnknown, consts.ZoneUnknown, consts.InvalidChargeParams, err
	}

	return appInfo, inputZone, outputZone, chargeParams, nil
}

// ValidatePreScheduleParams 验证预调度参数
func (h *Handler) ValidatePreScheduleParams(ctx context.Context, preScheduleID string) (*models.PreSchedule, bool, error) {
	logger := logging.GetLogger(ctx).With("func", "ValidatePreScheduleParams", "preScheduleID", preScheduleID)
	preSchedule := preScheduleID != ""
	if !preSchedule {
		return nil, preSchedule, nil
	}

	// preScheduleID 合法性校验
	if err := ValidateID(ctx, preScheduleID); err != nil {
		logger.Infof("preScheduleID is invalid: %v", err)
		return nil, preSchedule, errors.WithMessage(common.ErrInvalidPreScheduleID, err.Error())
	}

	// 预调度作业是否存在
	scheduleInfo, exist, err := h.jobSrv.GetPreSchedule(ctx, preScheduleID)
	if err != nil {
		logger.Infof("get preSchedule info error: %v", err)
		return nil, preSchedule, errors.WithMessage(common.ErrInternalServer, "get preSchedule info error")
	}

	if !exist {
		logger.Infof("preSchedule not found")
		return nil, preSchedule, errors.WithMessage(common.ErrPreScheduleNotFound, "preSchedule not found")
	}

	if scheduleInfo.Used {
		logger.Infof("preSchedule has been used")
		return nil, preSchedule, errors.WithMessage(common.ErrPreScheduleUsed, "preSchedule has been used")
	}
	return scheduleInfo, preSchedule, nil
}

// ValidateNoPreScheduleParams 验证非预调度参数
func ValidateNoPreScheduleParams(ctx context.Context, req *jobcreate.Request) error {
	logger := logging.GetLogger(ctx).With("func", "ValidateNoPreScheduleParams")
	if req.Params.Application.AppID == "" {
		logger.Infof("appID cannot be empty")
		return errors.WithMessagef(common.ErrInvalidAppID, "appID cannot be empty when not preSchedule")
	}

	if req.Params.Resource == nil {
		logger.Infof("resource cannot be empty")
		return errors.WithMessagef(common.ErrInvalidArgumentResource, "resource cannot be empty when not preSchedule")
	}

	if req.Params.Input == nil {
		logger.Infof("input cannot be empty")
		return errors.WithMessagef(common.ErrInvalidArgumentInput, "input cannot be empty when not preSchedule")
	}
	return nil
}

// validateApplication 验证app信息
func validateApplication(ctx context.Context, appID, command, zone string,
	env map[string]string, app *models.Application, userID snowflake.ID,
	appQuotaChecker AppQuotaChecker) (*models.Application, error) {
	logger := logging.GetLogger(ctx).With("func", "ValidateApplication", "appID", appID, "zone", zone)

	// app 发布状态校验
	if app.PublishStatus != string(update.Published) {
		// [403 ErrAppNotPublished app未发布]
		return nil, errors.WithMessagef(common.ErrAppNotPublished, "app not published, your input is '%s'", appID)
	}

	// 校验用户应用配额(是否能使用该应用)
	if !appQuotaChecker(ctx, app.ID) {
		// [403 ErrUserNoAppQuota 用户无应用配额]
		return nil, errors.WithMessagef(common.ErrUserNoAppQuota, "user no app quota, your input is '%s'", appID)
	}

	// 校验appimage和appbinpath是否都为空
	if app.Image == "" && app.BinPath == "" {
		// [500 InternalServer 服务器内部错误]
		logger.Warnf("appimage and appbinpath cannot be empty at the same time")
		return nil, errors.WithMessagef(common.ErrInternalServer,
			"appimage and appbinpath cannot be empty at the same time, your input is '%s'", appID)
	}

	// 校验command长度是否超过限制
	if len(command) > consts.MaxCommandLength {
		// [400 InvalidArgument.Command 参数错误,command内容过长]
		return nil, errors.WithMessagef(common.ErrInvalidArgumentCommand,
			"command is too long, the value must contain a maximum of %d characters.", consts.MaxCommandLength)
	}

	// 如果command为空, 则为非命令行作业, 使用app中的command
	noCommandJob := command == ""

	if noCommandJob {
		if app.Command == "" {
			// [400 InvalidArgument.Command 参数错误,command不能为空]
			logger.Infof("app command is empty")
			return nil, errors.WithMessagef(common.ErrInternalServer, "app command cannot be empty")
		}

		// 解析ExtentionParams
		extentionParams := make(map[string]schema.ExtentionParam)
		err := json.Unmarshal([]byte(app.ExtentionParams), &extentionParams)
		if err != nil {
			logger.Warnf("parse extentionParams error: %v", err)
			// [400 InvalidArgument.Command 参数错误,解析extentionParams错误]
			return nil, errors.WithMessagef(common.ErrInternalServer, "parse app extentionParams error")
		}

		// 检查必填项
		for k, v := range extentionParams {
			if v.Must {
				if _, ok := env[k]; !ok {
					// 缺少必填参数
					// [400 InvalidArgument.ErrInvalidArgumentEnv 参数错误,缺少必填环境变量]
					return nil, errors.WithMessagef(common.ErrInvalidArgumentEnv,
						"missing required environment variable '%s'", k)
				}
			}
		}

		// 检查参数是否合法
		for k, v := range env {
			if _, ok := extentionParams[k]; !ok {
				continue // 提交作业时会将extentionParams中不存在的key忽略，这里可以continue
			}

			// 检查参数类型
			switch extentionParams[k].Type {
			case schema.ExtentionParamTypeString:
			case schema.ExtentionParamTypeStringList:
				// 检查参数是否在允许的范围内
				if extentionParams[k].Values != nil && len(extentionParams[k].Values) > 0 &&
					!extentionParams[k].Values.Contains(v) {
					// [400 InvalidArgument.ErrInvalidArgumentEnv 参数错误,参数value不在允许的范围内]
					return nil, errors.WithMessagef(common.ErrInvalidArgumentEnv,
						"environment variable '%s' value '%s' is not in the allowed range [%v]",
						k, v, extentionParams[k].Values)
				}
			}
		}
	}
	return app, nil
}

// validateResource 验证资源
func validateResource(ctx context.Context, resource *jobcreate.Resource) error {
	logger := logging.GetLogger(ctx).With("func", "ValidateResource", "resource", resource)

	// 核数校验
	core := consts.DefaultCoreNum
	if resource.Cores != nil {
		core = *resource.Cores
	}
	if core < consts.MinCoreNum {
		// [400 QuotaExhausted.Resource 指定的资源数量超过限额]
		logger.Infof("core cannot less than %d, your input is %d", consts.MinCoreNum, core)
		return errors.WithMessagef(common.ErrQuotaExhaustedResource,
			"core cannot less than %d, your input is %d", consts.MinCoreNum, core)
	}
	resource.Cores = &core

	// 内存校验
	memory := consts.DefaultMemory
	if resource.Memory != nil {
		memory = *resource.Memory
	}
	if memory < consts.MinMemory {
		// [400 QuotaExhausted.Resource 指定的资源数量超过限额]
		logger.Infof("memory cannot less than %d, your input is %d", consts.MinMemory, memory)
		return errors.WithMessagef(common.ErrQuotaExhaustedResource,
			"memory cannot less than %d, your input is %d", consts.MinMemory, memory)
	}
	resource.Memory = &memory
	return nil
}

// validateFile 验证文件信息
func validateFile(ctx context.Context, input *jobcreate.Input, output *jobcreate.Output,
	zones schema.Zones, userID snowflake.ID) (inputZone consts.Zone, outputZone consts.Zone, err error) {
	inputZone = consts.ZoneUnknown
	outputZone = consts.ZoneUnknown
	inputZone, err = validateInput(ctx, input, zones, userID)
	if err != nil {
		return inputZone, outputZone, err
	}
	outputZone, err = validateOutput(ctx, output, zones, userID)
	if err != nil {
		return inputZone, outputZone, err
	}
	return inputZone, outputZone, nil
}

// validateInput 验证输入参数
func validateInput(ctx context.Context, input *jobcreate.Input, zones schema.Zones,
	userID snowflake.ID) (inputZone consts.Zone, err error) {
	logger := logging.GetLogger(ctx).With("func", "ValidateInput")

	inputZone = consts.ZoneUnknown

	// 校验输出文件
	if input == nil {
		// happened when preSchedule
		return inputZone, nil
	}

	inputType := consts.FileType(input.Type)
	// 判断input.Type是否为'hpc_storage'或'cloud_storage'
	if !inputType.EnableStorageType() {
		errorMessage := fmt.Sprintf("the file input type must in [%s], your input type is '%s'.",
			consts.EnableStorageTypeString(), input.Type)
		return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
	}

	inputSource := input.Source
	// 判断input.Source是否为空
	if inputSource == "" {
		return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput,
			"the file input path cannot be empty.")
	}

	// 判断input.Source是否是带域名的绝对路径
	// 构成: http[s]://文件区域域名/目标用户的ys_id/文件路径
	// 示例: https://jinan-storage.yuansuan.cn/ys_id/input/path1/
	if !strings.HasPrefix(inputSource, "http://") && !strings.HasPrefix(inputSource, "https://") {
		return inputZone, errors.WithMessagef(common.ErrInvalidArgumentInput,
			"the file input path must be absolute path with domain, like: '%s', your input path is '%s'.",
			demoInput, inputSource)
	}

	// 必须是配置文件存在的区域的域名
	if !zones.Exist(inputSource) {
		return inputZone, errors.WithMessagef(common.ErrInvalidArgumentInput,
			"the file input path must be in the zone list [%v], your input path is '%s'.", zones.List(), inputSource)
	}

	// 解析出ys_id,判断授权
	inputYsIDStr := util.ParseYsID(inputSource)
	if inputYsIDStr == "" {
		return inputZone, errors.WithMessagef(common.ErrInvalidArgumentInput,
			"the file input path must have ys_id, like: '%s', your input path is '%s'.", demoInput, inputSource)
	}

	inputYsID, err := snowflake.ParseString(inputYsIDStr)
	if err != nil {
		return inputZone, errors.WithMessagef(common.ErrInvalidArgumentInput,
			"the file input path ys_id is invalid, your input path ys_id is '%s'.", inputYsIDStr)
	}

	inputCanAccess, err := validateCanAccess(ctx, userID, inputYsID)
	if err != nil {
		logger.Warnf("the file input path ys_id is Unauthorized, error: %v, your input path ys_id is '%s'.",
			err, inputYsIDStr)
		return inputZone, err
	}

	if !inputCanAccess {
		errorMessage := fmt.Sprintf("the file input path ys_id is Unauthorized, userID is %s, "+
			"input path ys_id is '%s'.", userID, inputYsIDStr)
		logger.Warn(errorMessage)
		return inputZone, errors.WithMessage(common.ErrJobPathUnauthorized, errorMessage)
	}

	// 判断input.Source是否包含..
	if strings.Contains(inputSource, "..") {
		errorMessage := fmt.Sprintf("the file input path cannot contain '..', your input path is '%s'.", inputSource)
		logger.Warn(errorMessage)
		return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
	}

	// 构成: 目标用户的ys_id/文件路径 或 为空 ,不带前缀'/'
	// 示例: ys_id/input/path1/
	inputDest := input.Destination

	if inputDest != "" {
		// input.Destination不能以域名开头
		if strings.HasPrefix(inputDest, "http://") || strings.HasPrefix(inputDest, "https://") {
			errorMessage := fmt.Sprintf("the file input destination cannot start with 'http://' or"+
				" 'https://', your input destination is '%s'.", inputDest)
			logger.Warn(errorMessage)
			return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
		}

		// 判断input.Destination是否以'/'开头
		if strings.HasPrefix(inputDest, "/") {
			errorMessage := fmt.Sprintf("the file input destination cannot start with '/', should like '%s', "+
				"your input destination is '%s'.", demoDest, inputDest)
			logger.Warn(errorMessage)
			return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
		}

		// 判断input.Destination是否包含..
		if strings.Contains(inputDest, "..") {
			errorMessage := fmt.Sprintf("the file input destination cannot contain '..', should like '%s', "+
				"your input destination is '%s'.", demoDest, inputDest)
			logger.Warn(errorMessage)
			return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
		}

		// 判断input.Destination是否包含正确的ys_id
		destYsIDStr := util.ParseYsIDWithOutDomain(inputDest)
		if destYsIDStr != "" { // 理论上进入这里就不可能会为空
			destYsID, err := snowflake.ParseString(destYsIDStr)
			if err != nil {
				errorMessage := fmt.Sprintf("the file input destination path ys_id is invalid, should like '%s', "+
					"your input destination path ys_id is '%s'.", demoDest, destYsIDStr)
				logger.Warn(errorMessage)
				return inputZone, errors.WithMessage(common.ErrInvalidArgumentInput, errorMessage)
			}

			destCanAccess, err := validateCanAccess(ctx, userID, destYsID)
			if err != nil {
				logger.Warnf("the file input destination ys_id is Unauthorized, ValidateCanAccess error: %v, "+
					"your input destination ys_id is '%s'.", err, destYsIDStr)
				return inputZone, err
			}

			if !destCanAccess {
				errorMessage := fmt.Sprintf("the file input destination path ys_id is Unauthorized, should like '%s', "+
					"your input destination path ys_id is '%s'.", demoDest, destYsIDStr)
				logger.Warn(errorMessage)
				return inputZone, errors.WithMessage(common.ErrJobPathUnauthorized, errorMessage)
			}
		}

	}

	inputZone = consts.Zone(zones.GetZoneByEndpoint(inputSource))
	return inputZone, nil
}

// validateOutput 验证输出参数
func validateOutput(ctx context.Context, output *jobcreate.Output, zones schema.Zones,
	userID snowflake.ID) (outputZone consts.Zone, err error) {
	logger := logging.GetLogger(ctx).With("func", "ValidateOutput")

	outputZone = consts.ZoneUnknown
	// 校验输出文件
	if output == nil {
		return outputZone, nil
	}

	// [400 InvalidArgument.Output 输出路径不正确,例如前缀含文件路径等错误。带"../"等错误]
	outputType := consts.FileType(output.Type)
	// 判断output.Type是否为'hpc_storage'或'cloud_storage'
	if !outputType.EnableStorageType() {
		logger.Infof("the file output type must in [%v], your output type is '%s'.",
			consts.EnableStorageTypeString(), output.Type)
		return outputZone, errors.WithMessagef(common.ErrInvalidArgumentOutput,
			"the file output type must in [%v], your output type is '%s'.",
			consts.EnableStorageTypeString(), output.Type)
	}

	outputAddress := output.Address
	// 判断output.Address是否为空
	if outputAddress == "" {
		logger.Infof("the file output path cannot be empty.")
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput,
			"the file output path cannot be empty.")
	}

	// 判断output.Address是否是带域名的绝对路径
	// 构成: http[s]://文件区域域名/目标用户的ys_id/文件路径
	// 示例: https://jinan-storage.yuansuan.cn/ys_id/output/path1/
	if !strings.HasPrefix(outputAddress, "http://") && !strings.HasPrefix(outputAddress, "https://") {
		errorMessage := fmt.Sprintf("the file output path must be absolute path with domain,like: '%s', "+
			"your output path is '%s'.", demoOutput, outputAddress)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput, errorMessage)
	}

	// 必须是配置文件存在的区域的域名
	if !zones.Exist(outputAddress) {
		errorMessage := fmt.Sprintf("the file output path must be in the zone list [%v], "+
			"your output path is '%s'.", zones.List(), outputAddress)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput, errorMessage)
	}

	// 解析出ys_id,判断授权
	outputYsIDStr := util.ParseYsID(outputAddress)
	if outputYsIDStr == "" {
		errorMessage := fmt.Sprintf("the file output path must have ys_id, like: '%s', "+
			"your output path is '%s'.", demoOutput, outputAddress)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput, errorMessage)
	}

	outputYsID, err := snowflake.ParseString(outputYsIDStr)
	if err != nil {
		errorMessage := fmt.Sprintf("the file output path ys_id is invalid, "+
			"your output path ys_id is '%s'.", outputYsIDStr)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput, errorMessage)
	}

	outputCanAccess, err := validateCanAccess(ctx, userID, outputYsID)
	if err != nil {
		logger.Warnf("the file output path ys_id is Unauthorized, ValidateCanAccess error: %v, "+
			"userID is %s, output path ys_id is '%s'.", err, userID, outputYsIDStr)
		return outputZone, err
	}

	if !outputCanAccess {
		errorMessage := fmt.Sprintf("the file output path ys_id is Unauthorized, "+
			"your output path ys_id is '%s'.", outputYsIDStr)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrJobPathUnauthorized, errorMessage)
	}

	// 判断output.Address是否包含..
	if strings.Contains(outputAddress, "..") {
		errorMessage := fmt.Sprintf("the file output path cannot contain '..', "+
			"your output path is '%s'.", outputAddress)
		logger.Warn(errorMessage)
		return outputZone, errors.WithMessage(common.ErrInvalidArgumentOutput, errorMessage)
	}

	outputZone = consts.Zone(zones.GetZoneByEndpoint(outputAddress))
	return outputZone, nil
}

// validateCanAccess 验证用户是否有权限访问他人文件的权限
func validateCanAccess(ctx context.Context, userID snowflake.ID, ysID snowflake.ID) (bool, error) {
	// 是自己的，直接返回
	if ysID == userID {
		return true, nil
	}

	// 不是自己的要校验有没有权限
	// TODO: 日后可能有跨账号访问的需求，调用iam接口判断是否有权限，暂时不做

	// 返回false
	return false, nil
}

// validateCustomRule 验证自定义规则
func validateCustomRule(ctx context.Context, customRule *jobcreate.CustomStateRule) error {
	// 为nil，直接返回
	if customRule == nil {
		return nil
	}

	// 校验KeyStatement长度是否超限
	if len(customRule.KeyStatement) > consts.MaxCustomStateRuleLength {
		// [400 InvalidArgument.CustomStateRule.KeyStatement 参数错误,KeyStatement参数过长]
		return errors.WithMessagef(common.ErrInvalidArgumentCustomStateRuleKeyStatement,
			"KeyStatement is too long, the value must contain a maximum of %d characters.",
			consts.MaxCustomStateRuleLength)
	}

	// 校验ResultState是否合法
	if customRule.ResultState != jobcreate.ResultStateCompleted && customRule.ResultState != jobcreate.ResultStateFailed {
		// [400 InvalidArgument.CustomStateRule.ResultState 参数错误,ResultState参数错误]
		return errors.WithMessagef(common.ErrInvalidArgumentCustomStateRuleResultState,
			"ResultState is invalid, the value must in [%s,%s].",
			jobcreate.ResultStateCompleted, jobcreate.ResultStateFailed)
	}
	return nil
}

func validateChargeParams(ctx context.Context, chargeParams schema.ChargeParams) (schema.ChargeParams, error) {
	if !config.GetConfig().BillEnabled {
		return chargeParams, nil
	}
	defaultChargeType := schema.PostPaid
	defaultChargeParams := schema.ChargeParams{
		ChargeType: &defaultChargeType,
	}
	if chargeParams.ChargeType == nil {
		return defaultChargeParams, nil
	}
	if *chargeParams.ChargeType == "" {
		return defaultChargeParams, nil
	}
	if !chargeParams.ChargeType.IsValid() {
		return schema.ChargeParams{}, errors.WithMessagef(common.ErrInvalidChargeParams,
			"ChargeType can only be [PrePaid | PostPaid]")
	}
	// FIXME 暂时不支持PrePaid模式，后续支持了再将校验去除
	if *chargeParams.ChargeType == schema.PrePaid {
		return schema.ChargeParams{}, errors.WithMessagef(common.ErrInvalidChargeParams,
			"unsupported for [PrePaid]")
	}
	return chargeParams, nil
}

// validateOther 验证其他参数
func validateOther(ctx context.Context, name, comment, zone, command string) error {
	// 作业name长度校验
	if len(name) > consts.MaxNameLength {
		// [400 InvalidArgument.Name 参数错误,name参数过长]
		return errors.WithMessagef(common.ErrInvalidArgumentName,
			"name is too long, the value must contain a maximum of %d characters.", consts.MaxNameLength)
	}

	// 校验comment
	if len(comment) > consts.MaxCommentLength {
		// [400 InvalidArgument.Comment 参数错误,comment内容过长]
		return errors.WithMessagef(common.ErrInvalidArgumentComment,
			"comment is too long, the value must contain a maximum of %d characters.", consts.MaxCommentLength)
	}

	// 校验zone
	zones := config.GetConfig().Zones
	if zone != "" {
		if !zones.IsZone(zone) {
			// [400 InvalidArgument.Zone 参数错误,zone未知]
			return errors.WithMessagef(common.ErrInvalidArgumentZone,
				"zone is unknown.must in [%v], your zone is '%s'.", zones.List(), zone)
		}

		// 校验zone是否有HPCEndpoint
		if zones[zone].HPCEndpoint == "" {
			// [400 InvalidArgument.Zone 参数错误,zone hpc endpoint为空]
			return errors.WithMessagef(common.ErrInvalidArgumentZone,
				"zone hpc endpoint is empty, your zone is '%s'.", zone)
		}
	}
	return nil
}

// validateShared 校验是否是特殊用户能用Shared参数
func validateShared(ctx context.Context, shared bool, userID snowflake.ID, sharedChecker SharedChecker) error {
	logger := logging.GetLogger(ctx).With("func", "ValidShared", "shared", shared, "userID", userID)

	if !shared {
		return nil
	}

	if canShared, err := sharedChecker(ctx, userID); err != nil {
		logger.Infof("check user is ys product user error: %v", err)
		return err
	} else if !canShared {
		logger.Warnf("user is not ys product user")
		// [403 ErrJobAccessDenied 用户无权限]
		return errors.WithMessagef(common.ErrJobAccessDenied,
			"user is not ys product user, cannot use shared(noRound) parameter")
	}
	return nil
}

// validateAllocType 校验是否是特殊用户能用AllocType参数
func validateAllocType(ctx context.Context, allocType string,
	userID snowflake.ID, allocTypeChecker AllocTypeChecker) error {
	logger := logging.GetLogger(ctx).With("func", "ValidateAllocType", "allocType",
		allocType, "userID", userID)

	isYsProductUser, err := allocTypeChecker(ctx, userID)
	if err != nil {
		logger.Infof("check user is ys product user error: %v", err)
		return err
	}
	if allocType != "" {
		// 如果不是内部用户，但尝试使用 AllocType，返回错误
		if !isYsProductUser {
			logger.Warnf("non-YsProductUser attempted to use AllocType")
			return errors.WithMessagef(common.ErrJobAccessDenied,
				"user is not ys product user, cannot use AllocType parameter")
		}
		// 如果 allocType 有值，那么一定得是 "average"
		if allocType != "average" {
			logger.Warnf("Invalid AllocType value: %s", allocType)
			return errors.WithMessagef(common.ErrInvalidArgumentAllocType,
				"Invalid AllocType value: %s. The only valid value is 'average'", allocType)
		}
	}
	// 如果是内部用户或 allocType 为空，则通过验证
	return nil
}

// ValidateUserInfo ...
func ValidateUserInfo(c *gin.Context) (snowflake.ID, error) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return 0, errors.WithMessagef(common.ErrInvalidUserID, err.Error())
	}
	return userID, err
}

// ValidateBatchGet 校验批量获取作业
func ValidateBatchGet(ctx context.Context, req *jobbatchget.Request) error {
	if len(req.JobIDs) > consts.MaxBatchGetJobIDs {
		logging.GetLogger(ctx).Warnf("invalid jobIDs length: %v", len(req.JobIDs))
		return errors.WithMessagef(common.ErrInvalidArgumentJobIDs,
			"Too many JobID, the number of JobIDs must be less than %d.", consts.MaxBatchGetJobIDs)
	}
	return nil
}

// ValidateList 校验列表
func ValidateList(ctx context.Context, req *joblist.Request) error {
	if req.PageSize == nil {
		ps := consts.DefaultPageSize
		req.PageSize = &ps
	}

	if req.PageOffset == nil {
		po := consts.DefaultPageOffset
		req.PageOffset = &po
	}

	// 分页合法性判断 pageSize
	if *req.PageSize <= 0 || *req.PageSize > 1000 {
		logging.GetLogger(ctx).Warnf("invalid page size: %v", req.PageSize)
		return errors.WithMessage(common.ErrInvalidArgumentPageSize,
			"page size must be greater than 0 and less than or equal to 1000")
	}

	// 分页合法性判断 pageOffset
	if *req.PageOffset < 0 {
		logging.GetLogger(ctx).Warnf("invalid page offset: %v", req.PageOffset)
		return errors.WithMessage(common.ErrInvalidArgumentPageOffset, "page index can't be less than 0")

	}

	// zone 合法性校验
	zones := config.GetConfig().Zones
	if req.Zone != "" && !zones.IsZone(req.Zone) {
		logging.GetLogger(ctx).Warnf("invalid zone: %v", req.Zone)
		return errors.WithMessage(common.ErrInvalidArgumentZone, "zone is unknown")
	}

	// job state valid check
	if req.JobState != "" {
		ok := consts.ValidStringState(req.JobState)
		if !ok {
			err := fmt.Errorf("invalid job state: %v, valid job state: %v", req.JobState, consts.AllStateString())
			logging.GetLogger(ctx).Warnf(err.Error())
			return errors.WithMessage(common.ErrInvalidArgumentJobState, err.Error())
		}
	}
	return nil
}

// ValidateSnapshotImg 校验云图数据
func ValidateSnapshotImg(ctx context.Context, JobID, Path string) (*jobsnapshotget.Request, error) {
	req := &jobsnapshotget.Request{}
	if JobID == "" {
		return req, errors.WithMessage(common.ErrInvalidJobID, "job id empty")
	}
	if Path == "" {
		return req, errors.WithMessage(common.ErrInvalidPath, "path empty")
	}
	// path非法校验
	if strings.Contains(Path, "..") {
		return req, errors.WithMessage(common.ErrInvalidPath, "path cannot contain '..'")
	}
	req.JobID = JobID
	req.Path = Path
	return req, nil
}

// ValidateJobID 校验
func ValidateJobID(ctx context.Context, id string) error {
	// jobID 合法性校验
	if err := ValidateID(ctx, id); err != nil {
		return errors.WithMessage(common.ErrInvalidJobID, err.Error())
	}
	return nil
}

func ValidateUserID(ctx context.Context, id string) error {
	// userID 合法性校验
	if err := ValidateID(ctx, id); err != nil {
		return errors.WithMessagef(common.ErrInvalidUserID, err.Error())
	}
	return nil
}

func ValidateAppID(ctx context.Context, id string) error {
	// appID 合法性校验
	if err := ValidateID(ctx, id); err != nil {
		return errors.WithMessagef(common.ErrInvalidAppID, err.Error())
	}
	return nil
}

// ValidateID 校验ID
func ValidateID(ctx context.Context, id string) error {
	logger := logging.GetLogger(ctx).With("func", "ValidateID", "ID", id)

	_, err := snowflake.ParseString(id)
	if err != nil {
		logger.Infof("parse id error: %v", err)
		return errors.WithMessagef(err, "parse id error, your input is '%s'", id)
	}
	return nil
}

func ValidateAccountID(ctx context.Context, id string) error {
	// accountID 合法性校验
	if err := ValidateID(ctx, id); err != nil {
		return errors.WithMessagef(common.ErrInvalidAccountId, err.Error())
	}
	return nil
}

func (h *Handler) ValidatePreSchedule(ctx context.Context, req *jobpreschedule.Request,
	userID snowflake.ID, appQuotaChecker AppQuotaChecker, sharedChecker SharedChecker) (
	*models.Application, schema.Zones, error) {
	// 校验app信息
	appID := req.Params.Application.AppID
	appYSID, err := snowflake.ParseString(appID)
	if err != nil {
		return nil, nil, err
	}
	app, err := h.appSrv.Apps().GetApp(ctx, appYSID)
	if err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			// [404 AppIDNotFound 未找到app]
			return nil, nil, errors.WithMessagef(common.ErrAppIDNotFound,
				"app not found, your appID=%s, err=%v", appID, err)
		}
		return nil, nil, err
	}
	// 校验app信息
	appInfo, err := validateApplication(ctx, req.Params.Application.AppID, req.Params.Application.Command,
		"PreSchedule", req.Params.EnvVars, app, userID, appQuotaChecker)
	if err != nil {
		return nil, nil, err
	}
	// 校验资源
	err = ValidatePreScheduleResource(ctx, req.Params.Resource)
	if err != nil {
		return nil, nil, err
	}
	err = validateShared(ctx, req.Shared, userID, sharedChecker)
	if err != nil {
		return nil, nil, err
	}
	// 校验分区
	zones, err := ValidatePreScheduleZones(ctx, req.Fixed, req.Zones)
	if err != nil {
		return nil, nil, err
	}
	return appInfo, zones, nil
}

// ValidatePreScheduleResource 预调度资源校验
func ValidatePreScheduleResource(ctx context.Context, resource *jobpreschedule.Resource) error {
	logger := logging.GetLogger(ctx).With("func", "ValidateResource", "resource", resource)

	// 核数校验
	mincores := consts.DefaultCoreNum
	if resource.MinCores != nil {
		mincores = *resource.MinCores
	}
	if mincores < consts.MinCoreNum {
		// [400 QuotaExhausted.Resource 指定的资源数量超过限额]
		logger.Infof("core cannot less than %d, your input is %d", consts.MinCoreNum, mincores)
		return errors.WithMessagef(common.ErrQuotaExhaustedResource,
			"min core cannot less than %d, your input is %d", consts.MinCoreNum, mincores)
	}
	resource.MinCores = &mincores
	maxcores := consts.DefaultCoreNum
	if resource.MaxCores != nil {
		maxcores = *resource.MaxCores
	}
	if maxcores < mincores {
		// [400 QuotaExhausted.Resource 指定的资源数量超过限额]
		logger.Infof("max core cannot less than min core, your input is %d", maxcores)
		return errors.WithMessagef(common.ErrQuotaExhaustedResource,
			"max core cannot less than min core, your input is %d", maxcores)
	}

	// 内存校验
	memory := consts.DefaultMemory
	if resource.Memory != nil {
		memory = *resource.Memory
	}
	if memory < consts.MinMemory {
		// [400 QuotaExhausted.Resource 指定的资源数量超过限额]
		logger.Infof("memory cannot less than %d, your input is %d", consts.MinMemory, memory)
		return errors.WithMessagef(common.ErrQuotaExhaustedResource,
			"memory cannot less than %d, your input is %d", consts.MinMemory, memory)
	}
	resource.Memory = &memory
	return nil
}

func ValidListNeedSyncFileJobs(c context.Context, req *jobneedsyncfile.Request) error {
	if req.PageSize == nil {
		ps := consts.DefaultPageSize
		req.PageSize = &ps
	}
	if req.PageOffset == nil {
		po := consts.DefaultPageOffset
		req.PageOffset = &po
	}

	// 分页合法性判断 pageSize
	if *req.PageSize <= 0 || *req.PageSize > 1000 {
		logging.GetLogger(c).Warnf("invalid page size: %v", req.PageSize)
		return errors.WithMessage(common.ErrInvalidArgumentPageSize,
			"page size must be greater than 0 and less than or equal to 1000")
	}
	// 分页合法性判断 pageOffset
	if *req.PageOffset < 0 {
		logging.GetLogger(c).Warnf("invalid page offset: %v", req.PageOffset)
		return errors.WithMessage(common.ErrInvalidArgumentPageOffset, "page index can't be less than 0")
	}
	// 校验zone
	zones := config.GetConfig().Zones
	if req.Zone == "" {
		return errors.WithMessage(common.ErrInvalidArgumentZone, "zone should not be empty.")
	}
	if req.Zone != "" && !zones.IsZone(req.Zone) {
		logging.GetLogger(c).Warnf("invalid zone: %v", req.Zone)
		return errors.WithMessage(common.ErrInvalidArgumentZone, "zone is unknown")
	}
	return nil
}

func ValidSyncFileState(c context.Context, req *jobsyncfilestate.Request) error {
	logger := logging.GetLogger(c)
	if req.DownloadFileSizeCurrent < 0 {
		logger.Warnf("invalid file size: %v", req.DownloadFileSizeCurrent)
		return errors.WithMessage(common.ErrInvalidArgumentDownloadFileSizeCurrent,
			"current file size can not be less than 0")
	}

	if req.DownloadFileSizeTotal < 0 {
		logger.Warnf("invalid file size: %v", req.DownloadFileSizeTotal)
		return errors.WithMessage(common.ErrInvalidArgumentDownloadFileSizeTotal,
			"total file size can not be less than 0")
	}

	if req.DownloadFinishedTime != "" {
		_, err := util.ParseTime(req.DownloadFinishedTime, time.RFC3339)
		if err != nil {
			logger.Warnf("invalid download finished time: %v, err: %v", req.DownloadFinishedTime, err)
			return errors.WithMessage(common.ErrInvalidArgumentDownloadFinishedTime,
				"downloadFinishedTime is not RFC3339 format")
		}
	}

	if req.TransmittingTime != "" {
		_, err := util.ParseTime(req.TransmittingTime, time.RFC3339)
		if err != nil {
			logger.Warnf("invalid transmitting time: %v, err: %v", req.TransmittingTime, err)
			return errors.WithMessage(common.ErrInvalidArgumentDownloadFinishedTime,
				"downloadFinishedTime is not RFC3339 format")
		}
	}

	fileSyncState := consts.FileSyncState(req.FileSyncState)
	if !fileSyncState.IsValid() {
		logger.Warnf("invalid file sync state: %v", req.FileSyncState)
		return errors.WithMessage(common.ErrJobFileSyncStateUpdateFailed, "invalid file sync state!")
	}
	return nil
}

func ValidatePreScheduleZones(ctx context.Context, fixed bool, zones []string) (schema.Zones, error) {
	// 校验分区
	// fixed为true时，只从request.Zones的分区范围中选取，此时request.Zones不能为空
	// fixed为false时，从所有分区中选取，此时request.Zones可为空
	schemaZones := config.GetConfig().Zones
	if fixed {
		if len(zones) == 0 {
			return nil, errors.WithMessage(common.ErrInvalidArgumentZone, "fixed is true, zones cannot be empty")
		}

		resZones := make(schema.Zones)
		for _, zone := range zones {
			if !schemaZones.IsZone(zone) {
				continue
			}
			if schemaZones[zone].HPCEndpoint == "" {
				continue
			}
			resZones[zone] = schemaZones[zone]
		}

		if len(resZones) == 0 {
			return nil, errors.WithMessage(common.ErrInvalidArgumentZone, "fixed is true, no valid zone")
		}
		return resZones, nil
	}
	return schemaZones, nil
}
