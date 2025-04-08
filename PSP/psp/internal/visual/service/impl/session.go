package impl

import (
	"context"
	"encoding/csv"
	"fmt"
	openapivisual "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/visual"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	projectpb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/utils"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
)

type VisualService struct {
	sid *snowflake.Node
	api *openapi.OpenAPI

	sessionDao  dao.SessionDao
	hardwareDao dao.HardwareDao
	softwareDao dao.SoftwareDao
}

func NewVisualService() (*VisualService, error) {
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
	return &VisualService{
		sid:         node,
		api:         api,
		sessionDao:  dao.NewSessionDao(),
		hardwareDao: dao.NewHardwareDao(),
		softwareDao: dao.NewSoftwareDao(),
	}, nil
}

func (s *VisualService) ListSession(ctx context.Context, hardwareIDStrs, softwareIDStrs, projectIDStrs, statuses []string, username string, isAdmin bool, pageIndex, pageSize int, loginUserID snowflake.ID) ([]*dto.Session, int64, error) {
	projectIDs := make([]snowflake.ID, 0)
	if !isAdmin && len(projectIDStrs) == 0 {
		err := s.getProjectIdsByUserID(ctx, loginUserID, &projectIDs)
		if err != nil {
			return nil, 0, err
		}
	}
	if len(projectIDStrs) != 0 {
		projectIDs = snowflake.BatchParseStringToID(projectIDStrs)
	}

	sessions, total, err := s.sessionDao.ListSessionInfos(ctx, username, isAdmin, nil, statuses, snowflake.BatchParseStringToID(hardwareIDStrs),
		snowflake.BatchParseStringToID(softwareIDStrs), projectIDs, 0, 0, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}
	remoteApps, err := s.MapSoftwareIDToRemoteApps(ctx)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*dto.Session, 0)
	for _, v := range sessions {
		session := utils.ConvertSession(v)
		session.Software.RemoteApps = remoteApps[v.SoftwareID]
		if session.Software.RemoteApps == nil {
			session.Software.RemoteApps = make([]*dto.RemoteApp, 0)
		}
		list = append(list, session)
	}
	return list, total, nil
}

