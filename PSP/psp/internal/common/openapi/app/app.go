package openapiapp

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/add"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/delete"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/get"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/list"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/applist"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func ListCloudApp(api *openapi.OpenAPI) (*applist.Response, error) {
	response, err := api.Client.Job.ListAPP()

	if response != nil {
		logging.Default().Debugf("openapi list app request id: [%v]", response.RequestID)
	}

	return response, err
}

func ListApp(api *openapi.OpenAPI) (*list.Response, error) {
	response, err := api.Client.Job.AdminListAPP()

	if response != nil {
		logging.Default().Debugf("openapi list app request id: [%v]", response.RequestID)
	}

	return response, err
}

func GetApp(api *openapi.OpenAPI, appID string) (*get.Response, error) {
	if strutil.IsEmpty(appID) {
		return nil, ErrAppIDEmpty
	}

	response, err := api.Client.Job.AdminGetAPP(
		api.Client.Job.AdminGetAPP.AppID(appID),
	)

	if response != nil {
		logging.Default().Debugf("openapi get app request id: [%v], appID: [%v]", response.RequestID, appID)
	}

	return response, err
}

func AddApp(api *openapi.OpenAPI, appName, appType, appVersion, imageId, licenseManagerId, residualLogParser string, enableResidual, enableSnapshot bool, binPath map[string]string) (*add.Response, error) {
	if strutil.IsEmpty(appName) {
		return nil, ErrAppNameEmpty
	}
	if strutil.IsEmpty(appType) {
		return nil, ErrAppTypeEmpty
	}
	if strutil.IsEmpty(appVersion) {
		return nil, ErrAppVersionEmpty
	}

	if strutil.IsNotEmpty(licenseManagerId) {
		id, err := snowflake.ParseString(licenseManagerId)
		if err != nil || id < 0 {
			return nil, status.Errorf(errcode.ErrAppParamLicenseManagerIdInvalid, "add app param format err, licenseManagerId: "+licenseManagerId)
		}
	}
	if binPath != nil {
		binPathStr, err := json.Marshal(binPath)
		if err != nil {
			return nil, err
		}
		binPathStrRuneLength := utf8.RuneCountInString(string(binPathStr))
		if binPathStrRuneLength > common.StringParamLengthLimit255 {
			return nil, status.Errorf(errcode.ErrAppParamBinPathOverLengthLimit, "bin path length: [%v] over limit 255", binPathStrRuneLength)
		}
	}

	response, err := api.Client.Job.AdminAddAPP(
		api.Client.Job.AdminAddAPP.Name(appName),
		api.Client.Job.AdminAddAPP.Type(appType),
		api.Client.Job.AdminAddAPP.Version(appVersion),
		api.Client.Job.AdminAddAPP.Image(imageId),
		api.Client.Job.AdminAddAPP.LicManagerId(licenseManagerId),
		api.Client.Job.AdminAddAPP.ResidualEnable(enableResidual),
		api.Client.Job.AdminAddAPP.ResidualLogParser(residualLogParser),
		api.Client.Job.AdminAddAPP.SnapshotEnable(enableSnapshot),
		api.Client.Job.AdminAddAPP.BinPath(binPath),
		api.Client.Job.AdminAddAPP.Command(`#YS_COMMAND_PREPARED`),
		api.Client.Job.AdminAddAPP.ExtentionParams("{\"YS_MAIN_FILE\":{\"Type\":\"String\",\"ReadableName\":\"主文件\",\"Must\":true},\"YS_SKIP_MESH\":{\"Type\":\"String\",\"ReadableName\":\"跳过画网格\",\"Must\":false}}"),
	)

	if response != nil {
		logging.Default().Debugf("openapi add app request id: [%v], appName: [%v], appType: [%v], appVersion: [%v], "+
			"imageId: [%v], licenceManagerId: [%v], enableResidual: [%v], enableSnapshot: [%v], residualLogParser: [%v], binPath: [%+v]", response.RequestID,
			appName, appType, appVersion, imageId, licenseManagerId, enableResidual, enableSnapshot, residualLogParser, binPath)
	}

	// 默认从 Paas 发布, 保证后续可用
	if err == nil && response != nil && response.Data != nil && response.Data.AppID != "" {
		_, err = UpdateApp(api, response.Data.AppID, appName, appType, appVersion, imageId, licenseManagerId, residualLogParser, enableResidual, enableSnapshot, binPath, update.Published)
	}

	return response, err
}

