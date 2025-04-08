package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ListHardware
//
//	@Summary		硬件列表
//	@Description	硬件列表接口
//	@Tags			可视化-硬件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListHardwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListHardwareResponse
//	@Router			/vis/hardware [get]
func (s *RouteService) ListHardware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListHardwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	loginUserID := ginutil.GetUserID(ctx)

	hardwares, total, err := s.visualService.ListHardware(ctx, req.Name, req.HasUsed, req.IsAdmin,
		req.CPU, req.Mem, req.GPU, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list hardware err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListHardwareFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListHardwareResponse{Hardwares: hardwares, Total: total})
}

// AddHardware
//
//	@Summary		新增硬件
//	@Description	新增硬件接口
//	@Tags			可视化-硬件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AddHardwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddHardwareResponse
//	@Router			/vis/hardware [post]
func (s *RouteService) AddHardware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AddHardwareRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Name == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] add hardware req: [%+v]", ginutil.GetUserID(ctx), req))

	id, err := s.visualService.AddHardware(ctx, &dto.Hardware{
		Name:           req.Name,
		Desc:           req.Desc,
		Network:        req.Network,
		CPU:            req.CPU,
		Mem:            req.Mem,
		GPU:            req.GPU,
		CPUModel:       req.CPUModel,
		GPUModel:       req.GPUModel,
		InstanceType:   req.InstanceType,
		InstanceFamily: req.InstanceFamily,
	})
	if err != nil {
		logger.Errorf("add hardware err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualAddHardwareFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v添加实例[%v]", ginutil.GetUserName(ctx), req.Name))

	ginutil.Success(ctx, &dto.AddHardwareResponse{ID: id})
}

// UpdateHardware
//
//	@Summary		更新硬件
//	@Description	更新硬件接口
//	@Tags			可视化-硬件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UpdateHardwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.UpdateHardwareResponse
//	@Router			/vis/hardware [put]
func (s *RouteService) UpdateHardware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.UpdateHardwareRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" || req.Name == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] update hardware req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.visualService.UpdateHardware(ctx, &dto.Hardware{
		ID:             req.ID,
		Name:           req.Name,
		Desc:           req.Desc,
		Network:        req.Network,
		CPU:            req.CPU,
		Mem:            req.Mem,
		GPU:            req.GPU,
		CPUModel:       req.CPUModel,
		GPUModel:       req.GPUModel,
		InstanceType:   req.InstanceType,
		InstanceFamily: req.InstanceFamily,
	})
	if err != nil {
		logger.Errorf("update hardware err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualUpdateHardwareFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v编辑实例[%v]", ginutil.GetUserName(ctx), req.Name))

	ginutil.Success(ctx, &dto.UpdateHardwareResponse{})
}

// DeleteHardware
//
//	@Summary		删除硬件
//	@Description	删除硬件接口
//	@Tags			可视化-硬件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.DeleteHardwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.DeleteHardwareResponse
//	@Router			/vis/hardware [delete]
func (s *RouteService) DeleteHardware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DeleteHardwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] delete hardware req: [%+v]", ginutil.GetUserID(ctx), req))

	hardware, _ := s.visualService.GetHardware(ctx, snowflake.MustParseString(req.ID))
	err := s.visualService.DeleteHardware(ctx, req.ID)
	if err != nil {
		logger.Errorf("delete hardware err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualDeleteHardwareFailed)
		return
	}

	var hardwareName string
	if hardware != nil {
		hardwareName = hardware.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v删除实例[%v]", ginutil.GetUserName(ctx), hardwareName))

	ginutil.Success(ctx, &dto.DeleteHardwareResponse{})
}
