package dao

import (
	"context"
	"fmt"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type approveImpl struct{}

func (d approveImpl) GetRecord(ctx context.Context, ID snowflake.ID) (*model.ApproveRecord, error) {
	record := &model.ApproveRecord{}
	session := boot.MW.DefaultSession(ctx)
	ok, err := session.ID(ID).Get(record)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, status.Error(errcode.ErrApproveRecordNotExist, "")
	}
	return record, nil
}

func (d approveImpl) UpdateApproveRecord(ctx context.Context, record *model.ApproveRecord) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		record.UpdateTime = time.Now()
		_, err := session.ID(record.Id).Update(record)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d approveImpl) AllApproved(ctx context.Context, recordID snowflake.ID) (bool, error) {
	session := boot.MW.DefaultSession(ctx)

	count, err := session.Where("approve_record_id = ? and result != 0", recordID).Count(&model.ApproveUser{})
	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}
	return false, nil
}

func (d approveImpl) UpdateApproveUser(ctx context.Context, approveUser *model.ApproveUser) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		approveUser.UpdateTime = time.Now()
		_, err := session.ID(approveUser.Id).Update(approveUser)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d approveImpl) CheckSign(ctx context.Context, sign string) (bool, error) {
	session := boot.MW.DefaultSession(ctx)
	return session.Where("status = 1 and sign = ?", sign).Exist(&model.ApproveRecord{})
}

func (d approveImpl) AddApproveUser(ctx context.Context, approveUser *model.ApproveUser) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		approveUser.CreateTime = time.Now()
		_, err := session.Insert(approveUser)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d approveImpl) AddApproveRecord(ctx context.Context, record *model.ApproveRecord) error {
	return with.DefaultSession(ctx, func(session *xorm.Session) error {
		_, err := session.Insert(record)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewApproveDao() ApproveDao {
	return &approveImpl{}
}

func (a *approveImpl) CancelApprove(ctx context.Context, id int64) error {
	session := boot.MW.DefaultSession(ctx)

	var record model.ApproveRecord
	exist, err := session.ID(id).Get(&record)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("userid not exist")
	}

	record.Status = int8(dto.ApproveStatusCancel)
	record.UpdateTime = time.Now()
	record.ApproveTime = time.Now()
	_, err = session.ID(id).Update(record)
	if err != nil {
		return err
	}

	return nil
}

func (a *approveImpl) ApproveList(ctx context.Context, page *xtype.Page, condition *dto.ApproveListCondition) (list []*model.ApproveUserWithRecord, total int64, err error) {
	session := boot.MW.DefaultSession(ctx)

	session.Table(model.ApproveUserName).Alias("au").Join("INNER", "approve_record ar", "au.approve_record_id = ar.Id")

	if condition.ApplyId > 0 {
		session.Where("ar.apply_user_id = ?", condition.ApplyId)
	}

	if condition.RecordType > 0 {
		session.Where("ar.type = ?", condition.RecordType)
	}

	if condition.StartTime > 0 {
		session.Where("au.create_time > ?", time.Unix(condition.StartTime, 0))
	}

	if condition.EndTime > 0 {
		session.Where("au.create_time < ?", time.Unix(condition.EndTime, 0))
	}

	if condition.Status > 0 {
		session.In("ar.status", condition.Status)
	}

	if page.Index > 0 {
		session.Limit(int(page.Size), int((page.Index-1)*page.Size))
	}

	total, err = session.Desc("ar.approve_time").Desc("ar.create_time").FindAndCount(&list)

	return
}

func (a *approveImpl) ApplicationList(ctx context.Context, page *xtype.Page, condition *dto.ApplicationListCondition) (list []*model.ApproveUserWithRecord, total int64, err error) {
	session := boot.MW.DefaultSession(ctx)

	session.Table(model.ApproveUserName).Alias("au").Join("INNER", "approve_record ar", "au.approve_record_id = ar.Id")

	if condition.UserId > 0 {
		session.Where("au.approve_user_id = ?", condition.UserId)
	}

	if len(condition.ApplyName) > 0 {
		session.Where("ar.apply_user_name like ?", "%"+condition.ApplyName+"%")
	}

	if condition.RecordType > 0 {
		session.Where("ar.type = ?", condition.RecordType)
	}

	if condition.StartTime > 0 {
		session.Where("au.create_time > ?", time.Unix(condition.StartTime, 0))
	}

	if condition.EndTime > 0 {
		session.Where("au.create_time < ?", time.Unix(condition.EndTime, 0))
	}

	if len(condition.Status) != 0 {
		session.In("ar.status", condition.Status)
	}

	if page.Index > 0 {
		session.Limit(int(page.Size), int((page.Index-1)*page.Size))
	}
	total, err = session.Desc("ar.approve_time").Desc("ar.create_time").FindAndCount(&list)

	return
}

func (a approveImpl) CheckUnhandledApprove(ctx context.Context, userID int64) (bool, error) {
	session := boot.MW.DefaultSession(ctx)

	return session.Table(model.ApproveRecordName).Alias("ar").Join("LEFT OUTER", "approve_user au", "ar.status = 1 and au.approve_record_id = ar.Id").Where(
		"(ar.apply_user_id = ? or au.approve_user_id = ?) and status = 1", userID, userID).Exist(&model.ApproveRecord{})
}
