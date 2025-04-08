package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type SessionDao interface {
	ListSession(ctx context.Context, userName string, outStatuses, statuses []string, hardwareIds, softwareIds []snowflake.ID, startTime, endTime int64, offset, limit int) ([]*model.Session, int64, error)
	ListSessionInfos(ctx context.Context, userName string, isAdmin bool, outStatuses, statuses []string, hardwareIds, softwareIds, ProjectIds []snowflake.ID, startTime, endTime int64, offset, limit int) ([]*ListSessionResponse, int64, error)
	InsertSessions(ctx context.Context, sessions []*model.Session) error
	UpdateSession(ctx context.Context, sessionData *model.Session) error
	DeleteSession(ctx context.Context, sessionID snowflake.ID) error
	GetSession(ctx context.Context, sessionID snowflake.ID) (*model.Session, bool, error)
	HasUsedResource(ctx context.Context, projectID snowflake.ID, username string, hardwareID, softwareID snowflake.ID, statuses []string) ([]*model.Session, int64, error)
	GetUsedResource(ctx context.Context, projectIds []snowflake.ID, isAdmin bool, username string) ([]*model.Session, int64, error)
	ListUsedProjectNames(ctx context.Context, username string) ([]string, error)
	GetSessionList(ctx context.Context, projectIds []int64) ([]*model.Session, error)

	DurationStatistics(ctx context.Context, appIDs []snowflake.ID, startTime, endTime string) ([]*Statistics, error)
	ListHistory(ctx context.Context, appIDs []string, startTime, endTime string, pageIndex, pageSize int) (int64, []*History, error)
	SessionStatistics(ctx context.Context, startTime, endTime int64, reportType, dimensionType string) ([]*StatisticItem, error)
}

type HardwareDao interface {
	ListHardware(ctx context.Context, IDs []snowflake.ID, name string, cpu, mem, gpu, offset, limit int) ([]*model.Hardware, int64, error)
	InsertHardware(ctx context.Context, hardware *model.Hardware) error
	UpdateHardware(ctx context.Context, hardware *model.Hardware) error
	DeleteHardware(ctx context.Context, hardwareID snowflake.ID) error
	GetHardware(ctx context.Context, hardwareID snowflake.ID, name string) (*model.Hardware, bool, error)
}

type SoftwareDao interface {
	ListSoftware(ctx context.Context, IDs []snowflake.ID, name, platform, state string, offset, limit int) ([]*model.Software, int64, error)
	InsertSoftware(ctx context.Context, software *model.Software) error
	UpdateSoftware(ctx context.Context, software *model.Software) error
	DeleteSoftware(ctx context.Context, softwareID snowflake.ID) error
	PublishSoftware(ctx context.Context, id snowflake.ID, state string) error
	GetSoftware(ctx context.Context, softwareID snowflake.ID, name string) (*model.Software, bool, error)
	UsingStatuses(ctx context.Context, username string, softwareIds []snowflake.ID) ([]*UsingStatusesResponse, error)

	ListRemoteApp(ctx context.Context, offset, limit int) ([]*model.RemoteApp, int64, error)
	InsertRemoteApp(ctx context.Context, remoteApp *model.RemoteApp) error
	UpdateRemoteApp(ctx context.Context, remoteApp *model.RemoteApp) error
	DeleteRemoteApp(ctx context.Context, remoteAppID snowflake.ID) error
	DeleteRemoteAppWithSoftwareID(ctx context.Context, softwareID snowflake.ID) error
	GetRemoteApp(ctx context.Context, remoteAppID snowflake.ID) (*model.RemoteApp, bool, error)
	GetRemoteApps(ctx context.Context, softwareID snowflake.ID) ([]*model.RemoteApp, error)

	InsertSoftwarePresets(ctx context.Context, presets []*model.SoftwarePreset) error
	DeleteSoftwarePresets(ctx context.Context, softwareID snowflake.ID) error
	GetSoftwarePresets(ctx context.Context, softwareID snowflake.ID) ([]*model.SoftwarePreset, error)
}
