package module

import (
	"context"
	"fmt"
	"strings"

	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountlist"

	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/amountrefund"
	userBillList "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/billlist"
	userIdGet "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/idget"
	userResourceBillList "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/resourcebilllist"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/billlist"
	accountCreate "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/create"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditadd"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditquotamodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/frozenmodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idget"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idreduce"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentfreezeunfreeze"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentreduce"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidget"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidreduce"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accountbilltype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/accounttype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"

	"time"

	"github.com/pkg/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"google.golang.org/grpc/status"
)

type ReduceRequest struct {
	AccountID             string
	Amount                int64
	TradeID               string
	IdempotentID          string
	Comment               string
	MerchandiseID         string
	MerchandiseName       string
	UnitPrice             int64
	PriceDes              string
	Quantity              float64
	QuantityUnit          string
	ResourceID            string
	ProductName           string
	StartTime             time.Time
	EndTime               time.Time
	OptUserID             string
	AccountCashVoucherIDs string
	VoucherConsumeMode    int64
}

// Create ...
func Create(ctx context.Context, req *accountCreate.Request, optUserID string) (*accountCreate.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.Create params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	var accountName string
	var customerID snowflake.ID
	accountType := accounttype.AccountType(req.AccountType)

	// 验证企业账号名称重复
	if accountType == accounttype.COMPANY {
		// 验证名字是否有重名
		exist, err := dao.Account.ExistSameName(ctx, sess, req.AccountName)
		if err != nil {
			logger.Errorf("account.query_db error, err: %v", err)
			return nil, errors.New("query.account_db error!")
		}

		// 企业名称重复存在
		if exist {
			logger.Infof("account.create_company_account error, name exists: %s", req.AccountName)
			return nil, common.ErrAccountExists
		}

		accountName = req.AccountName
		if req.UserID != "" {
			customerID, err = snowflake.ParseString(req.UserID)
			if err != nil {
				logger.Errorf("parse company id error! company ID: %s, err: %v", req.UserID, err)
			}

			exist, err := dao.Account.ExistSameCustomerID(ctx, sess, customerID)
			if err != nil {
				logger.Errorf("account.query_account_db error, err: %v", err)
				return nil, errors.New("account.query_account_db error!")
			}

			// userID重复存在
			if exist {
				logger.Infof("account.create_user_account error, id exists error, id: %s", req.UserID)
				return nil, common.ErrAccountExists
			}
		}
	}

	// 验证用户Id 是否存在
	// sso 验证
	if accountType == accounttype.PERSONAL {
		userInfo, err := rpc.GetInstance().GetSSOUserByID(ctx, req.UserID)
		if err != nil {
			if status.Code(err) == consts.ErrHydraLcpDBUserNotExist {
				logger.Infof("user id not exists, user id: %s", req.UserID)
				return nil, common.ErrUserNotExists
			} else {
				logger.Error("account.query_sso_user_info error, err: %v", err)
				return nil, errors.New("account.query_sso_user_info error!")
			}
		}

		// 账号体系内唯一验证
		userId := snowflake.MustParseString(req.UserID)
		exist, err := dao.Account.ExistSameCustomerID(ctx, sess, userId)
		if err != nil {
			logger.Errorf("account.query_account_db error, err: %v", err)
			return nil, errors.New("account.query_account_db error!")
		}

		// userID重复存在
		if exist {
			logger.Infof("account.create_user_account error, id exists error, id: %s", req.UserID)
			return nil, common.ErrAccountExists
		}
		accountName = userInfo.Phone
		customerID = userId
	}

	newID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		logger.Errorf("gen snowflake error, err: %v", err)
		return nil, errors.New("account.create_account.gen_ID error!")
	}

	// 保存账户信息
	now := time.Now()
	newAccount := &models.Account{
		Id:              newID,
		CustomerId:      customerID,
		Name:            accountName,
		Currency:        "CNY",
		Status:          1,
		AccountBalance:  0,
		FreezedAmount:   0,
		NormalBalance:   0,
		AwardBalance:    0,
		WithdrawEnabled: 0,
		CreditQuota:     0,
		AccountType:     int32(accountType),
		CreateTime:      now,
		UpdateTime:      now,
	}

	_, err = dao.Account.Add(ctx, sess, newAccount)
	if err != nil {
		logger.Errorf("account.create.save_account error, err: %v", err)
		return nil, errors.New("account.create.insert_db_account error!")
	}

	return &accountCreate.Data{
		AccountID: newAccount.Id.String(),
	}, nil

}

// CreditAdd ...
func CreditAdd(ctx context.Context, req *creditadd.Request, optUserID string) (*creditadd.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.CreditAdd params: %+v, optUserID: %s", req, optUserID)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradeCredit
	sign := accountbilltype.AccountBillSignAdd

	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		addAmount := req.DeltaNormalBalance + req.DeltaAwardBalance
		idempotentID, err := defaultIdempotentID(ctx, req.IdempotentID)
		if err != nil {
			return nil, err
		}

		createBillReq := &CreateAccountBillRequest{
			AccountID:    accountID,
			TradeID:      req.TradeId,
			IdempotentID: idempotentID,
			OutTradeID:   consts.AmountZero, // 外部订单和内部订单
			Amount:       addAmount,
			TradeType:    tradeType,
			Sign:         sign,
			Comment:      req.Comment,
		}
		modelAccountBill, err := createAccountBill(ctx, sess, createBillReq)
		if err != nil {
			return nil, err
		}

		return modelAccountBill, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		modelAccount.Add(req.DeltaNormalBalance, req.DeltaAwardBalance)
		modelAccountBill.DeltaNormalBalance = req.DeltaNormalBalance
		modelAccountBill.DeltaAwardBalance = req.DeltaAwardBalance
		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount
		return true /*update account bill*/, nil
	}

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		return &AccountCashVoucherResult{}, nil
	}

	reqParamsStr := fmt.Sprintf("CreditAdd requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: optUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		sign:           sign,
	}

	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		// 触发幂等性错误， 返回账户详情字段
		if errors.Is(err, common.ErrAccountBillIdempotentIDRepeat) {
			logger.Warnf("account.account_credit_add.idempotentid_repeat, req params: %+v", req)
			account, err := getAccountReplyIfIdempotentIDRepeat(ctx, accountID)
			if err != nil {
				return nil, err
			}

			return &creditadd.Data{
				Account: &v20230530.Account{
					AccountID:         account.AccountID,
					AccountName:       account.AccountName,
					AccountBalance:    account.AccountBalance,
					NormalBalance:     account.NormalBalance,
					AwardBalance:      account.AwardBalance,
					FreezedAmount:     account.FreezedAmount,
					CreditQuotaAmount: account.CreditQuotaAmount,
					CashVoucherAmount: account.CashVoucherAmount,
				},
			}, nil
		}

		logger.Errorf("accountOperator error, err: %v", err)
		return nil, err
	}

	return &creditadd.Data{
		Account: &v20230530.Account{
			AccountID:         accountReply.Id.String(),
			AccountName:       accountReply.Name,
			AccountBalance:    accountReply.AccountBalance,
			NormalBalance:     accountReply.NormalBalance,
			AwardBalance:      accountReply.AwardBalance,
			FreezedAmount:     accountReply.FreezedAmount,
			CashVoucherAmount: accountReply.accountCashVoucherTotalAmount,
		},
	}, nil
}

