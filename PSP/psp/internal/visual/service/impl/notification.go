package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	pbproject "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

type SessionNotification struct {
	sessionDao dao.SessionDao
	rpc        *client.Client
}

func NewSessionNotification() (*SessionNotification, error) {
	rpc := client.GetInstance()

	return &SessionNotification{
		sessionDao: dao.NewSessionDao(),
		rpc:        rpc,
	}, nil
}

// SessionCheckTicker 检查
func (s *SessionNotification) SessionCheckTicker() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	currentTime := time.Now()
	dailyTime := config.GetConfig().SessionNotificationCheck.DailyTime
	if dailyTime == 0 {
		dailyTime = 9 // 默认9点
	}

	minTime := config.GetConfig().SessionNotificationCheck.MinTime

	tomorrow9AM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, dailyTime, minTime, 0, 0, currentTime.Location())
	duration := tomorrow9AM.Sub(currentTime)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := s.NotificationSessionExpire(ctx); err != nil {
				logger.Errorf("sesssion notification check err: %v", err)
			}
			currentTime := time.Now()
			tomorrow9AM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, dailyTime, minTime, 0, 0, currentTime.Location())
			duration := tomorrow9AM.Sub(currentTime)
			timer.Reset(duration)
		}
	}
}

// NotificationSessionExpire 通知会话到期信息
func (s *SessionNotification) NotificationSessionExpire(ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	start, end := timeutil.GetTomorrow()

	//1.查询明日 要结束的项目
	projects, err := s.rpc.Project.GetProjectsIdByTimePeriod(ctx, &pbproject.GetProjectsIdByTimePeriodRequest{StartTime: start, EndTime: end})
	if err != nil {
		logger.Errorf("get projects err: %v", err)
		return err
	}

	//2.查询会话
	sessions, err := s.sessionDao.GetSessionList(ctx, s.GetProjectId(projects))
	if err != nil {
		logger.Errorf("get sessions err: %v", err)
		return err
	}

	for _, session := range sessions {
		projectEndTime := s.GetProjectEndTimeById(int64(session.ProjectID), projects)
		err := s.SendNotification(ctx, session, projectEndTime)
		if err != nil {
			logger.Errorf("SendNotification err: %v", err)
			continue
		}
	}

	return nil
}

func (s *SessionNotification) GetProjectId(projects *pbproject.GetProjectsIdByTimePeriodResponse) []int64 {
	var projectIds []int64

	for _, project := range projects.Projects {
		projectIds = append(projectIds, project.ProjectId)
	}
	return projectIds
}

func (s *SessionNotification) GetProjectEndTimeById(ProjectID int64, projects *pbproject.GetProjectsIdByTimePeriodResponse) time.Time {

	for _, project := range projects.Projects {
		if project.ProjectId == ProjectID {
			return project.EndTime.AsTime()
		}
	}

	return time.Time{}
}

// SendNotification 发送到期信息
func (s *SessionNotification) SendNotification(ctx context.Context, session *model.Session, projectEndTime time.Time) error {
	logger := logging.GetLogger(ctx)

	content := fmt.Sprintf("会话[会话编号:%v 项目名:%v]，由于项目将要于%v结束，%v", session.OutSessionID, session.ProjectName, timeutil.FormatTime(projectEndTime, common.DatetimeFormat), consts.EndSessionContent)
	msg := &pbnotice.WebsocketMessage{
		UserId:  session.UserID.String(),
		Type:    common.SessionEventType,
		Content: content,
	}

	if _, err := s.rpc.Notice.SendWebsocketMessage(ctx, msg); err != nil {
		logger.Errorf("session notification send err: %v", err)
		return fmt.Errorf("session notification send err: %v", err)
	}
	return nil
}
