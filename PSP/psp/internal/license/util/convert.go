package util

import (
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"

	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
)

func ConvertToModuleConfig(moduleConfigs []*moduleconfig.GetModuleConfigResponseData) []*dto.ModuleConfigResponse {
	if len(moduleConfigs) == 0 {
		return []*dto.ModuleConfigResponse{}
	}

	resp := make([]*dto.ModuleConfigResponse, 0, len(moduleConfigs))
	for _, config := range moduleConfigs {
		moduleConfig := &dto.ModuleConfigResponse{
			Id:         config.Id,
			ModuleName: config.ModuleName,
			Total:      config.Total,
			UsedNum:    config.UsedNum,
			FreeNum:    config.Total - config.UsedNum,
		}

		resp = append(resp, moduleConfig)
	}
	return resp
}

func ConvertToLicenseInfo(licenseInfos []*licenseinfo.GetLicenseInfoResponseData) []*dto.LicenseInfoResponse {
	if len(licenseInfos) == 0 {
		return []*dto.LicenseInfoResponse{}
	}

	resp := make([]*dto.LicenseInfoResponse, 0, len(licenseInfos))
	for _, licenseInfo := range licenseInfos {
		licenseInfoResp := &dto.LicenseInfoResponse{
			Id:                    licenseInfo.Id,
			ManagerId:             licenseInfo.ManagerId,
			LicenseName:           licenseInfo.Provider,
			MacAddr:               licenseInfo.MacAddr,
			ToolPath:              licenseInfo.ToolPath,
			LicenseUrl:            licenseInfo.LicenseUrl,
			Port:                  licenseInfo.Port,
			LicenseNum:            licenseInfo.LicenseNum,
			Weight:                licenseInfo.Weight,
			BeginTime:             licenseInfo.BeginTime,
			EndTime:               licenseInfo.EndTime,
			Auth:                  licenseInfo.Auth,
			LicenseEnvVar:         licenseInfo.LicenseEnvVar,
			AllowableHpcEndpoints: licenseInfo.AllowableHpcEndpoints,
			CollectorType:         licenseInfo.CollectorType,
			ModuleConfigInfos:     ConvertToModuleConfig(licenseInfo.ModuleConfigs),
		}

		resp = append(resp, licenseInfoResp)
	}

	return resp
}

func ConvertToLicManagerListResp(licenseManagerList []*licmanager.GetLicManagerResponseData) []*dto.LicenseManagerResponse {
	if len(licenseManagerList) == 0 {
		return []*dto.LicenseManagerResponse{}
	}

	resp := make([]*dto.LicenseManagerResponse, 0, len(licenseManagerList))
	for _, licManager := range licenseManagerList {
		licManagerResp := &dto.LicenseManagerResponse{
			Id:          licManager.Id,
			LicenseType: licManager.AppType,
			Status:      licManager.Status,
			Os:          licManager.Os,
			Desc:        licManager.Desc,
			ComputeRule: licManager.ComputeRule,
			CreateTime:  licManager.CreateTime,
		}

		resp = append(resp, licManagerResp)
	}

	return resp
}

func Convert2LicManagerData(licenseManager *licmanager.GetLicManagerResponseData) *dto.LicenseManagerData {

	licManagerResp := &dto.LicenseManagerData{
		Id:           licenseManager.Id,
		AppType:      licenseManager.AppType,
		Os:           licenseManager.Os,
		Desc:         licenseManager.Desc,
		ComputeRule:  licenseManager.ComputeRule,
		CreateTime:   licenseManager.CreateTime,
		LicenseInfos: ConvertToLicenseInfo(licenseManager.LicenseInfos),
	}

	return licManagerResp
}
