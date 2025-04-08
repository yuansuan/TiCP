package dao

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type AccountDaoImpl struct {
}

func NewAccountDaoImpl() *AccountDaoImpl {
	return &AccountDaoImpl{}
}

func (c *AccountDaoImpl) Add(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error) {
	_, err := session.Insert(account)
	if err != nil {
		return consts.AmountZero, err
	}
	return account.Id, nil
}

func (c *AccountDaoImpl) Get(ctx context.Context, session *xorm.Session, id snowflake.ID, forUpdate bool) (*models.Account, error) {
	account := models.Account{}
	if forUpdate {
		session = session.ForUpdate()
	}

	exist, err := session.ID(id).Get(&account)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.Wrap(common.ErrAccountNotExists, "account get dao:")
	}

	return &account, nil
}

func (c *AccountDaoImpl) GetByCustomerID(ctx context.Context, session *xorm.Session, customerID snowflake.ID) (*models.Account, error) {
	accountModel := &models.Account{}
	exist, err := session.Where("customer_id = ?", customerID).OrderBy("create_time DESC").Limit(1).Get(accountModel)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.Wrap(common.ErrAccountNotExists, "account get by customerID dao:")
	}

	return accountModel, nil
}

func (c *AccountDaoImpl) EditAmount(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error) {
	session = session.ID(account.Id).MustCols("account_balance", "freezed_amount", "normal_balance", "award_balance", "update_time")
	_, err := session.Update(account)
	if err != nil {
		return consts.AmountZero, err
	}

	return account.Id, nil
}

func (c *AccountDaoImpl) EditFrozenStatus(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error) {
	_, err := session.ID(account.Id).MustCols("is_freeze", "update_time").Update(account)
	if err != nil {
		return consts.AmountZero, err
	}

	return account.Id, nil
}

func (c *AccountDaoImpl) EditCreditQuota(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error) {
	_, err := session.ID(account.Id).MustCols("credit_quota", "update_time").Update(account)
	if err != nil {
		return consts.AmountZero, err
	}

	return account.Id, nil
}

func (c *AccountDaoImpl) ExistSameName(ctx context.Context, session *xorm.Session, accountName string) (bool, error) {
	exist, err := session.Exist(&models.Account{
		Name: accountName,
	})

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (c *AccountDaoImpl) ExistSameCustomerID(ctx context.Context, session *xorm.Session, customerID snowflake.ID) (bool, error) {
	exist, err := session.Exist(&models.Account{
		CustomerId: customerID,
	})

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (c *AccountDaoImpl) SelectAccounts(ctx context.Context, session *xorm.Session, req *AccountListRequest) (accountList []*models.Account, total int64, err error) {
	if req.AccountID != snowflake.ID(0) {
		session.Where("id = ?", req.AccountID)
	}

	if req.AccountName != "" {
		session.Where("name like ?", "%"+req.AccountName+"%")
	}

	if req.CustomerID != snowflake.ID(0) {
		session.Where("customer_id = ?", req.CustomerID)
	}
	if req.FrozenStatus != nil {
		session.Where("is_freeze = ?", req.FrozenStatus)
	}

	session.Limit(int(req.PageSize), int((req.PageIndex-1)*req.PageSize))

	accountList = []*models.Account{}
	total, err = session.FindAndCount(&accountList)
	if err != nil {
		return nil, consts.AmountZero, err
	}

	return accountList, total, nil

}
