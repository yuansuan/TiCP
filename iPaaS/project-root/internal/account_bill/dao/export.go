package dao

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"xorm.io/xorm"
)

var (
	Account                    = NewAccountDaoImpl()
	AccountBill                = NewAccountBillDaoImpl()
	AccountLog                 = NewAccountLogDaoImpl()
	CashVoucher                = NewCashVoucherDaoImpl()
	AccountCashVoucherRelation = NewAccountCashVoucherDaoImpl()
	AccountCashVoucherLog      = NewAccountCashVoucherLogDaoImpl()
)

// AccountDao AccountDao interface
type AccountDao interface {
	Get(ctx context.Context, session *xorm.Session, id snowflake.ID, forUpdate bool) (*models.Account, error)

	GetByCustomerID(ctx context.Context, session *xorm.Session, customerID snowflake.ID) (*models.Account, error)

	Add(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error)

	EditAmount(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error)

	EditFrozenStatus(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error)

	EditCreditQuota(ctx context.Context, session *xorm.Session, account *models.Account) (snowflake.ID, error)

	ExistSameName(ctx context.Context, session *xorm.Session, accountName string) (bool, error)

	ExistSameCustomerID(ctx context.Context, session *xorm.Session, customerID snowflake.ID) (bool, error)

	SelectAccounts(ctx context.Context, session *xorm.Session, req *AccountListRequest) (accountList []*models.Account, total int64, err error)
}

// AccountBillDao AccountBillDao interface
type AccountBillDao interface {
	Get(ctx context.Context, session *xorm.Session, id snowflake.ID, forUpdate bool) (*models.AccountBill, error)

	GetAccountBillBySearchParams(ctx context.Context, session *xorm.Session, req *BillGetRequest) (*models.AccountBill, error)

	Add(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error)

	Edit(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error)

	EditAmount(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error)

	SelectAccountBills(ctx context.Context, session *xorm.Session, req *BillListRequest) (accountBillList []*models.AccountBill, total int64, err error)

	SumBillAmountByResourceID(ctx context.Context, session *xorm.Session, req *BillListRequest) ([]*models.AccountBill, int64, error)
}

// AccountLogDao AccountLogDao interface
type AccountLogDao interface {
	Add(ctx context.Context, session *xorm.Session, accountLog *models.AccountLog) (snowflake.ID, error)
}

// CashVoucherDao cashVoucher dao interface
type CashVoucherDao interface {
	Get(ctx context.Context, session *xorm.Session, cashVoucherID snowflake.ID, forUpdate bool) (*models.CashVoucher, error)

	Add(ctx context.Context, session *xorm.Session, cashVoucher *models.CashVoucher) (snowflake.ID, error)

	Edit(ctx context.Context, session *xorm.Session, cashVoucher *models.CashVoucher) (snowflake.ID, error)

	SelectCashVouchers(ctx context.Context, session *xorm.Session, in *models.SelectCashVouchersReq) (int64, []*models.CashVoucher, error)
}

// AccountCashVoucherRelationDao ...
type AccountCashVoucherRelationDao interface {
	Get(ctx context.Context, session *xorm.Session, accountVoucherID snowflake.ID, forUpdate bool) (*models.AccountCashVoucherRelation, error)

	Add(ctx context.Context, session *xorm.Session, accountVoucher *models.AccountCashVoucherRelation) (snowflake.ID, error)

	BatchAdd(ctx context.Context, session *xorm.Session, accountVouchers []*models.AccountCashVoucherRelation) ([]*models.AccountCashVoucherRelation, error)

	Edit(ctx context.Context, session *xorm.Session, accountVoucher *models.AccountCashVoucherRelation) (snowflake.ID, error)

	GetTotalAmountByAccountId(ctx context.Context, session *xorm.Session, accountId snowflake.ID) (int64, error)

	SelectByAccountId(ctx context.Context, session *xorm.Session, req *ListRequest) (accountCashVouchers []*models.AccountCashVoucherRelation, total int64, err error)
}

// AccountCashVoucherLogDao ...
type AccountCashVoucherLogDao interface {
	Add(ctx context.Context, session *xorm.Session, cashVoucherLog *models.AccountCashVoucherLog) (snowflake.ID, error)

	UpdateAccountBillIDByID(ctx context.Context, session *xorm.Session, accountBillID, id snowflake.ID) error
}
