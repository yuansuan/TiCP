package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	hardwareApi "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/validator"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/response"
)

func GetHardWare(ctx *gin.Context) {
	logger := trace.GetLogger(ctx).Base()

	userId, err := util.GetUserId(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(hardwareApi.APIGetRequest)
	err = bindGetHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind get hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIGetHardwareRequest(req)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate request failed, %v", err)
		return
	}

	state, err := util.GetState(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	isYSProduct, err := state.IamClient.IsYsProductUser(userId)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
		logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
		return
	}

	hardwareId := snowflake.MustParseString(*req.HardwareId)
	var hardware *models.Hardware
	var exist bool
	if isYSProduct {
		// YSProduct用户不联表查，默认全能看到
		hardware, exist, err = dao.GetHardware(ctx, hardwareId)
	} else {
		hardware, exist, err = dao.GetHardwareByUser(ctx, hardwareId, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get hardware from database failed")); err != nil {
		logger.Warnf("get hardware from database failed, %v", err)
		return
	}
	if !exist {
		err = fmt.Errorf("hardware not found")
		_ = response.NotFoundIfError(ctx, err, response.WrapErrorResp(common.HardwareNotFound, "hardware not found"))
		logger.Warn(err)
		return
	}

	data := &hardwareApi.APIGetResponseData{
		Hardware: v20230530.Hardware{
			HardwareId:     hardware.Id.String(),
			Zone:           hardware.Zone.String(),
			Name:           hardware.Name,
			Desc:           hardware.Desc,
			InstanceType:   hardware.InstanceType,
			InstanceFamily: hardware.InstanceFamily,
			Network:        int(hardware.Network),
			Cpu:            int(hardware.Cpu),
			Mem:            int(hardware.Mem),
			Gpu:            int(hardware.Gpu),
			GpuModel:       hardware.GpuModel,
			CpuModel:       hardware.CpuModel,
		},
	}

	response.RenderJson(data, ctx)
}

func bindGetHardwareRequest(req *hardwareApi.APIGetRequest, c *gin.Context) error {
	if err := c.ShouldBindUri(req); err != nil {
		return fmt.Errorf("bind uri failed, %w", err)
	}

	return nil
}

func ListHardWare(ctx *gin.Context) {
	logger := trace.GetLogger(ctx).Base()

	userId, err := util.GetUserId(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.UserNotExistsErrorCode, "invalid userId")); err != nil {
		logger.Warnf("get userId failed, %v", err)
		return
	}
	logger = logger.With("user-id", userId.String())

	req := new(hardwareApi.APIListRequest)
	err = bindListHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, util.InvalidArgumentErrResp); err != nil {
		logger.Warnf("bind list hardware request failed, %v", err)
		return
	}

	err, errResp := validator.ValidateAPIListHardwareRequest(req, ctx)
	if err = response.BadRequestIfError(ctx, err, errResp); err != nil {
		logger.Warnf("validate list hardware request failed, %v", err)
		return
	}
	pageOffset, pageSize := *req.PageOffset, *req.PageSize

	state, err := util.GetState(ctx)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "get state failed")); err != nil {
		logger.Warnf("get state from gin ctx failed, %v", err)
		return
	}

	isYSProduct, err := state.IamClient.IsYsProductUser(userId)
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "check user is YSProduct or not failed")); err != nil {
		logger.Warnf("check user is YSProductUser by iam client failed, %v", err)
		return
	}

	params := ensureListHardwareDaoParams(req)
	var hardwares []*models.Hardware
	var total int64
	if isYSProduct {
		// YSProduct用户不联表查，默认全能看到
		hardwares, total, err = dao.ListHardware(ctx, params)
	} else {
		hardwares, total, err = dao.ListHardwareByUser(ctx, params, userId)
	}
	if err = response.InternalErrorIfError(ctx, err, response.WrapErrorResp(common.InternalServerErrorCode, "list hardware from database failed")); err != nil {
		logger.Warnf("list hardware from database failed, %v", err)
		return
	}

	data := &hardwareApi.APIListResponseData{
		Hardware: make([]*v20230530.Hardware, 0, len(hardwares)),
		Offset:   pageOffset,
		Size:     pageSize,
		Total:    int(total),
	}

	if params.PageSize+params.PageOffset < int(total) {
		data.NextMarker = params.PageOffset + params.PageSize
	} else {
		data.NextMarker = -1
	}

	for _, hardware := range hardwares {
		data.Hardware = append(data.Hardware, &v20230530.Hardware{
			HardwareId:     hardware.Id.String(),
			Zone:           hardware.Zone.String(),
			Name:           hardware.Name,
			Desc:           hardware.Desc,
			InstanceType:   hardware.InstanceType,
			InstanceFamily: hardware.InstanceFamily,
			Network:        int(hardware.Network),
			Cpu:            int(hardware.Cpu),
			Mem:            int(hardware.Mem),
			Gpu:            int(hardware.Gpu),
			GpuModel:       hardware.GpuModel,
			CpuModel:       hardware.CpuModel,
		})
	}

	response.RenderJson(data, ctx)
}

func bindListHardwareRequest(req *hardwareApi.APIListRequest, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return fmt.Errorf("bind query failed, %w", err)
	}

	return nil
}

func ensureListHardwareDaoParams(req *hardwareApi.APIListRequest) *dao.ListHardwareParams {
	params := &dao.ListHardwareParams{}
	if req == nil {
		return params
	}

	if req.Zone != nil {
		params.Zone = zone.Zone(*req.Zone)
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Cpu != nil {
		params.Cpu = *req.Cpu
	}

	if req.Mem != nil {
		params.Mem = *req.Mem
	}

	if req.Gpu != nil {
		params.Gpu = *req.Gpu
	}

	if req.PageSize != nil {
		params.PageSize = *req.PageSize
	}

	if req.PageOffset != nil {
		params.PageOffset = *req.PageOffset
	}

	return params
}
