package validator

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

const (
	softwareNameMaxSize    = 64
	softwareDescMaxSize    = 255
	softwareIconMaxSize    = 255
	softwareImageIdMaxSize = 64
)

func ValidateAPIGetSoftwareRequest(req *software.APIGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get software request is nil")
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

	return nil, response.ErrorResp{}
}

func ValidateAPIListSoftwareRequest(req *software.APIListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list software request is nil")
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

	err, errResp = isSoftwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isSoftwarePlatformValid(req.Platform, true)
	if err != nil {
		return fmt.Errorf("Platform invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isSoftwareNameValid(name *string, allowEmpty bool) (error, response.ErrorResp) {
	if name == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *name == "" {
		err = fmt.Errorf("[Name] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareName, err.Error())
	}

	if len(*name) > softwareNameMaxSize {
		err = fmt.Errorf("[Name] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareName, err.Error())
	}

	return nil, response.ErrorResp{}
}

var validSoftwarePlatformList = []string{common.PlatformWindows, common.PlatformLinux}

func isSoftwarePlatformValid(platform *string, allowEmpty bool) (error, response.ErrorResp) {
	if platform == nil || (allowEmpty && *platform == "") {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *platform == "" {
		err = fmt.Errorf("[Platform] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwarePlatform, err.Error())
	}

	if !util.StringInSlice(*platform, validSoftwarePlatformList) {
		err = fmt.Errorf("[Platform] invalid")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwarePlatform, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPostSoftwareRequest(req *software.AdminPostRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post software request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isZoneValid(req.Zone, false)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isSoftwareNameValid(req.Name, false)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isSoftwareDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isSoftwareIconValid(req.Icon, true)
	if err != nil {
		return fmt.Errorf("Icon invalid, %w", err), errResp
	}

	err, errResp = isSoftwarePlatformValid(req.Platform, false)
	if err != nil {
		return fmt.Errorf("Platform invalid, %w", err), errResp
	}

	err, errResp = isSoftwareImageIdValid(req.ImageId, false)
	if err != nil {
		return fmt.Errorf("ImageId invalid, %w", err), errResp
	}

	err, errResp = isSoftwareInitScript(req.InitScript, true)
	if err != nil {
		return fmt.Errorf("InitScript invalid, %w", err), errResp
	}

	// no need to check GpuDesired
	return nil, response.ErrorResp{}
}

func isSoftwareDescValid(desc *string, allowEmpty bool) (error, response.ErrorResp) {
	if desc == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *desc == "" {
		err = fmt.Errorf("[Desc] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	if len(*desc) > softwareDescMaxSize {
		err = fmt.Errorf("[Desc] too long")
		return err, response.WrapErrorResp(common.InvalidDesc, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isSoftwareIconValid(icon *string, allowEmpty bool) (error, response.ErrorResp) {
	if icon == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *icon == "" {
		err = fmt.Errorf("[Icon] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentIcon, err.Error())
	}

	if len(*icon) > softwareIconMaxSize {
		err = fmt.Errorf("[Icon] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentIcon, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isSoftwareImageIdValid(imageId *string, allowEmpty bool) (error, response.ErrorResp) {
	if imageId == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *imageId == "" {
		err = fmt.Errorf("[ImageId] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentImageId, err.Error())
	}

	if len(*imageId) > softwareImageIdMaxSize {
		err = fmt.Errorf("[ImageId] too long")
		return err, response.WrapErrorResp(common.InvalidArgumentImageId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isSoftwareInitScript(initScript *string, allowEmpty bool) (error, response.ErrorResp) {
	if initScript == nil {
		return nil, response.ErrorResp{}
	}

	var err error
	if !allowEmpty && *initScript == "" {
		err = fmt.Errorf("[InitScript] cannot be empty")
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareInitScript, err.Error())
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminPutSoftwareRequest(req *software.AdminPutRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("put software request is nil")
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

	return ValidateAdminPostSoftwareRequest(&req.AdminPostRequest)
}

func ValidateAdminPatchSoftwareRequest(req *software.AdminPatchRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("patch software request is nil")
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

	err, errResp = isZoneValid(req.Zone, true)
	if err != nil {
		return fmt.Errorf("Zone invalid, %w", err), errResp
	}

	err, errResp = isSoftwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isSoftwareDescValid(req.Desc, true)
	if err != nil {
		return fmt.Errorf("Desc invalid, %w", err), errResp
	}

	err, errResp = isSoftwareIconValid(req.Icon, true)
	if err != nil {
		return fmt.Errorf("Icon invalid, %w", err), errResp
	}

	err, errResp = isSoftwarePlatformValid(req.Platform, true)
	if err != nil {
		return fmt.Errorf("Platform invalid, %w", err), errResp
	}

	err, errResp = isSoftwareImageIdValid(req.ImageId, true)
	if err != nil {
		return fmt.Errorf("ImageId invalid, %w", err), errResp
	}

	err, errResp = isSoftwareInitScript(req.InitScript, true)
	if err != nil {
		return fmt.Errorf("InitScript invalid, %w", err), errResp
	}

	// no need to check GpuDesired
	return nil, response.ErrorResp{}
}

func ValidateAdminGetSoftwareRequest(req *software.AdminGetRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("get software request is nil")
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

	return nil, response.ErrorResp{}
}

func ValidateAdminListSoftwareRequest(req *software.AdminListRequest, c *gin.Context) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("list software request is nil")
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

	err, errResp = isSoftwareNameValid(req.Name, true)
	if err != nil {
		return fmt.Errorf("Name invalid, %w", err), errResp
	}

	err, errResp = isSoftwarePlatformValid(req.Platform, true)
	if err != nil {
		return fmt.Errorf("Platform invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func ValidateAdminDeleteSoftwareRequest(req *software.AdminDeleteRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("delete software request is nil")
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

	return nil, response.ErrorResp{}
}

func ValidateAdminPostSoftwaresUsersRequest(req *software.AdminPostUsersRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("post softwares users request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSoftwareIdsValid(req.Softwares)
	if err != nil {
		return fmt.Errorf("Softwares invalid, %w", err), errResp
	}

	err, errResp = isUserIdsValid(req.Users)
	if err != nil {
		return fmt.Errorf("Users invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}

func isUserIdsValid(users []string) (error, response.ErrorResp) {
	if err := isSnowflakeIdsValid(users); err != nil {
		return err, response.WrapErrorResp(common.InvalidUserID, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isSoftwareIdsValid(softwares []string) (error, response.ErrorResp) {
	if err := isSnowflakeIdsValid(softwares); err != nil {
		return err, response.WrapErrorResp(common.InvalidArgumentSoftwareId, err.Error())
	}

	return nil, response.ErrorResp{}
}

func isSnowflakeIdsValid(ids []string) error {
	if ids == nil || len(ids) == 0 {
		return nil
	}

	var err error
	for _, id := range ids {
		if _, err = snowflake.ParseString(id); err != nil {
			return fmt.Errorf("parse [%s] to snowflake id failed", id)
		}
	}

	return nil
}

func ValidateAdminDeleteSoftwaresUsersRequest(req *software.AdminDeleteUsersRequest) (error, response.ErrorResp) {
	var err error
	if req == nil {
		err = fmt.Errorf("delete softwares users request is nil")
		return err, response.WrapErrorResp(common.InvalidArgumentErrorCode, err.Error())
	}

	err, errResp := CheckRequestFieldsRequired(req, reflect.TypeOf(*req))
	if err != nil {
		return fmt.Errorf("check request fields required failed, %w", err), errResp
	}

	err, errResp = isSoftwareIdsValid(req.Softwares)
	if err != nil {
		return fmt.Errorf("Softwares invalid, %w", err), errResp
	}

	err, errResp = isUserIdsValid(req.Users)
	if err != nil {
		return fmt.Errorf("Users invalid, %w", err), errResp
	}

	return nil, response.ErrorResp{}
}
