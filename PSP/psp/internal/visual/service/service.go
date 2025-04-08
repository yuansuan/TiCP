package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type VisualService interface {
	ListSession(ctx context.Context, hardwareIDStrs, softwareIDStrs, projectIDStrs, statuses []string, username string, isAdmin bool, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Session, int64, error)
	StartSession(ctx context.Context, projectIDStr, hardwareIDStr, softwareIDStr, username string, mounts []string, loginUserID snowflake.ID) (string, string, string, error)
	GetMountInfo(ctx context.Context, projectIDStr, username string, loginUserID snowflake.ID) (*dto.GetMountInfoResponse, error)
	PowerOffSession(ctx context.Context, sessionIDStr string) (string, error)
	PowerOnSession(ctx context.Context, sessionIDStr string) (string, error)
	CloseSession(ctx context.Context, sessionIDStr, exitReason string, admin bool) (string, error)
	ReadySession(ctx context.Context, sessionIDStr string) (bool, error)
	RebootSession(ctx context.Context, sessionIDStr, exitReason string, admin bool) (bool, error)
	GetRemoteAppURL(ctx context.Context, sessionIDStr, remoteAppName string) (string, error)
	ListUsedProjectNames(ctx context.Context, username string) ([]string, error)
	ExportSessionInfo(ctx *gin.Context, startTime, endTime int64) error

	ListHardware(ctx context.Context, name string, hasUsed, isAdmin bool, cpu, mem, gpu, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Hardware, int64, error)
	AddHardware(ctx context.Context, hardware *dto.Hardware) (string, error)
	UpdateHardware(ctx context.Context, hardware *dto.Hardware) error
	DeleteHardware(ctx context.Context, hardwareIDStr string) error

	ListSoftware(ctx context.Context, userId int64, name, platform, state, username string, hasPermission, hasUsed, isAdmin bool, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Software, int64, error)
	AddSoftware(ctx context.Context, software *dto.Software) (string, error)
	UpdateSoftware(ctx context.Context, software *dto.Software) error
	DeleteSoftware(ctx context.Context, softwareIDStr string) error
	PublishSoftware(ctx context.Context, idStr string, state string) error
	GetSoftwarePresets(ctx context.Context, softwareIDStr string) ([]*dto.Hardware, error)
	SetSoftwarePresets(ctx context.Context, softwareIDStr string, presets []*dto.SoftwarePreset) error
	SoftwareUseStatuses(ctx context.Context, userId int64, username string) ([]*dto.SoftwareUsingStatus, error)

	AddRemoteApp(ctx context.Context, softwareIDStr string, remoteApp *dto.RemoteApp) (string, error)
	UpdateRemoteApp(ctx context.Context, softwareIDStr string, remoteApp *dto.RemoteApp) error
	DeleteRemoteApp(ctx context.Context, remoteAppIDStr string) error

	DurationStatistic(ctx context.Context, appIDs []string, startTime, endTime string) ([]*dto.DurationStatistic, error)
	ListHistoryDuration(ctx context.Context, appIDs []string, startTime, endTime string, pageIndex, pageSize int) ([]*dto.HistoryDuration, int64, error)

	GetSoftware(ctx context.Context, softwareID snowflake.ID) (*dto.Software, error)
	GetHardware(ctx context.Context, hardwareID snowflake.ID) (*dto.Hardware, error)

	SessionUsageDurationStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionUsageDurationStatisticResponse, error)
	ExportUsageDurationStatistic(ctx *gin.Context, startTime, endTime int64) error
	SessionCreateNumberStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionCreateNumberStatisticResponse, error)
	SessionNumberStatusStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionNumberStatusStatisticResponse, error)
	SessionInfo(ctx context.Context, sessionID string) (*dto.Session, error)
}
