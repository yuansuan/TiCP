package dao

import (
	"context"
	"fmt"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type sessionDaoImpl struct{}

func NewSessionDao() SessionDao {
	return &sessionDaoImpl{}
}

func (d *sessionDaoImpl) InsertSessions(ctx context.Context, sessions []*model.Session) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.Insert(sessions)
	if err != nil {
		return err
	}
	return nil
}

func (d *sessionDaoImpl) ListSession(ctx context.Context, userName string, outStatuses, statuses []string, hardwareIds, softwareIds []snowflake.ID, startTime, endTime int64, offset, limit int) ([]*model.Session, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var response []*model.Session
	if len(outStatuses) > 0 {
		session.In("raw_status", outStatuses)
	}
	if len(statuses) > 0 {
		session.In("status", statuses)
	}
	if len(hardwareIds) > 0 {
		session.In("hardware_id", hardwareIds)
	}
	if len(softwareIds) > 0 {
		session.In("software_id", softwareIds)
	}
	if userName != "" {
		session.Where("user_name = ?", userName)
	}
	if startTime > 0 {
		session.Where("start_time >= ?", time.UnixMilli(startTime))
	}
	if endTime > 0 {
		session.Where("end_time <= ?", time.UnixMilli(endTime))
	}
	if offset > 0 {
		session.Limit(int(limit), int((offset-1)*limit))
	} else {
		session.Limit(int(limit))
	}

	total, err := session.Desc("update_time").FindAndCount(&response)
	if err != nil {
		return nil, 0, err
	}
	return response, total, nil
}

type ListSessionResponse struct {
	*model.Session  `xorm:"extends"`
	*model.Software `xorm:"extends"`
	*model.Hardware `xorm:"extends"`
}

func (d *sessionDaoImpl) ListSessionInfos(ctx context.Context, userName string, isAdmin bool, outStatuses, statuses []string, hardwareIds, softwareIds, projectIds []snowflake.ID, startTime, endTime int64, offset, limit int) ([]*ListSessionResponse, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var response []*ListSessionResponse
	session = session.Table(model.SessionTableName).Alias("sr").
		Join("LEFT", model.SoftwareTableName+" as s", "sr.software_id = s.id").
		Join("LEFT", model.HardwareTableName+" as h", "sr.hardware_id = h.id")

	if len(outStatuses) > 0 {
		session.In("sr.raw_status", outStatuses)
	}
	if len(statuses) > 0 {
		session.In("sr.status", statuses)
	}
	if len(hardwareIds) > 0 {
		session.In("sr.hardware_id", hardwareIds)
	}
	if len(softwareIds) > 0 {
		session.In("sr.software_id", softwareIds)
	}
	if len(projectIds) > 0 {
		session.In("project_id", projectIds)
	}
	if startTime > 0 {
		session.Where("sr.start_time >= ?", time.UnixMilli(startTime))
	}
	if endTime > 0 {
		session.Where("sr.end_time <= ?", time.UnixMilli(endTime))
	}
	if offset > 0 {
		session.Limit(limit, (offset-1)*limit)
	} else {
		session.Limit(limit)
	}

	if !isAdmin {
		session.Where("not (project_id = ? and user_name != ?)", common.PersonalProjectID, userName)
	}

	total, err := session.Select("sr.*, s.*, h.*").Desc("sr.create_time").FindAndCount(&response)
	if err != nil {
		return nil, 0, err
	}
	return response, total, nil
}

func (d *sessionDaoImpl) UpdateSession(ctx context.Context, sessionData *model.Session) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(sessionData.ID).UseBool().Update(sessionData)
	if err != nil {
		return err
	}
	return nil
}

func (d *sessionDaoImpl) DeleteSession(ctx context.Context, sessionID snowflake.ID) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	_, err := session.ID(sessionID).Delete(&model.Session{})
	if err != nil {
		return err
	}
	return nil
}

func (d *sessionDaoImpl) GetSession(ctx context.Context, sessionID snowflake.ID) (*model.Session, bool, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	sessionData := &model.Session{}
	exist, err := session.ID(sessionID).Get(sessionData)
	if err != nil {
		return nil, exist, err
	}
	return sessionData, exist, nil
}

