package validator

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

const (
	ScriptContentMaxLength = 65535
)

func ValidateAPIPostSessionsRequest(req *session.ApiPostRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post sessions request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isHardwareIdValid(req.HardwareId, false)
	if err != nil {
		return fmt.Errorf("HardwareId invalid, %w", err), errResp
	}

	err, errResp = isSoftwareIdValid(req.SoftwareId, false)
	if err != nil {
		return fmt.Errorf("SoftwareId invalid, %w", err), errResp
	}

	// no validate mount path here

	return nil, response.ErrorResp{}
}

func isSoftwareIdValid(softwareId *string, allowEmpty bool) (error, response.ErrorResp) {
	if softwareId == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *softwareId == "" {
		err = fmt.Errorf("[SoftwareId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareId, err.Error())
	}

	_, err = snowflake.ParseString(*softwareId)
	if err != nil {
		err = fmt.Errorf("parse [SoftwareId] \"%s\" to snowflake id failed, %w", *softwareId, err)
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAPIGetSessionRequest(req *session.ApiGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAPIListSessionRequest(req *session.ApiListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	// PageOffset PageSize 特殊处理一下，如果是查询参数 ?PageSize= ，gin的解析会将他解析为PageSize=0，而期望是赋默认值1000
	err, errResp = ensurePageOffset(&req.PageOffset, c)
	if err != nil {
		return fmt.Errorf("PageOffset invalid, %w", err), errResp
	}

	err, errResp = ensurePageSize(&req.PageSize, c)
	if err != nil {
		return fmt.Errorf("PageSize invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	// not validate status/sessionIds here
	return nil, response.ErrorResp{}
}

func ValidateAPICloseSessionRequest(req *session.ApiCloseRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("close session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAPISessionReadyRequest(req *session.ApiReadyRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("check session ready request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAPIDeleteSessionRequest(req *session.ApiDeleteRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("delete session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminCloseSessionRequest(req *session.AdminCloseRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("close session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	err, errResp = isSessionCloseReasonValid(req.Reason, false)
	if err != nil {
		return fmt.Errorf("Reason invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isSessionCloseReasonValid(reason *string, allowEmpty bool) (error, response.ErrorResp) {
	if reason == nil {
		return nil, response.ErrorResp{}
	}

	if !allowEmpty && *reason == "" {
		err := fmt.Errorf("[Reason] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSessionAdminCloseReason, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminListSessionRequest(req *session.AdminListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	// PageOffset PageSize 特殊处理一下，如果是查询参数 ?PageSize= ，gin的解析会将他解析为PageSize=0，而期望是赋默认值1000
	err, errResp = ensurePageOffset(&req.PageOffset, c)
	if err != nil {
		return fmt.Errorf("PageOffset invalid, %w", err), errResp
	}

	err, errResp = ensurePageSize(&req.PageSize, c)
	if err != nil {
		return fmt.Errorf("PageSize invalid, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	// not validate status/sessionIds here
	return nil, response.ErrorResp{}
}

func ValidateStartSessionRequest(req *session.PowerOnRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("start session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateStopSessionRequest(req *session.PowerOffRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("stop session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateRestartSessionRequest(req *session.RebootRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("restart session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAPIRestoreSessionRequest(req *session.ApiRestoreRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("api restore session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminRestoreSessionRequest(req *session.AdminRestoreRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("admin restore session request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	err, errResp = isUserIdValid(req.UserId, false)
	if err != nil {
		return fmt.Errorf("UserId invalid. %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isUserIdValid(userId *string, allowEmpty bool) (error, response.ErrorResp) {
	if userId == nil || (allowEmpty && *userId == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *userId == "" {
		err = fmt.Errorf("[UserId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidUserID, err.Error())
	}

	_, err = snowflake.ParseString(*userId)
	if err != nil {
		err = fmt.Errorf("parse [UserId] \"%s\" to snowflake id failed, %w", *userId, err)
		return err, response.WrapErrorResp(common.InvalidUserID, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateExecScriptRequest(req *session.ExecScriptRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("api session exec script request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	err, errResp = isScriptRunnerValid(req.ScriptRunner, true)
	if err != nil {
		return fmt.Errorf("ScriptRunner invalid, %w", err), errResp
	}

	err, errResp = isScriptContentValid(req.ScriptContent, false)
	if err != nil {
		return fmt.Errorf("ScriptContent invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

var validScriptRunnerList = []string{common.ScriptRunnerPowershell}

func isScriptRunnerValid(scriptRunner *string, allowEmpty bool) (error, response.ErrorResp) {
	if scriptRunner == nil || (allowEmpty && *scriptRunner == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *scriptRunner == "" {
		err = fmt.Errorf("[ScriptRunner] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentScriptRunner, err.Error())
	}

	if !util.StringInSlice(*scriptRunner, validScriptRunnerList) {
		err = fmt.Errorf("[ScriptRunner] invalid, must be in %s", validScriptRunnerList)
		return err, response.WrapErrorResp(common.InvalidArgumentScriptRunner, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isScriptContentValid(scriptContent *string, allowEmpty bool) (error, response.ErrorResp) {
	if scriptContent == nil || (allowEmpty && *scriptContent == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *scriptContent == "" {
		err = fmt.Errorf("[ScriptContent] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentScriptContent, err.Error())
	}

	if len(*scriptContent) > ScriptContentMaxLength {
		err = fmt.Errorf("[ScriptContent] cannot be longer than 65535")
		return err, response.WrapErrorResp(common.InvalidArgumentScriptContent, err.Error())
	}

	// check if in base64 encoded format
	if _, err = base64.StdEncoding.DecodeString(*scriptContent); err != nil {
		err = fmt.Errorf("[ScriptContent] not in base64 encoded format")
		return err, response.WrapErrorResp(common.InvalidArgumentScriptContent, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateMountRequest(req *session.MountRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("session mount request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	err, errResp = isShareDirectoryValid(req.ShareDirectory, true)
	if err != nil {
		return fmt.Errorf("ShareDirectory invalid, %w", err), errResp
	}

	err, errResp = isMountPointValid(req.MountPoint, false)
	if err != nil {
		return fmt.Errorf("MountPoint invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isShareDirectoryValid(shareDirectory *string, allowEmpty bool) (error, response.ErrorResp) {
	if shareDirectory == nil || (allowEmpty && *shareDirectory == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *shareDirectory == "" {
		err = fmt.Errorf("[ShareDirectory] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentShareDirectory, err.Error())
	}

	if util.IsAbsPath(*shareDirectory) {
		err = fmt.Errorf("[ShareDirectory] cannot be absolutely")
		return err, response.WrapErrorResp(common.InvalidArgumentShareDirectory, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isMountPointValid(mountPoint *string, allowEmpty bool) (error, response.ErrorResp) {
	if mountPoint == nil || (allowEmpty && *mountPoint == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *mountPoint == "" {
		err = fmt.Errorf("[MountPoint] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentMountPoint, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateUmountRequest(req *session.UmountRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("session umount request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSessionIdValid(req.SessionId, false)
	if err != nil {
		return fmt.Errorf("SessionId invalid, %w", err), errResp
	}

	err, errResp = isMountPointValid(req.MountPoint, false)
	if err != nil {
		return fmt.Errorf("MountPoint invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}
