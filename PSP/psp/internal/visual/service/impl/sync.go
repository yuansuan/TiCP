package impl

import (
	"context"
	"fmt"
	openapivisual "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/visual"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/utils"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type SessionScanner struct {
	sid *snowflake.Node
	api *openapi.OpenAPI

	sessionDao dao.SessionDao
}

func NewSessionScanner() (*SessionScanner, error) {
	logger := logging.Default()
	node, err := snowflake.GetInstance()
	if err != nil {
		logger.Errorf("new snowflake node err: %v", err)
		return nil, err
	}
	api, err := openapi.NewLocalAPI()
	if err != nil {
		logger.Errorf("new openapi err: %v", err)
		return nil, err
	}
	return &SessionScanner{
		sid:        node,
		api:        api,
		sessionDao: dao.NewSessionDao(),
	}, nil
}

var NormalSessionStatus = []string{
	consts.SessionStatusPending,
	consts.SessionStatusStarting,
	consts.SessionStatusStarted,
	consts.SessionStatusUnAvailable,
	consts.SessionStatusRebooting,
	consts.SessionStatusPoweringOff,
	consts.SessionStatusPowerOff,
	consts.SessionStatusPoweringOn,
}

var ActiveSessionStatus = []string{
	consts.SessionStatusPending,
	consts.SessionStatusStarting,
	consts.SessionStatusStarted,
	consts.SessionStatusUnAvailable,
	consts.SessionStatusClosing,
	consts.SessionStatusRebooting,
	consts.SessionStatusPoweringOff,
	consts.SessionStatusPowerOff,
	consts.SessionStatusPoweringOn,
}

var ReadySessionStatus = []string{
	consts.SessionStatusStarting,
	consts.SessionStatusStarted,
	consts.SessionStatusUnAvailable,
	consts.SessionStatusRebooting,
	consts.SessionStatusPoweringOff,
	consts.SessionStatusPowerOff,
	consts.SessionStatusPoweringOn,
}

// SyncSessionDataRoutine 同步Paas层和PSP层间会话数据信息
func (s *SessionScanner) SyncSessionDataRoutine() {
	logger := logging.Default()

	interval := config.GetConfig().SyncData.DataInterval
	logger.Infof("sync session data routine interval: [%v]", interval)
	if interval <= 0 {
		interval = consts.DefaultSyncDataInterval
	}

	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(interval) * time.Second):
			pageIndex, pageSize := common.DefaultPageIndex, consts.DefaultSyncNumber
			for {
				sessions, _, err := s.sessionDao.ListSession(ctx, "", ActiveSessionStatus, nil, nil, nil, 0, 0, pageIndex, pageSize)
				if err != nil {
					logger.Errorf("scan sessions err: %v", err)
					break
				}
				if len(sessions) == 0 {
					break
				}

				// 更新状态变化的会话数据
				s.updateSessionData(ctx, sessions)

				if len(sessions) < consts.DefaultSyncNumber {
					break
				}

				pageIndex++
			}
		}
	}
}

func (s *SessionScanner) updateSessionData(ctx context.Context, sessions []*model.Session) {
	logger := logging.GetLogger(ctx)

	ids := make([]string, 0)
	idToSessionMap := make(map[string]*model.Session)
	for _, v := range sessions {
		ids = append(ids, v.OutSessionID)
		idToSessionMap[v.OutSessionID] = v
	}

	idsStr := strings.Join(ids, ",")
	response, err := openapivisual.ListSession(s.api, idsStr, common.DefaultPageOffset, consts.DefaultSyncNumber)
	if err != nil {
		logger.Errorf("openapi get session err: %v, idsStr: [%v]", err, idsStr)
		return
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v], idsStr: [%v]", response, idsStr)
		return
	}
	if response.Data == nil {
		logger.Errorf("openapi response data nil, idsStr: [%v]", idsStr)
		return
	}

	for _, v := range response.Data.Sessions {
		session, ok := idToSessionMap[v.Id]
		if !ok {
			continue
		}
		if session.RawStatus == v.Status {
			continue
		}

		sessionData := utils.ConvertDBSession(v, session.ID)

		if sessionData.RawStatus != consts.SessionStatusStarted {
			sessionData.Status = sessionData.RawStatus
		}

		if sessionData.RawStatus == consts.SessionStatusClosed {
			sessionData.Duration = sessionData.EndTime.Unix() - sessionData.StartTime.Unix()
		}

		err = s.sessionDao.UpdateSession(ctx, sessionData)
		if err != nil {
			logger.Errorf("update session err: %v, session: [%+v]", err, sessionData)
			continue
		}
	}
}

