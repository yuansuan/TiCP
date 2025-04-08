package dao

import (
	"context"
	"fmt"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

var JobStatisticStatuses = []string{
	consts.JobStateRunning,
	consts.JobStateTerminated,
	consts.JobStateSuspended,
	consts.JobStateCompleted,
	consts.JobStateFailed,
}

func (d *jobDaoImpl) GetTop5ProjectByCPUTime(ctx context.Context, projectIds []snowflake.ID, start, end int64) ([]*dto.ProjectCPUTime, error) {
	session := boot.MW.DefaultSession(ctx)

	var top5Projects []*dto.ProjectCPUTime
	session = session.Table(&model.Job{}).In("state", JobStatisticStatuses).In("project_id", projectIds)
	if start > 0 {
		session.Where("submit_time >= ?", time.UnixMilli(start))
	}
	if end > 0 {
		session.Where("submit_time <= ?", time.UnixMilli(end))
	}
	session.GroupBy("project_id").Desc("cpu_time").Limit(5)
	err := session.Select("project_id, project_name, sum(cast(floor(cast(exec_duration * cpus_alloc / 3600 as decimal(26, 6)) * 100000) / 100000 as decimal(25, 5))) as cpu_time").Find(&top5Projects)
	if err != nil {
		return nil, err
	}

	return top5Projects, nil
}

func (d *jobDaoImpl) GetJobCountByProjectIds(ctx context.Context, projectIds []snowflake.ID, start, end int64) ([]*dto.ProjectJobCount, error) {
	session := boot.MW.DefaultSession(ctx)
	var projectJobCounts []*dto.ProjectJobCount
	session = session.Table(&model.Job{}).In("project_id", projectIds)
	if start > 0 {
		session.Where("submit_time >= ?", time.UnixMilli(start))
	}
	if end > 0 {
		session.Where("submit_time <= ?", time.UnixMilli(end))
	}
	session.GroupBy("project_id").OrderBy("project_id")
	err := session.Select("project_id, project_name, count(*) as count").Find(&projectJobCounts)
	if err != nil {
		return nil, err
	}

	return projectJobCounts, nil
}

func (d *jobDaoImpl) GetJobCPUTimeTotal(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64) (float64, error) {
	session := boot.MW.DefaultSession(ctx)

	var statisticsTotal float64
	session = session.Table(&model.Job{})
	wrapParamCondition(session, names, projectIds, "", queryType, computeType, startTime, endTIme)

	_, err := session.Select("sum(cast(floor(cast(exec_duration * cpus_alloc / 3600 as decimal(26, 6)) * 100000) / 100000 as decimal(25, 5))) as cpu_time").Get(&statisticsTotal)
	if err != nil {
		return 0.0, err
	}

	return statisticsTotal, nil
}

func (d *jobDaoImpl) GetJobStatisticsOverview(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64, pageIndex, pageSize int) ([]*model.StatisticsJob, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var overviews []*model.StatisticsJob
	session = session.Table(&model.Job{})
	wrapParamCondition(session, names, projectIds, "", queryType, computeType, startTime, endTIme)

	selectColumn := ""
	switch queryType {
	case consts.JobStatisticsQueryTypeApp:
		selectColumn = "app_id, app_name, type, project_name"
	case consts.JobStatisticsQueryTypeUser:
		selectColumn = "user_id, user_name, type, project_name"
	default:
		return nil, 0, fmt.Errorf("query type not match")
	}

	session.GroupBy(selectColumn)
	wrapPageCondition(session, pageIndex, pageSize)

	total, err := session.Select(selectColumn + ", sum(cast(floor(cast(exec_duration * cpus_alloc / 3600 as decimal(26, 6)) * 100000) / 100000 as decimal(25, 5))) as cpu_time").FindAndCount(&overviews)
	if err != nil {
		return nil, 0, err
	}

	return overviews, total, nil
}

func (d *jobDaoImpl) GetJobStatisticsDetail(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTIme int64, pageIndex, pageSize int) ([]*model.StatisticsJob, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var jobList []*model.StatisticsJob
	session = session.Table(&model.Job{})
	wrapParamCondition(session, names, projectIds, consts.JobStatisticsShowTypeDetail, queryType, computeType, startTime, endTIme)
	wrapPageCondition(session, pageIndex, pageSize)

	total, err := session.Select("id, name, type, project_name, app_name, user_name, submit_time, start_time, end_time, cast(floor(cast(exec_duration * cpus_alloc / 3600 as decimal(26, 6)) * 100000) / 100000 as decimal(25, 5)) as cpu_time").FindAndCount(&jobList)
	if err != nil {
		return nil, 0, err
	}

	return jobList, total, nil
}

func wrapParamCondition(session *xorm.Session, names, projectIds []string, showType, queryType string, computeType string, startTime int64, endTIme int64) {
	if len(names) > 0 {
		switch queryType {
		case consts.JobStatisticsQueryTypeApp:
			session.In("app_name", names).Asc("app_name")
		case consts.JobStatisticsQueryTypeUser:
			session.In("user_name", names).Asc("user_name")
		}
	}
	if len(projectIds) > 0 {
		session.In("project_id", snowflake.BatchParseStringToID(projectIds))
	}
	if computeType != "" {
		session.Where("type = ?", computeType)
	}
	if startTime > 0 {
		session.Where("submit_time >= ?", time.Unix(startTime, 0))
	}
	if endTIme > 0 {
		session.Where("submit_time <= ?", time.Unix(endTIme, 0))
	}

	// 作业统计: 只统计处于 "已完成状态" 的作业数据
	session.In("state", consts.JobStateTerminated, consts.JobStateCompleted, consts.JobStateFailed)

	// 排除执行时间为 0 的数据
	session.Where("exec_duration > 0")

	if consts.JobStatisticsShowTypeDetail == showType {
		// 优先展示最新提交的作业数据
		session.Desc("submit_time")
	}
}

func wrapPageCondition(session *xorm.Session, pageIndex int, pageSize int) {
	if pageIndex > 0 {
		session.Limit(pageSize, (pageIndex-1)*pageSize)
	} else {
		session.Limit(pageSize)
	}
}
