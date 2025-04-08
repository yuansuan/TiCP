package sessionaction

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

type config struct {
	userId snowflake.ID
	logger *logging.Logger
}

type Option interface {
	apply(c *config)
}

type optionFunc func(c *config)

func (f optionFunc) apply(c *config) {
	f(c)
}

func WithUserId(userId snowflake.ID) Option {
	return optionFunc(func(c *config) {
		c.userId = userId
	})
}

func WithLogger(logger *logging.Logger) Option {
	return optionFunc(func(c *config) {
		c.logger = logger
	})
}

func RestartSession(c *gin.Context, opts ...Option) {
	conf := new(config)
	for _, opt := range opts {
		opt.apply(conf)
	}

	logger := conf.logger
	if logger == nil {
		logger = trace.GetLogger(c).Base()
	}

	if conf.userId != 0 {
		logger = logger.With("user-id", conf.userId.String())
	}

	req := new(session.RebootRequest)
	err := bindRestartSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind reboot session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateRestartSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate reboot session request failed, %v", err)
		return
	}
	sessionId := snowflake.MustParseString(*req.SessionId)

	state, err := util.GetState(c)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	exist, allowed := true, true
	sessionWithInstance := new(models.SessionWithInstance)
	err = with.DefaultTransaction(c, func(ctx context.Context) error {
		var e error
		sessionWithInstance, exist, e = dao.GetSessionWithInstanceBySessionId(ctx, sessionId, conf.userId, true)
		if e != nil {
			return fmt.Errorf("get session with instance by sessionId [%s] failed, %w", sessionId, e)
		}
		if !exist {
			return nil
		}

		if sessionWithInstance.Session.Status != schema.SessionStarted {
			allowed = false
			return nil
		}

		e = dao.UpdateSessionStatus(ctx, sessionId, []schema.SessionStatus{schema.SessionStarted}, schema.SessionRebooting)
		if e != nil {
			return fmt.Errorf("update session [%s] to status [REBOOTING] failed, %w", sessionId, e)
		}

		return state.Cloud.RestartInstance(sessionWithInstance.Session.Zone, sessionWithInstance.Instance.InstanceId)
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("database error, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("session [%s] not found", sessionId)
		_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.SessionNotFound, err.Error()))
		logger.Warn(err)
		return
	}
	if !allowed {
		err = fmt.Errorf("session [%s] in state [%s] now allowed to restart", sessionWithInstance.Session.Id, sessionWithInstance.Session.Status)
		_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(ensureForbiddenRestartSessionErrorCode(isAdminApi(c.FullPath())), err.Error()))
		logger.Warn(err)
		return
	}

	response.RenderJson(nil, c)
}

func bindRestartSessionRequest(req *session.RebootRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ensureForbiddenRestartSessionErrorCode(isAdmin bool) string {
	if isAdmin {
		return common.ForbiddenSessionAdminRestart
	}

	return common.ForbiddenSessionUserRestart
}

func isAdminApi(path string) bool {
	return strings.HasPrefix(path, "/admin")
}
