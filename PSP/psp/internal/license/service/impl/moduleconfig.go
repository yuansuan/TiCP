package impl

import (
	"context"
	"fmt"

	licmoduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/license/moduleconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

type moduleConfigServiceImpl struct {
	localAPI *openapi.OpenAPI
}

func NewModuleConfigService() (service.ModuleConfigService, error) {
	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	return &moduleConfigServiceImpl{
		localAPI: localAPI,
	}, nil

}

func (s *moduleConfigServiceImpl) AddModuleConfig(ctx context.Context, req *dto.AddModuleConfigRequest) (*dto.AddModuleConfigResponse, error) {
	addReq := &licmoduleconfig.AddModuleConfigRequest{
		LicenseId:  req.LicenseId,
		ModuleName: req.ModuleName,
		Total:      req.Total,
	}

	moduleConfigListResp, err := moduleconfig.ListModuleConfig(s.localAPI, req.LicenseId)
	if err != nil {
		return nil, err
	}

	for _, config := range moduleConfigListResp.Data.ModuleConfigs {
		if config.ModuleName == req.ModuleName {
			return nil, status.Error(errcode.ErrFailedModuleNameRepeat, errcode.MsgFailedModuleNameRepeat)
		}
	}

	resp, err := moduleconfig.Add(s.localAPI, addReq)
	if err != nil {
		return nil, err
	}

	return &dto.AddModuleConfigResponse{
		Id: resp.Data.Id,
	}, nil
}

func (s *moduleConfigServiceImpl) EditModuleConfig(ctx context.Context, req *dto.EditModuleConfigRequest) error {

	editReq := &licmoduleconfig.PutModuleConfigRequest{
		Id:         req.Id,
		ModuleName: req.ModuleName,
		Total:      req.Total,
	}

	moduleConfigListResp, err := moduleconfig.ListModuleConfig(s.localAPI, req.LicenseId)
	if err != nil {
		return err
	}

	for _, config := range moduleConfigListResp.Data.ModuleConfigs {
		if config.ModuleName == req.ModuleName && req.Id != config.Id {
			return status.Error(errcode.ErrFailedModuleNameRepeat, errcode.MsgFailedModuleNameRepeat)
		}
	}

	_, err = moduleconfig.Edit(s.localAPI, editReq)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update moduleConfig success, params:[%v]", req))

	return nil
}

func (s *moduleConfigServiceImpl) DeleteModuleConfig(ctx context.Context, moduleConfigID string) error {

	_, err := moduleconfig.Delete(s.localAPI, moduleConfigID)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("delete moduleConfig success, moduleConfigID:[%v]", moduleConfigID))

	return nil
}

func (s *moduleConfigServiceImpl) ModuleConfigList(ctx context.Context, licenseID string) (*dto.ModuleConfigListResponse, error) {
	resp, err := moduleconfig.ListModuleConfig(s.localAPI, licenseID)
	if err != nil {
		return nil, err
	}

	var totalNum, usedNum int
	moduleConfigs := resp.Data.ModuleConfigs
	if len(moduleConfigs) == 0 {
		return &dto.ModuleConfigListResponse{}, nil
	}

	for _, moduleConfig := range moduleConfigs {
		totalNum += moduleConfig.Total
		usedNum += moduleConfig.UsedNum
	}

	moduleConfigResponses := util.ConvertToModuleConfig(moduleConfigs)
	var usedPercent string
	if totalNum <= 0 || usedNum <= 0 {
		usedPercent = "0.00%"
	} else {
		usedPercent = fmt.Sprintf("%.2f%%", float64(usedNum*100)/float64(totalNum))
	}

	return &dto.ModuleConfigListResponse{
		ModuleConfigInfos: moduleConfigResponses,
		UsedPercent:       usedPercent,
	}, nil
}
