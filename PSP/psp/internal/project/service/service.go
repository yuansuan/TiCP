package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type ProjectService interface {
	// ProjectSave 保存项目
	ProjectSave(ctx context.Context, req *dto.ProjectAddRequest, loginUserID snowflake.ID, loginUserName string) (*dto.ProjectAddResponse, error)
	// ProjectList 项目列表
	ProjectList(ctx context.Context, req *dto.ProjectListRequest, loginUserID snowflake.ID) (*dto.ProjectListResponse, error)
	// CurrentProjectList 当前项目列表
	CurrentProjectList(ctx context.Context, req *dto.CurrentProjectListRequest, loginUserID snowflake.ID) (*dto.CurrentProjectListResponse, error)
	// CurrentProjectListForParam 当前项目参数列表
	CurrentProjectListForParam(ctx context.Context, req *dto.CurrentProjectListForParamRequest, loginUserID snowflake.ID) (*dto.CurrentProjectListForParamResponse, error)
	// ProjectDetail 项目详情
	ProjectDetail(ctx context.Context, req *dto.ProjectDetailRequest, loginUserID snowflake.ID) (*dto.ProjectDetailResponse, error)
	// ProjectDelete  删除项目
	ProjectDelete(ctx *gin.Context, projectID string, loginUserID snowflake.ID) error
	// ProjectTerminate  终止项目
	ProjectTerminate(ctx context.Context, projectID string, loginUserID snowflake.ID) error
	// ProjectEdit 编辑项目
	ProjectEdit(ctx context.Context, req *dto.ProjectEditRequest, loginUserID snowflake.ID) error
	// ProjectModifyOwner 修改项目管理员
	ProjectModifyOwner(ctx context.Context, req *dto.ProjectModifyOwnerRequest, loginUserID snowflake.ID) error
}

type ProjectMemberService interface {
	// ProjectMemberSave 增加项目成员
	ProjectMemberSave(ctx context.Context, req *dto.ProjectMemberRequest, userID snowflake.ID) (*dto.ProjectMemberResponse, error)
	// CreatePersonalProjectPath 新增项目成员软链接
	CreatePersonalProjectPath(ctx context.Context, insertMembers []*model.ProjectMember, project *model.Project) error
	// DelPersonalProjectPath 删除成员软链接
	DelPersonalProjectPath(ctx context.Context, projectID snowflake.ID, userIDs []snowflake.ID)
}
