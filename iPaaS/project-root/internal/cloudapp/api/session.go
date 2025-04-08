package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/utils/payby"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"go.uber.org/multierr"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	internal_common "github.com/yuansuan/ticp/iPaaS/project-root/internal/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/sessionaction"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/sessionrestore"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/openapi"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rdp"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/template"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/webrtc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

var (
	defaultChargeType = schema.PostPaid
	maxMountPaths     = 20
)

var (
	ErrSoftwareNotExist = errors.New("Software not exist")
	ErrHardwareNotExist = errors.New("Hardware not exist")
)

// PostSessions FIXME 代码冗长，数据库操作有重复，需要整理
// 通过storage的openapi接口创建共享点然后再创建虚拟机，不在虚拟机中创建共享点
func PostSessions(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiPostRequest)
	err = bindPostSessionsRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind request body failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIPostSessionsRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api post sessions request failed, %v", err)
		return
	}

	s, err := util.GetState(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	isYSProduct, err := s.IamClient.IsYsProductUser(userId)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
		logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
		return
	}
	// 计费相关校验
	accountId := snowflake.ID(0)
	payByAccountId := snowflake.ID(0)
	chargeType := schema.ChargeType("")
	payByUserId := snowflake.ID(0)
	billEnabled := config.GetConfig().BillEnabled
	if billEnabled {
		var errMsg string
		req.ChargeParams, err, errMsg = checkAndEnsureChargeParams(req.ChargeParams)
		if err = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentChargeParams, errMsg)); err != nil {
			logger.Warnf("check and ensure charge params failed, %v", err)
			return
		}
		chargeType = *req.ChargeParams.ChargeType

		// 代支付
		if req.PayBy != nil {
			reqPayBy, err := payby.ParseToken(*req.PayBy)
			if err = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentPayBy, "parse payBy token failed")); err != nil {
				logger.Infof("parse payby token failed, req paby: %+v, err:%v", *req.PayBy, err)
				return
			}

			accessKeyID := reqPayBy.GetAccessKeyID()
			resourceTag := reqPayBy.GetResourceTag()
			timestamp := reqPayBy.GetTimestamp()
			// 获取accountId
			iamClient := iam_client.NewClient(config.GetConfig().OpenAPI.Endpoint, config.GetConfig().OpenAPI.AccessKeyId, config.GetConfig().OpenAPI.AccessKeySecret)
			resp, err := iamClient.GetSecret(&iam_api.GetSecretRequest{
				AccessKeyId: accessKeyID,
			})
			if err != nil {
				if strings.Contains(err.Error(), "secret not found") {
					_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentErrorCode, "accessKeyId not found"))
					return
				}
				logger.Warnf("get secret %s from iam server failed %s", accessKeyID, err.Error())
				_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get secret from iam server failed"))
				return
			}

			// 校验token是否一致
			newPayBy, err := payby.NewPayBy(accessKeyID, resp.AccessKeySecret, resourceTag, timestamp)
			if err != nil {
				logger.Infof("generate payBy failed, payBy: %+v, err: %v", newPayBy, err)
				_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentPayBy, "generate payBy failed"))
				return
			}

			if !newPayBy.SignEqualTo(reqPayBy) {
				err := errors.Errorf("payBy sign check failed, request payBy params: %+v, generate new sign: %+v", reqPayBy, newPayBy)
				logger.Infof("payBy sign check failed, request payBy params: %+v, generate new sign: %+v", reqPayBy, newPayBy)
				_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentPayBySignature, "invalid Token"))
				return
			}

			// token 是否过期
			if time.Now().Sub(time.UnixMilli(timestamp)) > time.Minute*5 {
				err := errors.Errorf("payBy sign check failed, request payBy params: %+v, generate new sign: %+v", reqPayBy, newPayBy)
				logger.Infof("payBy token expired, token: %v", req.PayBy)
				_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.PayByTokenExpire, "payBy token expired"))
				return
			}

			payByUserId, err = snowflake.ParseString(resp.YSId)
			if err != nil {
				logger.Errorf("parse payByUserId %s from snowflake failed %s", resp.YSId, err.Error())
				_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "parse payByUserId failed"))
				return
			}

			payByAccountId, err = util.CheckAccountInPostSession(logger, s, payByUserId)
			if err != nil {
				logger.Infof("check accountId failed, ysid: %s, err: %s", userId, err.Error())
				if errors.Is(err, internal_common.ErrRequestAccountByYSID) {
					_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountId) {
					_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalErrorInvalidAccountId, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountStatusNotEnoughBalance) {
					_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusNotEnoughBalance, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountStatusFrozen) {
					_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusFrozen, err.Error()))
					return
				}
				_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentUserIds, err.Error()))
				return
			}
		} else {
			accountId, err = util.CheckAccountInPostSession(logger, s, userId)
			if err != nil {
				logger.Infof("check accountId failed, ysid: %s, err: %s", userId, err.Error())
				if errors.Is(err, internal_common.ErrRequestAccountByYSID) {
					_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountId) {
					_ = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalErrorInvalidAccountId, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountStatusNotEnoughBalance) {
					_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusNotEnoughBalance, err.Error()))
					return
				}

				if errors.Is(err, internal_common.ErrInvalidAccountStatusFrozen) {
					_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusFrozen, err.Error()))
					return
				}
				_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentUserIds, err.Error()))
				return
			}
		}
	}

	// 软硬件校验
	softwareId, hardwareId := snowflake.MustParseString(*req.SoftwareId), snowflake.MustParseString(*req.HardwareId)
	var software *models.Software
	var hardware *models.Hardware
	if billEnabled && req.PayBy != nil && *req.PayBy != "" {
		isYSProduct, err := s.IamClient.IsYsProductUser(payByUserId)
		if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
			logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
			return
		}
		software, err = validateSoftware(c, softwareId, payByUserId, isYSProduct, s.OpenAPIClient)
		if err != nil {
			logger.Infof("validate software failed, %v", err)
			return
		}

		hardware, err = validateHardware(c, hardwareId, payByUserId, isYSProduct, s.OpenAPIClient)
		if err != nil {
			logger.Infof("validate hardware failed, %v", err)
			return
		}
	} else {
		software, err = validateSoftware(c, softwareId, userId, isYSProduct, s.OpenAPIClient)
		if err != nil {
			logger.Infof("validate software failed, %v", err)
			return
		}

		hardware, err = validateHardware(c, hardwareId, userId, isYSProduct, s.OpenAPIClient)
		if err != nil {
			logger.Infof("validate hardware failed, %v", err)
			return
		}
	}

	if software.Zone != hardware.Zone {
		err = response.BadRequestIfError(c, fmt.Errorf("software, hardware zone not equal"), response.WrapErrorResp(common.InvalidArgumentHardwareSoftwareZoneNotEqual, "software, hardware zone not equal"))
		logger.Infof("software zone: %s, hardware: %s, %v", software.Zone, hardware.Zone, err)
		return
	}

	// 云盘挂载校验
	scriptParams := &cloud.ScriptParams{
		SignalHost: config.GetConfig().CloudApp.SignalHost,
	}
	var shareUsernameList, sharePasswordList []string

	// MountPaths为空时，不挂载目录
	if req.MountPaths != nil {
		// validate MountPaths
		err, errResp = validateMountPaths(software.Platform, *req.MountPaths)
		if err = response.BadRequestIfError(c, err, errResp); err != nil {
			logger.Warnf("validate mount path failed, %v", err)
			return
		}

		// Windows -> "subA/xxx=X:,subB/xxx=Y:,subC/xxx=Z:"
		// Linux -> "subA/xxx=/mnt/data1,subB/xxx=/mnt/data2,subC/xxx=/mnt/data3"
		mountPathsStr, err := parseMountPaths(*req.MountPaths)
		if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, fmt.Sprintf("parse MountPaths failed, %v", err))); err != nil {
			logger.Warnf("parse mount paths to string failed, %v", err)
			return
		}

		scriptParams.ShareMountPaths = cloud.ShareMountPaths(mountPathsStr)

		// 解析mountPathsStr，分割出所有需要创建共享点的列表
		//for mountSrc := range *req.MountPaths {
		//	shareDirectory, err := util.CreateShareDirectory(s, util.CreateShareDirectoryArgs{
		//		Zone:         software.Zone.String(),
		//		AssumeUserId: userId,
		//		SubPath:      mountSrc,
		//	})
		//	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "create share directory failed")); err != nil {
		//		logger.Errorf("create share directory failed, %v", err)
		//		return
		//	}
		//
		//	shareUsernameList = append(shareUsernameList, shareDirectory.UserName)
		//	sharePasswordList = append(sharePasswordList, shareDirectory.Password)
		//}
	}

	scriptParams.ShareUsername = cloud.StringList(strings.Join(shareUsernameList, ","))
	scriptParams.SharePassword = cloud.StringList(strings.Join(sharePasswordList, ","))
	scriptParams.LoginPassword = utils.RandomPassword(24)
	var instance *models.Instance
	var sess *models.Session
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		// 用software查出该software下所有的remoteapp
		// 对每一个remoteapp创建一个随机密码
		remoteApps, err := dao.ListRemoteAppBySoftwareID(c, software.Id)
		if err != nil {
			return fmt.Errorf("list remoteApp by softwareId [%s] failed, %w", software.Id, err)
		}

		sessionId, err := rpc.GenID(ctx)
		if err != nil {
			return fmt.Errorf("generate sessionId failed, %w", err)
		}
		scriptParams.RemoteAppUserPasses, err = createRemoteAppUserPass(ctx, sessionId, remoteApps, scriptParams.LoginPassword, software.Platform)
		if err != nil {
			return fmt.Errorf("create remoteApp user pass pair failed, %w", err)
		}
		scriptParams.RoomId = sessionId.String()

		instance, err = createInstance(ctx, hardware, software, scriptParams)
		if err != nil {
			return fmt.Errorf("create instance failed, %w", err)
		}

		sess, err = createSession(ctx, sessionId, instance, userId, accountId, payByAccountId, chargeType)
		if err != nil {
			return fmt.Errorf("create session failed, %w", err)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("create instance/session in db failed, %v", err)
		return
	}

	rawInstanceId, err := s.Cloud.RunInstance(instance.Zone, software, instance, sess, hardware)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalErrorRunInstanceFailed, "run instance failed")); err != nil {
		logger.Warnf("run instance failed, %v", err)
		err = with.DefaultTransaction(c, func(ctx context.Context) error {
			return multierr.Append(
				dao.InstanceTerminated(ctx, instance.Id),
				dao.SessionClosed(ctx, sess.Id, err.Error()),
			)
		})
		if err != nil {
			logger.Error("run instance failed, update instance status to [terminated] and session status to [closed] failed")
		}
		return
	}

	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		return multierr.Append(
			dao.InstanceCreated(ctx, instance.Id, rawInstanceId),
			dao.SessionStarting(ctx, sess.Id),
		)
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warn("update instance status to [created] and session status to [starting] failed")
		return
	}

	sd, exist, err := dao.GetSessionDetailsBySessionID(c, userId, sess.Id)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get session details failed")); err != nil {
		logger.Warnf("get session details by sessionId [%s] failed, %v", sess.Id, err)
		return
	}
	if !exist {
		err = response.InternalErrorIfError(c, errors.New("session not exist"), response.WrapErrorResp(common.InternalServerErrorCode, "session not exist"))
		logger.Warnf("session details not found where session id [%s]", sess.Id)
		return
	}

	response.RenderJson(sd.ToDetailHTTPModel(), c)
}