// PaymentReduce ...
func PaymentReduce(ctx context.Context, req *paymentreduce.Request, OptUserID string) (*paymentreduce.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.paymentReduce params: %+v, optUserID: %s", req, OptUserID)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradePay
	sign := accountbilltype.AccountBillSignReduce
	startTime, endTime := formatReqTime(req.StartTime, req.EndTime)

	lastFreezeAmount := int64(0)
	reduceAmount := int64(0)
	idempotentID, err := defaultIdempotentID(ctx, "")
	if err != nil {
		return nil, err
	}

	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		modelAccountBill, err := dao.AccountBill.GetAccountBillBySearchParams(ctx, sess, &dao.BillGetRequest{
			AccountID: accountID,
			TradeID:   req.TradeID,
			TradeType: tradeType,
		})

		if err != nil {
			logger.Errorf("account.payment_reduce.found_account_bill err: %v", err)
			return nil, errors.New("account.payment_reduce.found_account_bill.db error")
		}

		if modelAccountBill.Id > 0 {
			// account bill exists
			if accountbilltype.AccountBillSign(modelAccountBill.Sign) != accountbilltype.AccountBillSignFreeze {
				logger.Warnf("account.payment_reduce.account_bill can not change sign %v to %v", accountbilltype.AccountBillSign(modelAccountBill.Sign), sign)
				return nil, errors.WithMessage(common.ErrAccountBillSignStatusInvalid, "account.payment_reduce.account_bill sign is not freezed status")
			}

			lastFreezeAmount = modelAccountBill.Amount
			modelAccountBill.Comment = req.Comment
			modelAccountBill.UpdateTime = time.Now()
			modelAccountBill.Sign = int(sign)
			modelAccountBill.MerchandiseId = req.MerchandiseID
			modelAccountBill.MerchandiseName = req.MerchandiseName
			modelAccountBill.UnitPrice = req.UnitPrice
			modelAccountBill.PriceDes = util.DefaultString(req.PriceDes, consts.UNIT_PRICE_DEFAULT)
			modelAccountBill.Quantity = req.Quantity
			modelAccountBill.QuantityUnit = req.QuantityUnit
			modelAccountBill.ResourceId = req.ResourceID
			modelAccountBill.StartTime = *startTime
			modelAccountBill.EndTime = *endTime
			modelAccountBill.IdempotentId = idempotentID
			_, err1 := dao.AccountBill.Edit(ctx, sess, modelAccountBill)
			if err1 != nil {
				logger.Errorf("account.payment_reduce.update_account_bill err: %v", err1)
				return nil, errors.New("account.payment_reduce.update_account_bill error")
			}
		} else {
			logger.Warnf("account.payment_reduce.found_account_bill account_bill not exists! accountID: %v, tradeID:%s", accountID, req.TradeID)
			return nil, errors.New("account.payment_reduce.found_account_bill account_bill not exists")
		}
		return modelAccountBill, nil
	}

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		modelAccountBill, err := dao.AccountBill.GetAccountBillBySearchParams(ctx, sess, &dao.BillGetRequest{
			AccountID: accountID,
			TradeID:   req.TradeID,
			TradeType: tradeType,
		})

		if err != nil {
			logger.Errorf("account.payment_reduce.found_account_bill err: %v", err)
			return nil, errors.New("account.payment_reduce.found_account_bill.db error")
		}

		if modelAccountBill.Id <= 0 {
			logger.Warnf("account.payment_reduce.found_account_bill account_bill not exists! accountID: %v, tradeID:%s", accountID, req.TradeID)
			return nil, common.ErrAccountBillNotExists
		}

		reduceAmount = modelAccountBill.Amount

		if util.IsBlank(req.AccountCashVoucherIDs) {
			return &AccountCashVoucherResult{}, nil
		}

		// check account cash voucher valid
		trimIDs := strings.TrimRight(req.AccountCashVoucherIDs, ",")
		voucherIDs := strings.Split(trimIDs, ",")
		if len(voucherIDs) > consts.CASH_VOUCHER_EXCEED_AMOUNT {
			logger.Warnf("account.payment_reduce.cash_voucher_exceed_amount.error!")
			return nil, common.ErrAccountVoucherExceed
		}

		voucherRelations := make([]*models.AccountCashVoucherRelation, 0, len(voucherIDs))
		for _, id := range voucherIDs {
			voucherId, err2 := snowflake.ParseString(id)
			if err2 != nil {
				logger.Errorf("account.payment_reduce.parse_voucher_id.error! accountCashVoucherId:%s, err:%v", id, err2)
				return nil, fmt.Errorf("account.payment_reduce.parse_voucher_id.error! accountCashVoucherId:%s", id)
			}

			voucherRelation, err2 := dao.AccountCashVoucherRelation.Get(ctx, sess, voucherId, true)
			if err2 != nil {
				if errors.Is(err2, common.ErrAccountVoucherIDNotFound) {
					return nil, errors.WithMessagef(common.ErrAccountVoucherIDNotFound, "voucher not exists, id:%s", id)
				}

				logger.Errorf("account.payment_reduce.query_voucher.error, accountCashVoucherId: %s, err:%v", id, err2)
				return nil, err2
			}

			if voucherexpiredtype.IsExpiredType(voucherRelation.IsExpired) == voucherexpiredtype.EXPIRED {
				logger.Warnf("account.payment_reduce.voucher_expired.error! accountCashVoucherId: %v", id)
				return nil, errors.WithMessagef(common.ErrAccountVoucherExpired, "account cash voucher expired, id: %v", id)
			}

			if voucherstatus.AccountCashVoucherStatus(voucherRelation.Status) == voucherstatus.DISABLED {
				logger.Warnf("account.payment_reduce.voucher_disabled.error! accountCashVoucherId: %v", id)
				return nil, errors.WithMessagef(common.ErrAccountVoucherDisabled, "account voucher disabled, id: %v", id)
			}

			voucherRelations = append(voucherRelations, voucherRelation)
		}

		// 扣减金额
		var deltaVoucherAmount int64
		var usedAccountCashVoucherIDs []string
		var accountCashVoucherLogsIDs []string

		for _, voucherRelation := range voucherRelations {
			historyVoucherRelation := voucherRelation.String()

			if reduceAmount <= 0 {
				break
			}

			// 剩余金额
			remainingAmount := voucherRelation.RemainingAmount
			// 代金券剩余金额为0， 则返回空对象默认值
			if remainingAmount <= 0 {
				continue
			}

			// 代金券还有剩余金额可以抵扣，即 remainingAmount > 0
			if remainingAmount >= reduceAmount {
				deltaVoucherAmount += reduceAmount
				voucherRelation.RemainingAmount = voucherRelation.RemainingAmount - reduceAmount
				voucherRelation.UsedAmount = voucherRelation.UsedAmount + reduceAmount
				reduceAmount = 0
				remainingAmount = voucherRelation.RemainingAmount
			} else {
				deltaVoucherAmount += voucherRelation.RemainingAmount
				reduceAmount = reduceAmount - voucherRelation.RemainingAmount
				voucherRelation.UsedAmount = voucherRelation.CashVoucherAmount
				voucherRelation.RemainingAmount = 0
				remainingAmount = 0
			}

			logger.Infof("current voucher info: %s", util.ToJsonString(ctx, voucherRelation))

			// 更新账户代金券表，保存代金券日志
			_, err2 := dao.AccountCashVoucherRelation.Edit(ctx, sess, voucherRelation)
			if err2 != nil {
				logger.Errorf("account.payment_reduce.update.accountCashVoucher.error! err:%v", err)
				return nil, errors.New("account.payment_reduce.update.accountCashVoucher.error")
			}

			// 增加代金券到数组
			usedAccountCashVoucherIDs = append(usedAccountCashVoucherIDs, voucherRelation.Id.String())

			newID, err2 := rpc.GetInstance().GenID(ctx)
			if err2 != nil {
				logger.Errorf("account.payment_reduce.genID.error! err:%v", err2)
				return nil, errors.New("account.payment_reduce.genID.error")
			}
			now := time.Now()
			cashVoucherLog := &models.AccountCashVoucherLog{
				Id:                   newID,
				AccountId:            voucherRelation.AccountId,
				CashVoucherId:        voucherRelation.CashVoucherId,
				AccountCashVoucherId: voucherRelation.Id,
				SignType:             consts.VOUCHER_LOG_SIGN_CONSUME,
				SourceInfo:           historyVoucherRelation,
				TargetInfo:           voucherRelation.String(),
				Comment:              req.Comment,
				OptUserId:            snowflake.MustParseString(OptUserID),
				CreateTime:           now,
				UpdateTime:           now,
			}
			_, err2 = dao.AccountCashVoucherLog.Add(ctx, sess, cashVoucherLog)
			if err2 != nil {
				logger.Errorf("account.payment_reduce.add.cashVoucherLog.error! err:%v", err2)
				return nil, errors.New("account.payment_reduce.add.cashVoucherLog.error!")
			}

			// 保存日志id
			accountCashVoucherLogsIDs = append(accountCashVoucherLogsIDs, newID.String())
		}

		return &AccountCashVoucherResult{
			AccountCashVoucherIDs:     usedAccountCashVoucherIDs,
			AccountCashVoucherLogsIDs: accountCashVoucherLogsIDs,
			DeltaVoucherAmount:        deltaVoucherAmount,
		}, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		if lastFreezeAmount != 0 {
			modelAccount.Unfreeze(lastFreezeAmount)
		}

		normal, award := modelAccount.Reduce(reduceAmount)
		modelAccountBill.DeltaNormalBalance = normal
		modelAccountBill.DeltaAwardBalance = award
		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount
		return true /*update account bill*/, nil
	}

	reqParamsStr := fmt.Sprintf("Payment reduce requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: OptUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		sign:           sign,
	}

	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		logger.Errorf("accountOperator error, err: %v", err)
		return nil, err
	}

	return &paymentreduce.Data{
		AccountID:         accountReply.Id.String(),
		AccountName:       accountReply.Name,
		AccountBalance:    accountReply.AccountBalance,
		NormalBalance:     accountReply.NormalBalance,
		AwardBalance:      accountReply.AwardBalance,
		FreezedAmount:     accountReply.FreezedAmount,
		CashVoucherAmount: accountReply.accountCashVoucherTotalAmount,
		CreditQuotaAmount: accountReply.CreditQuota,
	}, nil
}

func formatReqTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	startTime := &util.InvalidTime
	endTime := &util.InvalidTime
	if util.IsNotBlank(startTimeStr) {
		toTime, _ := util.StringToTime(startTimeStr)
		startTime = toTime
	}

	if util.IsNotBlank(endTimeStr) {
		toTime, _ := util.StringToTime(endTimeStr)
		endTime = toTime
	}
	return startTime, endTime
}

// AccountList ...
func AccountList(ctx context.Context, req *accountlist.Request, optUserID string) (*accountlist.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.BillList params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	// check account exists
	var accountID snowflake.ID
	var customerID snowflake.ID
	if req.AccountID != "" {
		accountID = snowflake.MustParseString(req.AccountID)
		_, err := IsAccountExists(ctx, accountID, sess)
		if err != nil {
			return nil, err
		}
	}

	if req.CustomerID != "" {
		customerID = snowflake.MustParseString(req.CustomerID)
	}

	accountReq := &dao.AccountListRequest{
		AccountID:    accountID,
		AccountName:  req.AccountName,
		CustomerID:   customerID,
		FrozenStatus: req.FrozenStatus,
		PageSize:     req.PageSize,
		PageIndex:    req.PageIndex,
	}

	accountList, total, err := dao.Account.SelectAccounts(ctx, sess, accountReq)
	if err != nil {
		logger.Errorf("account.accountlist.dao.error, err: %v", err)
		return &accountlist.Data{
			Total:    consts.AmountZero,
			Accounts: nil,
		}, errors.WithMessagef(common.ErrInternalServer, "account.bill_list.dao.error")
	}

	accountListResp := make([]*v20230530.AccountDetail, 0, len(accountList))
	for _, account := range accountList {
		cashVoucherAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, accountID)
		if err != nil {
			return nil, err
		}
		accountData := util.ModelToOpenApiAccount(account, cashVoucherAmount)
		accountListResp = append(accountListResp, accountData)
	}

	return &accountlist.Data{
		Total:    total,
		Accounts: accountListResp,
	}, nil
}

// BillList ...
func BillList(ctx context.Context, req *billlist.Request, optUserID string) (*billlist.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.BillList params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	// check account exists
	var accountID snowflake.ID
	if req.AccountID != "" {
		accountID = snowflake.MustParseString(req.AccountID)
		_, err := IsAccountExists(ctx, accountID, sess)
		if err != nil {
			return nil, err
		}
	}

	billReq := &dao.BillListRequest{
		AccountID:   accountID,
		StartTime:   util.InvalidTime,
		EndTime:     util.InvalidTime,
		TradeType:   req.TradeType,
		SignType:    req.SignType,
		ProductName: req.ProductName,
		SortByAsc:   req.SortByAsc,
		PageSize:    req.PageSize,
		PageIndex:   req.PageIndex,
	}

	if req.StartTime != "" {
		startTime, _ := util.StringToTime(req.StartTime)
		billReq.StartTime = *startTime
	}

	if req.EndTime != "" {
		endTIme, _ := util.StringToTime(req.EndTime)
		billReq.EndTime = *endTIme
	}

	billList, total, err := dao.AccountBill.SelectAccountBills(ctx, sess, billReq)
	if err != nil {
		logger.Errorf("account.bill_list.dao.error, err: %v", err)
		return &billlist.Data{
			Total:        consts.AmountZero,
			AccountBills: nil,
		}, errors.WithMessagef(common.ErrInternalServer, "account.bill_list.dao.error")
	}

	billListResp := make([]*v20230530.BillListData, 0, len(billList))
	for _, bill := range billList {
		billData := util.ModelToOpenApiBill(bill)
		billListResp = append(billListResp, billData)
	}

	return &billlist.Data{
		Total:        total,
		AccountBills: billListResp,
	}, nil
}

