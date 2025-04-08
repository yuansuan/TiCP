package impl

import (
	"context"
	"fmt"
	openapivisual "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/visual"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/utils"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *VisualService) AddRemoteApp(ctx context.Context, softwareIDStr string, remoteApp *dto.RemoteApp) (string, error) {
	softwareID := snowflake.MustParseString(softwareIDStr)
	software, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return "", err
	}
	if !exist {
		return "", status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}

	response, err := openapivisual.AddRemoteApp(s.api, software.OutSoftwareID, remoteApp)
	if err != nil {
		return "", err
	}
	if response == nil || response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	newID := s.sid.Generate()
	err = s.softwareDao.InsertRemoteApp(ctx, &model.RemoteApp{
		ID:             newID,
		OutRemoteAppID: response.Data.Id,
		SoftwareID:     softwareID,
		OutSoftwareID:  software.OutSoftwareID,
		Name:           remoteApp.Name,
		Desc:           remoteApp.Desc,
		Dir:            remoteApp.Dir,
		Args:           remoteApp.Args,
		Logo:           remoteApp.Logo,
		DisableGfx:     remoteApp.DisableGfx,
	})
	if err != nil {
		return "", err
	}
	return newID.String(), nil
}

func (s *VisualService) UpdateRemoteApp(ctx context.Context, softwareIDStr string, remoteApp *dto.RemoteApp) error {
	softwareID := snowflake.MustParseString(softwareIDStr)
	software, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareID)
	}
	remoteAppID := snowflake.MustParseString(remoteApp.ID)
	remoteAppData, exist, err := s.softwareDao.GetRemoteApp(ctx, remoteAppID)
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualRemoteAppNotFound, "remote app not found, remoteAppID: [%v]", remoteAppID)
	}

	response, err := openapivisual.UpdateRemoteApp(s.api, remoteAppData.OutRemoteAppID, software.OutSoftwareID, remoteApp)
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	err = s.softwareDao.UpdateRemoteApp(ctx, &model.RemoteApp{
		ID:            remoteAppID,
		SoftwareID:    softwareID,
		OutSoftwareID: software.OutSoftwareID,
		Name:          remoteApp.Name,
		Desc:          remoteApp.Desc,
		Dir:           remoteApp.Dir,
		Args:          remoteApp.Args,
		Logo:          remoteApp.Logo,
		DisableGfx:    remoteApp.DisableGfx,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *VisualService) DeleteRemoteApp(ctx context.Context, remoteAppIDStr string) error {
	remoteAppID := snowflake.MustParseString(remoteAppIDStr)
	remoteApp, exist, err := s.softwareDao.GetRemoteApp(ctx, remoteAppID)
	if err != nil {
		return err
	}
	if !exist {
		return status.Errorf(errcode.ErrVisualRemoteAppNotFound, "remote app not found, remoteAppID: [%v]", remoteAppID)
	}

	_, total, err := s.sessionDao.HasUsedResource(ctx, -1, "", 0, remoteApp.SoftwareID, ActiveSessionStatus)
	if err != nil {
		return err
	}
	if total > 0 {
		return status.Errorf(errcode.ErrVisualSoftwareHasUsed, "software has used, softwareID: [%v]", remoteApp.SoftwareID)
	}

	response, err := openapivisual.DeleteRemoteApp(s.api, remoteApp.OutRemoteAppID)
	if err != nil {
		return err
	}
	if response == nil || response.ErrorCode != "" {
		return fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	err = s.softwareDao.DeleteRemoteApp(ctx, remoteAppID)
	if err != nil {
		return err
	}
	return nil
}

// MapSoftwareIDToRemoteApps 返回软件ID和远程应用之间的映射关系
func (s *VisualService) MapSoftwareIDToRemoteApps(ctx context.Context) (map[snowflake.ID][]*dto.RemoteApp, error) {
	remoteApps, _, err := s.softwareDao.ListRemoteApp(ctx, common.DefaultPageIndex, common.DefaultMaxPageSize)
	if err != nil {
		return nil, err
	}
	softwareIDToRemoteApps := make(map[snowflake.ID][]*dto.RemoteApp)
	for _, v := range remoteApps {
		softwareIDToRemoteApps[v.SoftwareID] = append(softwareIDToRemoteApps[v.SoftwareID], utils.ConvertRemoteApp(v))
	}
	return softwareIDToRemoteApps, nil
}