func bindPostSessionsRequest(req *session.ApiPostRequest, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func parseMountPaths(mountPaths map[string]string) (string, error) {
	res := new(bytes.Buffer)
	for mountSrc, mountPoint := range mountPaths {
		// mountSrc进行加密
		mountSrc, err := hashid.EncodeStr(mountSrc)
		if err != nil {
			return "", err
		}
		_, err = res.WriteString(fmt.Sprintf("%s=%s,", mountSrc, mountPoint))
		if err != nil {
			return "", fmt.Errorf("bytes write string failed, %w", err)
		}
	}

	return strings.TrimSuffix(res.String(), ","), nil
}

func createSession(ctx context.Context, sessionId snowflake.ID, instance *models.Instance, userID, accountID, payByAccountID snowflake.ID, chargeType schema.ChargeType) (*models.Session, error) {
	// use sessionId to be roomId
	desktopUrl, err := webrtc.GenerateDesktopURLBase64(sessionId.String())
	if err != nil {
		return nil, fmt.Errorf("generate desktop url failed, %w", err)
	}

	sess := &models.Session{
		Id:             sessionId,
		Zone:           instance.Zone,
		UserId:         userID,
		AccountId:      accountID,
		PayByAccountId: payByAccountID,
		InstanceId:     instance.Id,
		Status:         schema.SessionPending,
		ChargeType:     chargeType,
		RoomId:         sessionId,
		DesktopUrl:     desktopUrl,
	}

	return sess, dao.CreateSession(ctx, sess)
}

// createInstance 创建一个实例模型
func createInstance(ctx context.Context, hardware *models.Hardware, software *models.Software, params *cloud.ScriptParams) (*models.Instance, error) {
	id, err := rpc.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate snowflake id failed, %w", err)
	}

	userScript, err := template.Render(software.InitScript, params)
	if err != nil {
		return nil, fmt.Errorf("render template failed, %w", err)
	}

	instance := &models.Instance{
		Id:             id,
		Zone:           software.Zone,
		HardwareId:     hardware.Id,
		SoftwareId:     software.Id,
		InitScript:     software.InitScript,
		UserParams:     utils.MustMarshalJson(params),
		UserScript:     userScript,
		InstanceId:     "-",
		InstanceData:   "{}",
		SshPassword:    params.LoginPassword,
		InstanceStatus: models.InstancePending,
	}

	return instance, dao.NewInstance(ctx, instance)
}

