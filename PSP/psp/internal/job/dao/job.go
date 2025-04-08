package dao

import (
	"context"
	"fmt"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"reflect"
	"time"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/dbutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type jobDaoImpl struct {
	sid *snowflake.Node
}

// NewJobDao 创建JobDao
func NewJobDao() (JobDao, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	return &jobDaoImpl{sid: node}, nil
}

// InsertJob 保存作业信息
func (d *jobDaoImpl) InsertJob(ctx context.Context, job *model.Job) (snowflake.ID, error) {
	session := boot.MW.DefaultSession(ctx)

	job.Id = d.sid.Generate()
	job.UpdateTime = time.Now()
	_, err := session.Insert(job)
	if err != nil {
		return 0, err
	}

	return job.Id, nil
}

// UpdateJob 更新作业信息
func (d *jobDaoImpl) UpdateJob(ctx context.Context, job *model.Job) error {
	session := boot.MW.DefaultSession(ctx)

	job.UpdateTime = time.Now()
	cols := []string{"app_id", "user_name", "update_time"}
	_, err := session.Where("id=?", job.Id).Cols(cols...).Update(job)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBurstJob 更新爆发作业信息
func (d *jobDaoImpl) UpdateBurstJob(ctx context.Context, job *model.Job) error {
	session := boot.MW.DefaultSession(ctx)

	updatedCols := make([]string, 0, reflect.TypeOf(*job).NumField())
	if job.State != "" {
		updatedCols = append(updatedCols, "state")
	}

	if job.DataState != "" {
		updatedCols = append(updatedCols, "data_state")
	}

	if job.UploadTaskId != "" {
		updatedCols = append(updatedCols, "upload_task_id")
	}

	if job.OutJobId != "" {
		updatedCols = append(updatedCols, "out_job_id")
	}

	if job.Type != "" {
		updatedCols = append(updatedCols, "type")
	}

	if job.Queue != "" {
		updatedCols = append(updatedCols, "queue")
	}

	if job.AppName != "" {
		updatedCols = append(updatedCols, "app_name")
	}

	if job.AppId != 0 {
		updatedCols = append(updatedCols, "app_id")
	}

	if job.BurstNum != 0 {
		updatedCols = append(updatedCols, "burst_num")
	}

	if len(job.VisAnalysis) > 0 {
		updatedCols = append(updatedCols, "vis_analysis")
	}

	_, err := session.Where("id=?", job.Id).Cols(updatedCols...).Update(job)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDataState 更新作业数据状态
func (d *jobDaoImpl) UpdateDataState(ctx context.Context, outJobID, dataState string) error {
	session := boot.MW.DefaultSession(ctx)

	job := &model.Job{}
	job.DataState = dataState
	job.UpdateTime = time.Now()
	cols := []string{"data_state", "update_time"}
	_, err := session.Where("out_job_id=?", outJobID).Cols(cols...).Update(job)
	if err != nil {
		return err
	}

	return nil
}

// UpdateJobWithCols 指定列更新作业信息
func (d *jobDaoImpl) UpdateJobWithCols(ctx context.Context, job *model.Job, cols []string) error {
	session := boot.MW.DefaultSession(ctx)

	job.UpdateTime = time.Now()
	cols = append(cols, "update_time")
	_, err := session.Where("id=?", job.Id).Cols(cols...).Update(job)
	if err != nil {
		return err
	}

	return nil
}

// GetJobUserNameList 获取作业用户名称列表
func (d *jobDaoImpl) GetJobUserNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var names []string
	session = session.Table(&model.Job{})

	if computeType != "" {
		session.Where("type = ?", computeType)
	}
	if len(projectIds) > 0 {
		session.In("project_id", snowflake.BatchParseStringToID(projectIds))
	}
	if !isAdmin {
		session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
	}

	if err := session.Distinct("user_name").Asc("user_name").Find(&names); err != nil {
		return nil, err
	}

	return names, nil
}

// GetJobComputeTypeList 获取作业计算类型列表
func (d *jobDaoImpl) GetJobComputeTypeList(ctx context.Context, isAdmin bool, loginUserID snowflake.ID) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var names []string
	if err := session.Table(&model.Job{}).Distinct("type").Asc("type").Find(&names); err != nil {
		return nil, err
	}
	if !isAdmin {
		session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
	}

	return names, nil
}

// GetJobSetNameList 获取作业集名称列表
func (d *jobDaoImpl) GetJobSetNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var names []string
	session = session.Table(&model.Job{})

	if computeType != "" {
		session.Where("type = ?", computeType)
	}
	if len(projectIds) > 0 {
		session.In("project_id", snowflake.BatchParseStringToID(projectIds))
	}
	if !isAdmin {
		session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
	}

	err := session.Where("job_set_name is not null and job_set_name != ''").Distinct("job_set_name").Asc("job_set_name").Find(&names)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// GetJobAppNameList 获取作业应用名称列表
func (d *jobDaoImpl) GetJobAppNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var names []string
	session = session.Table(&model.Job{})

	if computeType != "" {
		session.Where("type = ?", computeType)
	}
	if len(projectIds) > 0 {
		session.In("project_id", snowflake.BatchParseStringToID(projectIds))
	}
	if !isAdmin {
		session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
	}

	if err := session.Distinct("app_name").Asc("app_name").Find(&names); err != nil {
		return nil, err
	}

	return names, nil
}

