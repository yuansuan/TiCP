package dao

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type jobAttrDaoImpl struct {
	sid *snowflake.Node
}

// NewJobAttrDao 创建JobAttrDao
func NewJobAttrDao() (*jobAttrDaoImpl, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	return &jobAttrDaoImpl{sid: node}, nil
}

// GetJobAttrList 获取作业属性信息
func (d *jobAttrDaoImpl) GetJobAttrList(ctx context.Context, jobID snowflake.ID) ([]*model.JobAttr, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var attrsList []*model.JobAttr
	err := session.Where("job_id=?", jobID).Find(&attrsList)
	if err != nil {
		return nil, err
	}

	return attrsList, err
}

// GetJobAttrByKey 获取指定key的作业属性信息
func (d *jobAttrDaoImpl) GetJobAttrByKey(ctx context.Context, jobID snowflake.ID, key string) (bool, *model.JobAttr, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	attr := &model.JobAttr{}
	has, err := UniqueJobAttrWhere(session, jobID, key).Get(attr)
	if err != nil {
		return false, nil, err
	}

	return has, attr, nil
}

// UpdateJobAttr 更新作业属性信息
func (d *jobAttrDaoImpl) UpdateJobAttr(ctx context.Context, attr *model.JobAttr) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := UniqueJobAttrWhere(session, attr.JobId, attr.Key).Cols("value").Update(attr)
	if err != nil {
		return err
	}

	return nil
}

// InsertJobAttr 保存作业属性信息
func (d *jobAttrDaoImpl) InsertJobAttr(ctx context.Context, attr *model.JobAttr) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Insert(attr)
	if err != nil {
		return err
	}

	return nil
}

// UniqueJobAttrWhere unique job attr where condition
func UniqueJobAttrWhere(session *xorm.Session, jobID snowflake.ID, key string) *xorm.Session {
	return session.Where("job_id=? and `key`=?", jobID.Int64(), key)
}
