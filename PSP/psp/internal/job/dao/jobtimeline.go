package dao

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type jobTimelineDaoImpl struct {
	sid *snowflake.Node
}

// NewJobTimelineDao 创建JobTimelineDao
func NewJobTimelineDao() (*jobTimelineDaoImpl, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	return &jobTimelineDaoImpl{sid: node}, nil
}

// GetJobTimeline 获取指定作业的时间线信息
func (d *jobTimelineDaoImpl) GetJobTimeline(ctx context.Context, jobID snowflake.ID) ([]*model.JobTimeline, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var timeline []*model.JobTimeline
	if err := session.Where("job_id=?", jobID).Asc("event_time").Find(&timeline); err != nil {
		return nil, err
	}

	return timeline, nil
}

// InsertJobTimeline 保存作业时间线信息
func (d *jobTimelineDaoImpl) InsertJobTimeline(ctx context.Context, timeline *model.JobTimeline) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Insert(timeline)
	if err != nil {
		return err
	}

	return nil
}

// GetJobTimelineByName 获取指定名称的作业时间线信息
func (d *jobTimelineDaoImpl) GetJobTimelineByName(ctx context.Context, jobID snowflake.ID, eventName string) (bool, *model.JobTimeline, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	timeline := &model.JobTimeline{}
	has, err := UniqueJobTimelineWhere(session, jobID, eventName).Get(timeline)
	if err != nil {
		return false, nil, err
	}

	return has, timeline, nil
}

// UniqueJobTimelineWhere unique job timeline where condition
func UniqueJobTimelineWhere(session *xorm.Session, jobID snowflake.ID, eventName string) *xorm.Session {
	return session.Where("job_id=? and event_name=?", jobID, eventName)
}
