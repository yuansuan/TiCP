package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/get"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/license/licensemanager"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/app"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

type licenseManagerServiceImpl struct {
	localAPI *openapi.OpenAPI
}

func NewLicenseManagerService() (service.LicenseManagerService, error) {
	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	return &licenseManagerServiceImpl{
		localAPI: localAPI,
	}, nil
}

func (s *licenseManagerServiceImpl) LicenseManagerList(ctx context.Context, licenseType string) (*dto.LicenseManagerListResponse, error) {
	logger := logging.GetLogger(ctx)

	licManagerList, err := licensemanager.ListLicenseManage(s.localAPI)
	if err != nil {
		logger.Errorf("list license manager list err: %v", err)
		return nil, err
	}

	//在内存中过滤
	var resultItems []*licmanager.GetLicManagerResponseData
	for _, item := range licManagerList.Data.Items {
		if strings.Contains(strings.ToLower(item.AppType), strings.ToLower(licenseType)) {
			resultItems = append(resultItems, item)
		}
	}

	licenseManagers := util.ConvertToLicManagerListResp(resultItems)

	return &dto.LicenseManagerListResponse{
		LicenseManagers: licenseManagers,
		Total:           len(resultItems),
	}, nil
}

func (s *licenseManagerServiceImpl) LicenseManagerInfo(ctx context.Context, managerID string) (*dto.LicenseManagerData, error) {
	logger := logging.GetLogger(ctx)

	getLicenseManager := s.localAPI.Client.License.GetLicenseManager
	options := []get.Option{
		getLicenseManager.Id(managerID),
	}
	response, err := getLicenseManager(options...)
	if err != nil {
		logger.Errorf("get license manager info err: %v", err)
		return nil, err
	}

	licenseManager := util.Convert2LicManagerData(response.Data)

	return licenseManager, nil
}

func (s *licenseManagerServiceImpl) LicenseTypeList(ctx context.Context) (*dto.LicenseTypeListResponse, error) {
	logger := logging.GetLogger(ctx)

	licManagerList, err := s.localAPI.Client.License.ListLicenseManager()
	if err != nil {
		logger.Errorf("list license manager list err: %v", err)
		return nil, err
	}

	//在内存中过滤
	var resultItems []*dto.LicenseTypeInfo
	for _, item := range licManagerList.Data.Items {
		licenseManager := &dto.LicenseTypeInfo{
			Id:          item.Id,
			LicenseType: item.AppType,
		}
		resultItems = append(resultItems, licenseManager)
	}

	return &dto.LicenseTypeListResponse{
		LicenseTypeInfos: resultItems,
	}, nil
}

func (s *licenseManagerServiceImpl) AddLicenseManager(ctx context.Context, req *dto.AddLicenseManagerRequest) (*dto.AddLicenseManagerResponse, error) {
	logger := logging.GetLogger(ctx)

	//校验license类型是否重复
	licenseManage, err := licensemanager.ListLicenseManage(s.localAPI)
	if err != nil {
		logger.Errorf("get license manager list err: %v", err)
		return nil, err
	}
	for _, item := range licenseManage.Data.Items {
		if item.AppType == req.AppType {
			logger.Warnf("add license manager failed,AppType repeat")
			return nil, status.Error(errcode.ErrFailedAppTypeRepeat, errcode.MsgFailedAppTypeRepeat)
		}
	}

	//新增license信息
	licReq := &licmanager.AddLicManagerRequest{
		AppType:     req.AppType,
		Os:          req.Os,
		Desc:        req.Desc,
		ComputeRule: req.ComputeRule,
	}

	resp, err := licensemanager.Add(s.localAPI, licReq)
	if err != nil {
		return nil, err
	}

	// 修改license状态为发布
	editReq := &licmanager.PutLicManagerRequest{
		Id:     resp.Data.Id,
		Status: consts.Publish,
		AddLicManagerRequest: licmanager.AddLicManagerRequest{
			AppType:     req.AppType,
			Os:          req.Os,
			Desc:        req.Desc,
			ComputeRule: req.ComputeRule,
		},
	}

	_, err = licensemanager.Edit(s.localAPI, editReq)
	if err != nil {
		return nil, err
	}

	return &dto.AddLicenseManagerResponse{
		Id: resp.Data.Id,
	}, nil
}

func (s *licenseManagerServiceImpl) EditLicenseManager(ctx context.Context, req *dto.EditLicenseManagerRequest) error {
	logger := logging.GetLogger(ctx)

	var licManagerStatus int
	//校验license类型是否重复
	licenseManage, err := licensemanager.ListLicenseManage(s.localAPI)
	if err != nil {
		logger.Errorf("get license manager list err: %v", err)
		return err
	}
	for _, item := range licenseManage.Data.Items {
		if item.AppType == req.AppType && item.Id != req.Id {
			logger.Warnf("edit license manager failed,AppType repeat")
			return status.Error(errcode.ErrFailedAppTypeRepeat, errcode.MsgFailedAppTypeRepeat)
		}
		if item.Id == req.Id {
			licManagerStatus = item.Status
		}
	}

	licReq := &licmanager.PutLicManagerRequest{
		Id:     req.Id,
		Status: licManagerStatus,
		AddLicManagerRequest: licmanager.AddLicManagerRequest{
			ComputeRule: req.ComputeRule,
			AppType:     req.AppType,
			Os:          req.Os,
			Desc:        req.Desc,
		},
	}

	_, err = licensemanager.Edit(s.localAPI, licReq)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update licenseManager success, params:[%v]", req))

	return nil
}

func (s *licenseManagerServiceImpl) DeleteLicenseManager(ctx context.Context, managerID string) error {
	req := &app.CheckLicenseManagerIdUsedRequest{
		LicenseManagerId: managerID,
	}

	resp, err := client.GetInstance().App.CheckLicenseManagerIdUsed(ctx, req)
	if err != nil {
		return errors.Wrap(err, "check license manager id used error")
	}

	if resp.IsUsed {
		return status.Error(errcode.ErrFailedLicenseManagerDeleteBind, errcode.MsgFailedLicenseManagerDeleteBind)
	}

	_, err = licensemanager.Delete(s.localAPI, managerID)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("delete licenseManager success, managerID:[%v]", managerID))

	return nil
}
