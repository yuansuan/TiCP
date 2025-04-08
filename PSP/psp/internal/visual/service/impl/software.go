package impl

import (
	"context"
	"fmt"
	openapivisual "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/visual"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/utils"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/serializeutil"
)

func (s *VisualService) ListSoftware(ctx context.Context, userId int64, name, platform, state, username string, hasPermission, hasUsed, isAdmin bool, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Software, int64, error) {
	resourceIds := make([]snowflake.ID, 0)
	softwareList := make([]*dto.Software, 0)

	if !hasUsed && hasPermission {
		permissions, err := client.GetInstance().RBAC.Permission.ListObjectResources(ctx, &rbac.ListObjectResourcesRequest{
			Id: &rbac.ObjectID{
				Id:   snowflake.ID(userId).String(),
				Type: rbac.ObjectType_USER,
			},
			ResourceType: []string{common.PermissionResourceTypeVisualSoftware},
		})
		if err != nil {
			return nil, 0, err
		}

		for _, v := range permissions.Perms {
			resourceIds = append(resourceIds, snowflake.ID(v.ResourceId))
		}

		tracelog.Info(ctx, fmt.Sprintf("user: [%d] has permission software, resourceIds: %+v", userId, resourceIds))
	}

	if !hasPermission && hasUsed {
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
			resourceIds = append(resourceIds, v.SoftwareID)
		}

		tracelog.Info(ctx, fmt.Sprintf("user: [%d] has used software, resourceIds: %+v", userId, resourceIds))
	}

	if hasPermission || hasUsed {
		if len(resourceIds) == 0 {
			return softwareList, 0, nil
		}

		// 重置分页参数
		pageIndex = 0
		pageSize = len(resourceIds)
	}

	softwares, total, err := s.softwareDao.ListSoftware(ctx, resourceIds, name, platform, state, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	presetsMap, err := s.GetPresets(ctx)
	if err != nil {
		return nil, 0, err
	}
	remoteApps, err := s.MapSoftwareIDToRemoteApps(ctx)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range softwares {
		software := utils.ConvertSoftware(v)

		software.Presets = presetsMap[v.ID]
		if software.Presets == nil {
			software.Presets = make([]*dto.Hardware, 0)
		}
		software.RemoteApps = remoteApps[v.ID]
		if software.RemoteApps == nil {
			software.RemoteApps = make([]*dto.RemoteApp, 0)
		}
		softwareList = append(softwareList, software)
	}
	return softwareList, total, nil
}

func (s *VisualService) AddSoftware(ctx context.Context, software *dto.Software) (string, error) {
	_, exist, err := s.softwareDao.GetSoftware(ctx, 0, software.Name)
	if err != nil {
		return "", err
	}
	if exist {
		return "", status.Errorf(errcode.ErrVisualSoftwareHasExist, "software: [%v] has exist", software.Name)
	}

	response, err := openapivisual.AddSoftware(s.api, software, config.GetZone())
	if err != nil {
		return "", err
	}
	if response == nil || response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	newID := s.sid.Generate()
	softwareModelData := &model.Software{
		ID:            newID,
		OutSoftwareID: response.Data.SoftwareId,
		Name:          software.Name,
		Desc:          software.Desc,
		Platform:      software.Platform,
		State:         common.Unpublished,
		ImageID:       software.ImageID,
		InitScript:    software.InitScript,
		Icon:          software.Icon,
		GpuDesired:    software.GPUDesired,
		Zone:          config.GetZone(),
	}
	err = s.softwareDao.InsertSoftware(ctx, softwareModelData)
	if err != nil {
		return "", err
	}

	err = s.AddVisualPermission(ctx, common.PermissionResourceTypeVisualSoftware, softwareModelData)
	if err != nil {
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("add software: [%v]", serializeutil.GetStringForTraceLog(utils.GetSimpleSoftwareValue(softwareModelData))))

	return newID.String(), nil
}

