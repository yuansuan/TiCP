package utils

import (
	"strconv"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

func ConvertDBSession(value *v20230530.Session, ID snowflake.ID) *model.Session {
	if value == nil {
		return nil
	}
	session := &model.Session{
		ID:         ID,
		RawStatus:  value.Status,
		StreamURL:  value.StreamUrl,
		ExitReason: value.ExitReason,
		Zone:       value.Zone,
	}
	if value.Hardware != nil {
		session.OutHardwareID = value.Hardware.HardwareId
	}
	if value.Software != nil {
		session.OutSoftwareID = value.Software.SoftwareId
	}
	if value.StartTime != nil && !value.StartTime.IsZero() {
		session.StartTime = *value.StartTime
	}
	if value.EndTime != nil && !value.EndTime.IsZero() {
		session.EndTime = *value.EndTime
	}
	return session
}

func ConvertSession(value *dao.ListSessionResponse) *dto.Session {
	if value.Session == nil {
		return nil
	}
	session := &dto.Session{
		ID:          value.Session.ID.String(),
		OutAppID:    value.OutSessionID,
		UserName:    value.Session.UserName,
		ProjectName: value.ProjectName,
		Status:      value.Session.Status,
		StreamURL:   value.Session.StreamURL,
		ExitReason:  value.Session.ExitReason,
		Duration:    strconv.Itoa(int(value.Duration)),
		Hardware:    ConvertHardware(value.Hardware),
		Software:    ConvertSoftware(value.Software),
	}
	if !value.Session.StartTime.IsZero() {
		session.StartTime = value.Session.StartTime
	}
	if !value.Session.EndTime.IsZero() {
		session.EndTime = timeutil.DefaultFormatTime(value.Session.EndTime)
	}
	if !value.Session.CreateTime.IsZero() {
		session.CreateTime = value.Session.CreateTime
	}
	if !value.Session.UpdateTime.IsZero() {
		session.UpdateTime = value.Session.UpdateTime
	}
	return session
}

func Convert2Session(value *model.Session) *dto.Session {
	if value == nil {
		return nil
	}
	session := &dto.Session{
		ID:          value.ID.String(),
		OutAppID:    value.OutSessionID,
		UserName:    value.UserName,
		ProjectName: value.ProjectName,
		Status:      value.Status,
		StreamURL:   value.StreamURL,
		ExitReason:  value.ExitReason,
		Duration:    strconv.Itoa(int(value.Duration)),
	}
	if !value.StartTime.IsZero() {
		session.StartTime = value.StartTime
	}
	if !value.EndTime.IsZero() {
		session.EndTime = timeutil.DefaultFormatTime(value.EndTime)
	}
	if !value.CreateTime.IsZero() {
		session.CreateTime = value.CreateTime
	}
	if !value.UpdateTime.IsZero() {
		session.UpdateTime = value.UpdateTime
	}
	return session
}

func ConvertHardware(value *model.Hardware) *dto.Hardware {
	if value == nil {
		return nil
	}
	hardware := &dto.Hardware{
		ID:             value.ID.String(),
		Name:           value.Name,
		Desc:           value.Desc,
		Network:        value.Network,
		CPU:            value.CPU,
		Mem:            value.Mem,
		GPU:            value.Gpu,
		CPUModel:       value.CPUModel,
		GPUModel:       value.GpuModel,
		InstanceType:   value.InstanceType,
		InstanceFamily: value.InstanceFamily,
		DefaultPreset:  false,
	}
	if !value.CreateTime.IsZero() {
		hardware.CreateTime = value.CreateTime
	}
	if !value.UpdateTime.IsZero() {
		hardware.UpdateTime = value.UpdateTime
	}
	return hardware
}

func ConvertSoftware(value *model.Software) *dto.Software {
	if value == nil {
		return nil
	}
	software := &dto.Software{
		ID:         value.ID.String(),
		Name:       value.Name,
		Desc:       value.Desc,
		Platform:   value.Platform,
		ImageID:    value.ImageID,
		State:      value.State,
		InitScript: value.InitScript,
		Icon:       value.Icon,
		GPUDesired: value.GpuDesired,
		RemoteApps: nil,
	}
	if !value.CreateTime.IsZero() {
		software.CreateTime = value.CreateTime
	}
	if !value.UpdateTime.IsZero() {
		software.UpdateTime = value.UpdateTime
	}
	return software
}

func ConvertRemoteApp(value *model.RemoteApp) *dto.RemoteApp {
	if value == nil {
		return nil
	}
	remoteApp := &dto.RemoteApp{
		ID:         value.ID.String(),
		Name:       value.Name,
		Desc:       value.Desc,
		Dir:        value.Dir,
		Args:       value.Args,
		Logo:       value.Logo,
		DisableGfx: value.DisableGfx,
	}
	if !value.CreateTime.IsZero() {
		remoteApp.CreateTime = value.CreateTime
	}
	if !value.UpdateTime.IsZero() {
		remoteApp.UpdateTime = value.UpdateTime
	}
	return remoteApp
}

func GetSimpleSoftwareValue(value *model.Software) *model.Software {
	if value != nil {
		value.Icon = ""
		value.InitScript = ""
	}
	return value
}
