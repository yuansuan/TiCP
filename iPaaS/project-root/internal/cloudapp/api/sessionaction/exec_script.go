package sessionaction

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func ExecScript(c *gin.Context, opts ...Option) {
	//conf := new(config)
	//for _, opt := range opts {
	//	opt.apply(conf)
	//}
	//
	//logger := conf.logger
	//if logger == nil {
	//	logger = trace.GetLogger(c).Base()
	//}
	//
	//req := new(session.ExecScriptRequest)
	//err := bindSessionExecScriptRequest(req, c)
	//if err = response.BadRequestIfError(c, err, util.InvalidArgumentErrResp); err != nil {
	//	logger.Warnf("bind restore session request failed, %v", err)
	//	return
	//}
	//
	//err, errResp := validator.ValidateExecScriptRequest(req)
	//if err = response.BadRequestIfError(c, err, errResp); err != nil {
	//	logger.Warnf("validate api session exec script request failed, %v", err)
	//	return
	//}
	//
	//sessionId := snowflake.MustParseString(*req.SessionId)
	//logger = logger.With("session-id", sessionId.String())
	//
	//// check if session exist
	//sessionDetail, exist, err := dao.GetSessionDetailsBySessionID(c, conf.userId, sessionId)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, fmt.Sprintf("get session detail from database failed, sessionId [%s]", sessionId))); err != nil {
	//	logger.Warnf("get session detail by sessionId [%s] from db failed, %v", sessionId, err)
	//	return
	//}
	//if !exist {
	//	err = fmt.Errorf("session not found where id = [%s]", sessionId)
	//	_ = response.NotFoundIfError(c, err, response.WrapErrorResp(common.SessionNotFound, err.Error()))
	//	logger.Warn(err)
	//	return
	//}
	//
	//// check if session ready to ensure session can exec script
	//isSessionReady := util.IsSessionReadyFromSignalServer(sessionDetail.RoomId, logger)
	//if !isSessionReady {
	//	err = fmt.Errorf("session not ready wehre id = [%s]", sessionId)
	//	_ = response.ForbiddenIfError(c, err, response.WrapErrorResp(common.ForbiddenSessionNotReady, err.Error()))
	//	logger.Warn(err)
	//	return
	//}
	//
	//s, err := util.GetState(c)
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
	//	logger.Warnf("get state from gin ctx failed, %v", err)
	//	return
	//}
	//
	//execScriptResp, err := util.ExecScript(logger, s, sessionDetail, req, trace.GetRequestId(c))
	//if err = response.InternalErrorIfError(c, err, response.WrapErrorResp(common.InternalServerErrorCode, "exec script failed")); err != nil {
	//	logger.Warnf("exec script failed, %v", err)
	//	return
	//}
	//
	//respData := &session.ExecScriptResponseData{
	//	ExitCode: execScriptResp.Data.ExitCode,
	//	Stdout:   execScriptResp.Data.Stdout,
	//	Stderr:   execScriptResp.Data.Stderr,
	//}
	response.RenderJson(nil, c)
}

func bindSessionExecScriptRequest(req *session.ExecScriptRequest, c *gin.Context) error {
	var err error
	if err = c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	if err = c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("bind json failed, %w", err)
	}

	return nil
}