func (s *VisualService) UpdateSoftware(ctx context.Context, software *dto.Software) error {
	softwareID := snowflake.MustParseString(software.ID)
	softwareData, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}

	response, err := openapivisual.UpdateSoftware(s.api, softwareData.OutSoftwareID, software, config.GetZone())
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	softwareModel := &model.Software{
		ID:            softwareID,
		OutSoftwareID: softwareData.OutSoftwareID,
		Name:          software.Name,
		Desc:          software.Desc,
		Platform:      software.Platform,
		ImageID:       software.ImageID,
		InitScript:    software.InitScript,
		Icon:          software.Icon,
		GpuDesired:    software.GPUDesired,
		Zone:          config.GetZone(),
	}
	err = s.softwareDao.UpdateSoftware(ctx, softwareModel)
	if err != nil {
		return err
	}

	err = s.UpdateVisualPermission(ctx, common.PermissionResourceTypeVisualSoftware, softwareData, softwareModel)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update software: [%v]", serializeutil.GetStringForTraceLog(utils.GetSimpleSoftwareValue(softwareModel))))

	return nil
}

func (s *VisualService) DeleteSoftware(ctx context.Context, softwareIDStr string) error {
	logger := logging.GetLogger(ctx)
	softwareID := snowflake.MustParseString(softwareIDStr)
	software, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}

	_, total, err := s.sessionDao.HasUsedResource(ctx, -1, "", 0, softwareID, ActiveSessionStatus)
	if err != nil {
		return err
	}
	if total > 0 {
		return status.Errorf(errcode.ErrVisualSoftwareHasUsed, "software has been used, softwareID: [%v]", softwareID)
	}

	if software.State == common.Published {
		return status.Errorf(errcode.ErrVisualSoftwareHasPublishedFailed, "software has published, softwareID: [%v]", softwareID)
	}

	response, err := openapivisual.DeleteSoftware(s.api, software.OutSoftwareID)
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	// 删除和软件相关联的资源数据: 远程应用、 软件预设
	remoteApp, err := s.softwareDao.GetRemoteApps(ctx, softwareID)
	if err != nil {
		return status.Errorf(errcode.ErrVisualRemoteAppNotFound, "get remote app err: %v, softwareID: [%v]", err, softwareID)
	}
	for _, app := range remoteApp {
		err = s.DeleteRemoteApp(ctx, app.OutRemoteAppID)
		if err != nil {
			logger.Errorf("delete remote app err: %v, softwareID: [%v]", err, softwareID)
			continue
		}
	}
	err = s.softwareDao.DeleteRemoteAppWithSoftwareID(ctx, softwareID)
	if err != nil {
		return err
	}
	err = s.softwareDao.DeleteSoftwarePresets(ctx, softwareID)
	if err != nil {
		return err
	}

	err = s.DeleteVisualPermission(ctx, common.PermissionResourceTypeVisualSoftware, &model.Software{ID: softwareID})
	if err != nil {
		return err
	}

	err = s.softwareDao.DeleteSoftware(ctx, softwareID)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("delete software: [%v]", serializeutil.GetStringForTraceLog(utils.GetSimpleSoftwareValue(software))))

	return nil
}

func (s *VisualService) PublishSoftware(ctx context.Context, idStr string, state string) error {
	softwareID := snowflake.MustParseString(idStr)
	_, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}

	_, total, err := s.sessionDao.HasUsedResource(ctx, -1, "", 0, softwareID, NormalSessionStatus)
	if err != nil {
		return err
	}
	if total > 0 {
		return status.Errorf(errcode.ErrVisualSoftwareHasUsedForPublish, "software has been used and can not unpublished softwareID: [%v]", softwareID)
	}

	err = s.softwareDao.PublishSoftware(ctx, softwareID, state)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("publish software: [%v]", softwareID))

	return nil
}

func (s *VisualService) GetSoftwarePresets(ctx context.Context, softwareIDStr string) ([]*dto.Hardware, error) {
	softwareID := snowflake.MustParseString(softwareIDStr)
	_, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}

	presetsMap, err := s.GetPresets(ctx)
	if err != nil {
		return nil, err
	}

	presets := presetsMap[softwareID]
	if presets == nil {
		presets = make([]*dto.Hardware, 0)
	}

	return presets, nil
}