func createRemoteAppUserPass(ctx context.Context, sessionId snowflake.ID, remoteApps []*models.RemoteApp,
	defaultPassword string, platform models.Platform) (string, error) {
	userPasses := make([]string, 0)
	remoteAppsUserPass := make([]*models.RemoteAppUserPass, 0)
	for _, remoteApp := range remoteApps {
		remoteAppUserPass := &models.RemoteAppUserPass{
			SessionId:     sessionId,
			RemoteAppName: remoteApp.Name,
		}

		if remoteApp.LoginUser == "" {
			// 使用默认密码不需要在虚机内部二次改动密码了，则不需要写到user-data中
			remoteAppUserPass.Username = rdp.GetDefaultUsernameByPlatform(platform)
			remoteAppUserPass.Password = defaultPassword
		} else {
			remoteAppUserPass.Username = remoteApp.LoginUser
			remoteAppUserPass.Password = utils.RandomPassword(24)

			userPasses = append(userPasses, fmt.Sprintf("%s@%s", remoteAppUserPass.Username, remoteAppUserPass.Password))
		}

		remoteAppsUserPass = append(remoteAppsUserPass, remoteAppUserPass)
	}

	return strings.Join(userPasses, ","), dao.BatchInsertRemoteAppUserPass(ctx, remoteAppsUserPass)
}