// GetJobQueueNameList 获取作业队列名称列表
func (d *jobDaoImpl) GetJobQueueNameList(ctx context.Context, projectIds []string, computeType string, isAdmin bool, loginUserID snowflake.ID) ([]string, error) {
	session := boot.MW.DefaultSession(ctx)

	var queues []string
	session = session.Table(&model.Job{})

	if computeType != "" {
		session.Where("type = ?", computeType)
	}
	if len(projectIds) > 0 {
		session.In("project_id", snowflake.BatchParseStringToID(projectIds))
	}
	if !isAdmin {
		session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
	}

	if err := session.Distinct("queue").Asc("queue").Find(&queues); err != nil {
		return nil, err
	}

	return queues, nil
}

// GetJobByOutID 获取作业详细信息
func (d *jobDaoImpl) GetJobByOutID(ctx context.Context, outJobID, jobType string) (bool, *model.Job, error) {
	session := boot.MW.DefaultSession(ctx)

	job := &model.Job{}
	if jobType != "" {
		session.Where("type=?", jobType)
	}

	has, err := session.Where("out_job_id=?", outJobID).Get(job)
	if err != nil {
		return false, nil, err
	}

	return has, job, err
}

// GetJobListByOutJobID 获取作业信息列表
func (d *jobDaoImpl) GetJobListByOutJobID(ctx context.Context, outJobIDList []string) ([]*model.Job, error) {
	session := boot.MW.DefaultSession(ctx)

	var jobList []*model.Job
	if err := session.In("out_job_id", outJobIDList).Find(&jobList); err != nil {
		return nil, err
	}

	return jobList, nil
}

// GetJobDetail 获取作业详细信息
func (d *jobDaoImpl) GetJobDetail(ctx context.Context, jobID snowflake.ID) (bool, *model.Job, error) {
	session := boot.MW.DefaultSession(ctx)

	job := &model.Job{}
	has, err := session.Where("id=?", jobID).Get(job)
	if err != nil {
		return false, nil, err
	}

	return has, job, err
}