func (s *VisualService) StartSession(ctx context.Context, projectIDStr, hardwareIDStr, softwareIDStr, username string, mounts []string, loginUserID snowflake.ID) (string, string, string, error) {
	mountDirectory := config.GetConfig().MountDirectory
	projectId, projectName, projectFilePath := common.PersonalProjectID, common.PersonalProjectName, ""
	if projectIDStr != projectId.String() {
		project, err := client.GetInstance().Project.GetProjectDetailById(ctx, &projectpb.GetProjectDetailByIdRequest{ProjectId: projectIDStr})
		if err != nil {
			return "", "", "", err
		}

		existResponse, err := client.GetInstance().Project.ExistsProjectMember(ctx, &projectpb.ExistsProjectMemberRequest{UserId: loginUserID.String(), ProjectId: project.ProjectId})
		if err != nil {
			return "", "", "", err
		}
		if !existResponse.IsExist {
			return "", "", "", status.Errorf(errcode.ErrVisualCurrentProjectNotAccess, fmt.Sprintf("the project: [%v] not found the user: [%v]", project.ProjectId, loginUserID.String()))
		}

		if project.State != common.ProjectRunning {
			return "", "", "", status.Errorf(errcode.ErrVisualCurrentProjectNotRunning, fmt.Sprintf("the project: [%v] not running", project.ProjectId))
		}

		projectId = snowflake.MustParseString(project.ProjectId)
		projectName = project.ProjectName
		projectFilePath = project.FilePath
	}

	hardwareID := snowflake.MustParseString(hardwareIDStr)
	hardware, exist, err := s.hardwareDao.GetHardware(ctx, hardwareID, "")
	if err != nil {
		return "", "", "", err
	}
	if !exist {
		return "", "", "", status.Errorf(errcode.ErrVisualHardwareNotFound, "hardware not found, hardwareID: [%v]", hardwareIDStr)
	}
	softwareID := snowflake.MustParseString(softwareIDStr)
	software, exist, err := s.softwareDao.GetSoftware(ctx, softwareID, "")
	if err != nil {
		return "", "", "", err
	}
	if !exist {
		return "", "", "", status.Errorf(errcode.ErrVisualSoftwareNotFound, "software not found, softwareID: [%v]", softwareIDStr)
	}

	// 限制当前用户对一种软件只能有一个session
	sessions, total, err := s.sessionDao.HasUsedResource(ctx, projectId, username, hardwareID, softwareID, NormalSessionStatus)
	if err != nil {
		return "", "", "", err
	}
	if total > 0 {
		return sessions[0].OutSessionID, projectName, sessions[0].OutSessionID, status.Errorf(errcode.ErrVisualSessionRepeatStart, "project: [%v] session has used, hardwareID: [%v], softwareID: [%v], username: [%v]", projectIDStr, hardwareIDStr, softwareIDStr, username)
	}

	if len(mountDirectory.DriveNames) < mountDirectory.LimitNum {
		return "", "", "", fmt.Errorf("drive names len: [%v] less than mount limit num: [%v], please check settings", mountDirectory.DriveNames, mountDirectory.LimitNum)
	}

	mountDriveList := mountDirectory.DriveNames
	mountDrives := make(map[string]string, len(mountDriveList))
	mountDrives[username] = fmt.Sprintf("%v%v", mountDriveList[0], common.Colon)
	if projectFilePath != "" {
		mountDrives[projectFilePath[1:]] = fmt.Sprintf("%v%v", mountDriveList[len(mountDrives)], common.Colon)
	}
	if mountDirectory.EnablePublicDirectory {
		mountDrives[common.PublicFolderPath] = fmt.Sprintf("%v%v", mountDriveList[len(mountDrives)], common.Colon)
	}
	if len(mountDrives)+len(mounts) > mountDirectory.LimitNum {
		return "", "", "", fmt.Errorf("over mount num limit: [%v],drive names len: [%v], mounts len: [%v]", mountDirectory.LimitNum, len(mountDrives), len(mounts))
	}

	if len(mounts) > 0 {
		projectsResponse, err := client.GetInstance().Project.GetProjectsDetailByIds(ctx, &projectpb.GetProjectsDetailByIdsRequest{ProjectIds: mounts})
		if err != nil {
			return "", "", "", err
		}
		mountDriveLength := len(mountDrives)
		for i, v := range projectsResponse.Projects {
			mountDrives[v.FilePath[1:]] = fmt.Sprintf("%v%v", mountDriveList[mountDriveLength+i], common.Colon)
		}
	}

	if software.Platform == consts.PlatformTypeLinux {
		for k := range mountDrives {
			mountDrives[k] = fmt.Sprintf("%v/%v", mountDirectory.LinuxMountRootPath, k)
		}
	}

	response, err := openapivisual.StartSession(s.api, hardware.OutHardwareID, software.OutSoftwareID, mountDrives)
	if err != nil {
		return "", "", "", err
	}
	if response == nil || response.ErrorCode != "" {
		return "", "", "", fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	newID := s.sid.Generate()
	err = s.sessionDao.InsertSessions(ctx, []*model.Session{
		{
			ID:           newID,
			OutSessionID: response.Data.Id,
			ProjectID:    projectId,
			ProjectName:  projectName,
			UserID:       loginUserID,
			UserName:     username,
			RawStatus:    response.Data.Status,
			Status:       response.Data.Status,
			StreamURL:    response.Data.StreamUrl,
			HardwareID:   hardwareID,
			SoftwareID:   softwareID,
			Zone:         config.GetZone(),
			StartTime:    time.Now(),
		},
	})
	if err != nil {
		return "", "", "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("start session: [%v], outSessionId: [%v], projectId: [%v], projectName: [%v]", newID, response.Data.Id, projectId, projectName))

	return newID.String(), "", response.Data.Id, nil
}

func (s *VisualService) GetMountInfo(ctx context.Context, projectIDStr, username string, loginUserID snowflake.ID) (*dto.GetMountInfoResponse, error) {
	mountDirectory := config.GetConfig().MountDirectory
	defaultMounts := make([]*dto.MountInfo, 0, 2)
	defaultMounts = append(defaultMounts, &dto.MountInfo{ID: common.PersonalProjectID.String(), Name: common.PersonalProjectName})

	if snowflake.MustParseString(projectIDStr) != common.PersonalProjectID {
		projectResponse, err := client.GetInstance().Project.GetProjectDetailById(ctx, &projectpb.GetProjectDetailByIdRequest{ProjectId: projectIDStr})
		if err != nil {
			return nil, err
		}
		if projectResponse == nil || projectResponse.ProjectName == "" {
			return nil, fmt.Errorf("project name is empty")
		}

		defaultMounts = append(defaultMounts, &dto.MountInfo{ID: projectResponse.ProjectId, Name: projectResponse.ProjectName})
	}

	if mountDirectory.EnablePublicDirectory {
		defaultMounts = append(defaultMounts, &dto.MountInfo{Name: common.PublicFolderPath})
	}

	projectsResponse, err := client.GetInstance().Project.GetMemberProjectsByUserId(ctx, &projectpb.GetMemberProjectsByUserIdRequest{UserId: loginUserID.String(), IncludeDefault: false})
	if err != nil {
		return nil, err
	}

	projectNames := make([]*dto.MountInfo, 0, len(projectsResponse.Projects))
	for _, v := range projectsResponse.Projects {
		if v.State == common.ProjectRunning && v.ProjectId != projectIDStr {
			projectNames = append(projectNames, &dto.MountInfo{ID: v.ProjectId, Name: v.ProjectName})
		}
	}

	return &dto.GetMountInfoResponse{
		DefaultMounts: defaultMounts,
		SelectMounts:  projectNames,
		SelectLimit:   mountDirectory.LimitNum - len(defaultMounts),
	}, nil
}

func (s *VisualService) CloseSession(ctx context.Context, sessionIDStr, exitReason string, admin bool) (string, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}

	var response v20230530.Response
	if admin {
		result, errTmp := openapivisual.AdminCloseSession(s.api, session.OutSessionID, exitReason)
		response, err = result.Response, errTmp
	} else {
		result, errTmp := openapivisual.UserCloseSession(s.api, session.OutSessionID)
		response, err = result.Response, errTmp
	}
	if err != nil {
		return "", err
	}
	if response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response)
	}

	err = s.sessionDao.UpdateSession(ctx, &model.Session{
		ID:        sessionID,
		RawStatus: consts.SessionStatusClosing,
		Status:    consts.SessionStatusClosing,
	})
	if err != nil {
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("close session: [%v], outSessionId: [%v], projectId: [%v], projectName: [%v]", session.ID, session.OutSessionID, session.ProjectID, session.ProjectName))

	return session.OutSessionID, nil
}