func validateMountPaths(softwarePlatform models.Platform, mountPaths map[string]string) (error, response.ErrorResp) {
	if len(mountPaths) == 0 {
		return nil, response.ErrorResp{}
	}

	if len(mountPaths) > maxMountPaths {
		return fmt.Errorf("MountPaths should less than %d", maxMountPaths), response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("MountPaths should less than %d", maxMountPaths))
	}

	for mountSrc, mountPoint := range mountPaths {
		switch softwarePlatform {
		//case models.Windows:
		//	if !util.StringInSlice(mountPoint, util.WindowsMountPathPermit) {
		//		return fmt.Errorf("invalid mount point %s, should be in %v", mountPoint, util.WindowsMountPathPermit),
		//			response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("invalid mount point %s, should be in %v", mountPoint, util.WindowsMountPathPermit))
		//	}
		case models.Linux:
			if !util.IsAbsPath(mountPoint) {
				return fmt.Errorf("invalid mount point %s, should be absolute", mountPoint),
					response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("invalid mount point %s, should be absolute", mountPoint))
			}
		default:
			return fmt.Errorf("unsupported software platform"), response.WrapErrorResp(common.InternalServerErrorCode, "unsupported software platform")
		}

		if util.IsAbsPath(mountSrc) {
			return fmt.Errorf("invalid mount src %s, should be relative", mountSrc),
				response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("invalid mount src %s, should be relative", mountSrc))
		}

		err, errResp := validateMountString(mountSrc)
		if err != nil {
			return err, errResp
		}

		err, errResp = validateMountString(mountPoint)
		if err != nil {
			return err, errResp
		}
	}

	return nil, response.ErrorResp{}
}