func UpdateApp(api *openapi.OpenAPI, appID, appName, appType, appVersion, imageId, licenseManagerId, residualLogParser string, enableResidual, enableSnapshot bool, binPath map[string]string, state update.Status) (*update.Response, error) {
	if strutil.IsEmpty(appID) {
		return nil, ErrAppIDEmpty
	}
	if strutil.IsEmpty(appName) {
		return nil, ErrAppNameEmpty
	}
	if strutil.IsEmpty(appType) {
		return nil, ErrAppTypeEmpty
	}
	if strutil.IsEmpty(appVersion) {
		return nil, ErrAppVersionEmpty
	}

	if strutil.IsNotEmpty(licenseManagerId) {
		id, err := snowflake.ParseString(licenseManagerId)
		if err != nil || id < 0 {
			return nil, status.Errorf(errcode.ErrAppParamLicenseManagerIdInvalid, "add app param format err, licenseManagerId: "+licenseManagerId)
		}
	}
	if binPath != nil {
		binPathStr, err := json.Marshal(binPath)
		if err != nil {
			return nil, err
		}
		binPathStrRuneLength := utf8.RuneCountInString(string(binPathStr))
		if binPathStrRuneLength > common.StringParamLengthLimit255 {
			return nil, status.Errorf(errcode.ErrAppParamBinPathOverLengthLimit, "bin path length: [%v] over limit 255", binPathStrRuneLength)
		}
	}

	response, err := api.Client.Job.AdminUpdateAPP(
		api.Client.Job.AdminUpdateAPP.AppID(appID),
		api.Client.Job.AdminUpdateAPP.Name(appName),
		api.Client.Job.AdminUpdateAPP.Type(appType),
		api.Client.Job.AdminUpdateAPP.Version(appVersion),
		api.Client.Job.AdminUpdateAPP.Image(imageId),
		api.Client.Job.AdminUpdateAPP.LicManagerId(licenseManagerId),
		api.Client.Job.AdminUpdateAPP.ResidualEnable(enableResidual),
		api.Client.Job.AdminUpdateAPP.ResidualLogParser(residualLogParser),
		api.Client.Job.AdminUpdateAPP.SnapshotEnable(enableSnapshot),
		api.Client.Job.AdminUpdateAPP.BinPath(binPath),
		api.Client.Job.AdminUpdateAPP.PublishStatus(state),
		api.Client.Job.AdminUpdateAPP.Command(`#YS_COMMAND_PREPARED`),
		api.Client.Job.AdminUpdateAPP.ExtentionParams("{\"YS_MAIN_FILE\":{\"Type\":\"String\",\"ReadableName\":\"主文件\",\"Must\":true},\"YS_SKIP_MESH\":{\"Type\":\"String\",\"ReadableName\":\"跳过画网格\",\"Must\":false}}"),
	)

	if response != nil {
		logging.Default().Debugf("openapi update app request id: [%v], appID: [%v], appName: [%v], appType: [%v], "+
			"appVersion: [%v], imageId: [%v], licenceManagerId: [%v], enableResidual: [%v], enableSnapshot: [%v], residualLogParser: [%v], binPath: [%+v], "+
			"state: [%v]", response.RequestID, appID, appName, appType, appVersion, imageId, licenseManagerId, enableResidual, enableSnapshot, residualLogParser, binPath, state)
	}

	return response, err
}

func DeleteApp(api *openapi.OpenAPI, appID string) (*delete.Response, error) {
	if strutil.IsEmpty(appID) {
		return nil, ErrAppIDEmpty
	}

	response, err := api.Client.Job.AdminDeleteAPP(
		api.Client.Job.AdminDeleteAPP.AppID(appID),
	)

	if response != nil {
		logging.Default().Debugf("openapi delete app request id: [%v], appID: [%v]", response.RequestID, appID)
	}

	return response, err
}