func (s *VisualService) PowerOffSession(ctx context.Context, sessionIDStr string) (string, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}

	response, err := openapivisual.AdminPowerOffSession(s.api, session.OutSessionID)
	if err != nil {
		return "", err
	}
	if response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response)
	}

	err = s.sessionDao.UpdateSession(ctx, &model.Session{
		ID:        sessionID,
		RawStatus: consts.SessionStatusPoweringOff,
		Status:    consts.SessionStatusPoweringOff,
	})
	if err != nil {
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("power off session: [%v], outSessionId: [%v], projectId: [%v], projectName: [%v]", session.ID, session.OutSessionID, session.ProjectID, session.ProjectName))

	return session.OutSessionID, nil
}

func (s *VisualService) PowerOnSession(ctx context.Context, sessionIDStr string) (string, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}

	response, err := openapivisual.AdminPowerOnSession(s.api, session.OutSessionID)
	if err != nil {
		return "", err
	}
	if response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response)
	}

	err = s.sessionDao.UpdateSession(ctx, &model.Session{
		ID:        sessionID,
		RawStatus: consts.SessionStatusPoweringOn,
		Status:    consts.SessionStatusPoweringOn,
	})
	if err != nil {
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("power on session: [%v], outSessionId: [%v], projectId: [%v], projectName: [%v]", session.ID, session.OutSessionID, session.ProjectID, session.ProjectName))

	return session.OutSessionID, nil
}

func (s *VisualService) RebootSession(ctx context.Context, sessionIDStr, Reason string, admin bool) (bool, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}

	var response v20230530.Response
	result, errTmp := openapivisual.AdminRebootSession(s.api, session.OutSessionID)
	response, err = result.Response, errTmp
	if err != nil {
		return false, err
	}
	if response.ErrorCode != "" {
		return false, fmt.Errorf("openapi response nil or response err: [%+v]", response)
	}

	err = s.sessionDao.UpdateSession(ctx, &model.Session{
		ID:        sessionID,
		RawStatus: consts.SessionStatusRebooting,
		Status:    consts.SessionStatusRebooting,
	})
	if err != nil {
		return false, err
	}

	tracelog.Info(ctx, fmt.Sprintf("reboot session: [%v], outSessionId: [%v], projectId: [%v], projectName: [%v]", session.ID, session.OutSessionID, session.ProjectID, session.ProjectName))

	return true, nil
}

