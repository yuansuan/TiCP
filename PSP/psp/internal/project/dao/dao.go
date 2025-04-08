package dao

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type ProjectDao interface {
	// InsertProject 保存项目
	InsertProject(ctx context.Context, project *model.Project) (string, error)
	// ExistSameProjectName 是否存在相同项目名称
	ExistSameProjectName(ctx context.Context, projectName string) (bool, error)
	// UpdateProject 更新项目
	UpdateProject(ctx context.Context, project *model.Project) error
	// UpdateProjectStatus 更新项目
	UpdateProjectStatus(ctx context.Context, start, end int64, state string) error
	// UpdateProjectWithCols 根据指定列更新项目
	UpdateProjectWithCols(ctx context.Context, project *model.Project, cols []string) error
	// GetProjectDetailById 根据id查询项目
	GetProjectDetailById(ctx context.Context, projectId snowflake.ID) (*model.Project, bool, error)
	// GetProjectsDetailByIds 根据ids查询项目
	GetProjectsDetailByIds(ctx context.Context, projectIds []snowflake.ID) ([]*model.Project, int64, error)
	// GetProjectList 获取项目列表
	GetProjectList(ctx context.Context, req *dto.ProjectListRequest, userID snowflake.ID, isSysRole bool) ([]*model.Project, int64, error)
	// CurrentProjectList 获取当前项目列表
	CurrentProjectListForParam(ctx context.Context, req *dto.CurrentProjectListForParamRequest, userID snowflake.ID, starttime, endTime time.Time) ([]*model.Project, int64, error)
	// GetProjectListByTimePeriod 获取某个时间段内的项目列表
	GetProjectListByTimePeriod(ctx context.Context, start, end int64) ([]*model.Project, error)
	// GetRunningProjectIdsByTime 获取某个时间段内的运行中的项目id列表
	GetRunningProjectIdsByTime(ctx context.Context, timePoint time.Time) ([]int64, error)
}

type ProjectMemberDao interface {
	// BatchInsertProjectMember 批量保存项目成员
	BatchInsertProjectMember(ctx context.Context, projectMembers []*model.ProjectMember) error
	// BatchDeleteProjectMember 批量移除项目成员
	BatchDeleteProjectMember(ctx context.Context, projectID snowflake.ID, ids []snowflake.ID) error
	// GetProjectMembersByProjectId 根据项目id获取项目成员
	GetProjectMembersByProjectId(ctx context.Context, projectID snowflake.ID) ([]*model.ProjectMember, error)
	// GetProjectsByUserId 根据成员id获取所属项目列表
	GetProjectsByUserId(ctx context.Context, states []string, userID snowflake.ID) ([]*dto.ProjectMemberPbResp, error)
	// ExistsProjectMember 指定项目是否存在userID 用户
	ExistsProjectMember(ctx context.Context, projectID, userID snowflake.ID) (bool, error)
	// GetProjectMembersByProjectIdAndUserIds 根据项目id和用户id获取项目成员信息
	GetProjectMembersByProjectIdAndUserIds(ctx context.Context, projectID snowflake.ID, userIDs []snowflake.ID) ([]*model.ProjectMember, int64, error)
	// GetProjectMemberCountByProjectId 根据项目id获取项目成员数量
	GetProjectMemberCountByProjectId(ctx context.Context, projectIds []snowflake.ID) ([]*dto.ProjectMemberCount, error)
}
