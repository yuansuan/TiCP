package module

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accountbilltype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"xorm.io/xorm"
)

type CreateAccountBillRequest struct {
	AccountID       snowflake.ID
	TradeID         string
	IdempotentID    string
	OutTradeID      snowflake.ID
	Amount          int64
	TradeType       accountbilltype.AccountBillTradeType
	Sign            accountbilltype.AccountBillSign
	Comment         string
	MerchandiseID   string
	MerchandiseName string
	UnitPrice       int64
	PriceDes        string
	Quantity        float64
	QuantityUnit    string
	ResourceID      string
	ProductName     string
	StartTime       time.Time
	EndTime         time.Time
}

type BillListRequest struct {
	AccountID string
	StartTime time.Time
	EndTime   time.Time
	PageIndex int64
	PageSize  int64
}

func createAccountBill(ctx context.Context, sess *xorm.Session, req *CreateAccountBillRequest) (*models.AccountBill, error) {
	// account bill not exists
	logger := logging.GetLogger(ctx)
	newID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		logger.Errorf("account bill gen newID error, err: %v", err)
		return nil, errors.New("account bill gen newID error")
	}

	now := time.Now()
	modelAccountBill := &models.AccountBill{
		Id:              newID,
		AccountId:       req.AccountID,
		Sign:            int(req.Sign),
		Amount:          req.Amount,
		TradeType:       int(req.TradeType),
		TradeId:         req.TradeID,
		IdempotentId:    req.IdempotentID,
		Comment:         req.Comment,
		OutTradeId:      req.OutTradeID,
		MerchandiseId:   req.MerchandiseID,
		MerchandiseName: req.MerchandiseName,
		UnitPrice:       req.UnitPrice,
		PriceDes:        req.PriceDes,
		Quantity:        req.Quantity,
		QuantityUnit:    req.QuantityUnit,
		ResourceId:      req.ResourceID,
		ProductName:     req.ProductName,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		CreateTime:      now,
		UpdateTime:      now,
	}

	_, err = dao.AccountBill.Add(ctx, sess, modelAccountBill)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// 判断是否为唯一键冲突的错误类型
			if mysqlErr.Number == 1062 {
				return nil, common.ErrAccountBillIdempotentIDRepeat
			}
		}

		logger.Errorf("account bill insert error, err: %v", err)
		return nil, errors.New("account bill insert error")
	}

	return modelAccountBill, nil
}
