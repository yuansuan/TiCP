package validator

import (
	"fmt"
	"reflect"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"
	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

const (
	remoteAppNameMaxSize      = 64
	remoteAppDescMaxSize      = 255
	remoteAppDirMaxSize       = 255
	remoteAppLogoMaxSize      = 1024
	remoteAppLoginUserMaxSize = 64
)

func ValidateApiGetRemoteAppRequest(req *remoteapp.ApiGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get remoteapp request is nil")
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

	err, errResp = isRemoteAppNameValid(req.RemoteAppName, false)
	if err != nil {
		return fmt.Errorf("RemoteAppName invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isSessionIdValid(sessionId *string, allowEmpty bool) (error, response.ErrorResp) {
	if sessionId == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *sessionId == "" {
		err = fmt.Errorf("[SessionId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSessionId, err.Error())
	}

	_, err = snowflake.ParseString(*sessionId)
	if err != nil {
		err = fmt.Errorf("parse [SessionId] \"%s\" to snowflake id failed, %w", *sessionId, err)
		return err, response.WrapErrorResp(common.InvalidArgumentSessionId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isRemoteAppNameValid(remoteAppName *string, allowEmpty bool) (error, response.ErrorResp) {
	if remoteAppName == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *remoteAppName == "" {
		err = fmt.Errorf("[RemoteAppName] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppName, err.Error())
	}

	if len(*remoteAppName) > remoteAppNameMaxSize {
		err = fmt.Errorf("[RemoteAppName] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppName, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPostRemoteAppsRequest(req *remoteapp.AdminPostRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post remoteapp request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSoftwareIdValid(req.SoftwareId, false)
	if err != nil {
		return fmt.Errorf("SoftwareId invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppNameValid(req.Name, false)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDirValid(req.Dir, true)
	if err != nil {
		return fmt.Errorf("Dir invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLogoValid(req.Logo, true)
	if err != nil {
		return fmt.Errorf("Logo invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLoginUserValid(req.LoginUser, true)
	if err != nil {
		return fmt.Errorf("LoginUser invalid, %w", err), errResp
	}

	// no need to check disableGfx
	return nil, response.ErrorResp{}
}

func isRemoteAppDirValid(dir *string, allowEmpty bool) (error, response.ErrorResp) {
	if dir == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *dir == "" {
		err = fmt.Errorf("[Dir] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppDir, err.Error())
	}

	if len(*dir) > remoteAppDirMaxSize {
		err = fmt.Errorf("[Dir] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppDir, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isRemoteAppLogoValid(logo *string, allowEmpty bool) (error, response.ErrorResp) {
	if logo == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *logo == "" {
		err = fmt.Errorf("[Logo] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentLogo, err.Error())
	}

	if len(*logo) > remoteAppLogoMaxSize {
		err = fmt.Errorf("[Logo] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentLogo, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isRemoteAppLoginUserValid(loginUser *string, allowEmpty bool) (error, response.ErrorResp) {
	if loginUser == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *loginUser == "" {
		err = fmt.Errorf("[LoginUser] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppLoginUser, err.Error())
	}

	if len(*loginUser) > remoteAppLoginUserMaxSize {
		err = fmt.Errorf("[LoginUser] contains more than 64 characters")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppLoginUser, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isRemoteAppDescValid(desc *string, allowEmpty bool) (error, response.ErrorResp) {
	if desc == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *desc == "" {
		err = fmt.Errorf("[Desc] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	if len(*desc) > remoteAppDescMaxSize {
		err = fmt.Errorf("[Desc] too long")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPutRemoteAppRequest(req *remoteapp.AdminPutRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("put remoteapp request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isRemoteAppIdValid(req.RemoteAppId, false)
	if err != nil {
		return fmt.Errorf("RemoteAppId invalid, %w", err), errResp
	}

	err, errResp = isSoftwareIdValid(req.SoftwareId, false)
	if err != nil {
		return fmt.Errorf("SoftwareId invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppNameValid(req.Name, false)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDirValid(req.Dir, true)
	if err != nil {
		return fmt.Errorf("Dir invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLogoValid(req.Logo, true)
	if err != nil {
		return fmt.Errorf("Logo invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLoginUserValid(req.LoginUser, true)
	if err != nil {
		return fmt.Errorf("LoginUser invalid, %w", err), errResp
	}

	// no need to check disableGfx
	return nil, response.ErrorResp{}
}

func isRemoteAppIdValid(remoteAppId *string, allowEmpty bool) (error, response.ErrorResp) {
	if remoteAppId == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *remoteAppId == "" {
		err = fmt.Errorf("[RemoteAppId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppId, err.Error())
	}

	_, err = snowflake.ParseString(*remoteAppId)
	if err != nil {
		err = fmt.Errorf("parse [RemoteAppId] \"%s\" to snowflake id failed, %w", *remoteAppId, err)
		return err, response.WrapErrorResp(common.InvalidArgumentRemoteAppId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPatchRemoteAppRequest(req *remoteapp.AdminPatchRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get remoteapp request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isRemoteAppIdValid(req.RemoteAppId, false)
	if err != nil {
		return fmt.Errorf("RemoteAppId invalid, %w", err), errResp
	}

	err, errResp = isSoftwareIdValid(req.SoftwareId, true)
	if err != nil {
		return fmt.Errorf("SoftwareId invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppDirValid(req.Dir, true)
	if err != nil {
		return fmt.Errorf("Dir invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLogoValid(req.Logo, true)
	if err != nil {
		return fmt.Errorf("Logo invalid, %w", err), errResp
	}

	err, errResp = isRemoteAppLoginUserValid(req.LoginUser, true)
	if err != nil {
		return fmt.Errorf("LoginUser invalid, %w", err), errResp
	}

	// no need to check disableGfx
	return nil, response.ErrorResp{}
}

func ValidateAdminDeleteRemoteAppRequest(req *remoteapp.AdminDeleteRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get remoteapp request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isRemoteAppIdValid(req.RemoteAppId, false)
	if err != nil {
		return fmt.Errorf("RemoteAppId invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}