func validateMountString(src string) (error, response.ErrorResp) {
	if strings.Contains(src, "=") || strings.Contains(src, ",") {
		return fmt.Errorf("invalid mount path %s, cannot contains =,", src),
			response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("invalid mount path %s, cannot contains =,", src))
	}

	if util.ContainsChinese(src) {
		return fmt.Errorf("mount path [%s] cannot contain Chinese", src),
			response.WrapErrorResp(common.InvalidArgumentMountPaths, fmt.Sprintf("invalid mount path %s, cannot contains Chinese", src))
	}

	return nil, response.ErrorResp{}
}

func validateSoftware(c *gin.Context, softwareId, userId snowflake.ID, isYSProductUser bool, client *openapi.Client) (*models.Software, error) {
	var (
		software *models.Software
		exist    bool
		err      error
	)
	if isYSProductUser {
		// YSProduct用户不联表查，默认全能看到
		software, exist, err = dao.GetSoftware(c, softwareId)
	} else {
		software, exist, err = dao.GetSoftwareByUser(c, softwareId, userId)
	}
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get software from db failed")); err != nil {
		return nil, fmt.Errorf("get software from db failed, %w", err)
	}
	if !exist {
		err = ErrSoftwareNotExist
		return nil, response.BadRequestIfError(c, err, response.WrapErrorResp(common.SoftwareNotFound, "software not exist"))
	}

	return software, nil
}

func validateHardware(c *gin.Context, hardwareId, userId snowflake.ID, isYSProductUser bool, client *openapi.Client) (*models.Hardware, error) {
	var (
		hardware *models.Hardware
		exist    bool
		err      error
	)
	if isYSProductUser {
		// YSProduct用户不联表查，默认全能看到
		hardware, exist, err = dao.GetHardware(c, hardwareId)
	} else {
		hardware, exist, err = dao.GetHardwareByUser(c, hardwareId, userId)
	}
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get hardware from db failed")); err != nil {
		return nil, fmt.Errorf("get hardware from db failed, %w", err)
	}
	if !exist {
		err = ErrHardwareNotExist
		_ = response.BadRequestIfError(c, err, response.WrapErrorResp(common.HardwareNotFound, "hardware not exist"))
		return nil, err
	}

	return hardware, nil
}

func checkAndEnsureChargeParams(chargeParams *schema.ChargeParams) (*schema.ChargeParams, error, string) {
	defaultChargeParams := &schema.ChargeParams{
		ChargeType: &defaultChargeType,
	}

	if chargeParams == nil {
		return defaultChargeParams, nil, ""
	}

	if chargeParams.ChargeType == nil {
		return defaultChargeParams, nil, ""
	}

	if *chargeParams.ChargeType == "" {
		return defaultChargeParams, nil, "'"
	}

	if !chargeParams.ChargeType.IsValid() {
		msg := "invalid ChargeParams.ChargeType, should be in [PrePaid | PostPaid]"
		return nil, errors.New(msg), msg
	}

	// FIXME 暂时不支持PrePaid模式，后续支持了再将校验去除
	if *chargeParams.ChargeType == schema.PrePaid {
		msg := "unsupport for [PrePaid]"
		return nil, errors.New(msg), msg
	}

	return chargeParams, nil, ""
}

func GetSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiGetRequest)
	err = bindGetSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIGetSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api get session request failed, %v", err)
		return
	}

	sessionId := snowflake.MustParseString(*req.SessionId)
	sessionHTTPModel := &schema.Session{}
	exist := true
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		sess, exists, err := dao.GetSessionDetailsBySessionID(c, userId, sessionId)
		if err != nil {
			return fmt.Errorf("get session details by SessionId [%s] failed, %w", sessionId, err)
		}
		if !exists {
			exist = false
			return nil
		}

		sessionHTTPModel = sess.ToDetailHTTPModel()
		remoteApps, err := dao.ListRemoteAppBySoftwareID(c, sess.SoftwareId)
		if err != nil {
			return fmt.Errorf("list remote app by softwareId [%s] failed, %w", sessionId, err)
		}

		for _, remoteApp := range remoteApps {
			sessionHTTPModel.RemoteApps = append(sessionHTTPModel.RemoteApps, remoteApp.ToHTTPModel())
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("database error, %v", err)
		return
	}
	if !exist {
		err = response.NotFoundIfError(c, errors.New("Session not found"), response.WrapErrorResp(common.SessionNotFound, "Session not found"))
		logger.Warnf("Session [%s] not found", sessionId)
		return
	}

	response.RenderJson(sessionHTTPModel, c)
}