func UserBillList(ctx context.Context, req *userBillList.Request, userID string) (*userBillList.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.UserBillList params: %+v, optUserID: %s", req, userID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	// 获取账号信息
	account, err := getAccountInfoByUserID(ctx, userID, "account.UserBillList")
	if err != nil {
		return nil, err
	}

	billReq := &dao.BillListRequest{
		AccountID:   account.Id,
		StartTime:   util.InvalidTime,
		EndTime:     util.InvalidTime,
		TradeType:   req.TradeType,
		SignType:    req.SignType,
		ProductName: req.ProductName,
		SortByAsc:   req.SortByAsc,
		PageSize:    req.PageSize,
		PageIndex:   req.PageIndex,
	}

	if req.StartTime != "" {
		startTime, _ := util.StringToTime(req.StartTime)
		billReq.StartTime = *startTime
	}

	if req.EndTime != "" {
		endTIme, _ := util.StringToTime(req.EndTime)
		billReq.EndTime = *endTIme
	}

	billList, total, err := dao.AccountBill.SelectAccountBills(ctx, sess, billReq)
	if err != nil {
		logger.Errorf("account.user_bill_list.dao.error, err: %v", err)
		return &userBillList.Data{
			Total:        consts.AmountZero,
			AccountBills: nil,
		}, errors.WithMessagef(common.ErrInternalServer, "account.user_bill_list.dao.error")
	}

	billListResp := make([]*v20230530.BillListData, 0, len(billList))
	for _, bill := range billList {
		billData := util.ModelToOpenApiBill(bill)
		billListResp = append(billListResp, billData)
	}

	return &userBillList.Data{
		Total:        total,
		AccountBills: billListResp,
	}, nil
}

// UserResourceBillList ...
func UserResourceBillList(ctx context.Context, req *userResourceBillList.Request, userID string) (*userResourceBillList.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.UserResourceBillList params: %+v, optUserID: %s", req, userID)
	account, err := getAccountInfoByUserID(ctx, userID, "account.ResourceBillList")
	if err != nil {
		return nil, err
	}

	billReq := &dao.BillListRequest{
		AccountID:   account.Id,
		TradeType:   req.TradeType,
		SignType:    req.SignType,
		ProductName: req.ProductName,
		SortByAsc:   req.SortByAsc,
		StartTime:   util.InvalidTime,
		PageSize:    req.PageSize,
		PageIndex:   req.PageIndex,
	}

	if req.StartTime != "" {
		startTime, _ := util.StringToTime(req.StartTime)
		billReq.StartTime = *startTime
	}

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	accountBillList, total, err := dao.AccountBill.SumBillAmountByResourceID(ctx, sess, billReq)
	if err != nil {
		logger.Errorf("Account.UserResourceBillList err: %v", err)
		return nil, errors.New("Account.UserResourceBillList.calculate error!")
	}

	resourceBillListResp := make([]*userResourceBillList.AccountResourceBillListData, 0, len(accountBillList))

	for _, bill := range accountBillList {
		// 查找资源商品信息
		billReq.ResourceID = bill.ResourceId
		billReq.SortByAsc = false
		billReq.PageIndex = 1
		billReq.PageSize = 1
		billReq.StartTime = util.InvalidTime
		billReq.EndTime = util.InvalidTime
		selAccountBillList, _, err := dao.AccountBill.SelectAccountBills(ctx, sess, billReq)
		if err != nil {
			logger.Errorf("Account.UserResourceBillList.query_db_err, err: %v", err)
			return nil, errors.New("Account.UserResourceBillList.query_db_err!")
		}
		if len(selAccountBillList) == 0 {
			msg := "Account.UserResourceBillList.selectAccountBills.respAccountBillList empty!"
			logger.Warnf(msg)
			return nil, errors.New(msg)
		}

		latestAccountBill := selAccountBillList[0]
		// 转换openApi格式数据
		resourceBillListInfo := util.ModelToOpenApiResourceBillListInfo(bill)
		resourceBillListInfo.AccountID = account.Id.String()
		resourceBillListInfo.MerchandiseID = latestAccountBill.MerchandiseId
		resourceBillListInfo.MerchandiseName = latestAccountBill.MerchandiseName
		resourceBillListInfo.UnitPrice = latestAccountBill.UnitPrice
		resourceBillListInfo.PriceDes = latestAccountBill.PriceDes
		resourceBillListInfo.QuantityUnit = latestAccountBill.QuantityUnit
		resourceBillListInfo.ProductName = latestAccountBill.ProductName
		// 退款字段特殊处理
		if req.TradeType == int64(accountbilltype.AccountBillTradeRefund) {
			resourceBillListInfo.TotalRefundAmount = resourceBillListInfo.TotalAmount
		}
		// 组装数组数据
		resourceBillListResp = append(resourceBillListResp, resourceBillListInfo)
	}

	return &userResourceBillList.Data{
		AccountBills: resourceBillListResp,
		Total:        total,
	}, nil
}

