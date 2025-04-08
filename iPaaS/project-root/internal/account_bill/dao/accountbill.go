package dao

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accountbilltype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type AccountBillDaoImpl struct {
}

func NewAccountBillDaoImpl() *AccountBillDaoImpl {
	return &AccountBillDaoImpl{}
}

func (c *AccountBillDaoImpl) Add(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error) {
	_, err := session.Insert(accountBill)
	if err != nil {
		return consts.AmountZero, err
	}

	return accountBill.Id, nil
}

func (c *AccountBillDaoImpl) Get(ctx context.Context, session *xorm.Session, id snowflake.ID, forUpdate bool) (*models.AccountBill, error) {
	accountBill := models.AccountBill{}

	if forUpdate {
		session = session.ForUpdate()
	}

	exist, err := session.ID(id).Get(&accountBill)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, common.ErrAccountBillNotExists
	}

	return &accountBill, nil
}

func (c *AccountBillDaoImpl) GetAccountBillBySearchParams(ctx context.Context, session *xorm.Session, req *BillGetRequest) (*models.AccountBill, error) {
	modelAccountBill := &models.AccountBill{}
	if req.AccountID != 0 {
		session.Where("account_id = ?", req.AccountID)
	}

	if req.TradeType != 0 {
		session.Where("trade_type = ?", req.TradeType)
	}

	if !util.IsBlank(req.TradeID) {
		session.Where("trade_id = ?", req.TradeID)
	}

	if !util.IsBlank(req.ResourceID) {
		session.Where("resource_id = ?", req.ResourceID)
	}

	_, err := session.Get(modelAccountBill)
	if err != nil {
		return nil, err
	}

	return modelAccountBill, nil
}

func (c *AccountBillDaoImpl) Edit(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error) {
	_, err := session.ID(accountBill.Id).Update(accountBill)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// 判断是否为唯一键冲突的错误类型
			if mysqlErr.Number == 1062 {
				return consts.AmountZero, common.ErrAccountBillIdempotentIDRepeat
			}
		}

		return consts.AmountZero, err
	}

	return accountBill.Id, nil
}

func (c *AccountBillDaoImpl) EditAmount(ctx context.Context, session *xorm.Session, accountBill *models.AccountBill) (snowflake.ID, error) {
	session = session.ID(accountBill.Id).MustCols("amount", "account_balance", "freezed_amount", "normal_balance", "award_balance", "delta_voucher_balance", "account_voucher_ids")
	_, err := session.Update(accountBill)
	if err != nil {
		return consts.AmountZero, err
	}
	return accountBill.Id, nil
}

func (c *AccountBillDaoImpl) SelectAccountBills(ctx context.Context, session *xorm.Session, req *BillListRequest) (accountBillList []*models.AccountBill, total int64, err error) {
	if req.AccountID != snowflake.ID(0) {
		session.Where("account_id = ?", req.AccountID)
	}
	if !req.StartTime.Equal(util.InvalidTime) && !req.StartTime.IsZero() {
		session.Where("create_time >= ?", req.StartTime)
	}

	if !req.EndTime.Equal(util.InvalidTime) && !req.EndTime.IsZero() {
		session.Where("create_time < ?", req.EndTime)
	}

	if req.TradeType != consts.AmountZero {
		session.Where("trade_type = ?", req.TradeType)
	}

	if req.ProductName != "" {
		session.Where("product_name = ?", req.ProductName)
	}

	// 只取加、扣款；
	if req.SignType != consts.AmountZero {
		session.Where("sign = ?", req.SignType)
	} else {
		session.Where("sign = 1 OR sign = 2")
	}

	if !req.SortByAsc {
		session.OrderBy("create_time DESC")
	}

	session.Limit(int(req.PageSize), int((req.PageIndex-1)*req.PageSize))

	accountBillList = []*models.AccountBill{}
	total, err = session.FindAndCount(&accountBillList)
	if err != nil {
		return nil, consts.AmountZero, err
	}

	return accountBillList, total, nil

}

func (c *AccountBillDaoImpl) SumBillAmountByResourceID(ctx context.Context, session *xorm.Session, req *BillListRequest) (accountBillList []*models.AccountBill, total int64, err error) {
	accountBillList = []*models.AccountBill{}

	session.Table(models.AccountBill{}).
		Where("account_id = ? and resource_id != '' and trade_type = ?", req.AccountID, req.TradeType)

	if req.SignType != 0 {
		session.Where("sign = ?", req.SignType)
	}

	if !util.IsBlank(req.ProductName) {
		session.Where("product_name = ?", req.ProductName)
	}

	if !req.StartTime.IsZero() {
		session.Where("create_time > ?", req.StartTime)
	}

	total, err = session.Select(
		"resource_id, sum(amount) as amount, sum(delta_normal_balance) as delta_normal_balance, "+
			"sum(delta_award_balance) as delta_award_balance, sum(delta_voucher_balance) as delta_voucher_balance, "+
			"sum(truncate(quantity, 5)) as quantity, sum(freezed_amount) as freezed_amount, min(start_time) as start_time, "+
			"max(end_time) as end_time, max(create_time) as create_time").
		GroupBy("resource_id").
		OrderBy("create_time asc").
		Limit(int(req.PageSize), int((req.PageIndex-1)*req.PageSize)).
		FindAndCount(&accountBillList)
	if err != nil {
		return nil, 0, err
	}

	return accountBillList, total, nil
}

type BillGetRequest struct {
	AccountID    snowflake.ID
	TradeID      string
	IdempotentID string
	TradeType    accountbilltype.AccountBillTradeType
	ResourceID   string
}

type BillListRequest struct {
	AccountID   snowflake.ID
	StartTime   time.Time
	EndTime     time.Time
	TradeType   int64
	SignType    int64
	ProductName string
	ResourceID  string
	SortByAsc   bool
	PageIndex   int64
	PageSize    int64
}

type AccountListRequest struct {
	AccountID    snowflake.ID
	AccountName  string
	CustomerID   snowflake.ID
	FrozenStatus *bool
	PageIndex    int64
	PageSize     int64
}
