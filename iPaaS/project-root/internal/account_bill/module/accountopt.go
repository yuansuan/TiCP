package module

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accountbilltype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"xorm.io/xorm"
)

// OperatorParams OperatorParams
type OperatorParams struct {
	AccountID            string
	TradeID              string
	Comment              string
	OutTradeID           string
	OperatorUserID       string
	AccountCashVoucherID string
	VoucherConsumeMode   int64
	ReduceAmount         int64
	ReqPrams             string // 请求参数序列字符串
	sign                 accountbilltype.AccountBillSign
}

type AccountCashVoucherResult struct {
	AccountCashVoucherIDs     []string
	AccountCashVoucherLogsIDs []string
	DeltaVoucherAmount        int64
}

type AccountOperatorReply struct {
	*models.Account
	accountCashVoucherTotalAmount int64
}

// HandlerFoundAccountBill HandlerFoundAccountBill
type HandlerFoundAccountBill func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*models.AccountBill, error)

// HandlerUpdateAccount HandlerUpdateAccount
type HandlerUpdateAccount func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, optParams *OperatorParams) (bool, error)

// HandlerAccountCashVoucher ...
type HandlerAccountCashVoucher func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error)

// accountOperator
func accountOperator(ctx context.Context, handlerFoundAccountBill HandlerFoundAccountBill, handlerUpdateAccount HandlerUpdateAccount, handlerAccountCashVoucher HandlerAccountCashVoucher, optParams *OperatorParams) (*AccountOperatorReply, error) {
	modelAccount := &models.Account{}
	modelAccountBill := &models.AccountBill{}
	var voucherTotalAmount int64
	accountID := snowflake.MustParseString(optParams.AccountID)

	_, err := boot.MW.DefaultTransaction(ctx, func(sess *xorm.Session) (interface{}, error) {
		// get account
		accountReply, err := dao.Account.Get(ctx, sess, accountID, true)
		if err != nil {
			return nil, err
		}

		modelAccount = accountReply

		oldAccount := (*modelAccount)

		cashVoucherResult, err := handlerAccountCashVoucher(ctx, sess, optParams)
		if err != nil {
			return nil, err
		}

		// get account bill
		modelAccountBill, err = handlerFoundAccountBill(ctx, sess, optParams)
		if err != nil {
			return nil, err
		}

		updateAccountBill, err := handlerUpdateAccount(ctx, sess, modelAccount, modelAccountBill, optParams)
		if err != nil {
			return nil, err
		}

		// update AccountBill
		if updateAccountBill {
			if len(cashVoucherResult.AccountCashVoucherIDs) != 0 {
				modelAccountBill.DeltaVoucherBalance = cashVoucherResult.DeltaVoucherAmount
				//modelAccountBill.Amount += modelAccountBill.DeltaVoucherBalance
				// 更新用户代金券日志表对账单id
				for _, id := range cashVoucherResult.AccountCashVoucherLogsIDs {
					parsedId := snowflake.MustParseString(id)
					err1 := dao.AccountCashVoucherLog.UpdateAccountBillIDByID(ctx, sess, modelAccountBill.Id, parsedId)
					if err1 != nil {
						return nil, err1
					}
				}

				// 更新accountBill 代金券id
				var tmpVoucherIdString []string
				for _, id := range cashVoucherResult.AccountCashVoucherIDs {
					parsedId := int64(snowflake.MustParseString(id))
					tmpVoucherIdString = append(tmpVoucherIdString, strconv.FormatInt(parsedId, 10))
				}
				modelAccountBill.AccountVoucherIds = strings.Join(tmpVoucherIdString, ",")
			}

			// 更新账单记录
			now := time.Now()
			//modelAccountBill.CreateTime = now
			modelAccountBill.UpdateTime = now
			_, err1 := dao.AccountBill.EditAmount(ctx, sess, modelAccountBill)
			if err1 != nil {
				return nil, fmt.Errorf("account.operator err: %v", err1)
			}
		}

		// update Account
		now := time.Now()
		modelAccount.UpdateTime = now
		_, err1 := dao.Account.EditAmount(ctx, sess, modelAccount)
		if err1 != nil {
			return nil, fmt.Errorf("account.operator err: %v", err1)
		}

		// 冻结操作日志不记录
		if optParams.sign != accountbilltype.AccountBillSignFreeze {
			err = addAccountLog(ctx, sess, &oldAccount, modelAccount, optParams)
			if err != nil {
				return nil, err
			}
		}

		// get sum voucher amount
		accountCashVoucherTotalAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, modelAccount.Id)
		if err != nil {
			return nil, fmt.Errorf("query account voucher err:%v", err)
		}
		voucherTotalAmount = accountCashVoucherTotalAmount

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	accountReply := &AccountOperatorReply{
		accountCashVoucherTotalAmount: voucherTotalAmount,
	}
	accountReply.Account = modelAccount
	return accountReply, nil
}

func addAccountLog(ctx context.Context, sess *xorm.Session, old, updated *models.Account, optParams *OperatorParams) error {
	operatorUserID := snowflake.MustParseString(optParams.OperatorUserID)
	accountID := snowflake.MustParseString(optParams.AccountID)

	newID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		return fmt.Errorf("account.operator.add_account_log err: %v", err)
	}

	model := &models.AccountLog{
		Id:          newID,
		AccountId:   accountID,
		OperatorUid: operatorUserID,
		Params:      optParams.ReqPrams,
		Old:         old.String(),
		Updated:     updated.String(),
		CreateTime:  time.Now(),
	}

	_, err = dao.AccountLog.Add(ctx, sess, model)
	if err != nil {
		logging.GetLogger(ctx).Warnf("account.operator.add_account_log err: %v", err)
		return errors.New("account.operator.add_account_log err")
	}
	return nil
}