func (d *sessionDaoImpl) HasUsedResource(ctx context.Context, projectID snowflake.ID, username string, hardwareID, softwareID snowflake.ID, statuses []string) ([]*model.Session, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	sessions := make([]*model.Session, 0)
	if projectID >= 0 {
		session.Where("project_id = ?", projectID)
	}
	if username != "" {
		session.Where("user_name = ?", username)
	}
	if hardwareID > 0 {
		session.Where("hardware_id = ?", hardwareID)
	}
	if softwareID > 0 {
		session.Where("software_id = ?", softwareID)
	}
	if len(statuses) > 0 {
		session.In("status", statuses)
	}

	total, err := session.FindAndCount(&sessions)
	if err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

func (d *sessionDaoImpl) GetUsedResource(ctx context.Context, projectIds []snowflake.ID, isAdmin bool, username string) ([]*model.Session, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	sessions := make([]*model.Session, 0)

	if username != "" {
		session.Where("user_name = ?", username)
	}
	if len(projectIds) > 0 {
		session.In("project_id", projectIds)
	}

	if !isAdmin {
		session.Where("not (project_id = ? and user_name != ?)", common.PersonalProjectID, username)
	}

	total, err := session.Distinct("hardware_id, software_id").FindAndCount(&sessions)
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

type Statistics struct {
	SoftwareID   snowflake.ID `xorm:"software_id"`
	SoftwareName string       `xorm:"name"`
	Duration     int64        `xorm:"duration"`
}

func (d *sessionDaoImpl) DurationStatistics(ctx context.Context, appIDs []snowflake.ID, startTime, endTime string) ([]*Statistics, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	var statistics []*Statistics
	session = session.Table(model.SessionTableName).Alias("sr").
		Join("LEFT", model.SoftwareTableName+" as s", "sr.software_id = s.id")
	if len(appIDs) > 0 {
		session.In("sr.software_id", appIDs)
	}
	if startTime != "" && endTime != "" {
		start, err := time.Parse(common.DatetimeFormat, startTime)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse(common.DatetimeFormat, endTime)
		if err != nil {
			return nil, err
		}
		session.And("sr.start_time >= ?", start).And("sr.end_time <= ?", end)
	}

	session.And("sr.status = ?", consts.SessionStatusClosed).GroupBy("sr.software_id")
	err := session.Select("sr.software_id, s.name, sum(sr.duration) as duration").Find(&statistics)
	if err != nil {
		return nil, err
	}
	return statistics, nil
}

type History struct {
	ID           snowflake.ID `xorm:"id"`
	SoftwareName string       `xorm:"name"`
	Duration     int64        `xorm:"duration"`
	StartTime    *time.Time   `xorm:"start_time"`
	EndTime      *time.Time   `xorm:"end_time"`
}

func (d *sessionDaoImpl) ListHistory(ctx context.Context, appIDs []string, startTime, endTime string, pageIndex, pageSize int) (int64, []*History, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	var histories []*History
	session = session.Table(model.SessionTableName).Alias("sr").
		Join("LEFT", model.SoftwareTableName+" as s", "sr.software_id = s.id")
	if len(appIDs) > 0 {
		session.In("sr.software_id", appIDs)
	}
	if startTime != "" && endTime != "" {
		start, err := time.Parse(common.DatetimeFormat, startTime)
		if err != nil {
			return 0, nil, err
		}
		end, err := time.Parse(common.DatetimeFormat, endTime)
		if err != nil {
			return 0, nil, err
		}
		session.And("sr.start_time >= ?", start).And("sr.end_time <= ?", end)
	}
	if pageIndex > 0 && pageSize > 0 {
		session.Limit(pageSize, (pageIndex-1)*pageSize)
	}

	session.And("sr.status = ?", consts.SessionStatusClosed).Desc("sr.create_time")
	total, err := session.Select("sr.id, s.name, sr.duration, sr.start_time, sr.end_time").FindAndCount(&histories)
	if err != nil {
		return 0, nil, err
	}
	return total, histories, nil
}

type StatisticItem struct {
	Key   string  `xorm:"key"`
	Value float64 `xorm:"value"`
}

func (d *sessionDaoImpl) SessionStatistics(ctx context.Context, startTime, endTime int64, reportType, dimensionType string) ([]*StatisticItem, error) {
	session := boot.MW.DefaultSession(ctx)

	statistics := make([]*StatisticItem, 0)
	session = session.Table(model.SessionTableName).Alias("sr").
		Join("LEFT", model.SoftwareTableName+" as s", "sr.software_id = s.id").
		Where("sr.status = ?", consts.SessionStatusClosed)

	if startTime > 0 {
		session.Where("sr.start_time >= ?", time.UnixMilli(startTime))
	}

	if endTime > 0 {
		session.Where("sr.end_time <= ?", time.UnixMilli(endTime))
	}

	switch {
	case consts.DimensionTypeSoftware == dimensionType && consts.ReportTypeSessionUsageDuration == reportType:
		session.GroupBy("s.name").Select("s.name `key`, sum(round(sr.duration/3600, 6)) `value`")
	case consts.DimensionTypeUser == dimensionType && consts.ReportTypeSessionUsageDuration == reportType:
		session.GroupBy("sr.user_name").Select("sr.user_name `key`, sum(round(sr.duration/3600, 6)) `value`")
	case consts.DimensionTypeSoftware == dimensionType && consts.ReportTypeSessionCreateNumber == reportType:
		session.GroupBy("s.name").Select("s.name `key`, count(*) `value`")
	case consts.DimensionTypeUser == dimensionType && consts.ReportTypeSessionCreateNumber == reportType:
		session.GroupBy("sr.user_name").Select("sr.user_name `key`, count(*) `value`")
	default:
		return nil, fmt.Errorf("unsupported report type: [%v] or dimension type: [%v]", reportType, dimensionType)
	}

	err := session.Find(&statistics)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (d *sessionDaoImpl) ListUsedProjectNames(ctx context.Context, username string) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var names []string
	session.Table(model.SessionTableName)

	if username != "" {
		session.Where("user_name = ?", username)
	}

	err := session.Distinct("project_name").Asc("project_name").Find(&names)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (d *sessionDaoImpl) GetSessionList(ctx context.Context, projectIds []int64) ([]*model.Session, error) {
	session := boot.MW.DefaultSession(ctx)

	var sessionList []*model.Session
	session = session.Table(model.SessionTableName)
	session.In("project_id", projectIds)
	err := session.Find(&sessionList)
	if err != nil {
		return nil, err
	}

	return sessionList, nil
}