// GetUnfinishedJobList 分页获取未结束作业信息列表
func (d *jobDaoImpl) GetUnfinishedJobList(ctx context.Context, page *xtype.Page) ([]*model.Job, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var jobs []*model.Job
	index, size := page.Index, page.Size
	offset, err := xtype.GetPageOffset(index, size)
	if err != nil {
		return nil, 0, err
	}

	total, err := session.In("state", util.UnFinishedStates).Limit(int(size), int(offset)).FindAndCount(&jobs)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetFinishedJobList 分页获取已结束作业信息列表
//func (d *jobDaoImpl) GetFinishedJobList(ctx context.Context, page *ptype.Page) ([]*model.Job, int64, error) {
//	session := boot.MW.DefaultSession(ctx)
//
//	var jobs []*model.Job
//	limit, offset := page.LimitOffset()
//	total, err := session.In("state", FinishedStates).Limit(int(limit), int(offset)).FindAndCount(&jobs)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return jobs, total, nil
//}

// GetJobList 分页获取作业列表
func (d *jobDaoImpl) GetJobList(ctx context.Context, filter *dto.JobFilter, page *xtype.Page, orderSort *xtype.OrderSort, isAdmin bool, loginUserID snowflake.ID) ([]*model.Job, int64, error) {
	session := boot.MW.DefaultSession(ctx)

	var jobs []*model.Job
	index, size := page.Index, page.Size
	offset, err := xtype.GetPageOffset(index, size)
	if err != nil {
		return nil, 0, err
	}

	total, err := wrapListSession(session, filter, orderSort, isAdmin, loginUserID).Limit(int(size), int(offset)).FindAndCount(&jobs)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetAppJobNum 获取应用作业数量
func (d *jobDaoImpl) GetAppJobNum(ctx context.Context, start, end int64) ([]*dto.AppJobInfo, error) {
	session := boot.MW.DefaultSession(ctx)
	statistics := make([]*dto.AppJobInfo, 0)
	session = session.Table(&model.Job{})
	if start > 0 {
		session.Where("submit_time >= ?", time.UnixMilli(start))
	}

	if end > 0 {
		session.Where("submit_time <= ?", time.UnixMilli(end))
	}
	session.GroupBy("app_name")
	session.OrderBy("num desc").Limit(5)
	session.Select("count(id) as num, app_name ")
	if err := session.Find(&statistics); err != nil {
		return nil, err
	}
	return statistics, nil
}

// GetJobStatusNum 作业作业状态数量
func (d *jobDaoImpl) GetJobStatusNum(ctx context.Context, start int64, states []string) ([]*dto.JobStatus, error) {
	session := boot.MW.DefaultSession(ctx)
	var jobStatusNums []*dto.JobStatus
	session = session.Table(&model.Job{})
	if start > 0 {
		session.Where("submit_time >= ?", time.Unix(start, 0))
	}
	session.In("state", states)

	session.GroupBy("state")
	session.Select("count(id) as num, state")
	if err := session.Find(&jobStatusNums); err != nil {
		return nil, err
	}
	return jobStatusNums, nil
}

// GetUserJobNum 获取用户作业数量
func (d *jobDaoImpl) GetUserJobNum(ctx context.Context, start, end int64) ([]*dto.UserJobInfo, error) {
	session := boot.MW.DefaultSession(ctx)
	statistics := make([]*dto.UserJobInfo, 0)
	session = session.Table(&model.Job{})
	if start > 0 {
		session.Where("submit_time >= ?", time.UnixMilli(start))
	}

	if end > 0 {
		session.Where("submit_time <= ?", time.UnixMilli(end))
	}
	session.GroupBy("user_name")
	session.OrderBy("num desc").Limit(5)
	session.Select("count(id) as num, user_name ")
	if err := session.Find(&statistics); err != nil {
		return nil, err
	}
	return statistics, nil
}

// GetJobCPUTimeMetric 作业核时运行指标统计
func (d *jobDaoImpl) GetJobCPUTimeMetric(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string, states []string) ([]*dto.JobCPUTimeQueryMetric, error) {
	session := boot.MW.DefaultSession(ctx)

	job := &model.Job{}
	var metrics []*dto.JobCPUTimeQueryMetric
	session.Table(job.TableName()).Select(groupByCol + " as group_col, sum(cpus_alloc * exec_duration) / 3600 as cpu_time")
	if filter.StartTime > 0 {
		session.Where("create_time > ?", time.UnixMilli(filter.StartTime))
	}

	if filter.EndTime > 0 {
		session.Where("create_time < ?", time.UnixMilli(filter.EndTime))
	}

	session.In("state", states)
	session.GroupBy("group_col").OrderBy("cpu_time desc").Limit(filter.TopSize)
	err := session.Find(&metrics)

	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// GetJobCountMetric 获取应用和用户数量统计指标
func (d *jobDaoImpl) GetJobCountMetric(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string, states []string) ([]*dto.JobQueryResultMetric, error) {
	session := boot.MW.DefaultSession(ctx)

	job := model.Job{}
	var metrics []*dto.JobQueryResultMetric
	session.Table(job.TableName()).Select(fmt.Sprintf("%s as item, count(0) as count ", groupByCol))
	if filter.StartTime > 0 {
		session.Where("create_time > ?", time.UnixMilli(filter.StartTime))
	}

	if filter.EndTime > 0 {
		session.Where("create_time < ?", time.UnixMilli(filter.EndTime))
	}

	if len(states) > 0 {
		session.In("state", states)
	}
	session.GroupBy("item").OrderBy("count desc")
	err := session.Find(&metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (d *jobDaoImpl) GetJobWaitStatistic(ctx context.Context, filter *dto.JobMetricFiler, statisticType string, states []string) ([]*dto.JobQueryResultMetric, error) {
	session := boot.MW.DefaultSession(ctx)

	selectSql := "ROUND(UNIX_TIMESTAMP(DATE_FORMAT(create_time,'%Y-%m-%d')) * 1000, 0) as item"

	switch statisticType {
	case consts.JobWaitTimeStatisticAvg:
		selectSql = selectSql + ", ROUND(AVG(TIME_TO_SEC(TIMEDIFF(start_time,pend_time)))/3600, 2) as count"
	case consts.JobWaitTimeStatisticMax:
		selectSql = selectSql + ", ROUND(MAX(TIME_TO_SEC(TIMEDIFF(start_time,pend_time)))/3600, 2) as count"
	case consts.JobWaitTimeStatisticTotal:
		selectSql = selectSql + ", ROUND(SUM(TIME_TO_SEC(TIMEDIFF(start_time,pend_time)))/3600, 2) as count"
	case consts.JobWaitNumStatistic:
		selectSql = selectSql + ", COUNT(0) as count"
	}

	job := model.Job{}
	var metrics []*dto.JobQueryResultMetric

	session.Table(job.TableName()).Select(selectSql)
	session.Where("pend_time < start_time")
	if filter.StartTime > 0 {
		session.Where("create_time > ?", time.UnixMilli(filter.StartTime))
	}

	if filter.EndTime > 0 {
		session.Where("create_time < ?", time.UnixMilli(filter.EndTime))
	}

	if len(states) > 0 {
		session.In("state", states)
	}
	session.GroupBy("item").OrderBy("item asc")
	err := session.Find(&metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil

}

// GetJobDeliverCount 作业提交数量指标
func (d *jobDaoImpl) GetJobDeliverCount(ctx context.Context, filter *dto.JobMetricFiler, groupByCol string) ([]*dto.JobQueryResultMetric, error) {
	session := boot.MW.DefaultSession(ctx)

	var metrics []*dto.JobQueryResultMetric
	job := model.Job{}
	var querySql string
	if groupByCol == consts.JobDeliveryCountUser {
		querySql = fmt.Sprintf(
			"select t.day as item, count(t.totalCount) as count "+
				"from (select ROUND(UNIX_TIMESTAMP(DATE_FORMAT(create_time,'%%Y-%%m-%%d')) * 1000, 0) as day, count(user_id) as totalCount "+
				"      from %s "+
				"      where (create_time > '%v') AND (create_time < '%v') "+
				"      group by day, user_id"+
				") t "+
				"group by item "+
				"order by item asc ",
			job.TableName(),
			time.UnixMilli(filter.StartTime), time.UnixMilli(filter.EndTime))
	} else if groupByCol == consts.JobDeliveryCountJob {
		querySql = fmt.Sprintf(
			"select ROUND(UNIX_TIMESTAMP(DATE_FORMAT(create_time,'%%Y-%%m-%%d')) * 1000, 0) as item, count(id) as count "+
				"from %s "+
				"where (create_time > '%v') AND (create_time < '%v') "+
				"group by item "+
				"order by item asc ",
			job.TableName(),
			time.UnixMilli(filter.StartTime), time.UnixMilli(filter.EndTime))
	}

	err := session.SQL(querySql).Find(&metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func wrapListSession(session *xorm.Session, filter *dto.JobFilter, orderSort *xtype.OrderSort, isAdmin bool, loginUserID snowflake.ID) *xorm.Session {
	filterSession := wrapFilterSession(session, filter, isAdmin, loginUserID)
	return dbutil.WrapSortSession(filterSession, orderSort)
}

func wrapFilterSession(session *xorm.Session, filter *dto.JobFilter, isAdmin bool, loginUserID snowflake.ID) *xorm.Session {
	if filter != nil {
		if !strutil.IsEmpty(filter.JobID) {
			session.Where("id = ?", snowflake.MustParseString(filter.JobID))
		}

		if !strutil.IsEmpty(filter.JobSetID) {
			session.Where("job_set_id = ?", snowflake.MustParseString(filter.JobSetID))
		}

		if len(filter.ProjectIDs) > 0 {
			session.In("project_id", snowflake.BatchParseStringToID(filter.ProjectIDs))
		}

		if len(filter.JobTypes) > 0 {
			session.In("type", filter.JobTypes)
		}

		if len(filter.UserNames) > 0 {
			session.In("user_name", filter.UserNames)
		}

		if len(filter.Queues) > 0 {
			session.In("queue", filter.Queues)
		}

		if len(filter.States) > 0 {
			session.In("state", filter.States)
		}

		if !isAdmin {
			session.Where("not (project_id = ? and user_id != ?)", common.PersonalProjectID, loginUserID)
		}

		if len(filter.AppNames) > 0 {
			session.In("app_name", filter.AppNames)
		}

		if len(filter.JobSetNames) > 0 {
			session.In("job_set_name", filter.JobSetNames)
		}

		if !strutil.IsEmpty(filter.JobName) {
			session.Where("name like ?", "%"+filter.JobName+"%")
		}

		if filter.StarTime > 0 {
			session.Where("submit_time >= ?", time.Unix(int64(filter.StarTime), 0))
		}

		if filter.EndTime > 0 {
			session.Where("submit_time <= ?", time.Unix(int64(filter.EndTime), 0))
		}
	}

	return session
}
