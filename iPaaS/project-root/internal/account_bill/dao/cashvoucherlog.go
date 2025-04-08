package dao

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type AccountCashVoucherLogDaoImpl struct {
}

func NewAccountCashVoucherLogDaoImpl() *AccountCashVoucherLogDaoImpl {
	return &AccountCashVoucherLogDaoImpl{}
}

func (c *AccountCashVoucherLogDaoImpl) Add(ctx context.Context, session *xorm.Session, accountCashVoucherLog *models.AccountCashVoucherLog) (snowflake.ID, error) {
	_, err := session.Insert(accountCashVoucherLog)
	if err != nil {
		return 0, errors.Wrap(err, "dao:")
	}
	return accountCashVoucherLog.Id, nil
}

func (c *AccountCashVoucherLogDaoImpl) BatchAdd(ctx context.Context, session *xorm.Session, accountCashVoucherLog []*models.AccountCashVoucherLog) ([]*models.AccountCashVoucherLog, error) {
	affected, err := session.Insert(accountCashVoucherLog)
	if err != nil {
		logging.GetLogger(ctx).Errorf("accountCashVoucherLog BatchAdd db err, error: %v", err)
		return nil, common.ErrInternalServer
	}
	if affected != int64(len(accountCashVoucherLog)) {
		return nil, errors.New("insert account cash voucher log error!")
	}

	return accountCashVoucherLog, nil
}

func (c *AccountCashVoucherLogDaoImpl) UpdateAccountBillIDByID(ctx context.Context, session *xorm.Session, accountBillID, id snowflake.ID) error {
	accountCashVoucherLog := &models.AccountCashVoucherLog{
		AccountBillId: accountBillID,
	}

	session = session.ID(id).Cols("account_bill_id")
	_, err := session.Update(accountCashVoucherLog)

	if err != nil {
		logging.GetLogger(ctx).Errorf("accountCashVoucherLog udpate AccountBillID error, logID:%v, err: %v", id, err)
		return fmt.Errorf("accountCashVoucherLog udpate AccountBillID error, logID:%v", id)
	}

	return nil
}