func bindGetSessionRequest(req *session.ApiGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ListSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiListRequest)
	err = bindListSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIListSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api list session request failed, %v", err)
		return
	}

	listParam := &dao.ListSessionDetailParams{
		UserID:      userId,
		PageOffset:  *req.PageOffset,
		PageSize:    *req.PageSize,
		WithDeleted: false,
	}

	statusList := req.Status
	if statusList != nil && *statusList != "" {
		statuses := strings.Split(*statusList, ",")

		queryList := make([]string, 0)
		for _, s := range statuses {
			if !models.SessionStatusExist(schema.SessionStatus(s)) {
				err = response.BadRequestIfError(c, errors.New("invalid Status"),
					response.WrapErrorResp(common.InvalidArgumentSessionStatus, "invalid Status"))
				logger.Warnf("invalid status: %s", s)
				return
			}

			queryList = append(queryList, s)
		}

		listParam.Statuses = queryList
	}

	if req.Zone != nil {
		listParam.Zone = zone.Zone(*req.Zone)
	}

	sessionIdsStr := req.SessionIds
	if sessionIdsStr != nil && *sessionIdsStr != "" {
		sessionIds, err := util.ParseSnowflakeIds(*sessionIdsStr)
		if err = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentSessionIds, "invalid SessionIds")); err != nil {
			logger.Warnf("parse SessionId %s failed, %v", *sessionIdsStr, err)
			return
		}

		listParam.SessionIDs = sessionIds
	}

	sessions, total, err := dao.ListSessionDetail(c, listParam)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "list sessions from db failed")); err != nil {
		logger.Warnf("list sessions from db failed, %v", err)
		return
	}

	data := &session.ApiListResponseData{
		Sessions: make([]*schema.Session, 0),
		Offset:   listParam.PageOffset,
		Size:     listParam.PageSize,
		Total:    int(total),
	}

	// FIXME not elegant query way
	appsMap := make(map[snowflake.ID][]*models.RemoteApp)
	for _, sess := range sessions {
		apps, exist := appsMap[sess.SoftwareId]
		if !exist {
			apps, err = dao.ListRemoteAppBySoftwareID(c, sess.SoftwareId)
			if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "list remote app by SoftwareID")); err != nil {
				logger.Warnf("list remote app by SoftwareID failed, %v", err)
				return
			}
			appsMap[sess.SoftwareId] = apps
		}

		sessionHTTPModel := sess.ToDetailHTTPModel()
		for _, v := range apps {
			sessionHTTPModel.RemoteApps = append(sessionHTTPModel.RemoteApps, v.ToHTTPModel())
		}

		data.Sessions = append(data.Sessions, sessionHTTPModel)
	}

	if listParam.PageSize+listParam.PageOffset < int(total) {
		data.NextMarker = listParam.PageOffset + listParam.PageSize
	} else {
		data.NextMarker = -1
	}

	response.RenderJson(data, c)
}

func bindListSessionRequest(req *session.ApiListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func CloseSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiCloseRequest)
	err = bindCloseSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind close session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPICloseSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api close session request failed, %v", err)
		return
	}

	//s, err := util.GetState(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	sessionId := snowflake.MustParseString(*req.SessionId)
	exist, allowed := true, true
	sess := &models.SessionWithDetail{}
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		var e error
		sess, exist, e = dao.GetSessionDetailsBySessionIDWithLock(ctx, userId, sessionId)
		if e != nil {
			return fmt.Errorf("get session failed, %w", e)
		}
		if !exist {
			return nil
		}

		// check close allowed or not by status
		if sess.Status != schema.SessionStarting && sess.Status != schema.SessionStarted {
			allowed = false
			return nil
		}

		// call agent to clean custom things in instance in case to reuse boot volume
		// ignore error for compatible
		//util.OnBeforeCloseSession(logger, s, sess)

		_, e = dao.SessionUserClosing(ctx, userId, sessionId)
		if e != nil {
			return fmt.Errorf("mark session user closing failed, %w", e)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("session user closing failed, %v", err)
		return
	}
	if !exist {
		_ = response.NotFoundIfError(c, fmt.Errorf("session not found"), response.WrapErrorResp(common.SessionNotFound, "Session not found"))
		logger.Warnf("Session [%s] not found", sessionId)
		return
	}
	if !allowed {
		err = fmt.Errorf("session status not allowed to close")
		_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionUserClose, err.Error()))
		logger.Warnf("Session [%s] not allowed to close, status is [%s]", sessionId, sess.Status)
		return
	}

	response.RenderJson(nil, c)
}