// SyncSessionStatusRoutine 同步更新 PSP 会话状态
func (s *SessionScanner) SyncSessionStatusRoutine() {
	logger := logging.Default()

	interval := config.GetConfig().SyncData.StatusInterval
	logger.Infof("sync session status routine interval: [%v]", interval)
	if interval <= 0 {
		interval = consts.DefaultSyncStatusInterval
	}

	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(interval) * time.Second):
			userProjectMap := make(map[string]struct{})
			pageIndex, pageSize := common.DefaultPageIndex, common.DefaultMaxPageSize

			for {
				sessions, _, err := s.sessionDao.ListSession(ctx, "", nil, ReadySessionStatus, nil, nil, 0, 0, pageIndex, pageSize)
				if err != nil {
					logger.Errorf("scan sessions err: %v", err)
					break
				}
				if len(sessions) == 0 {
					break
				}

				go s.resolveSessionStatus(ctx, sessions, userProjectMap)

				if len(sessions) < common.DefaultMaxPageSize {
					break
				}

				pageIndex++
			}
		}
	}
}

func (s *SessionScanner) resolveSessionStatus(ctx context.Context, sessions []*model.Session, userProjectMap map[string]struct{}) {
	for _, v := range sessions {
		// 当远程会话 Ready 后, 实时更改状态
		if v.RawStatus == consts.SessionStatusStarted {
			go s.updateSessionStatus(ctx, v)
		}

		// 当离开项目时, 会话实时关闭; 不影响个人项目
		if common.PersonalProjectID != v.ProjectID && common.PersonalProjectName != v.ProjectName {
			userProjectKey := s.checkUserProjectMap(ctx, v, userProjectMap)
			if _, ok := userProjectMap[userProjectKey]; userProjectKey != "" && !ok {
				go s.autoCloseActiveSession(ctx, v)
			}
		}
	}
}

func (s *SessionScanner) updateSessionStatus(ctx context.Context, v *model.Session) {
	logger := logging.GetLogger(ctx)

	response, err := openapivisual.ReadySession(s.api, v.OutSessionID)
	if err != nil {
		logger.Errorf("openapi get session ready status err: %v, outSessionID: [%v]", err, v.OutSessionID)
		return
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v], outSessionID: [%v]", response, v.OutSessionID)
		return
	}
	if response.Data == nil {
		logger.Errorf("openapi response data nil, outSessionID: [%v]", v.OutSessionID)
		return
	}

	ready := response.Data.Ready

	switch {
	case ready && v.Status == consts.SessionStatusStarting:
		v.Status = consts.SessionStatusStarted
	case !ready && v.Status == consts.SessionStatusStarting:
		return
	case ready && v.Status == consts.SessionStatusStarted:
		return
	case !ready && v.Status == consts.SessionStatusStarted:
		v.Status = consts.SessionStatusUnAvailable
	case ready && v.Status == consts.SessionStatusUnAvailable:
		v.Status = consts.SessionStatusStarted
	case !ready && v.Status == consts.SessionStatusUnAvailable:
		return
	case ready && v.Status == consts.SessionStatusRebooting:
		v.Status = consts.SessionStatusStarted
	case !ready && v.Status == consts.SessionStatusRebooting:
		return
	case ready && v.Status == consts.SessionStatusPoweringOn:
		v.Status = consts.SessionStatusStarted
	case !ready && v.Status == consts.SessionStatusPoweringOn:
		return
	default:
		logger.Errorf("visual session ready status: [%v] & status: [%v] not match", ready, v.Status)
		return
	}
	err = s.sessionDao.UpdateSession(ctx, v)
	if err != nil {
		logger.Errorf("update session err: %v, session: [%+v]", err, v)
	}
}

func (*SessionScanner) checkUserProjectMap(ctx context.Context, v *model.Session, userProjectMap map[string]struct{}) string {
	logger := logging.Default()

	// 标记选择运行中的项目
	userProjectKey := fmt.Sprintf("%v_%v_%v", v.UserID, v.ProjectID, common.ProjectRunning)

	if _, ok := userProjectMap[v.UserID.String()]; !ok {
		response, err := client.GetInstance().Project.GetMemberProjectsByUserId(ctx, &project.GetMemberProjectsByUserIdRequest{UserId: v.UserID.String()})
		if err != nil {
			logger.Errorf("get projects by user id err: %v", err)
			return ""
		}

		for _, p := range response.Projects {
			userProjectInputKey := fmt.Sprintf("%v_%v_%v", v.UserID, p.ProjectId, p.State)
			userProjectMap[userProjectInputKey] = struct{}{}
		}

		// 设置 map 存储标记
		userProjectMap[v.UserID.String()] = struct{}{}
	}

	return userProjectKey
}

func (s *SessionScanner) autoCloseActiveSession(ctx context.Context, session *model.Session) {
	logger := logging.GetLogger(ctx)

	response, err := openapivisual.AdminCloseSession(s.api, session.OutSessionID, consts.DefaultUserExistProjecReason)
	if err != nil {
		logger.Errorf("openapi auto close session err: %v, outSessionID: [%v]", err, session.OutSessionID)
		return
	}
	if response == nil || response.ErrorCode != "" {
		logger.Errorf("openapi response nil or response err: [%+v], outSessionID: [%v]", response, session.OutSessionID)
		return
	}

	session.RawStatus = consts.SessionStatusClosing
	session.Status = consts.SessionStatusClosing

	err = s.sessionDao.UpdateSession(ctx, session)
	if err != nil {
		logger.Errorf("update session err: %v, session: [%+v]", err, session)
	}
}
