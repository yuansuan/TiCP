package impl

import (
	"context"
	"fmt"
	openapivisual "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/visual"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/utils"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/serializeutil"
)

func (s *VisualService) ListHardware(ctx context.Context, name string, hasUsed, isAdmin bool, cpu, mem, gpu, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Hardware, int64, error) {
	resourceIds := make([]snowflake.ID, 0)
	hardwareList := make([]*dto.Hardware, 0)

	if hasUsed {
		projectIds := make([]snowflake.ID, 0)
		if !isAdmin {
			err := s.getProjectIdsByUserID(ctx, loginUserID, &projectIds)
			if err != nil {
				return nil, 0, err
			}
		}

		sessions, _, err := s.sessionDao.GetUsedResource(ctx, projectIds, isAdmin, "")
		if err != nil {
			return nil, 0, err
		}

		for _, v := range sessions {
			resourceIds = append(resourceIds, v.HardwareID)
		}
	}

	if hasUsed {
		if len(resourceIds) == 0 {
			return hardwareList, 0, nil
		}

		// 重置分页参数
		pageIndex = 0
		pageSize = len(resourceIds)
	}

	hardwares, total, err := s.hardwareDao.ListHardware(ctx, resourceIds, name, cpu, mem, gpu, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range hardwares {
		hardwareList = append(hardwareList, utils.ConvertHardware(v))
	}

	return hardwareList, total, nil
}

func (s *VisualService) AddHardware(ctx context.Context, hardware *dto.Hardware) (string, error) {
	_, exist, err := s.hardwareDao.GetHardware(ctx, 0, hardware.Name)
	if err != nil {
		return "", err
	}
	if exist {
		return "", status.Errorf(errcode.ErrVisualHardwareHasExist, "hardware: [%v] has exist", hardware.Name)
	}

	response, err := openapivisual.AddHardware(s.api, hardware, config.GetZone())
	if err != nil {
		return "", err
	}
	if response == nil || response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	newID := s.sid.Generate()
	hardwareModelData := &model.Hardware{
		ID:             newID,
		OutHardwareID:  response.Data.HardwareId,
		Name:           hardware.Name,
		Desc:           hardware.Desc,
		Network:        hardware.Network,
		CPU:            hardware.CPU,
		Mem:            hardware.Mem,
		Gpu:            hardware.GPU,
		CPUModel:       hardware.CPUModel,
		GpuModel:       hardware.GPUModel,
		InstanceType:   hardware.InstanceType,
		InstanceFamily: hardware.InstanceFamily,
		Zone:           config.GetZone(),
	}
	err = s.hardwareDao.InsertHardware(ctx, hardwareModelData)
	if err != nil {
		return "", err
	}

	err = s.AddVisualPermission(ctx, common.PermissionResourceTypeVisualHardware, hardwareModelData)
	if err != nil {
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("add hardware: [%v]", serializeutil.GetStringForTraceLog(hardwareModelData)))

	return newID.String(), nil
}

func (s *VisualService) UpdateHardware(ctx context.Context, hardware *dto.Hardware) error {
	hardwareID := snowflake.MustParseString(hardware.ID)
	hardwareData, exist, err := s.hardwareDao.GetHardware(ctx, hardwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualHardwareNotFound, "hardware not found, hardwareID: [%v]", hardwareID)
	}

	response, err := openapivisual.UpdateHardware(s.api, hardwareData.OutHardwareID, hardware, config.GetZone())
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	hardwareModel := &model.Hardware{
		ID:             hardwareID,
		OutHardwareID:  hardwareData.OutHardwareID,
		Name:           hardware.Name,
		Desc:           hardware.Desc,
		Network:        hardware.Network,
		CPU:            hardware.CPU,
		Mem:            hardware.Mem,
		Gpu:            hardware.GPU,
		CPUModel:       hardware.CPUModel,
		GpuModel:       hardware.GPUModel,
		InstanceType:   hardware.InstanceType,
		InstanceFamily: hardware.InstanceFamily,
		Zone:           config.GetZone(),
	}
	err = s.hardwareDao.UpdateHardware(ctx, hardwareModel)
	if err != nil {
		return err
	}

	err = s.UpdateVisualPermission(ctx, common.PermissionResourceTypeVisualHardware, hardwareData, hardwareModel)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update hardware: [%v]", serializeutil.GetStringForTraceLog(hardwareModel)))

	return nil
}

func (s *VisualService) DeleteHardware(ctx context.Context, hardwareIDStr string) error {
	hardwareID := snowflake.MustParseString(hardwareIDStr)
	hardware, exist, err := s.hardwareDao.GetHardware(ctx, hardwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualHardwareNotFound, "hardware not found, hardwareID: [%v]", hardwareID)
	}

	_, total, err := s.sessionDao.HasUsedResource(ctx, -1, "", hardwareID, 0, ActiveSessionStatus)
	if err != nil {
		return err
	}
	if total > 0 {
		return status.Errorf(errcode.ErrVisualHardwareHasUsed, "hardware has been used, hardwareID: [%v]", hardwareID)
	}

	response, err := openapivisual.DeleteHardware(s.api, hardware.OutHardwareID)
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	err = s.DeleteVisualPermission(ctx, common.PermissionResourceTypeVisualHardware, &model.Hardware{ID: hardwareID})
	if err != nil {
		return err
	}

	err = s.hardwareDao.DeleteHardware(ctx, hardwareID)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("delete hardware: [%v]", serializeutil.GetStringForTraceLog(hardware)))

	return nil
}

func (s *VisualService) GetHardware(ctx context.Context, id snowflake.ID) (*dto.Hardware, error) {
	hardware, exist, err := s.hardwareDao.GetHardware(ctx, id, "")
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, status.Errorf(errcode.ErrVisualHardwareNotFound, "")
	}

	return utils.ConvertHardware(hardware), nil
}
