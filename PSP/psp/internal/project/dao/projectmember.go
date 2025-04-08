package dao

import (
	"context"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
)

type newProjectMemberDaoImpl struct {
}

func NewProjectMemberDao() *newProjectMemberDaoImpl {
	return &newProjectMemberDaoImpl{}
}

func (d *newProjectMemberDaoImpl) BatchInsertProjectMember(ctx context.Context, projectMembers []*model.ProjectMember) error {
	err := with.DefaultSession(ctx, func(session *xorm.Session) error {
		_, err := session.Insert(projectMembers)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *newProjectMemberDaoImpl) BatchDeleteProjectMember(ctx context.Context, projectID snowflake.ID, userIds []snowflake.ID) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		now := time.Now()
		_, err := session.In("user_id", userIds).And("project_id=?", projectID).
			Update(&model.ProjectMember{IsDelete: 1, UpdateTime: now})
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *newProjectMemberDaoImpl) GetProjectMembersByProjectId(ctx context.Context, projectID snowflake.ID) ([]*model.ProjectMember, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectMembers []*model.ProjectMember
	session = session.Table(model.ProjectMemberTableName)

	if projectID != 0 {
		session.Where("project_id = ? and is_delete = ?", projectID, 0)
	}

	err := session.Find(&projectMembers)
	if err != nil {
		return nil, err
	}

	return projectMembers, nil
}

func (d *newProjectMemberDaoImpl) GetProjectsByUserId(ctx context.Context, states []string, userID snowflake.ID) ([]*dto.ProjectMemberPbResp, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectList []*dto.ProjectMemberPbResp
	session = session.Table(model.ProjectMemberTableName).Alias("m").
		Join("INNER", []string{model.ProjectTableName, "p"}, "p.id = m.project_id").
		Where("p.is_delete = ? and m.is_delete = ? ", 0, 0)

	if userID > 0 {
		session.Where("m.user_id = ?", userID)
	}
	if len(states) > 0 {
		session.In("state", states)
	}

	session.Select("p.id, p.project_name, p.state, m.link_path")
	err := session.Find(&projectList)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (d *newProjectMemberDaoImpl) ExistsProjectMember(ctx context.Context, projectID, userID snowflake.ID) (bool, error) {
	session := boot.MW.DefaultSession(ctx)

	has, err := session.Where("project_id = ? and user_id=? and is_delete = ?", projectID, userID, common.Normal).Exist(&model.ProjectMember{})
	if err != nil {
		return false, err
	}

	return has, nil
}

func (d *newProjectMemberDaoImpl) GetProjectMembersByProjectIdAndUserIds(ctx context.Context, projectID snowflake.ID, userIDs []snowflake.ID) ([]*model.ProjectMember, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectMembers []*model.ProjectMember

	total, err := session.Table(model.ProjectMemberTableName).
		Where("project_id = ? and is_delete = ?", projectID, common.Normal).In("user_id", userIDs).FindAndCount(&projectMembers)
	if err != nil {
		return nil, 0, err
	}

	return projectMembers, total, nil
}

func (d *newProjectMemberDaoImpl) GetProjectMemberCountByProjectId(ctx context.Context, projectIds []snowflake.ID) ([]*dto.ProjectMemberCount, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectMemberCount []*dto.ProjectMemberCount
	session = session.Table(model.ProjectMemberTableName)
	if len(projectIds) > 0 {
		session.In("project_id", projectIds)
	}
	session.Select("project_id, count(*) as count").GroupBy("project_id")
	err := session.Find(&projectMemberCount)
	if err != nil {
		return nil, err
	}

	return projectMemberCount, nil
}
