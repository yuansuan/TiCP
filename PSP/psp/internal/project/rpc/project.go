package rpc

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// GetProjectDetailById 获取项目详情
func (s *GRPCService) GetProjectDetailById(ctx context.Context, in *project.GetProjectDetailByIdRequest) (*project.GetProjectByIdResponse, error) {
	projectID := snowflake.MustParseString(in.ProjectId)
	projectDetail, exist, err := s.projectDao.GetProjectDetailById(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}
	projectNumberCountMap := make(map[snowflake.ID]int64)

	resp := util.Convert2ProtoProjectDetail(projectDetail, projectNumberCountMap)
	return resp, nil
}

// GetProjectsDetailByIds 根据projectId 批量获取项目详情
func (s *GRPCService) GetProjectsDetailByIds(ctx context.Context, in *project.GetProjectsDetailByIdsRequest) (*project.GetProjectsByIdsResponse, error) {
	if len(in.ProjectIds) == 0 {
		return nil, status.Error(errcode.ErrInvalidParam, errcode.MsgInvalidParam)
	}

	projectIds := lo.Map[string, snowflake.ID](in.ProjectIds, func(id string, _ int) snowflake.ID {
		return snowflake.MustParseString(id)
	})

	projects, total, err := s.projectDao.GetProjectsDetailByIds(ctx, projectIds)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return nil, status.Error(errcode.ErrProjectNotFound, errcode.ProjectCodeMsg[errcode.ErrProjectNotFound])
	}

	projectNumberCountMap := make(map[snowflake.ID]int64)
	if in.IncludeMemberCount {
		projectMemberCount, err := s.projectMemberDao.GetProjectMemberCountByProjectId(ctx, projectIds)
		if err != nil {
			return nil, err
		}
		for _, v := range projectMemberCount {
			projectNumberCountMap[v.ProjectId] = v.Count
		}
	}

	respProjects := lo.Map[*model.Project, *project.GetProjectByIdResponse](projects, func(item *model.Project, _ int) *project.GetProjectByIdResponse {
		return util.Convert2ProtoProjectDetail(item, projectNumberCountMap)
	})

	return &project.GetProjectsByIdsResponse{Projects: respProjects}, nil
}

// GetProjectsIdByTimePeriod 根据起始时间和终止时间查询项目
func (s *GRPCService) GetProjectsIdByTimePeriod(ctx context.Context, in *project.GetProjectsIdByTimePeriodRequest) (*project.GetProjectsIdByTimePeriodResponse, error) {
	if in.GetStartTime() <= 0 || in.GetEndTime() <= 0 {
		return nil, status.Error(errcode.ErrInvalidParam, errcode.MsgInvalidParam)
	}

	projects, err := s.projectDao.GetProjectListByTimePeriod(ctx, in.StartTime, in.EndTime)
	if err != nil {
		return nil, err
	}

	respProjects := lo.Map[*model.Project, *project.GetProjectIdByTimePeriodResponse](projects, func(item *model.Project, _ int) *project.GetProjectIdByTimePeriodResponse {
		return util.Convert2ProtoProjectId(item)
	})

	return &project.GetProjectsIdByTimePeriodResponse{Projects: respProjects}, nil
}

// GetRunningProjectIdsByTime 根据起止时间查询正在进行的项目
func (s *GRPCService) GetRunningProjectIdsByTime(ctx context.Context, in *project.GetRunningProjectIdsByTimeRequest) (*project.GetRunningProjectIdsByTimeResponse, error) {
	projectIds, err := s.projectDao.GetRunningProjectIdsByTime(ctx, in.TimePoint.AsTime())
	if err != nil {
		return nil, err
	}

	return &project.GetRunningProjectIdsByTimeResponse{ProjectIds: projectIds}, nil
}
