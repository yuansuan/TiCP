package sessionrestore

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/multierr"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/template"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/webrtc"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func Restore(c *gin.Context, logger *logging.Logger, userId, sessionId snowflake.ID) (*models.SessionWithDetail, error) {
	s, err := util.GetState(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return nil, err
	}

	//计费相关校验
	accountId := snowflake.ID(0)
	if config.GetConfig().BillEnabled {
		account, err := s.OpenAPIClient.GetAccountByUserId(userId)
		if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, fmt.Sprintf("get account by userId [%s] failed", userId))); err != nil {
			logger.Warnf("get account by userId [%s] failed, %v", userId, err)
			return nil, err
		}

		// accountId必须有且合法
		accountId, err = snowflake.ParseString(account.AccountID)
		if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalErrorInvalidAccountId, fmt.Sprintf("parse accountId %s to snowflakeId failed", account.AccountID))); err != nil {
			logger.Warnf("parse accountId %s to snowflakeId failed, %v", account.AccountID, err)
			return nil, err
		}

		// 账户欠费
		if util.IsAccountInArrears(account) {
			err = fmt.Errorf("account has not enough balance, userId = %s", userId)
			_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusNotEnoughBalance, "account has not enough balance"))
			logger.Warn(err)
			return nil, err
		}

		// 账户被冻结
		if util.IsAccountFrozen(account) {
			err = fmt.Errorf("account has been frozen, userId = %s", userId)
			_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.InvalidAccountStatusFrozen, "account has been frozen"))
			logger.Warn(err)
			return nil, err
		}
	}

	sd, exist, err := dao.GetSessionDetailsBySessionID(c, userId, sessionId)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, fmt.Sprintf("get session detail from database failed, sessionId [%s]", sessionId.String()))); err != nil {
		logger.Warnf("get session detail from database failed, sessionId [%s]", sessionId.String())
		return nil, err
	}
	if !exist {
		err = fmt.Errorf("session not found where id = [%s]", sessionId.String())
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.SessionNotFound, err.Error()))
		logger.Warn(err)
		return nil, err
	}

	// 只有关掉的虚拟机允许被重建
	if sd.Status != schema.SessionClosed {
		err = fmt.Errorf("only CLOSED session can be restored while current status is %s", sd.Status)
		_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionRestore, err.Error()))
		logger.Warn(err)
		return nil, err
	}

	// 没有记录启动盘的，不允许被重建
	if sd.Instance.BootVolumeId == "" {
		err = fmt.Errorf("session cannot be restored while session bootVolumeId not recorded")
		_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionRestore, err.Error()))
		logger.Warn(err)
		return nil, err
	}

	// 同样bootVolumeId的session状态只能为CLOSED,instance状态只能为TERMINATED，否则不允许被重建
	bootVolumeOccupied, occupiedSessionId, err := dao.IsBootVolumeOccupied(c, sd.Instance.BootVolumeId)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "check boot volume is occupied by database failed")); err != nil {
		logger.Warnf("check boot volume is occupied by database failed, %v", err)
		return nil, err
	}
	if bootVolumeOccupied {
		err = fmt.Errorf("the boot volume of the session [%s] you want to restore is occupied by another session [%s]", sessionId, occupiedSessionId)
		_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionRestore, err.Error()))
		logger.Warn(err)
		return nil, err
	}

	var instance *models.Instance
	var sess *models.Session
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		newSessionId, e := rpc.GenID(ctx)
		if e != nil {
			return fmt.Errorf("generate session id failed, %w", e)
		}

		scriptParams := new(cloud.ScriptParams)
		if e = jsoniter.UnmarshalFromString(sd.UserParams, scriptParams); e != nil {
			return fmt.Errorf("unmarshal from string failed, %w", e)
		}
		scriptParams.RoomId = newSessionId.String()

		userScript, e := template.Render(sd.Software.InitScript, scriptParams)
		if e != nil {
			return fmt.Errorf("render template failed, %w", e)
		}

		instanceId, e := rpc.GenID(ctx)
		if e != nil {
			return fmt.Errorf("generate instance id failed, %w", e)
		}

		instance = &models.Instance{
			Id:             instanceId,
			Zone:           sd.Session.Zone,
			HardwareId:     sd.HardwareId,
			SoftwareId:     sd.SoftwareId,
			InitScript:     sd.Software.InitScript,
			UserParams:     utils.MustMarshalJson(scriptParams),
			UserScript:     userScript,
			SshPassword:    sd.Instance.SshPassword,
			InstanceStatus: models.InstancePending,
			BootVolumeId:   sd.Instance.BootVolumeId,
		}
		if e = dao.NewInstance(ctx, instance); e != nil {
			return fmt.Errorf("new instance failed, %w", e)
		}

		desktopUrl, e := webrtc.GenerateDesktopURLBase64(newSessionId.String())
		if e != nil {
			return fmt.Errorf("generate desktop url failed, %w", e)
		}

		sess = &models.Session{
			Id:         newSessionId,
			Zone:       sd.Session.Zone,
			UserId:     userId,
			AccountId:  accountId,
			InstanceId: instanceId,
			Status:     schema.SessionPending,
			ChargeType: sd.ChargeType,
			RoomId:     newSessionId,
			DesktopUrl: desktopUrl,
		}
		if e = dao.CreateSession(ctx, sess); e != nil {
			return fmt.Errorf("create session failed, %w", e)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("create instance/session in db failed, %v", err)
		return nil, err
	}

	rawInstanceId, err := s.Cloud.RestoreInstance(sess.Zone, &sd.Software, instance, sess, &sd.Hardware)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalErrorRunInstanceFailed, "restore instance failed")); err != nil {
		logger.Warnf("restore instance failed, %v", err)
		err = with.DefaultTransaction(c, func(ctx context.Context) error {
			return multierr.Append(
				dao.InstanceTerminated(ctx, instance.Id),
				dao.SessionClosed(ctx, sess.Id, err.Error()),
			)
		})
		if err != nil {
			logger.Error("restore instance failed, update instance status to [terminated] and session status to [closed] failed")
		}
		return nil, err
	}

	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		return multierr.Append(
			dao.InstanceCreated(ctx, instance.Id, rawInstanceId),
			dao.SessionStarting(ctx, sess.Id),
		)
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warn("update instance status to [created] and session status to [starting] failed")
		return nil, err
	}

	sd, exist, err = dao.GetSessionDetailsBySessionID(c, userId, sess.Id)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get session details failed")); err != nil {
		logger.Warnf("get session details by sessionId [%s] failed, %v", sess.Id, err)
		return nil, err
	}
	if !exist {
		err = response.InternalErrorIfError(c, errors.New("session not exist"), response.WrapErrorResp(common.InternalServerErrorCode, "session not exist"))
		logger.Warnf("session details not found where session id [%s]", sess.Id)
		return nil, err
	}

	return sd, nil
}
