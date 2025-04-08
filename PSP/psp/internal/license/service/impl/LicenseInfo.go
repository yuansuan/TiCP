package impl

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	licinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/license/licenseinfo"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/license/licensemanager"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

type licenseInfoServiceImpl struct {
	localAPI *openapi.OpenAPI
}

func NewLicenseInfoService() (service.LicenseInfoService, error) {
	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	return &licenseInfoServiceImpl{
		localAPI: localAPI,
	}, nil
}

func (s *licenseInfoServiceImpl) AddLicenseInfo(ctx context.Context, req *dto.LicenseInfoAddRequest) (*dto.LicenseInfoAddResponse, error) {
	logger := logging.GetLogger(ctx)
	//校验许可证名称是否重复
	managerResponse, err := licensemanager.Get(s.localAPI, req.ManagerId)
	if err != nil {
		return nil, err
	}
	for _, info := range managerResponse.Data.LicenseInfos {
		if info.Provider == req.LicenseName {
			logger.Warnf("add license info failed,LicenseName repeat")
			return nil, status.Error(errcode.ErrFailedLicenseNameRepeat, errcode.MsgFailedLicenseNameRepeat)
		}
	}

	hpcEndpoint := config.GetConfig().Local.Settings.HPCEndpoint

	//执行许可证新增
	licReq := &licinfo.AddLicenseInfoRequest{
		ManagerId:             req.ManagerId,
		Provider:              req.LicenseName,
		MacAddr:               req.MacAddr,
		ToolPath:              req.ToolPath,
		LicenseUrl:            req.LicenseUrl, //license 服务器地址
		Port:                  req.Port,
		LicenseNum:            req.LicenseNum,
		Weight:                req.Weight, //调度优先级
		BeginTime:             req.StartTime,
		EndTime:               req.EndTime,
		Auth:                  req.Auth,              //是否授权
		LicenseEnvVar:         req.LicenseEnvVar,     //环境变量
		AllowableHpcEndpoints: []string{hpcEndpoint}, //允许的提交作业的节点地址
		CollectorType:         req.CollectorType,     //license 类型
		LicenseType:           consts.Owned,          //license 默认自有
		HpcEndpoint:           hpcEndpoint,
	}

	resp, err := licenseinfo.Add(s.localAPI, licReq)
	if err != nil {
		logger.Errorf("add license info failed, err: %v", err)
		return nil, err
	}

	return &dto.LicenseInfoAddResponse{
		Id: resp.Data.Id,
	}, nil
}

func (s *licenseInfoServiceImpl) EditLicenseInfo(ctx context.Context, req *dto.LicenseInfoEditRequest) error {
	logger := logging.GetLogger(ctx)

	//校验许可证名称是否重复
	managerResponse, err := licensemanager.Get(s.localAPI, req.ManagerId)
	if err != nil {
		return err
	}
	for _, info := range managerResponse.Data.LicenseInfos {
		if info.Provider == req.LicenseName && info.Id != req.Id {
			logger.Warnf("edit license info failed,LicenseName repeat")
			return status.Error(errcode.ErrFailedLicenseNameRepeat, errcode.MsgFailedLicenseNameRepeat)
		}
	}

	hpcEndpoint := config.GetConfig().Local.Settings.HPCEndpoint

	licReq := &licinfo.PutLicenseInfoRequest{
		Id: req.Id,
		AddLicenseInfoRequest: licinfo.AddLicenseInfoRequest{
			ManagerId:             req.ManagerId,
			Provider:              req.LicenseName,
			MacAddr:               req.MacAddr,
			ToolPath:              req.ToolPath,
			LicenseUrl:            req.LicenseUrl,
			Port:                  req.Port,
			LicenseNum:            req.LicenseNum,
			Weight:                req.Weight,
			BeginTime:             req.StartTime,
			EndTime:               req.EndTime,
			Auth:                  req.Auth,
			LicenseEnvVar:         req.LicenseEnvVar,
			AllowableHpcEndpoints: []string{hpcEndpoint},
			HpcEndpoint:           hpcEndpoint,
			CollectorType:         req.CollectorType,
			LicenseType:           consts.Owned,
		},
	}

	_, err = licenseinfo.Edit(s.localAPI, licReq)
	if err != nil {
		logger.Errorf("edit license info failed, err: %v", err)
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update licenseInfo success, params:[%v]", req))
	return nil
}

func (s *licenseInfoServiceImpl) DeleteLicenseInfo(ctx context.Context, licenseID string) error {
	logger := logging.GetLogger(ctx)
	_, err := licenseinfo.Delete(s.localAPI, licenseID)
	if err != nil {
		logger.Errorf("delete license info failed, err: %v", err)
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("delete licenseInfo success, licenseID:[%v]", licenseID))
	return nil
}