// getAccountInfoByUserID 根据 sso userID 获取account 信息
func getAccountInfoByUserID(ctx context.Context, userID, interfaceMsg string) (*models.Account, error) {
	logger := logging.GetLogger(ctx)
	// check sso user
	_, err := rpc.GetInstance().GetSSOUserByID(ctx, userID)
	if err != nil {
		if status.Code(err) == consts.ErrHydraLcpDBUserNotExist {
			return nil, common.ErrUserNotExists
		} else {
			logger.Errorf("%s error, err: %v", interfaceMsg, err)
			return nil, errors.Errorf("%s.query_sso_user_info error!", interfaceMsg)
		}
	}
	// get account info by user id
	uID := snowflake.MustParseString(userID)
	account, err := getOrCreateAccountByUserID(ctx, uID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// accountReduceByAccountID ...
func accountReduceByAccountID(ctx context.Context, req *ReduceRequest) (*AccountOperatorReply, error) {
	logger := logging.GetLogger(ctx)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradePay
	sign := accountbilltype.AccountBillSignReduce
	reduceAmount := req.Amount

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		if util.IsBlank(req.AccountCashVoucherIDs) {
			return &AccountCashVoucherResult{}, nil
		}

		// check account cash voucher valid
		trimIDs := strings.TrimRight(req.AccountCashVoucherIDs, ",")
		voucherIDs := strings.Split(trimIDs, ",")
		if len(voucherIDs) > consts.CASH_VOUCHER_EXCEED_AMOUNT {
			logger.Warnf("account.account_reduce.cash_voucher_exceed_amount.error!")
			return nil, common.ErrAccountVoucherExceed
		}

		voucherRelations := make([]*models.AccountCashVoucherRelation, 0, len(voucherIDs))
		for _, id := range voucherIDs {
			voucherId, err := snowflake.ParseString(id)
			if err != nil {
				logger.Errorf("account.account_reduce.parse_voucher_id.error! accountCashVoucherId:%s, err:%v", id, err)
				return nil, fmt.Errorf("account.account_reduce.parse_voucher_id.error! accountCashVoucherId:%s", id)
			}

			voucherRelation, err := dao.AccountCashVoucherRelation.Get(ctx, sess, voucherId, true)
			if err != nil {
				if errors.Is(err, common.ErrAccountVoucherIDNotFound) {
					return nil, errors.WithMessagef(common.ErrAccountVoucherIDNotFound, "voucher not exists, id:%s", id)
				}

				logger.Errorf("account.account_reduce.query_voucher.error, accountCashVoucherId: %s, err:%v", id, err)
				return nil, err
			}

			if voucherexpiredtype.IsExpiredType(voucherRelation.IsExpired) == voucherexpiredtype.EXPIRED {
				logger.Warnf("account.account_reduce.voucher_expired.error! accountCashVoucherId: %v", id)
				return nil, errors.WithMessagef(common.ErrAccountVoucherExpired, "account cash voucher expired, id: %v", id)
			}

			if voucherstatus.AccountCashVoucherStatus(voucherRelation.Status) == voucherstatus.DISABLED {
				logger.Warnf("account.account_reduce.voucher_disabled.error! accountCashVoucherId: %v", id)
				return nil, errors.WithMessagef(common.ErrAccountVoucherDisabled, "account voucher disabled, id: %v", id)
			}

			voucherRelations = append(voucherRelations, voucherRelation)
		}

		// 扣减金额
		var deltaVoucherAmount int64
		var usedAccountCashVoucherIDs []string
		var accountCashVoucherLogsIDs []string

		for _, voucherRelation := range voucherRelations {
			historyVoucherRelation := voucherRelation.String()

			if reduceAmount <= 0 {
				break
			}

			// 剩余金额
			remainingAmount := voucherRelation.RemainingAmount
			// 代金券剩余金额为0， 则返回空对象默认值
			if remainingAmount <= 0 {
				continue
			}

			// 代金券还有剩余金额可以抵扣，即 remainingAmount > 0
			if remainingAmount >= reduceAmount {
				deltaVoucherAmount += reduceAmount
				voucherRelation.RemainingAmount = voucherRelation.RemainingAmount - reduceAmount
				voucherRelation.UsedAmount = voucherRelation.UsedAmount + reduceAmount
				reduceAmount = 0
				remainingAmount = voucherRelation.RemainingAmount
			} else {
				deltaVoucherAmount += voucherRelation.RemainingAmount
				reduceAmount = reduceAmount - voucherRelation.RemainingAmount
				voucherRelation.UsedAmount = voucherRelation.CashVoucherAmount
				voucherRelation.RemainingAmount = 0
				remainingAmount = 0
			}

			logger.Infof("current voucher info: %s", util.ToJsonString(ctx, voucherRelation))

			// 更新账户代金券表，保存代金券日志
			_, err := dao.AccountCashVoucherRelation.Edit(ctx, sess, voucherRelation)
			if err != nil {
				logger.Errorf("account.account_reduce.update.accountCashVoucher.error! err:%v", err)
				return nil, errors.New("account.account_reduce.update.accountCashVoucher.error")
			}

			// 增加代金券到数组
			usedAccountCashVoucherIDs = append(usedAccountCashVoucherIDs, voucherRelation.Id.String())

			newID, err := rpc.GetInstance().GenID(ctx)
			if err != nil {
				logger.Errorf("account.payment_reduce.genID.error! err:%v", err)
				return nil, errors.New("account.account_reduce.genID.error")
			}
			now := time.Now()
			cashVoucherLog := &models.AccountCashVoucherLog{
				Id:                   newID,
				AccountId:            voucherRelation.AccountId,
				CashVoucherId:        voucherRelation.CashVoucherId,
				AccountCashVoucherId: voucherRelation.Id,
				SignType:             consts.VOUCHER_LOG_SIGN_CONSUME,
				SourceInfo:           historyVoucherRelation,
				TargetInfo:           voucherRelation.String(),
				Comment:              req.Comment,
				OptUserId:            snowflake.MustParseString(req.OptUserID),
				CreateTime:           now,
				UpdateTime:           now,
			}
			_, err = dao.AccountCashVoucherLog.Add(ctx, sess, cashVoucherLog)
			if err != nil {
				logger.Errorf("account.account_reduce.add.cashVoucherLog.error! err:%v", err)
				return nil, errors.New("account.account_reduce.add.cashVoucherLog.error!")
			}

			// 保存日志id
			accountCashVoucherLogsIDs = append(accountCashVoucherLogsIDs, newID.String())
		}

		return &AccountCashVoucherResult{
			AccountCashVoucherIDs:     usedAccountCashVoucherIDs,
			AccountCashVoucherLogsIDs: accountCashVoucherLogsIDs,
			DeltaVoucherAmount:        deltaVoucherAmount,
		}, nil
	}

	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		idempotentID, err := defaultIdempotentID(ctx, req.IdempotentID)
		if err != nil {
			return nil, err
		}

		createBillReq := &CreateAccountBillRequest{
			AccountID:       accountID,
			TradeID:         req.TradeID,
			IdempotentID:    idempotentID,
			OutTradeID:      0, // 外部订单暂时没有使用
			Amount:          reduceAmount,
			TradeType:       tradeType,
			Sign:            sign,
			Comment:         req.Comment,
			MerchandiseID:   req.MerchandiseID,
			MerchandiseName: req.MerchandiseName,
			UnitPrice:       req.UnitPrice,
			PriceDes:        req.PriceDes,
			Quantity:        req.Quantity,
			QuantityUnit:    req.QuantityUnit,
			ResourceID:      req.ResourceID,
			ProductName:     req.ProductName,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
		}

		// account bill not exists
		modelAccountBill, err := createAccountBill(ctx, sess, createBillReq)
		if err != nil {
			return nil, err
		}

		return modelAccountBill, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		// 账户是否冻结
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		// 超出消费余额， 账户余额 + 授信额度 > 0
		//if modelAccount.AccountBalance+modelAccount.CreditQuota < reduceAmount {
		//	return false, common.ErrInsufficientBalance
		//}

		normal, award := modelAccount.Reduce(reduceAmount)
		modelAccountBill.DeltaNormalBalance = normal
		modelAccountBill.DeltaAwardBalance = award
		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount
		return true /*update account bill*/, nil
	}

	reqParamsStr := fmt.Sprintf("Payment reduce requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: req.OptUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		ReduceAmount:   req.Amount,
		sign:           sign,
	}

	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		// 触发幂等性错误， 返回账户详情字段
		if errors.Is(err, common.ErrAccountBillIdempotentIDRepeat) {
			logger.Warnf("account.account_reduce.idempotentid_repeat, req params: %+v", req)
			account, err := getAccountReplyIfIdempotentIDRepeat(ctx, accountID)
			if err != nil {
				return nil, err
			}

			return &AccountOperatorReply{
				Account: &models.Account{
					Id:             snowflake.MustParseString(account.AccountID),
					Name:           account.AccountName,
					AccountBalance: account.AccountBalance,
					NormalBalance:  account.NormalBalance,
					AwardBalance:   account.AwardBalance,
					FreezedAmount:  account.FreezedAmount,
					CreditQuota:    account.CreditQuotaAmount,
				},
				accountCashVoucherTotalAmount: account.CashVoucherAmount,
			}, nil
		}
		logger.Errorf("account.account_reduce.accountOperator error, err: %v", err)
		return nil, err
	}

	return accountReply, nil
}