func bindCloseSessionRequest(req *session.ApiCloseRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func SessionReady(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiReadyRequest)
	err = bindSessionReadyRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind session ready request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPISessionReadyRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api session ready failed, %v", err)
		return
	}

	//sessionId := snowflake.MustParseString(*req.SessionId)
	//sess, exist, err := dao.GetSession(c, userId, sessionId, false)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "check session exist failed")); err != nil {
	//	logger.Warnf("check session exist in db failed, %v", err)
	//	return
	//}
	//if !exist {
	//	err = response.NotFoundIfError(c, errors.New("session not found"), response.WrapErrorResp(common.SessionNotFound, "session not found"))
	//	logger.Warnf("session not found")
	//	return
	//}

	//isSessionReady := util.IsSessionReadyFromSignalServer(sess.RoomId, logger)
	//if !isSessionReady {
	//	logger.Infof("signal-server room [%s] not ready yet, sessionId [%s]", sess.RoomId, sessionId)
	//}

	response.RenderJson(session.ApiReadyResponseData{
		//Ready: isSessionReady,
	}, c)
}

func bindSessionReadyRequest(req *session.ApiReadyRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func DeleteSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiDeleteRequest)
	err = bindDeleteSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind delete session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIDeleteSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api delete session request failed, %v", err)
		return
	}

	sessionId := snowflake.MustParseString(*req.SessionId)

	allowed := true
	exist := true
	sess := &models.Session{}
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		var e error
		sess, exist, e = dao.GetSession(ctx, userId, sessionId, true)
		if e != nil {
			return fmt.Errorf("get session failed, %w", e)
		}
		if !exist {
			exist = false
			return nil
		}

		if sess.Status != schema.SessionClosed {
			allowed = false
			return nil
		}

		_, e = dao.DeleteSession(ctx, sess.Id)
		if e != nil {
			return fmt.Errorf("delete session failed, %w", e)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("database error, %v", err)
		return
	}
	if !exist {
		_ = response.NotFoundIfError(c, fmt.Errorf("session not found"), response.WrapErrorResp(common.SessionNotFound, "session not found"))
		logger.Warnf("session not found where id = %s", sess.Id)
		return
	}
	if !allowed {
		_ = response.ForbiddenIfError(c, fmt.Errorf("session state not allowd to be deleted"),
			response.WrapErrorResp(common.ForbiddenSessionUserDelete, "Only session at CLOSED status allowed to delete"))
		logger.Warnf("Only session at CLOSED status allowed to delete, current status: %s", sess.Status)
		return
	}

	response.RenderJson(nil, c)
}

func bindDeleteSessionRequest(req *session.ApiDeleteRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func StartSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}

	sessionaction.StartSession(
		c,
		sessionaction.WithUserId(userId),
		sessionaction.WithLogger(logger),
	)
}

func StopSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()
	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}

	sessionaction.StopSession(
		c,
		sessionaction.WithUserId(userId),
		sessionaction.WithLogger(logger),
	)
}

func RestartSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()
	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}

	sessionaction.RestartSession(
		c,
		sessionaction.WithUserId(userId),
		sessionaction.WithLogger(logger),
	)
}

// RestoreSession restore from boot volume
func RestoreSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(session.ApiRestoreRequest)
	err = bindRestoreSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind restore session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIRestoreSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate api restore session request failed, %v", err)
		return
	}

	sessionId := snowflake.MustParseString(*req.SessionId)
	logger = logger.With("session-id", sessionId.String())

	sd, err := sessionrestore.Restore(c, logger, userId, sessionId)
	if err != nil {
		return
	}

	response.RenderJson(sd.ToDetailHTTPModel(), c)
}

func bindRestoreSessionRequest(req *session.ApiRestoreRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func SessionExecScript(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	sessionaction.ExecScript(c,
		sessionaction.WithLogger(logger),
		sessionaction.WithUserId(userId),
	)
}

func SessionMount(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	sessionaction.Mount(c,
		sessionaction.WithLogger(logger),
		sessionaction.WithUserId(userId),
	)
}

func SessionUmount(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	sessionaction.Umount(c,
		sessionaction.WithLogger(logger),
		sessionaction.WithUserId(userId),
	)
}
