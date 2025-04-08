package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	remoteappmodel "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"
	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/rdp"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func GetRemoteApp(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	userId, err := util.GetUserId(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(remoteappmodel.ApiGetRequest)
	err = bindGetRemoteAppRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get remote app request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateApiGetRemoteAppRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate get remote app request failed, %v", err)
		return
	}

	sessionId := snowflake.MustParseString(*req.SessionId)
	remoteAppName := *req.RemoteAppName
	remoteApp, exist, err := dao.GetRemoteAppByName(c, userId, sessionId, remoteAppName)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get remote app failed")); err != nil {
		logger.Warnf("get remote app by name failed, %v", err)
		return
	}
	if !exist {
		err = response.NotFoundIfError(c, fmt.Errorf("remote app not found"), response.WrapErrorResp(common.RemoteAppNotFound, "remote app not found"))
		logger.Warnf("remote app not found where sessionId: [%d], remote app name: [%s]", sessionId.Int64(), remoteAppName)
		return
	}

	state, err := util.GetState(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	sessionDetail, exist, err := dao.GetSessionDetailsBySessionID(c, userId, sessionId)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get session from database failed")); err != nil {
		logger.Warnf("get session details by session ID failed, %v", err)
		return
	}
	if !exist {
		err = response.NotFoundIfError(c, fmt.Errorf("session not found"), response.WrapErrorResp(common.SessionNotFound, "session not found"))
		logger.Warnf("session not found where session id = [%s], user id = [%s]", sessionId, userId)
		return
	}

	remoteAppUserPass, exist, err := dao.GetRemoteAppUserPass(c, sessionId, remoteAppName)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get remote app user pass failed")); err != nil {
		logger.Warnf("get remote app user pass failed, %w", err)
		return
	}
	if !exist {
		// 兼容老数据，如果不存在，则使用默认的用户密码登陆
		err = fmt.Errorf("remote app user pass not found")
		logger.Warn(err)
		remoteAppUserPass = &models.RemoteAppUserPass{
			SessionId:     sessionId,
			RemoteAppName: remoteAppName,
			Username:      rdp.GetDefaultUsernameByPlatform(sessionDetail.Software.Platform),
			Password:      sessionDetail.Instance.SshPassword,
		}
	}

	remoteAppURL, err := rdp.GenerateRemoteAppURLBase64(sessionDetail, remoteAppUserPass, remoteApp, state.Cloud)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "generate remote app url failed")); err != nil {
		logger.Warnf("generate remote app url failed, %v", err)
		return
	}

	response.RenderJson(remoteappmodel.ApiGetResponseData{
		Url: remoteAppURL,
	}, c)
}

func bindGetRemoteAppRequest(req *remoteappmodel.ApiGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}