// AccountIDReduce ...
func AccountIDReduce(ctx context.Context, req *idreduce.Request, optUserID string) (*idreduce.Data, error) {
	logging.GetLogger(ctx).Infof("Account.AccountIDReduce params: %+v, optUserID: %s", req, optUserID)
	startTime, endTime := formatReqTime(req.StartTime, req.EndTime)

	accountReply, err := accountReduceByAccountID(ctx, &ReduceRequest{
		AccountID:             req.AccountID,
		Amount:                req.Amount,
		TradeID:               req.TradeID,
		IdempotentID:          req.IdempotentID,
		Comment:               req.Comment,
		AccountCashVoucherIDs: req.AccountCashVoucherIDs,
		VoucherConsumeMode:    req.VoucherConsumeMode,
		OptUserID:             optUserID,
		MerchandiseID:         req.MerchandiseID,
		MerchandiseName:       req.MerchandiseName,
		UnitPrice:             req.UnitPrice,
		PriceDes:              util.DefaultString(req.PriceDes, consts.UNIT_PRICE_DEFAULT),
		Quantity:              req.Quantity,
		QuantityUnit:          req.QuantityUnit,
		ResourceID:            req.ResourceID,
		ProductName:           req.ProductName,
		StartTime:             *startTime,
		EndTime:               *endTime,
	})

	if err != nil {
		return nil, err
	}

	return &idreduce.Data{
		Account: &v20230530.Account{
			AccountID:         accountReply.Id.String(),
			AccountName:       accountReply.Name,
			AccountBalance:    accountReply.AccountBalance,
			NormalBalance:     accountReply.NormalBalance,
			AwardBalance:      accountReply.AwardBalance,
			FreezedAmount:     accountReply.FreezedAmount,
			CashVoucherAmount: accountReply.accountCashVoucherTotalAmount,
			CreditQuotaAmount: accountReply.CreditQuota,
		},
	}, nil
}

// AccountYsIDReduce ...
func AccountYsIDReduce(ctx context.Context, req *ysidreduce.Request, optUserID string) (*ysidreduce.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.AccountYsIDReduce params: %+v, optUserID: %s", req, optUserID)
	startTime, endTime := formatReqTime(req.StartTime, req.EndTime)

	// get account info by user id
	userID := snowflake.MustParseString(req.UserID)
	account, err := getOrCreateAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	accountReply, err := accountReduceByAccountID(ctx, &ReduceRequest{
		AccountID:             account.Id.String(),
		Amount:                req.Amount,
		TradeID:               req.TradeID,
		IdempotentID:          req.IdempotentID,
		Comment:               req.Comment,
		AccountCashVoucherIDs: req.AccountCashVoucherIDs,
		VoucherConsumeMode:    req.VoucherConsumeMode,
		MerchandiseID:         req.MerchandiseID,
		MerchandiseName:       req.MerchandiseName,
		UnitPrice:             req.UnitPrice,
		PriceDes:              util.DefaultString(req.PriceDes, consts.UNIT_PRICE_DEFAULT),
		Quantity:              req.Quantity,
		QuantityUnit:          req.QuantityUnit,
		ResourceID:            req.ResourceID,
		ProductName:           req.ProductName,
		StartTime:             *startTime,
		EndTime:               *endTime,
	})

	if err != nil {
		return nil, err
	}

	return &ysidreduce.Data{
		Account: &v20230530.Account{
			AccountID:         accountReply.Id.String(),
			AccountName:       accountReply.Name,
			AccountBalance:    accountReply.AccountBalance,
			NormalBalance:     accountReply.NormalBalance,
			AwardBalance:      accountReply.AwardBalance,
			FreezedAmount:     accountReply.FreezedAmount,
			CashVoucherAmount: accountReply.accountCashVoucherTotalAmount,
			CreditQuotaAmount: accountReply.CreditQuota,
		},
	}, nil

}

// 通过UID获取资金账户信息，如果不存在则创建，在创建前会检查远算账号系统中该用户是否存在
func getOrCreateAccountByUserID(ctx context.Context, userID snowflake.ID) (*models.Account, error) {
	if userID == snowflake.ID(0) {
		return nil, common.ErrAccountNotExists
	}
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	account, err := dao.Account.GetByCustomerID(ctx, sess, userID)
	if err != nil {
		if errors.Is(err, common.ErrAccountNotExists) {
			logging.GetLogger(ctx).Infof("account of %s not existed, will create", userID)
			req := accountCreate.Request{
				AccountName: userID.String(),
				UserID:      userID.String(),
				AccountType: 2,
			}
			if acc, err := Create(ctx, &req, "create by system(when not existed)"); err != nil {
				return nil, err
			} else {
				accId, _ := snowflake.ParseString(acc.AccountID)
				account, err = dao.Account.Get(ctx, sess, accId, false)
				if err != nil {
					return nil, err
				}
			}
		} else {
			logging.GetLogger(ctx).Errorf("account.query_account error, err: %v", err)
			return nil, errors.New("account.query_account_db error")
		}
	}

	return account, nil
}

// IsAccountExists ...
func IsAccountExists(ctx context.Context, accountId snowflake.ID, sess *xorm.Session) (*models.Account, error) {
	logger := logging.GetLogger(ctx)
	if accountId == 0 {
		return nil, common.ErrAccountNotExists
	}

	account, err := dao.Account.Get(ctx, sess, accountId, false)
	if err != nil {
		if errors.Is(err, common.ErrAccountNotExists) {
			return nil, err
		} else {
			logger.Errorf("account.query_account_ error, err: %v", err)
			return nil, errors.New("account.query_account_db error")
		}
	}

	return account, nil
}

// AccountGetByID ...
func AccountGetByID(ctx context.Context, req *idget.Request, optUserID string) (*idget.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.AccountGetByID params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	accountID := snowflake.MustParseString(req.AccountID)
	account, err := IsAccountExists(ctx, accountID, sess)
	if err != nil {
		return nil, err
	}

	cashVoucherAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, accountID)
	if err != nil {
		return nil, err
	}

	return &idget.Data{
		AccountDetail: util.ModelToOpenApiAccount(account, cashVoucherAmount),
	}, nil
}

// AccountGetByYsID ...
func AccountGetByYsID(ctx context.Context, req *ysidget.Request, optUserID string) (*ysidget.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.AccountGetByYsID params: %+v, optUserID: %s", req, optUserID)
	userID := snowflake.MustParseString(req.UserID)
	account, err := getOrCreateAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	cashVoucherAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, account.Id)
	if err != nil {
		return nil, err
	}

	return &ysidget.Data{
		AccountDetail: util.ModelToOpenApiAccount(account, cashVoucherAmount),
	}, nil
}

func AccountGetByUserID(ctx context.Context, userID string) (*userIdGet.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.AccountGetByUserID params: %+v, optUserID: %s", userID, userID)

	// 根据userID 获取账户信息
	account, err := getAccountInfoByUserID(ctx, userID, "account.AccountGetByUserID")
	if err != nil {
		return nil, err
	}

	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	cashVoucherAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, account.Id)
	if err != nil {
		return nil, err
	}

	return &userIdGet.Data{
		AccountDetail: util.ModelToOpenApiAccount(account, cashVoucherAmount),
	}, nil
}

