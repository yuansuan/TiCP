package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	appdto "github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/util"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/app"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *GRPCService) GetAppInfoById(ctx context.Context, req *app.GetAppInfoByIdRequest) (*app.GetAppInfoByIdResponse, error) {
	logger := logging.GetLogger(ctx)

	id := snowflake.MustParseString(req.AppId)
	appInfo, err := s.appService.GetAppInfo(ctx, &appdto.GetAppInfoServiceRequest{ID: id})
	if err != nil {
		logger.Errorf("get app info err: %v, req: [%+v]", err, req)
		return nil, err
	}

	return &app.GetAppInfoByIdResponse{
		App: util.ConvertRPCAppInfo(appInfo),
	}, nil
}

func (s *GRPCService) GetAppInfoByOutAppId(ctx context.Context, req *app.GetAppInfoByOutAppIdRequest) (*app.GetAppInfoByOutAppIdResponse, error) {
	logger := logging.GetLogger(ctx)

	appInfo, err := s.appService.GetAppInfo(ctx, &appdto.GetAppInfoServiceRequest{OutAppID: req.OutAppId})
	if err != nil {
		logger.Errorf("get app info err: %v, req: [%+v]", err, req)
		return nil, err
	}

	return &app.GetAppInfoByOutAppIdResponse{
		App: util.ConvertRPCAppInfo(appInfo),
	}, nil
}

func (s *GRPCService) GetAppInfoByPrams(ctx context.Context, req *app.GetAppInfoByPramsRequest) (*app.GetAppInfoByPramsResponse, error) {
	logger := logging.GetLogger(ctx)

	if req.Type == "" || req.Version == "" || req.ComputeType == "" {
		return nil, status.Errorf(errcode.ErrInvalidParam, "param [type | version | compute_type] has empty")
	}

	appInfo, err := s.appService.GetAppInfo(ctx, &appdto.GetAppInfoServiceRequest{
		AppType:     req.Type,
		Version:     req.Version,
		ComputeType: req.ComputeType,
	})
	if err != nil {
		logger.Errorf("get app info err: %v, req: [%+v]", err, req)
		return nil, err
	}

	return &app.GetAppInfoByPramsResponse{
		App: util.ConvertRPCAppInfo(appInfo),
	}, nil
}

func (s *GRPCService) GetAppTotalNum(ctx context.Context, req *app.GetAppTotalNumRequest) (*app.GetAppTotalNumResponse, error) {
	logger := logging.GetLogger(ctx)

	count, err := s.appService.GetAppTotalNum(ctx)
	if err != nil {
		logger.Errorf("get app total err: %v, req: [%+v]", err, req)
		return nil, err
	}

	return &app.GetAppTotalNumResponse{
		Total: count,
	}, nil
}

func (s *GRPCService) CheckLicenseManagerIdUsed(ctx context.Context, req *app.CheckLicenseManagerIdUsedRequest) (*app.CheckLicenseManagerIdUsedResponse, error) {
	logger := logging.GetLogger(ctx)

	isUsed, err := s.appService.CheckLicenseManagerIdUsed(ctx, req.LicenseManagerId)
	if err != nil {
		logger.Errorf("check license manager id:[%+v] used err: [%v]", req.LicenseManagerId, err)
		return nil, err
	}

	return &app.CheckLicenseManagerIdUsedResponse{
		IsUsed: isUsed,
	}, nil
}
