package util

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

func Convert2ProjectList(projectList []*model.Project, loginUserID snowflake.ID) []*dto.ProjectListInfo {
	if len(projectList) == 0 {
		return []*dto.ProjectListInfo{}
	}

	resp := make([]*dto.ProjectListInfo, 0, len(projectList))
	for _, project := range projectList {
		projectListInfo := Convert2ProjectInfo(project, loginUserID)
		if projectListInfo == nil {
			continue
		}

		resp = append(resp, projectListInfo)
	}

	return resp
}

func Convert2ProjectInfo(project *model.Project, loginUserID snowflake.ID) *dto.ProjectListInfo {
	if project == nil {
		return nil
	}

	projectListInfo := &dto.ProjectListInfo{
		ID:             project.Id.String(),
		ProjectName:    project.ProjectName,
		ProjectOwnerID: project.ProjectOwner.String(),
		State:          project.State,
		StartTime:      timeutil.DefaultFormatTime(project.StartTime),
		EndTime:        timeutil.DefaultFormatTime(project.EndTime),
		Comment:        project.Comment,
		CreateTime:     timeutil.DefaultFormatTime(project.CreateTime),
		IsProjectOwner: project.ProjectOwner == loginUserID,
	}

	return projectListInfo
}

func Convert2ProtoProjectDetail(project *model.Project, projectNumberCountMap map[snowflake.ID]int64) *pb.GetProjectByIdResponse {
	if project == nil {
		return nil
	}

	return &pb.GetProjectByIdResponse{
		ProjectId:    project.Id.String(),
		ProjectName:  project.ProjectName,
		ProjectOwner: project.ProjectOwner.String(),
		State:        project.State,
		FilePath:     project.FilePath,
		MemberCount:  projectNumberCountMap[project.Id],
		StartTime:    timestamppb.New(project.StartTime),
		EndTime:      timestamppb.New(project.EndTime),
		CreateTime:   timestamppb.New(project.CreateTime),
	}

}

func Convert2ProtoProjectId(project *model.Project) *pb.GetProjectIdByTimePeriodResponse {
	if project == nil {
		return nil
	}

	return &pb.GetProjectIdByTimePeriodResponse{
		ProjectId: int64(project.Id),
		EndTime:   timestamppb.New(project.EndTime),
	}

}

func Convert2ProtoProjectMember(projectMember *model.ProjectMember) *pb.GetProjectMemberByProjectIdAndUserIdResponse {
	if projectMember == nil {
		return nil
	}

	return &pb.GetProjectMemberByProjectIdAndUserIdResponse{
		Id:         projectMember.Id.String(),
		ProjectId:  projectMember.ProjectId.String(),
		UserId:     projectMember.UserId.String(),
		FilePath:   projectMember.LinkPath,
		CreateTime: timestamppb.New(projectMember.CreateTime),
	}
}