// FrozenModify ...
func FrozenModify(ctx context.Context, req *frozenmodify.Request, optUserID string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.FrozenModify params: %+v, optUserID: %s", req, optUserID)
	accountID := snowflake.MustParseString(req.AccountID)
	isFreezed := req.FrozenState

	err := with.DefaultSession(ctx, func(sess *xorm.Session) error {
		account, err := IsAccountExists(ctx, accountID, sess)
		if err != nil {
			logger.Errorf("freeze account err: %v", err)
			return err
		}

		account.IsFreeze = isFreezed
		account.UpdateTime = time.Now()
		_, err = dao.Account.EditFrozenStatus(ctx, sess, account)
		if err != nil {
			logger.Errorf("account.frozen_modify_db error, err: %v", err)
			return errors.New("account.frozen_modify_db error")
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, common.ErrAccountNotExists) {
			return err
		}
		logger.Errorf("account.frozen_modify error, err: %v", err)
		return errors.New("account.frozen_modify error")
	}

	return nil
}

// CreditQuotaModify ...
func CreditQuotaModify(ctx context.Context, req *creditquotamodify.Request, optUserID string) (*creditquotamodify.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.CreditQuotaModify params: %+v, optUserID: %s", req, optUserID)

	accountID := snowflake.MustParseString(req.AccountID)
	account := &models.Account{}
	_, err := boot.MW.DefaultTransaction(ctx, func(sess *xorm.Session) (interface{}, error) {
		// get account
		accountReply, err := dao.Account.Get(ctx, sess, accountID, true)
		if err != nil {
			if errors.Is(err, common.ErrAccountNotExists) {
				return nil, err
			} else {
				logger.Errorf("account.credit_quota_modify.query_account_db err: %v", err)
				return nil, errors.New("account.credit_quota_modify.query_account_db err!")
			}
		}

		account = accountReply

		// 账户是否冻结
		if account.IsFreeze {
			return nil, common.ErrFreezedAccount
		}

		// 超出消费余额， 账户余额 + 授信额度 <= 0
		if req.CreditQuotaAmount+account.AccountBalance <= 0 {
			return nil, common.ErrCreditQuotaExhausted
		}

		account.CreditQuota = req.CreditQuotaAmount
		account.UpdateTime = time.Now()
		_, err = dao.Account.EditCreditQuota(ctx, sess, account)
		if err != nil {
			logger.Errorf("account.credit_quota_modify.db error, err: %v", err)
			return nil, errors.New("account.credit_quota_modify.db error")
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &creditquotamodify.Data{
		AccountID:         account.Id.String(),
		CreditQuotaAmount: account.CreditQuota,
	}, nil
}

// PaymentFreezeUnfreeze ...
func PaymentFreezeUnfreeze(ctx context.Context, req *paymentfreezeunfreeze.Request, optUserID string) (*paymentfreezeunfreeze.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.paymentFreezeUnfreeze params: %+v, optUserID: %s", req, optUserID)

	idempotentID, err := defaultIdempotentID(ctx, "")
	if err != nil {
		return nil, err
	}

	// account freeze
	if req.IsFreezed {
		account, err := paymentFreeze(ctx, &paymentFreezeRequest{
			AccountID:    req.AccountID,
			Amount:       req.Amount,
			TradeID:      req.TradeID,
			IdempotentID: idempotentID,
			Comment:      req.Comment,
			OptUserID:    optUserID,
		})

		if err != nil {
			return nil, err
		}

		return util.ModelToOpenApiPaymentFreezedAccount(account.Account), nil
	} else {
		// 解冻操作
		account, err := paymentUnfreeze(ctx, &paymentUnfreezeRequest{
			AccountID:    req.AccountID,
			TradeID:      req.TradeID,
			IdempotentID: idempotentID,
			Comment:      req.Comment,
			OptUserID:    optUserID,
		})
		if err != nil {
			return nil, err
		}

		return util.ModelToOpenApiPaymentFreezedAccount(account.Account), nil
	}
}

type paymentFreezeRequest struct {
	AccountID    string
	Amount       int64
	TradeID      string
	IdempotentID string
	Comment      string
	OptUserID    string
}

// paymentFreeze ...
func paymentFreeze(ctx context.Context, req *paymentFreezeRequest) (*AccountOperatorReply, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.PaymentFreeze %+v", req)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradePay
	sign := accountbilltype.AccountBillSignFreeze

	freezeDelta := int64(0)
	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		modelAccountBill, err := dao.AccountBill.GetAccountBillBySearchParams(ctx, sess, &dao.BillGetRequest{
			AccountID: accountID,
			TradeID:   req.TradeID,
			TradeType: tradeType,
		})
		if err != nil {
			logger.Errorf("account.payment_freeze.found_account_bill err: %v", err)
			return nil, errors.New("account.payment_freeze.found_account_bill err")
		}

		if modelAccountBill.Id > 0 {
			// account bill exists
			if accountbilltype.AccountBillSign(modelAccountBill.Sign) != accountbilltype.AccountBillSignFreeze {
				msg := fmt.Sprintf("account.payment_freeze.found_account_bill can not change sign %v to %v", accountbilltype.AccountBillSign(modelAccountBill.Sign), sign)
				logger.Warnf(msg)
				return nil, errors.New(msg)
			}

			freezeDelta = req.Amount - modelAccountBill.Amount
			modelAccountBill.Amount = req.Amount
			modelAccountBill.Comment = req.Comment
			modelAccountBill.Sign = int(sign)
			modelAccountBill.IdempotentId = req.IdempotentID
			_, err1 := dao.AccountBill.Edit(ctx, sess, modelAccountBill)
			if err1 != nil {
				logger.Errorf("account.payment_freeze.update_account_bill err: %v", err1)
				return nil, errors.New("account.payment_freeze.update_account_bill err")
			}
		} else {
			freezeDelta = req.Amount
			// account bill not exists
			createBillReq := &CreateAccountBillRequest{
				AccountID:    accountID,
				TradeID:      req.TradeID,
				IdempotentID: req.IdempotentID,
				TradeType:    tradeType,
				Sign:         sign,
				Amount:       req.Amount,
				Comment:      req.Comment,
			}
			modelAccountBill, err = createAccountBill(ctx, sess, createBillReq)
			if err != nil {
				return nil, err
			}
		}
		return modelAccountBill, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		if freezeDelta != 0 {
			modelAccount.Freeze(freezeDelta)
		}
		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount

		return true /*not to update account_bill*/, nil
	}

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		return &AccountCashVoucherResult{}, nil
	}

	reqParamsStr := fmt.Sprintf("paymentFreeze requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: req.OptUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		sign:           sign,
	}
	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		logger.Errorf("account.operator error, err: %v", err)
		return nil, err
	}

	return accountReply, nil
}

type paymentUnfreezeRequest struct {
	AccountID    string
	TradeID      string
	IdempotentID string
	Comment      string
	OptUserID    string
}

// paymentUnfreeze ...
func paymentUnfreeze(ctx context.Context, req *paymentUnfreezeRequest) (*AccountOperatorReply, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account unfreeze: %+v", req)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradePay
	sign := accountbilltype.AccountBillSignUnfreeze

	lastFreezeAmount := int64(0)

	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		modelAccountBill, err := dao.AccountBill.GetAccountBillBySearchParams(ctx, sess, &dao.BillGetRequest{
			AccountID: accountID,
			TradeID:   req.TradeID,
			TradeType: tradeType,
		})

		if err != nil {
			logger.Errorf("account.payment_unfreeze.found_account_bill err: %v", err)
			return nil, errors.New("account.payment_unfreeze.found_account_bill err")
		}
		if modelAccountBill.Id <= 0 {
			logger.Warnf("account.payment_unfreeze.found_account_bill account_bill not exists! accountID: %v, tradeID:%s", accountID, req.TradeID)
			return nil, common.ErrAccountBillNotExists
		}

		if accountbilltype.AccountBillSign(modelAccountBill.Sign) != accountbilltype.AccountBillSignFreeze {
			msg := fmt.Sprintf("account.payment_unfreeze.found_account_bill can not change sign %v to %v", accountbilltype.AccountBillSign(modelAccountBill.Sign), sign)
			logger.Warnf(msg)
			return nil, errors.WithMessage(common.ErrAccountBillSignStatusInvalid, "account_bill sign is not freezed status")
		}

		lastFreezeAmount = modelAccountBill.Amount
		modelAccountBill.Comment = req.Comment
		modelAccountBill.UpdateTime = time.Now()
		modelAccountBill.Sign = int(sign)
		modelAccountBill.IdempotentId = req.IdempotentID
		_, err = dao.AccountBill.Edit(ctx, sess, modelAccountBill)
		if err != nil {
			logger.Errorf("account.payment_unfreeze.update_account_bill err: %v", err)
			return nil, errors.New("account.payment_unfreeze.update_account_bill err")
		}

		return modelAccountBill, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		if lastFreezeAmount != 0 {
			modelAccount.Unfreeze(lastFreezeAmount)
		}

		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount
		return true /*update account bill*/, nil
	}

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		return &AccountCashVoucherResult{}, nil
	}

	reqParamsStr := fmt.Sprintf("paymentUnfreeze requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: req.OptUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		sign:           sign,
	}
	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		logger.Errorf("accountOperator error, err: %v", err)
		return nil, err
	}

	return accountReply, nil
}

