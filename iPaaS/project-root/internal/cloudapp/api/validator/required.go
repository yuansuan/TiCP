package validator

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/yuansuan/ticp/common/project-root-api/common"

	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

var requiredErrorCodeMap = map[string]string{
	"HardwareId":    common.InvalidArgumentHardwareId,
	"SessionId":     common.InvalidArgumentSessionId,
	"RemoteAppName": common.InvalidArgumentRemoteAppName,
	"SoftwareId":    common.InvalidArgumentSoftwareId,
	"MountPaths":    common.InvalidArgumentMountPaths,
	"PageOffset":    common.InvalidPageOffset,
	"PageSize":      common.InvalidPageSize,
	"Status":        common.InvalidArgumentSessionStatus,
	"SessionIds":    common.InvalidArgumentSessionIds,
	"Zone":          common.InvalidArgumentZone,
	"Name":          common.InvalidArgumentName,
	"Platform":      common.InvalidArgumentSoftwarePlatform,
	"Desc":          common.InvalidDesc,
	"InstanceType":  common.InvalidArgumentInstanceType,
	"Network":       common.InvalidArgumentNetwork,
	"Cpu":           common.InvalidArgumentCpu,
	"Gpu":           common.InvalidArgumentGpu,
	"Reason":        common.InvalidArgumentSessionAdminCloseReason,
	"UserId":        common.InvalidUserID,
	"ScriptContent": common.InvalidArgumentScriptContent,
	"MountPoint":    common.InvalidArgumentMountPoint,
}

func CheckRequestFieldsRequired(req interface{}, reqType reflect.Type) (error, response.ErrorResp) {
	var err error
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)

		// ignore error
		required, _ := strconv.ParseBool(field.Tag.Get("required"))
		if !required {
			continue
		}

		if reflect.ValueOf(req).Elem().Field(i).IsNil() {
			err = fmt.Errorf("[%s] is required", field.Name)
			errorCode, exist := requiredErrorCodeMap[field.Name]
			if !exist {
				errorCode = common.InvalidArgumentErrorCode
			}

			return err, response.WrapErrorResp(errorCode, err.Error())
		}
	}

	return nil, response.ErrorResp{}
}