func (s *VisualService) GetPresets(ctx context.Context) (map[snowflake.ID][]*dto.Hardware, error) {
	softwareIDPresetsMap := make(map[snowflake.ID][]*dto.Hardware)

	presets, err := s.softwareDao.GetSoftwarePresets(ctx, 0)
	if err != nil {
		return nil, err
	}
	if presets == nil {
		return softwareIDPresetsMap, nil
	}

	defaultPresetMap := make(map[snowflake.ID]snowflake.ID)
	softwareIDHardwareIDsMap := make(map[snowflake.ID][]snowflake.ID)
	for _, v := range presets {
		if hardwareIDList, ok := softwareIDHardwareIDsMap[v.SoftwareID]; !ok {
			hardwareIDList = make([]snowflake.ID, 0)
			hardwareIDList = append(hardwareIDList, v.HardwareID)
			softwareIDHardwareIDsMap[v.SoftwareID] = hardwareIDList
		} else {
			hardwareIDList = append(hardwareIDList, v.HardwareID)
			softwareIDHardwareIDsMap[v.SoftwareID] = hardwareIDList
		}
		if v.Defaulted {
			defaultPresetMap[v.SoftwareID] = v.HardwareID
		}
	}

	hardwareList, _, err := s.hardwareDao.ListHardware(ctx, nil, "", 0, 0, 0, common.DefaultPageIndex, common.DefaultMaxPageSize)
	if err != nil {
		return nil, err
	}

	hardwareIDHardwareMap := make(map[snowflake.ID]*model.Hardware)
	for _, v := range hardwareList {
		hardwareIDHardwareMap[v.ID] = v
	}

	for softwareID, hardwareIDList := range softwareIDHardwareIDsMap {
		presetList := make([]*dto.Hardware, 0, len(hardwareIDList))
		for _, hardwareID := range hardwareIDList {
			if hardware, ok := hardwareIDHardwareMap[hardwareID]; ok {
				hardwareData := utils.ConvertHardware(hardware)
				presetList = append(presetList, hardwareData)

				if hardware.ID != 0 && hardware.ID == defaultPresetMap[softwareID] {
					hardwareData.DefaultPreset = true
				}
			}
		}
		softwareIDPresetsMap[softwareID] = presetList
	}

	return softwareIDPresetsMap, nil
}

func (s *VisualService) SetSoftwarePresets(ctx context.Context, softwareIDStr string, presets []*dto.SoftwarePreset) error {
	softwareID := snowflake.MustParseString(softwareIDStr)
	_, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}
	err = s.softwareDao.DeleteSoftwarePresets(ctx, softwareID)
	if err != nil {
		return err
	}

	newIDs := s.sid.Generates(int64(len(presets)))
	presetData := make([]*model.SoftwarePreset, 0)
	for i, v := range presets {
		presetData = append(presetData, &model.SoftwarePreset{
			ID:         newIDs[i],
			SoftwareID: snowflake.MustParseString(softwareIDStr),
			HardwareID: snowflake.MustParseString(v.HardwareID),
			Defaulted:  v.Default,
		})
	}

	if len(presets) > 0 {
		err = s.softwareDao.InsertSoftwarePresets(ctx, presetData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *VisualService) SoftwareUseStatuses(ctx context.Context, userId int64, username string) ([]*dto.SoftwareUsingStatus, error) {
	permissions, err := client.GetInstance().RBAC.Permission.ListObjectResources(ctx, &rbac.ListObjectResourcesRequest{
		Id: &rbac.ObjectID{
			Id:   snowflake.ID(userId).String(),
			Type: rbac.ObjectType_USER,
		},
		ResourceType: []string{common.PermissionResourceTypeVisualSoftware},
	})
	if err != nil {
		return nil, err
	}

	resourceIds := make([]snowflake.ID, 0)
	for _, v := range permissions.Perms {
		resourceIds = append(resourceIds, snowflake.ID(v.ResourceId))
	}

	usingStatuses := make([]*dto.SoftwareUsingStatus, 0)
	if len(resourceIds) == 0 {
		return usingStatuses, nil
	}

	usingStatusList, err := s.softwareDao.UsingStatuses(ctx, username, resourceIds)
	if err != nil {
		return nil, err
	}

	for _, v := range usingStatusList {
		usingStatus := &dto.SoftwareUsingStatus{
			Id:   v.Software.ID.String(),
			Name: v.Software.Name,
			Icon: v.Icon,
		}

		if v.Session != nil {
			usingStatus.SessionId = v.Session.ID.String()
			usingStatus.Status = v.Status
			usingStatus.StreamURL = v.StreamURL
		}

		usingStatuses = append(usingStatuses, usingStatus)
	}

	return usingStatuses, nil
}

func (s *VisualService) GetSoftware(ctx context.Context, id snowflake.ID) (*dto.Software, error) {
	software, exist, err := s.softwareDao.GetSoftware(ctx, id, "")
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, status.Errorf(errcode.ErrVisualSoftwareNotFound, "")
	}

	return utils.ConvertSoftware(software), nil
}
