package impl

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
)

type projectServiceImpl struct {
	projectDao           dao.ProjectDao
	projectMemberDao     dao.ProjectMemberDao
	sid                  *snowflake.Node
	projectMemberService service.ProjectMemberService
}

func NewProjectService() (*projectServiceImpl, error) {
	projectDao, err := dao.NewProjectDao()
	if err != nil {
		return nil, err
	}

	projectMemberDao := dao.NewProjectMemberDao()

	sid, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	projectMemberService, err := NewProjectMemberService()
	if err != nil {
		return nil, err
	}

	return &projectServiceImpl{
		projectDao:           projectDao,
		projectMemberDao:     projectMemberDao,
		sid:                  sid,
		projectMemberService: projectMemberService,
	}, nil
}

func (s *projectServiceImpl) ProjectSave(ctx context.Context, req *dto.ProjectAddRequest, loginUserID snowflake.ID, loginUserName string) (*dto.ProjectAddResponse, error) {
	// 校验当前userID权限, 需要项目管理员权限角色
	isPMR, err := util.CheckProjectAdminRole(ctx, loginUserID, true)
	if err != nil {
		return nil, err
	}

	if !isPMR {
		return nil, status.Error(errcode.ErrProjectAccessPermission, errcode.ProjectCodeMsg[errcode.ErrProjectAccessPermission])
	}

	// 项目名称不能为个人默认项目名称 default
	if req.ProjectName == common.PersonalProjectName {
		return nil, status.Errorf(errcode.ErrProjectNameIsDefault, errcode.ProjectCodeMsg[errcode.ErrProjectNameIsDefault])
	}

	var projectID string
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		// 判断项目名称是否重复
		exist, err := s.projectDao.ExistSameProjectName(ctx, req.ProjectName)
		if err != nil {
			return err
		}
		if exist {
			return status.Error(errcode.ErrProjectSameName, errcode.ProjectCodeMsg[errcode.ErrProjectSameName])
		}

		projectFilePath := path.Join("/", common.ProjectFolderPath, req.ProjectName)
		project := &model.Project{
			ProjectName:  req.ProjectName,
			ProjectOwner: snowflake.MustParseString(req.ProjectOwner),
			State:        common.ProjectInit,
			StartTime:    time.Unix(req.StartTime, 0),
			EndTime:      time.Unix(req.EndTime, 0),
			Comment:      req.Comment,
			FilePath:     projectFilePath,
		}
		id, err := s.projectDao.InsertProject(ctx, project)
		if err != nil {
			return err
		}
		projectID = id

		members := make([]*model.ProjectMember, 0)
		fileName := req.ProjectName + "_" + timeutil.FormatTime(time.Now(), common.YearMonthDayFormat)
		for _, memberID := range req.Members {
			userInfo, err := client.GetInstance().User.GetIncludeDeleted(ctx, &user.UserIdentity{Id: memberID})
			if err != nil {
				return err
			}

			if userInfo == nil {
				return fmt.Errorf("userID:[%v] not found err", memberID)
			}

			projectMember := &model.ProjectMember{
				Id:        s.sid.Generate(),
				ProjectId: snowflake.MustParseString(id),
				UserId:    snowflake.MustParseString(memberID),
				LinkPath:  path.Join("/", userInfo.Name, common.WorkspaceFolderPath, fileName),
			}
			members = append(members, projectMember)
		}

		err = s.projectMemberDao.BatchInsertProjectMember(ctx, members)
		if err != nil {
			return err
		}

		// 创建项目文件目录
		createDirReq := &storage.CreateDirReq{
			Path:     projectFilePath,
			UserName: loginUserName,
			Cross:    true,
			IsCloud:  false,
		}

		_, err = client.GetInstance().Storage.CreateDir(ctx, createDirReq)
		if err != nil {
			return err
		}

		// 创建对应的个人目录软链接
		err = s.projectMemberService.CreatePersonalProjectPath(ctx, members, project)
		if err != nil {
			// 创建个人目录失败，删除已创建项目文件目录和员工文件目录
			s.delProjectAndMemberFile(ctx, projectFilePath, members)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ProjectAddResponse{
		ID: projectID,
	}, nil
}

func (s *projectServiceImpl) delProjectAndMemberFile(ctx context.Context, projectFilePath string, members []*model.ProjectMember) {
	logger := logging.GetLogger(ctx)

	for _, member := range members {
		membersDirRmReq := &storage.RmReq{
			Paths:   []string{member.LinkPath},
			Cross:   true,
			IsCloud: false,
		}

		resp, err := client.GetInstance().Storage.Rm(ctx, membersDirRmReq)
		if err != nil {
			logger.Debugf("delete member dir resp:[%v], err:[%v]", resp, err)
		}
	}

	projectDirRmReq := &storage.RmReq{
		Paths:   []string{projectFilePath},
		Cross:   true,
		IsCloud: false,
	}

	resp, err := client.GetInstance().Storage.Rm(ctx, projectDirRmReq)
	if err != nil {
		logger.Errorf("delete project dir resp: [%v], err: [%v]", resp, err)
		return
	}

}

func (s *projectServiceImpl) ProjectList(ctx context.Context, req *dto.ProjectListRequest, loginUserID snowflake.ID) (*dto.ProjectListResponse, error) {
	isPMR := false
	if req.IsSysMenu {
		// 系统管理员权限角色成员 可以查看所有项目信息
		checkPMR, err := util.CheckProjectAdminRole(ctx, loginUserID, true)
		if err != nil {
			return nil, err
		}

		if checkPMR {
			isPMR = true
		}
	}

	list, total, err := s.projectDao.GetProjectList(ctx, req, loginUserID, isPMR)
	if err != nil {
		return nil, err
	}

	projectList := util.Convert2ProjectList(list, loginUserID)
	// 解析管理员名称
	for _, project := range projectList {
		userInfo, err := client.GetInstance().User.GetIncludeDeleted(ctx, &user.UserIdentity{Id: project.ProjectOwnerID})
		if err != nil {
			return nil, errors.Wrap(err, "get user info err")
		}

		project.ProjectOwnerName = userInfo.Name

		// 项目成员数据组装
		projectID := snowflake.MustParseString(project.ID)
		members, err := s.getMemberDetail(ctx, projectID)
		if err != nil {
			return nil, err
		}
		project.Members = members
	}

	return &dto.ProjectListResponse{
		ProjectList: projectList,
		Total:       total,
	}, nil
}

func (s *projectServiceImpl) CurrentProjectList(ctx context.Context, req *dto.CurrentProjectListRequest, loginUserID snowflake.ID) (*dto.CurrentProjectListResponse, error) {
	logger := logging.GetLogger(ctx)

	states := make([]string, 0)
	if req.State != "" {
		states = append(states, req.State)
	}

	projectList, err := s.projectMemberDao.GetProjectsByUserId(ctx, states, loginUserID)
	if err != nil {
		logger.Errorf("get current project list error, err: %v", err)
		return nil, err
	}

	projects := make([]*dto.CurrentProjectInfo, 0)
	if req.State == common.ProjectRunning || req.State == "" {
		projects = append(projects, &dto.CurrentProjectInfo{
			Id:   common.PersonalProjectID.String(),
			Name: common.PersonalProjectName,
		})
	}

	for _, v := range projectList {
		project := &dto.CurrentProjectInfo{
			Id:   v.ProjectID.String(),
			Name: v.ProjectName,
		}
		projects = append(projects, project)
	}

	return &dto.CurrentProjectListResponse{Projects: projects}, nil
}

func (s *projectServiceImpl) CurrentProjectListForParam(ctx context.Context, req *dto.CurrentProjectListForParamRequest, loginUserID snowflake.ID) (*dto.CurrentProjectListForParamResponse, error) {
	var startTime, endTime time.Time
	listLimit := config.GetConfig().SelectorListLimit
	if listLimit.Enable {
		endTime = time.Now()
		startTime = endTime.AddDate(0, -listLimit.MaxMonths, 0)
	}

	projectList, _, err := s.projectDao.CurrentProjectListForParam(ctx, req, loginUserID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	projects := make([]*dto.CurrentProjectInfo, 0)
	projects = append(projects, &dto.CurrentProjectInfo{
		Id:   common.PersonalProjectID.String(),
		Name: common.PersonalProjectName,
	})

	for _, v := range projectList {
		projects = append(projects, &dto.CurrentProjectInfo{
			Id:   v.Id.String(),
			Name: v.ProjectName,
		})
	}

	return &dto.CurrentProjectListForParamResponse{Projects: projects}, nil
}

func (s *projectServiceImpl) ProjectDetail(ctx context.Context, req *dto.ProjectDetailRequest, loginUserID snowflake.ID) (*dto.ProjectDetailResponse, error) {
	projectID := snowflake.MustParseString(req.ProjectID)

	// 校验当前userID权限, 当前是系统管理员角色或者项目组成员，可以查看所有的项目
	err := s.hasProjectAccessPermission(ctx, projectID, loginUserID)
	if err != nil {
		return nil, err
	}

	projectInfo, exist, err := s.projectDao.GetProjectDetailById(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	members, err := s.getMemberDetail(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 获取owner 名字
	project := util.Convert2ProjectInfo(projectInfo, loginUserID)
	owner, err := client.GetInstance().User.GetIncludeDeleted(ctx, &user.UserIdentity{Id: project.ProjectOwnerID})
	if err != nil {
		return nil, err
	}
	project.ProjectOwnerName = owner.Name

	resp := &dto.ProjectDetailResponse{}
	resp.ProjectListInfo = project
	resp.Members = members
	return resp, nil
}

// getMemberDetail 解析项目成员
func (s *projectServiceImpl) getMemberDetail(ctx context.Context, projectID snowflake.ID) ([]*dto.ProjectMembersInfo, error) {
	projectMembers, err := s.projectMemberDao.GetProjectMembersByProjectId(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 项目成员数据组装
	userIdentities := lo.Map[*model.ProjectMember, *user.UserIdentity](projectMembers, func(m *model.ProjectMember, _ int) *user.UserIdentity {
		return &user.UserIdentity{
			Id: m.UserId.String(),
		}
	})

	users, err := client.GetInstance().User.BatchGetUser(ctx, &user.UserIdentities{UserIdentities: userIdentities})
	if err != nil {
		return nil, err
	}

	members := lo.Map[*user.UserObj, *dto.ProjectMembersInfo](users.UserObj, func(u *user.UserObj, _ int) *dto.ProjectMembersInfo {
		return &dto.ProjectMembersInfo{
			UserID:   u.Id,
			UserName: u.Name,
		}
	})

	return members, nil
}

func (s *projectServiceImpl) ProjectDelete(ctx *gin.Context, projectID string, loginUserID snowflake.ID) error {
	id := snowflake.MustParseString(projectID)

	// 权限校验 是否项目管理员或者系统管理员
	err := s.hasProjectEditPermission(ctx, id, loginUserID)
	if err != nil {
		return err
	}

	projectInfo, exist, err := s.projectDao.GetProjectDetailById(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		return status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	projectMembers, err := s.projectMemberDao.GetProjectMembersByProjectId(ctx, id)
	if err != nil {
		return err
	}

	if len(projectMembers) > 1 || projectMembers[0].UserId != projectInfo.ProjectOwner {
		return status.Error(errcode.ErrProjectDelBeforeExistMembers, errcode.ProjectCodeMsg[errcode.ErrProjectDelBeforeExistMembers])
	}

	if projectInfo.State == common.ProjectRunning {
		return status.Error(errcode.ErrProjectDelete, errcode.ProjectCodeMsg[errcode.ErrProjectDelete])
	}

	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		// 删除人员
		err = s.projectMemberDao.BatchDeleteProjectMember(ctx, projectMembers[0].ProjectId, []snowflake.ID{projectMembers[0].UserId})
		if err != nil {
			return err
		}

		// 删除软链
		s.projectMemberService.DelPersonalProjectPath(ctx, id, []snowflake.ID{projectMembers[0].UserId})

		// 删除项目
		projectInfo.IsDelete = 1
		err = s.projectDao.UpdateProjectWithCols(ctx, projectInfo, []string{"is_delete"})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_PROJECT_MANAGER, fmt.Sprintf("用户%v删除项目[%v]", ginutil.GetUserName(ctx), projectInfo.ProjectName))
	return nil
}

func (s *projectServiceImpl) ProjectTerminate(ctx context.Context, projectID string, loginUserID snowflake.ID) error {
	id := snowflake.MustParseString(projectID)

	// 权限校验 是否项目管理员或者系统管理员
	err := s.hasProjectEditPermission(ctx, id, loginUserID)
	if err != nil {
		return err
	}

	projectInfo, exist, err := s.projectDao.GetProjectDetailById(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		return status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	if projectInfo.State != common.ProjectRunning &&
		projectInfo.State != common.ProjectInit {
		return status.Error(errcode.ErrProjectTerminateState, errcode.ProjectCodeMsg[errcode.ErrProjectTerminateState])
	}

	projectInfo.State = common.ProjectTerminated
	err = s.projectDao.UpdateProjectWithCols(ctx, projectInfo, []string{"state"})
	return err
}

func (s *projectServiceImpl) ProjectEdit(ctx context.Context, req *dto.ProjectEditRequest, loginUserID snowflake.ID) error {
	projectID := snowflake.MustParseString(req.ProjectID)

	// 权限校验 是否项目管理员或者系统管理员
	err := s.hasProjectEditPermission(ctx, projectID, loginUserID)
	if err != nil {
		return err
	}

	projectInfo, exist, err := s.projectDao.GetProjectDetailById(ctx, projectID)
	if err != nil {
		return err
	}

	if !exist {
		return status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	// 项目终止或者完成状态不能被编辑
	if projectInfo.State == common.ProjectCompleted ||
		projectInfo.State == common.ProjectTerminated {
		return status.Error(errcode.ErrProjectEditState, errcode.ProjectCodeMsg[errcode.ErrProjectEditState])
	}

	updateProject := &model.Project{
		Id:           projectInfo.Id,
		ProjectOwner: snowflake.MustParseString(req.ProjectOwner),
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		Comment:      req.Comment,
	}

	return with.DefaultTransaction(ctx, func(ctx context.Context) error {
		// 更新项目
		err = s.projectDao.UpdateProjectWithCols(ctx, updateProject, []string{"project_owner", "start_time", "end_time", "`comment`"})
		if err != nil {
			return err
		}

		// 更新成员
		members := &dto.ProjectMemberRequest{
			ProjectId: req.ProjectID,
			UserIds:   req.Members,
		}
		_, err = s.projectMemberService.ProjectMemberSave(ctx, members, loginUserID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *projectServiceImpl) ProjectModifyOwner(ctx context.Context, req *dto.ProjectModifyOwnerRequest, loginUserID snowflake.ID) error {
	// 权限校验 系统管理员角色才可转移权限
	isPMR, err := util.CheckProjectAdminRole(ctx, loginUserID, false)
	if err != nil {
		return err
	}

	if !isPMR {
		return status.Error(errcode.ErrProjectAccessPermission, errcode.ProjectCodeMsg[errcode.ErrProjectAccessPermission])
	}

	return with.DefaultTransaction(ctx, func(ctx context.Context) error {
		// 已结束项目不能修改项目管理员
		for _, ID := range req.ProjectIDs {
			projectID := snowflake.MustParseString(ID)
			projectMembers, err := s.projectMemberDao.GetProjectMembersByProjectId(ctx, projectID)
			if err != nil {
				return err
			}

			memberIDs := lo.Map[*model.ProjectMember, string](projectMembers, func(m *model.ProjectMember, _ int) string {
				return m.UserId.String()
			})

			contains := lo.Contains[string](memberIDs, req.TargetProjectOwnerID)
			if !contains {
				memberIDs = append(memberIDs, req.TargetProjectOwnerID)
				projectMemberReq := &dto.ProjectMemberRequest{
					ProjectId: ID,
					UserIds:   memberIDs,
				}
				_, err = s.projectMemberService.ProjectMemberSave(ctx, projectMemberReq, loginUserID)
				if err != nil {
					return err
				}
			}

			ownerId := snowflake.MustParseString(req.TargetProjectOwnerID)
			updateProject := &model.Project{
				Id:           projectID,
				ProjectOwner: ownerId,
			}

			err = s.projectDao.UpdateProjectWithCols(ctx, updateProject, []string{"project_owner"})
			if err != nil {
				return err
			}

		}

		return nil
	})
}

// 检查项目读取权限
func (s *projectServiceImpl) hasProjectAccessPermission(ctx context.Context, projectID, userID snowflake.ID) error {
	isPMR, err := util.CheckProjectAdminRole(ctx, userID, false)
	if err != nil {
		return err
	}

	has, err := s.projectMemberDao.ExistsProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}

	if !isPMR && !has {
		return status.Error(errcode.ErrProjectAccessPermission, errcode.ProjectCodeMsg[errcode.ErrProjectAccessPermission])
	}

	return nil
}

// 检查项目编辑权限
func (s *projectServiceImpl) hasProjectEditPermission(ctx context.Context, projectID, userID snowflake.ID) error {
	isPMR, err := util.CheckProjectAdminRole(ctx, userID, true)
	if err != nil {
		return err
	}

	detail, has, err := s.projectDao.GetProjectDetailById(ctx, projectID)
	if err != nil {
		return err
	}

	if !has {
		return status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	// 既没有系统管理员，也不是当前项目管理员，则不能编辑项目
	if !isPMR && detail.ProjectOwner != userID {
		return status.Error(errcode.ErrProjectAccessPermission, errcode.ProjectCodeMsg[errcode.ErrProjectAccessPermission])
	}

	return nil
}
