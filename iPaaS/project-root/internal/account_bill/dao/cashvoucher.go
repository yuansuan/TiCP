package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type CashVoucherDaoImpl struct {
}

func NewCashVoucherDaoImpl() *CashVoucherDaoImpl {
	return &CashVoucherDaoImpl{}
}

func (c *CashVoucherDaoImpl) Get(ctx context.Context, session *xorm.Session, cashVoucherID snowflake.ID, forUpdate bool) (*models.CashVoucher, error) {
	cashVoucher := models.CashVoucher{}

	if forUpdate {
		session = session.ForUpdate()
	}

	b, err := session.ID(cashVoucherID).Get(&cashVoucher)
	if err != nil {
		return nil, err
	}
	if !b {
		logging.GetLogger(ctx).Warnf("cashVoucher id not exists!, cashVoucher id: %v", cashVoucherID)
		return nil, common.ErrCashVoucherNotExists
	}

	return &cashVoucher, nil
}

func (c *CashVoucherDaoImpl) Add(ctx context.Context, session *xorm.Session, cashVoucher *models.CashVoucher) (snowflake.ID, error) {
	_, err := session.Insert(cashVoucher)
	if err != nil {
		return consts.AmountZero, errors.Wrap(err, "dao:")
	}
	return cashVoucher.Id, nil
}

func (c *CashVoucherDaoImpl) Edit(ctx context.Context, session *xorm.Session, cashVoucher *models.CashVoucher) (snowflake.ID, error) {
	session = session.ID(cashVoucher.Id).Cols("availability_status", "update_time")

	_, err := session.Update(cashVoucher)
	if err != nil {
		return consts.AmountZero, err
	}

	return cashVoucher.Id, nil
}

func (c *CashVoucherDaoImpl) SelectCashVouchers(ctx context.Context, session *xorm.Session, in *models.SelectCashVouchersReq) (int64, []*models.CashVoucher, error) {
	var cashVouchers []*models.CashVoucher

	session.Where("is_deleted = ?", consts.NoDeleted).Where("opt_user_id = ?", int64(in.OptUserId))

	if in.AvailabilityStatus != "" {
		b, status := voucherstatus.ValidAvailabilityStatusString(in.AvailabilityStatus)
		if b {
			session.Where("availability_status = ?", int64(status))
		}
	}

	if in.Id != 0 {
		session.Where("id = ?", in.Id)
	}

	if in.Name != "" {
		session.Where("name = ?", in.Name)
	}

	session.Where("is_expired = ?", in.IsExpired)
	if !in.StartTime.Equal(util.InvalidTime) {
		session.Where("create_time >= ?", in.StartTime)
	}

	if !in.EndTime.Equal(util.InvalidTime) {
		session.Where("create_time <= ?", in.EndTime)
	}

	session.OrderBy("create_time DESC")
	session.Limit(int(in.Size), int((in.Index-1)*in.Size))

	total, err := session.FindAndCount(&cashVouchers)
	if err != nil {
		return consts.AmountZero, nil, err
	}

	return total, cashVouchers, nil
}

func (c *CashVoucherDaoImpl) UpdateStatusOfExpired(ctx context.Context, session *xorm.Session) error {
	cashVoucher := &models.CashVoucher{IsExpired: consts.Expired}
	session.Table("cash_voucher")
	session.Where("is_deleted = ?", consts.NoDeleted)
	session.And("is_expired = ?", consts.NoExpired)
	session.And("expired_type = ?", consts.AbsExpired)
	session.And("abs_expired_time < ?", time.Now())
	_, err := session.Update(cashVoucher)
	if err != nil {
		logging.GetLogger(ctx).Error("err_cash_voucher_update_db", "error", err, "update_data", cashVoucher)
		return err
	}

	return nil
}
