package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	projectpb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// GetMemberProjectsByUserId 根据成员id获取所属项目列表
func (s *GRPCService) GetMemberProjectsByUserId(ctx context.Context, in *projectpb.GetMemberProjectsByUserIdRequest) (*projectpb.GetMemberProjectsByUserIdResponse, error) {
	logger := logging.GetLogger(ctx)
	uId := snowflake.MustParseString(in.UserId)
	if uId == snowflake.ID(0) {
		return &projectpb.GetMemberProjectsByUserIdResponse{}, nil
	}

	projectList, err := s.projectMemberDao.GetProjectsByUserId(ctx, nil, uId)
	if err != nil {
		logger.Errorf("get project list error, err: %v", err)
		return nil, err
	}

	projects := make([]*projectpb.MemberProjectObject, 0)

	if in.IncludeDefault {
		projects = append(projects, &projectpb.MemberProjectObject{
			ProjectId:   common.PersonalProjectID.String(),
			ProjectName: common.PersonalProjectName,
		})
	}

	for _, project := range projectList {
		projects = append(projects, &projectpb.MemberProjectObject{
			ProjectId:   project.ProjectID.String(),
			ProjectName: project.ProjectName,
			State:       project.State,
			LinkPath:    project.LinkPath,
		})
	}

	return &projectpb.GetMemberProjectsByUserIdResponse{
		Projects: projects,
	}, nil
}

// ExistsProjectMember 判断项目成员是否存在
func (s *GRPCService) ExistsProjectMember(ctx context.Context, in *projectpb.ExistsProjectMemberRequest) (*projectpb.ExistsProjectMemberResponse, error) {
	logger := logging.GetLogger(ctx)
	isExist, err := s.projectMemberDao.ExistsProjectMember(ctx, snowflake.MustParseString(in.ProjectId), snowflake.MustParseString(in.UserId))
	if err != nil {
		logger.Errorf("check project member error, err: %v", err)
		return nil, err
	}

	return &projectpb.ExistsProjectMemberResponse{
		IsExist: isExist,
	}, nil
}

// GetProjectMemberByProjectIdAndUserId 根据项目id 和用户id 获取项目成员信息
func (s *GRPCService) GetProjectMemberByProjectIdAndUserId(ctx context.Context, in *projectpb.GetProjectMemberByProjectIdAndUserIdRequest) (*projectpb.GetProjectMemberByProjectIdAndUserIdResponse, error) {
	logger := logging.GetLogger(ctx)

	projectId := snowflake.MustParseString(in.ProjectId)
	userIds := []snowflake.ID{snowflake.MustParseString(in.UserId)}
	projectMembers, total, err := s.projectMemberDao.GetProjectMembersByProjectIdAndUserIds(ctx, projectId, userIds)
	if err != nil {
		logger.Errorf("get project member error, err: %v", err)
		return nil, err
	}

	if total == 0 {
		return nil, status.Error(errcode.ErrProjectMemberNotExist, errcode.ProjectCodeMsg[errcode.ErrProjectMemberNotExist])
	}

	resp := util.Convert2ProtoProjectMember(projectMembers[0])
	return resp, nil
}
