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
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/dbutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type projectDaoImpl struct {
	sid *snowflake.Node
}

func NewProjectDao() (*projectDaoImpl, error) {
	sid, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	return &projectDaoImpl{sid: sid}, nil
}

func (d *projectDaoImpl) InsertProject(ctx context.Context, project *model.Project) (string, error) {
	project.Id = d.sid.Generate()
	err := with.DefaultSession(ctx, func(session *xorm.Session) error {
		now := time.Now()
		project.CreateTime = now
		project.UpdateTime = now
		_, err := session.Insert(project)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return project.Id.String(), nil
}

func (d *projectDaoImpl) ExistSameProjectName(ctx context.Context, projectName string) (bool, error) {
	var result = false
	err := with.DefaultSession(ctx, func(session *xorm.Session) error {
		exist, err := session.Exist(&model.Project{
			ProjectName: projectName,
		})

		result = exist
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return result, nil
}

func (d *projectDaoImpl) UpdateProjectWithCols(ctx context.Context, project *model.Project, cols []string) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		project.UpdateTime = time.Now()
		cols = append(cols, "update_time")
		_, err := session.ID(project.Id).Cols(cols...).Update(project)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *projectDaoImpl) UpdateProject(ctx context.Context, project *model.Project) error {
	//TODO implement me
	panic("implement me")
}

func (d *projectDaoImpl) UpdateProjectStatus(ctx context.Context, start, end int64, state string) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	project := &model.Project{State: state}
	if start > 0 {
		session.Where("start_time <= ?", time.Unix(start, 0)).And("end_time >= ?", time.Unix(start, 0))
		session.And("state = ?", common.ProjectInit)
	}
	if end > 0 {
		session.Where("end_time <= ?", time.Unix(end, 0)).
			In("state", []string{common.ProjectInit, common.ProjectRunning})
	}
	_, err := session.Cols("state").Update(project)
	if err != nil {
		return err
	}
	return nil
}

func (d *projectDaoImpl) GetProjectDetailById(ctx context.Context, projectId snowflake.ID) (*model.Project, bool, error) {
	session := boot.MW.DefaultSession(ctx)

	project := &model.Project{}
	exist, err := session.ID(projectId).Get(project)
	if err != nil {
		return nil, false, err
	}
	return project, exist, nil
}

func (d *projectDaoImpl) GetProjectsDetailByIds(ctx context.Context, projectIds []snowflake.ID) ([]*model.Project, int64, error) {
	if len(projectIds) == 0 {
		return nil, 0, nil
	}

	session := boot.MW.DefaultSession(ctx)

	var projectList []*model.Project

	total, err := session.Where("is_delete = ? ", common.Normal).In("id", projectIds).FindAndCount(&projectList)
	if err != nil {
		return nil, 0, err
	}

	return projectList, total, nil
}

func (d *projectDaoImpl) GetProjectList(ctx context.Context, req *dto.ProjectListRequest, userID snowflake.ID, isSysRole bool) ([]*model.Project, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	index, size := req.Page.Index, req.Page.Size
	offset, err := xtype.GetPageOffset(index, size)
	if err != nil {
		return nil, 0, err
	}

	var projectList []*model.Project
	session = session.Table(model.ProjectTableName).
		Where("is_delete = ?", common.Normal)

	if strutil.IsNotEmpty(req.ProjectName) {
		session.Where("project_name like ?", "%"+req.ProjectName+"%")
	}

	if len(req.State) > 0 {
		session.In("state", req.State)
	}

	if req.StartTime > 0 && req.EndTime > 0 {
		startTime := time.Unix(req.StartTime, 0)
		endTime := time.Unix(req.EndTime, 0)
		session.Where("start_time <= ? and end_time >= ?", endTime, startTime)
	}

	if !isSysRole {
		session.Where("id in (select distinct(project_id) from project_member where user_id = ?  and is_delete = ? )", userID, common.Normal)
	}

	session = dbutil.WrapSortSession(session, req.OrderSort).Limit(int(size), int(offset))

	total, err := session.FindAndCount(&projectList)
	if err != nil {
		return nil, 0, err
	}

	return projectList, total, nil
}

func (d *projectDaoImpl) CurrentProjectListForParam(ctx context.Context, req *dto.CurrentProjectListForParamRequest, userID snowflake.ID, starttime, endTime time.Time) ([]*model.Project, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectList []*model.Project
	session.Table(model.ProjectTableName).Where("is_delete = ?", 0)

	if !starttime.IsZero() && !endTime.IsZero() {
		session.Where("create_time >= ? and create_time <= ?", starttime, endTime)
	}
	if !req.IsAdmin {
		session.Where("id in (select distinct(project_id) from project_member where user_id = ? and is_delete= ?)", userID, common.Normal)
	}

	// 需要展示数据中所有正在进行中的项目信息
	session.Or("state = ? and is_delete = ? and id in (select distinct(project_id) from project_member where user_id = ? and is_delete= ?)", common.ProjectRunning, 0, userID, common.Normal)

	total, err := session.FindAndCount(&projectList)
	if err != nil {
		return nil, 0, err
	}

	return projectList, total, nil
}

func (d *projectDaoImpl) GetProjectListByTimePeriod(ctx context.Context, start, end int64) ([]*model.Project, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectList []*model.Project
	session = session.Table(model.ProjectTableName).
		Where("is_delete = ?", 0)

	if start > 0 && end > 0 {
		session.And("end_time >= ?", time.Unix(start, 0))
		session.And("end_time <= ?", time.Unix(end, 0))
	}

	err := session.Find(&projectList)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (d *projectDaoImpl) GetRunningProjectIdsByTime(ctx context.Context, timePoint time.Time) ([]int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var projectIds []int64
	session.Table(model.ProjectTableName).Where("is_delete = ?", 0).
		Where("state = ?", common.ProjectRunning)
	if !timePoint.IsZero() {
		session.Where("end_time > ?", timePoint)
	}
	err := session.Cols("id").Find(&projectIds)
	if err != nil {
		return nil, err
	}

	return projectIds, nil
}
