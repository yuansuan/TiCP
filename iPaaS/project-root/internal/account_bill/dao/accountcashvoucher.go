package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type AccountCashVoucherDaoImpl struct {
}

func NewAccountCashVoucherDaoImpl() *AccountCashVoucherDaoImpl {
	return &AccountCashVoucherDaoImpl{}
}

func (c *AccountCashVoucherDaoImpl) Get(ctx context.Context, session *xorm.Session, accountVoucherID snowflake.ID, forUpdate bool) (*models.AccountCashVoucherRelation, error) {

	accountCashVoucher := models.AccountCashVoucherRelation{}

	if forUpdate {
		session = session.ForUpdate()
	}

	b, err := session.ID(accountVoucherID).Get(&accountCashVoucher)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation GetQuota db err, error: %v", err)
		return nil, common.ErrInternalServer
	}
	if !b {
		logging.GetLogger(ctx).Warnf("accountCashVoucher id not exists!, accountCashVoucher id: %v", accountVoucherID)
		return nil, common.ErrAccountVoucherIDNotFound
	}

	return &accountCashVoucher, nil
}

func (c *AccountCashVoucherDaoImpl) Add(ctx context.Context, session *xorm.Session, accountVoucher *models.AccountCashVoucherRelation) (snowflake.ID, error) {
	_, err := session.Insert(accountVoucher)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation Add db err, error: %v", err)
		return 0, common.ErrInternalServer
	}
	return accountVoucher.Id, nil
}

func (c *AccountCashVoucherDaoImpl) BatchAdd(ctx context.Context, session *xorm.Session, accountVouchers []*models.AccountCashVoucherRelation) ([]*models.AccountCashVoucherRelation, error) {
	affected, err := session.Insert(accountVouchers)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation BatchAdd db err, error: %v", err)
		return nil, common.ErrInternalServer
	}
	if affected != int64(len(accountVouchers)) {
		return nil, errors.New("insert account cash voucher error!")
	}

	return accountVouchers, nil
}

func (c *AccountCashVoucherDaoImpl) Edit(ctx context.Context, session *xorm.Session, accountVoucher *models.AccountCashVoucherRelation) (snowflake.ID, error) {
	accountVoucher.UpdateTime = time.Now()
	session = session.ID(accountVoucher.Id).Cols("status", "used_amount", "remaining_amount")

	_, err := session.Update(accountVoucher)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation Edit db err, error: %v", err)
		return 0, common.ErrInternalServer
	}

	return accountVoucher.Id, nil
}

func (c *AccountCashVoucherDaoImpl) GetTotalAmountByAccountId(ctx context.Context, session *xorm.Session, accountId snowflake.ID) (int64, error) {
	relation := &models.AccountCashVoucherRelation{
		AccountId: accountId,
	}
	sum, err := session.Where("account_id = ? and is_deleted = ? and is_expired = ? and status = ?",
		accountId, consts.NoDeleted, voucherexpiredtype.NORMAL, voucherstatus.ENABLED).
		Sum(relation, "remaining_amount")
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation GetTotalAmountByAccountId err, accountId: %v, error: %v", accountId, err)
		return 0, errors.Errorf("account_cash_voucher_relation GetTotalAmountByAccountId err, accountId: %v", accountId)
	}

	return int64(sum), nil
}

type ListRequest struct {
	AccountID snowflake.ID
	OptUserID snowflake.ID
	StartTime time.Time
	EndTime   time.Time
	PageIndex int64
	PageSize  int64
}

func (c *AccountCashVoucherDaoImpl) SelectByAccountId(ctx context.Context, session *xorm.Session, req *ListRequest) (accountCashVouchers []*models.AccountCashVoucherRelation, total int64, err error) {

	session.Where("opt_user_id = ?", req.OptUserID)

	if req.AccountID.NotZero() {
		session.Where("account_id = ?", req.AccountID)
	}

	if !req.StartTime.Equal(util.InvalidTime) {
		session.Where("create_time >= ?", req.StartTime)
	}

	if !req.EndTime.Equal(util.InvalidTime) {
		session.Where("create_time < ?", req.EndTime)
	}

	session.OrderBy("create_time DESC")
	session.Limit(int(req.PageSize), int((req.PageIndex-1)*req.PageSize))

	accountCashVouchers = []*models.AccountCashVoucherRelation{}

	total, err = session.FindAndCount(&accountCashVouchers)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation SelectByAccountId db err, error: %v", err)
		return nil, 0, common.ErrInternalServer
	}

	return accountCashVouchers, total, nil
}

func (c *AccountCashVoucherDaoImpl) UpdateStatusOfExpired(sess *xorm.Session, ids []snowflake.ID, status int) (int64, error) {
	model := &models.AccountCashVoucherRelation{
		IsExpired:  status,
		UpdateTime: time.Now(),
	}
	num, err := sess.Cols("is_expired", "update_time").Table("account_cash_voucher_relation").In("id", ids).Update(model)
	return num, err
}

func (c *AccountCashVoucherDaoImpl) QueryExpired(ctx context.Context, session *xorm.Session) (accountCashVouchers []*models.AccountCashVoucherRelation, err error) {
	session.Where("is_deleted = ?", consts.NoDeleted)
	session.And("is_expired = ?", consts.NoExpired)
	session.And("expired_time < ?", time.Now())
	err = session.Find(&accountCashVouchers)
	if err != nil {
		logging.GetLogger(ctx).Errorf("account_cash_voucher_relation SelectExpired db err, error: %v", err)
		return nil, common.ErrInternalServer
	}

	return accountCashVouchers, nil
}