func (s *VisualService) ReadySession(ctx context.Context, sessionIDStr string) (bool, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}

	response, err := openapivisual.ReadySession(s.api, session.OutSessionID)
	if err != nil {
		return false, err
	}
	if response == nil || response.ErrorCode != "" {
		return false, fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	ready := false
	if response.Data != nil {
		ready = response.Data.Ready
	}
	return ready, nil
}

func (s *VisualService) GetRemoteAppURL(ctx context.Context, sessionIDStr, remoteAppName string) (string, error) {
	sessionID := snowflake.MustParseString(sessionIDStr)
	session, exist, err := s.sessionDao.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", status.Errorf(errcode.ErrVisualSessionNotFound, "session not found, sessionID: [%v]", sessionIDStr)
	}
	if session.Status != consts.SessionStatusStarted {
		return "", status.Errorf(errcode.ErrVisualSessionNotStarted, "the session not started, sessionID: [%v]", sessionIDStr)
	}

	response, err := openapivisual.GetRemoteAppURL(s.api, session.OutSessionID)
	if err != nil {
		return "", err
	}
	if response == nil || response.ErrorCode != "" {
		return "", fmt.Errorf("openapi response nil or response err: [%+v]", response.Response)
	}

	streamURL := ""
	if response.Data != nil {
		streamURL = response.Data.StreamUrl
	}
	return streamURL, nil
}

func (s *VisualService) ListUsedProjectNames(ctx context.Context, username string) ([]string, error) {
	names, err := s.sessionDao.ListUsedProjectNames(ctx, username)
	if err != nil {
		return make([]string, 0), err
	}

	return names, nil
}

func (s *VisualService) getProjectIdsByUserID(ctx context.Context, loginUserID snowflake.ID, projectIDs *[]snowflake.ID) error {
	response, err := client.GetInstance().Project.GetMemberProjectsByUserId(ctx, &projectpb.GetMemberProjectsByUserIdRequest{UserId: loginUserID.String(), IncludeDefault: true})
	if err != nil {
		return err
	}

	if projectIDs == nil {
		*projectIDs = make([]snowflake.ID, 0)
	}

	for _, v := range response.Projects {
		id := snowflake.MustParseString(v.ProjectId)
		*projectIDs = append(*projectIDs, id)
	}

	return nil
}

func (s *VisualService) ExportSessionInfo(ctx *gin.Context, startTime, endTime int64) error {
	logger := logging.GetLogger(ctx)

	start := time.Now()
	csvHeaders := []string{"会话编号", "镜像名称", "项目名称", "操作平台", "实例名称", "实例类型", "会话状态", "创建者", "创建时间", "开始时间", "结束时间"}
	err := csvutil.ExportCSVFile(ctx, &csvutil.ExportCSVFileInfo{
		CSVFileName: "3D可视化-会话详情",
		CSVHeaders:  csvHeaders,
		FillCSVData: func(w *csv.Writer) error {
			pageIndex, pageSize := common.DefaultPageIndex, common.CSVExportNumber

			for {
				exportData, _, err := s.sessionDao.ListSessionInfos(ctx, "", true, nil, nil, nil, nil, nil, startTime, endTime, pageIndex, pageSize)
				if err != nil {
					return err
				}

				if len(exportData) == 0 {
					break
				}

				for _, v := range exportData {
					rowData := make([]string, 0, len(csvHeaders))

					if v.Session == nil || v.Software == nil || v.Hardware == nil {
						logger.Errorf("session or software or hardware is nil, value: [%+v]", v)
						continue
					}

					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.ID.String()))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Software.Name))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.ProjectName))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Software.Platform))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Hardware.Name))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Hardware.InstanceType))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.Status))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.UserName))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.CreateTime.Format(common.DatetimeFormat)))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.StartTime.Format(common.DatetimeFormat)))
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Session.EndTime.Format(common.DatetimeFormat)))

					_ = w.Write(rowData)
				}
				w.Flush()

				if len(exportData) < common.CSVExportNumber {
					break
				}

				pageIndex++
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	logger.Infof("export session info cost time: %v", time.Since(start))

	return nil
}

func (s *VisualService) SessionInfo(ctx context.Context, sessionID string) (*dto.Session, error) {
	session, _, err := s.sessionDao.GetSession(ctx, snowflake.MustParseString(sessionID))
	return utils.Convert2Session(session), err
}
