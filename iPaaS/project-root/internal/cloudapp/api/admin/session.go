package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/sessionaction"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/sessionrestore"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func ListSession(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(session.AdminListRequest)
	err := bindListSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminListSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin list session request failed, %v", err)
		return
	}

	listParam := &dao.ListSessionDetailParams{
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
		listParam.Zone = config.Zone(*req.Zone)
	}

	if req.WithDeleted == true {
		listParam.WithDeleted = true
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

	userIdsStr := req.UserIds
	if userIdsStr != nil && *userIdsStr != "" {
		userIds, err := util.ParseSnowflakeIds(*userIdsStr)
		if err = response.BadRequestIfError(c, err, response.WrapErrorResp(common.InvalidArgumentUserIds, "invalid UserIds")); err != nil {
			logger.Warnf("parse UserIds %s failed, %v", *userIdsStr, err)
			return
		}

		listParam.UserIDs = userIds
	}

	sessions, total, err := dao.ListSessionDetail(c, listParam)
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "list sessions from db failed")); err != nil {
		logger.Warnf("list sessions from db failed, %v", err)
		return
	}

	data := &session.AdminListResponseData{
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

		sessionHTTPModel := sess.ToAdminDetailHTTPModel()
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

func bindListSessionRequest(req *session.AdminListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func CloseSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	req := new(session.AdminCloseRequest)
	err := bindCloseSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind close session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminCloseSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin close session request failed, %v", err)
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
		sess, exist, e = dao.GetSessionDetailsBySessionIDWithLock(ctx, snowflake.ID(0), sessionId)
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

		e = dao.SessionSysClosing(ctx, sessionId, *req.Reason)
		if e != nil {
			return fmt.Errorf("mark session user closing failed, %w", e)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "database error")); err != nil {
		logger.Warnf("session sys closing failed, %v", err)
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

func bindCloseSessionRequest(req *session.AdminCloseRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func StartSession(c *gin.Context) {
	sessionaction.StartSession(c)
}

func StopSession(c *gin.Context) {
	sessionaction.StopSession(c)
}

func RestartSession(c *gin.Context) {
	sessionaction.RestartSession(c)
}

func RestoreSession(c *gin.Context) {
	logger := trace.GetLogger(c).Base()

	req := new(session.AdminRestoreRequest)
	err := bindRestoreSessionRequest(req, c)
	if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind restore session request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAdminRestoreSessionRequest(req)
	if err = response.BadRequestIfError(c, err, errResp); err != nil {
		logger.Warnf("validate admin restore session request failed, %v", err)
		return
	}

	userId := snowflake.MustParseString(*req.UserId)
	sessionId := snowflake.MustParseString(*req.SessionId)
	logger = logger.With("user-id", userId.String(), "session-id", sessionId.String())

	sd, err := sessionrestore.Restore(c, logger, userId, sessionId)
	if err != nil {
		return
	}

	response.RenderJson(sd.ToAdminDetailHTTPModel(), c)
}

func bindRestoreSessionRequest(req *session.AdminRestoreRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}

func ExecScript(c *gin.Context) {
	sessionaction.ExecScript(c)
}

func SessionMount(c *gin.Context) {
	sessionaction.Mount(c)
}

func SessionUmount(c *gin.Context) {
	sessionaction.Umount(c)
}
