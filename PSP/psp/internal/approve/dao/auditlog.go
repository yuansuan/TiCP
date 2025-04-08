package dao

import (
	"context"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type auditLogImpl struct{}

func (d *auditLogImpl) ListAll(ctx context.Context, page *xtype.Page, userId int64, userName, ipAddress string, operateType string, startTime, endTime string, operateUserType dto.OperateUserType) (logList []*model.AuditLog, total int64, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if operateUserType > 0 {
		session.Where("operate_user_type = ?", operateUserType)
	}

	if strutil.IsNotEmpty(userName) {
		session.Where("user_name like ?", "%"+userName+"%")
	}

	if strutil.IsNotEmpty(ipAddress) {
		session.Where("`ip_address` = ?", ipAddress)
	}

	if strutil.IsNotEmpty(operateType) {
		session.Where("`operate_type` = ?", operateType)
	}

	if strutil.IsNotEmpty(startTime) {
		session.Where("operate_time > ?", startTime)
	}

	if strutil.IsNotEmpty(endTime) {
		session.Where("operate_time < ?", endTime)
	}

	if page.Index > 0 {
		session.Limit(int(page.Size), int((page.Index-1)*page.Size))
	}

	total, err = session.Desc("operate_time").FindAndCount(&logList)

	return
}

func (d *auditLogImpl) List(ctx context.Context, page *xtype.Page, userId int64, userName, ipAddress, operateType string, startTime, endTime string) (logList []*model.AuditLog, total int64, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if userId > 0 {
		session.Where("`user_id` = ?", userId)
	}

	if strutil.IsNotEmpty(userName) {
		session.Where("`user_name` = ?", userName)
	}
	if strutil.IsNotEmpty(ipAddress) {
		session.Where("`ip_address` = ?", ipAddress)
	}

	if strutil.IsNotEmpty(operateType) {
		session.Where("`operate_type` = ?", operateType)
	}

	if strutil.IsNotEmpty(startTime) {
		session.Where("operate_time > ?", startTime)
	}

	if strutil.IsNotEmpty(endTime) {
		session.Where("operate_time < ?", endTime)
	}

	if page.Index > 0 {
		session.Limit(int(page.Size), int((page.Index-1)*page.Size))
	}

	total, err = session.Desc("operate_time").FindAndCount(&logList)

	return
}

func NewLogDao() AuditLogDao {
	return &auditLogImpl{}
}

func (d *auditLogImpl) Add(ctx context.Context, log *model.AuditLog) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	_, err := session.Insert(log)
	if err != nil {
		return err
	}

	return nil
}