// AmountRefund ...
func AmountRefund(ctx context.Context, req *amountrefund.Request, optUserID string) (*amountrefund.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Account.AmountRefund params: %+v, optUserID: %s", req, optUserID)
	accountID := snowflake.MustParseString(req.AccountID)
	tradeType := accountbilltype.AccountBillTradeRefund
	sign := accountbilltype.AccountBillSignAdd

	handlerFoundAccountBill := func(ctx context.Context, sess *xorm.Session, params *OperatorParams) (*models.AccountBill, error) {
		idempotentID, err := defaultIdempotentID(ctx, req.IdempotentID)
		if err != nil {
			return nil, err
		}

		// 获取账单对应的业务类型
		var productName string
		if !util.IsBlank(req.ResourceID) {
			accountBill, err := dao.AccountBill.GetAccountBillBySearchParams(ctx, sess, &dao.BillGetRequest{
				AccountID:  accountID,
				ResourceID: req.ResourceID,
				TradeType:  accountbilltype.AccountBillTradePay,
			})
			if err != nil {
				logger.Errorf("account.account_refund.getAccountBill err: %v", err)
				return nil, errors.New("account.account_refund.getAccountBill err")
			}

			if accountBill != nil {
				productName = accountBill.ProductName
			}
		}

		createBillReq := &CreateAccountBillRequest{
			AccountID:    accountID,
			TradeID:      req.RefundID,
			IdempotentID: idempotentID,
			Amount:       req.Amount,
			TradeType:    tradeType,
			Sign:         sign,
			ResourceID:   req.ResourceID,
			ProductName:  productName,
			Comment:      req.Comment,
		}
		modelAccountBill, err := createAccountBill(ctx, sess, createBillReq)
		if err != nil {
			return nil, err
		}

		return modelAccountBill, nil
	}

	handlerUpdateAccount := func(ctx context.Context, sess *xorm.Session, modelAccount *models.Account, modelAccountBill *models.AccountBill, params *OperatorParams) (bool, error) {
		if modelAccount.IsFreeze {
			return false, common.ErrFreezedAccount
		}

		modelAccount.Add(req.Amount, consts.AmountZero)
		modelAccountBill.DeltaNormalBalance = req.Amount
		modelAccountBill.AccountBalance = modelAccount.AccountBalance
		modelAccountBill.FreezedAmount = modelAccount.FreezedAmount
		return true /*update account bill*/, nil
	}

	handlerAccountCashVoucher := func(ctx context.Context, sess *xorm.Session, optParams *OperatorParams) (*AccountCashVoucherResult, error) {
		return &AccountCashVoucherResult{}, nil
	}

	reqParamsStr := fmt.Sprintf("CreditAdd requst params: %v", util.ToJsonString(ctx, req))
	optParams := &OperatorParams{
		OperatorUserID: optUserID,
		ReqPrams:       reqParamsStr,
		AccountID:      req.AccountID,
		sign:           sign,
	}

	accountReply, err := accountOperator(ctx, handlerFoundAccountBill, handlerUpdateAccount, handlerAccountCashVoucher, optParams)
	if err != nil {
		// 触发幂等性错误， 返回账户详情字段
		if errors.Is(err, common.ErrAccountBillIdempotentIDRepeat) {
			logger.Warnf("account.account_refund.getAccountBill IdempotentId repeate, req params: %v", req)
			account, err := getAccountReplyIfIdempotentIDRepeat(ctx, accountID)
			if err != nil {
				return nil, err
			}

			return &amountrefund.Data{
				Account: &v20230530.Account{
					AccountID:         account.AccountID,
					AccountName:       account.AccountName,
					AccountBalance:    account.AccountBalance,
					NormalBalance:     account.NormalBalance,
					AwardBalance:      account.AwardBalance,
					FreezedAmount:     account.FreezedAmount,
					CreditQuotaAmount: account.CreditQuotaAmount,
					CashVoucherAmount: account.CashVoucherAmount,
				},
			}, nil
		}

		logger.Errorf("accountOperator error, err: %v", err)
		return nil, err
	}

	return &amountrefund.Data{
		Account: &v20230530.Account{
			AccountID:         accountReply.Id.String(),
			AccountName:       accountReply.Name,
			AccountBalance:    accountReply.AccountBalance,
			NormalBalance:     accountReply.NormalBalance,
			AwardBalance:      accountReply.AwardBalance,
			FreezedAmount:     accountReply.FreezedAmount,
			CreditQuotaAmount: accountReply.CreditQuota,
			CashVoucherAmount: accountReply.accountCashVoucherTotalAmount,
		},
	}, nil
}

func defaultIdempotentID(ctx context.Context, id string) (string, error) {
	if !util.IsBlank(id) {
		return id, nil
	}

	genID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		logging.GetLogger(ctx).Errorf("defaultIdempotentID generate id err %v", err)
		return "", errors.New("defaultIdempotentID generate id err")
	}

	return genID.String(), nil
}

func getAccountReplyIfIdempotentIDRepeat(ctx context.Context, accountID snowflake.ID) (*v20230530.Account, error) {
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	account, err := IsAccountExists(ctx, accountID, sess)
	if err != nil {
		return nil, err
	}

	cashVoucherAmount, err := dao.AccountCashVoucherRelation.GetTotalAmountByAccountId(ctx, sess, accountID)
	if err != nil {
		return nil, err
	}

	return &v20230530.Account{
		AccountID:         account.Id.String(),
		AccountName:       account.Name,
		AccountBalance:    account.AccountBalance,
		NormalBalance:     account.NormalBalance,
		AwardBalance:      account.AwardBalance,
		FreezedAmount:     account.FreezedAmount,
		CreditQuotaAmount: account.CreditQuota,
		CashVoucherAmount: cashVoucherAmount,
	}, nil
}
