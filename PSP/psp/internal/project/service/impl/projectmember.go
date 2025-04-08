package impl

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/samber/lo"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	userMgr "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
)

type projectMemberServiceImpl struct {
	projectMemberDao dao.ProjectMemberDao
	projectDao       dao.ProjectDao
	sid              *snowflake.Node
	rpc              *client.GRPC
}

func NewProjectMemberService() (*projectMemberServiceImpl, error) {
	projectMemberDao := dao.NewProjectMemberDao()
	projectDao, err := dao.NewProjectDao()
	if err != nil {
		return nil, err
	}
	rpc := client.GetInstance()
	instance, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("create snowflake instance err: %v", err)
		return nil, err
	}
	return &projectMemberServiceImpl{
		projectMemberDao: projectMemberDao,
		projectDao:       projectDao,
		sid:              instance,
		rpc:              rpc,
	}, nil
}

func (s *projectMemberServiceImpl) ProjectMemberSave(ctx context.Context, req *dto.ProjectMemberRequest, userID snowflake.ID) (*dto.ProjectMemberResponse, error) {
	err := with.DefaultTransaction(ctx, func(ctx context.Context) error {
		projectId := snowflake.MustParseString(req.ProjectId)

		project, has, err := s.projectDao.GetProjectDetailById(ctx, projectId)
		if err != nil {
			return err
		}

		if !has {
			return status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
		}

		//1.查询出该项目下的所有成员
		members, err := s.projectMemberDao.GetProjectMembersByProjectId(ctx, projectId)
		if err != nil {
			return err
		}

		//2.过滤出新增的
		var nowUserIds []string
		for _, member := range members {
			nowUserIds = append(nowUserIds, member.UserId.String())
		}
		insertUserIds := stringSliceDifference(req.UserIds, nowUserIds) //要新增的

		//3.批量添加项目成员
		now := time.Now()
		insertMembers := lo.Map[snowflake.ID, *model.ProjectMember](insertUserIds, func(userID snowflake.ID, _ int) *model.ProjectMember {
			return &model.ProjectMember{
				Id:         s.sid.Generate(),
				ProjectId:  projectId,
				UserId:     userID,
				CreateTime: now,
				UpdateTime: now,
			}
		})
		// 4.批量添加项目成员
		if len(insertMembers) != 0 {
			// 组装文件链接
			formatTime := timeutil.FormatTime(now, common.YearMonthDayFormat)
			userIdentities := lo.Map[*model.ProjectMember, *userMgr.UserIdentity](insertMembers, func(m *model.ProjectMember, _ int) *userMgr.UserIdentity {
				return &userMgr.UserIdentity{
					Id: m.UserId.String(),
				}
			})

			users, err := client.GetInstance().User.BatchGetUser(ctx, &userMgr.UserIdentities{UserIdentities: userIdentities})
			if err != nil {
				return err
			}

			userMaps := lo.KeyBy[snowflake.ID, *userMgr.UserObj](users.UserObj, func(u *userMgr.UserObj) snowflake.ID {
				return snowflake.MustParseString(u.Id)
			})

			for _, member := range insertMembers {
				user, ok := userMaps[member.UserId]
				if ok {
					member.LinkPath = path.Join("/", user.Name, common.WorkspaceFolderPath, project.ProjectName+"_"+formatTime)
				}
			}

			err = s.projectMemberDao.BatchInsertProjectMember(ctx, insertMembers)
			if err != nil {
				return err
			}

			// 调用创建storage用户链接接口
			err = s.CreatePersonalProjectPath(ctx, insertMembers, project)
			if err != nil {
				return err
			}
		}

		//5.过滤出删除的
		deleteUserIds := stringSliceDifference(nowUserIds, req.UserIds) //要删除的

		//6.批量删除项目成员
		if len(deleteUserIds) != 0 {
			err = s.projectMemberDao.BatchDeleteProjectMember(ctx, projectId, deleteUserIds)
			if err != nil {
				return err
			}

			// 删除用户对应的软链接目录
			go s.DelPersonalProjectPath(ctx, projectId, deleteUserIds)
		}

		//7.发送消息
		//发送消息(新增的成员)
		for _, userId := range insertUserIds {
			content := fmt.Sprintf("您已加入项目[编号:%v 项目名:%v]", project.Id, project.ProjectName)
			sendMessage(userId, content, s.rpc, ctx)
		}

		//发送消息(移出的成员)
		for _, userId := range deleteUserIds {
			content := fmt.Sprintf("您已退出项目[编号:%v 项目名:%v]", project.Id, project.ProjectName)
			sendMessage(userId, content, s.rpc, ctx)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ProjectMemberResponse{
		ID: req.ProjectId,
	}, nil
}

func (s *projectMemberServiceImpl) CreatePersonalProjectPath(ctx context.Context, insertMembers []*model.ProjectMember, project *model.Project) error {
	if len(insertMembers) == 0 {
		return nil
	}

	logger := logging.GetLogger(ctx)

	srcPath := path.Join("/", common.ProjectFolderPath, project.ProjectName)
	for _, member := range insertMembers {
		symLinkReq := &storage.SymLinkReq{
			SrcPath:   srcPath,
			DstPath:   member.LinkPath,
			Overwrite: true,
			Cross:     true,
			IsCloud:   false,
		}

		symLinkResp, err := client.GetInstance().Storage.SymLink(ctx, symLinkReq)
		if err != nil {
			logger.Errorf("create member:[%v] sym link err:[%v], resp: [%v]", member, err, symLinkResp)
			return err
		}
	}

	return nil
}

func (s *projectMemberServiceImpl) DelPersonalProjectPath(ctx context.Context, projectID snowflake.ID, userIDs []snowflake.ID) {
	logger := logging.GetLogger(ctx)

	projectMembers, total, err := s.projectMemberDao.GetProjectMembersByProjectIdAndUserIds(ctx, projectID, userIDs)
	if err != nil {
		logger.Errorf("DelPersonalProjectPath get members data err: %v", err)
		return
	}

	if total > 0 {
		for _, member := range projectMembers {
			// 删除对应的文件软链接
			rmReq := &storage.RmReq{
				Paths:   []string{member.LinkPath},
				Cross:   true,
				IsCloud: false,
			}
			rmResp, err := client.GetInstance().Storage.Rm(ctx, rmReq)
			if err != nil {
				logger.Errorf("del memeber:[%v], path:[%v], err:[%v], resp:[%v]", member.Id, member.LinkPath, err, rmResp)
				continue
			}
		}
	}
}

func sendMessage(userId snowflake.ID, content string, rpc *client.GRPC, ctx context.Context) error {
	msg := &pbnotice.WebsocketMessage{
		UserId:  userId.String(),
		Type:    common.ProjectEventType,
		Content: content,
	}

	if _, err := rpc.Notice.SendWebsocketMessage(ctx, msg); err != nil {
		logging.GetLogger(ctx).Errorf("project end send err: %v", err)
		return err
	}
	return nil
}

func stringSliceDifference(newSlice, oldSlice []string) []snowflake.ID {
	m := make(map[string]bool)

	for _, item := range oldSlice {
		m[item] = true
	}

	var diffs = make([]snowflake.ID, 0)

	for _, item := range newSlice {
		if _, exists := m[item]; !exists {
			diffs = append(diffs, snowflake.MustParseString(item))
		}
	}

	return diffs
}
