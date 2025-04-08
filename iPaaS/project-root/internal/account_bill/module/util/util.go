package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	resourceBillList "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/apiuser/resourcebilllist"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentfreezeunfreeze"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"

	//"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// InvalidTime ...
var InvalidTime, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00+08:00")

func ModelToOpenApiBill(bill *models.AccountBill) *v20230530.BillListData {
	startTimeStr, endTimeStr := "", ""
	if !InvalidTime.Equal(bill.StartTime) {
		startTimeStr = ModelTimeToString(bill.StartTime)
	}

	if !InvalidTime.Equal(bill.EndTime) {
		endTimeStr = ModelTimeToString(bill.EndTime)
	}

	return &v20230530.BillListData{
		ID:                    bill.Id.String(),
		AccountID:             bill.AccountId.String(),
		Amount:                bill.Amount,
		AccountBalance:        bill.AccountBalance,
		FreezedAmount:         bill.FreezedAmount,
		DeltaNormalBalance:    bill.DeltaNormalBalance,
		DeltaAwardBalance:     bill.DeltaAwardBalance,
		DeltaDeductionBalance: bill.DeltaVoucherBalance,
		TradeID:               ToStringIfSnowflake(bill.TradeId),
		SignType:              int64(bill.Sign),
		TradeType:             int64(bill.TradeType),
		Comment:               bill.Comment,
		TradeTime:             ModelTimeToString(bill.CreateTime),
		MerchandiseID:         bill.MerchandiseId,
		MerchandiseName:       bill.MerchandiseName,
		UnitPrice:             bill.UnitPrice,
		PriceDes:              bill.PriceDes,
		Quantity:              bill.Quantity,
		QuantityUnit:          bill.QuantityUnit,
		ResourceID:            bill.ResourceId,
		ProductName:           bill.ProductName,
		StartTime:             startTimeStr,
		EndTime:               endTimeStr,
	}
}

func ModelToOpenApiAccountVoucher(req *models.AccountCashVoucherRelation) *v20230530.AccountCashVoucher {
	return &v20230530.AccountCashVoucher{
		AccountCashVoucherID: req.Id.String(),
		AccountID:            req.AccountId.String(),
		CashVoucherID:        req.CashVoucherId.String(),
		Amount:               req.CashVoucherAmount,
		UsedAmount:           req.UsedAmount,
		RemainingAmount:      req.RemainingAmount,
		Status:               req.Status,
		ExpiredTime:          req.ExpiredTime.String(),
		IsExpired:            req.IsExpired,
		CreateTime:           req.CreateTime.String(),
	}
}

func ModelToOpenApiPaymentFreezedAccount(account *models.Account) *paymentfreezeunfreeze.Data {
	return &paymentfreezeunfreeze.Data{
		AccountID:         account.Id.String(),
		AccountName:       account.Name,
		AccountBalance:    account.AccountBalance,
		NormalBalance:     account.NormalBalance,
		AwardBalance:      account.AwardBalance,
		FreezedAmount:     account.FreezedAmount,
		CreditQuotaAmount: account.CreditQuota,
	}
}

func ModelToOpenApiCashVoucherDetail(cashVoucher *models.CashVoucher) *v20230530.CashVoucher {
	if cashVoucher == nil {
		return nil
	}

	availabilityStatus := voucherstatus.AvailabilityStatusType(cashVoucher.AvailabilityStatus)
	_, status := voucherstatus.GetAvailabilityStatusType(availabilityStatus)

	isExpiredType := voucherexpiredtype.IsExpiredType(cashVoucher.IsExpired)
	_, isExpiredTypeResp := voucherexpiredtype.GetIsExpiredType(isExpiredType)

	return &v20230530.CashVoucher{
		CashVoucherID:      cashVoucher.Id.String(),
		CashVoucherName:    cashVoucher.Name,
		AvailabilityStatus: strings.ToLower(status.String()),
		Amount:             cashVoucher.Amount,
		IsExpired:          strings.ToLower(isExpiredTypeResp.String()),
		ExpiredType:        cashVoucher.ExpiredType,
		AbsExpiredTime:     ModelTimeToString(cashVoucher.AbsExpiredTime),
		RelExpiredTime:     cashVoucher.RelExpiredTime,
		Comment:            cashVoucher.Comment,
		CreateTime:         ModelTimeToString(cashVoucher.CreateTime),
	}
}

func ModelToOpenApiAccount(account *models.Account, cashVoucherAmount int64) *v20230530.AccountDetail {
	if account == nil {
		return nil
	}

	return &v20230530.AccountDetail{
		AccountID:         account.Id.String(),
		AccountName:       account.Name,
		CustomerID:        account.CustomerId.String(),
		Currency:          account.Currency,
		AccountBalance:    account.AccountBalance,
		NormalBalance:     account.NormalBalance,
		AwardBalance:      account.AwardBalance,
		CreditQuotaAmount: account.CreditQuota,
		FreezedAmount:     account.FreezedAmount,
		CashVoucherAmount: cashVoucherAmount,
		FrozenStatus:      account.IsFreeze,
		CreateTime:        ModelTimeToString(account.CreateTime),
		UpdateTime:        ModelTimeToString(account.UpdateTime),
	}
}

func ModelToOpenApiResourceBillListInfo(accountBill *models.AccountBill) *resourceBillList.AccountResourceBillListData {
	if accountBill == nil {
		return nil
	}

	return &resourceBillList.AccountResourceBillListData{
		TotalAmount:          accountBill.Amount,
		TotalNormalAmount:    accountBill.DeltaNormalBalance,
		TotalAwardAmount:     accountBill.DeltaAwardBalance,
		TotalFreezedAmount:   accountBill.FreezedAmount,
		TotalDeductionAmount: accountBill.DeltaVoucherBalance,
		TotalDiscountAmount:  consts.AmountZero,
		Quantity:             accountBill.Quantity,
		ResourceID:           accountBill.ResourceId,
		StartTime:            ModelTimeToString(accountBill.StartTime),
		EndTime:              ModelTimeToString(accountBill.EndTime),
		LatestTradeTime:      ModelTimeToString(accountBill.CreateTime),
	}
}

// ModelTimeToString ...
func ModelTimeToString(timeInput time.Time) string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	formatTime := timeInput.In(cstSh).Format(time.RFC3339)
	return formatTime
}

// StringToTime ...
func StringToTime(timeString string) (*time.Time, bool) {
	if timeString == "" {
		return nil, false
	}
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	parseTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, cstSh)
	if err != nil {
		return nil, false
	}
	return &parseTime, true
}

func GetUserID(c *gin.Context) (string, error) {
	userIDStr := c.GetHeader(common.UserInfoKey)
	if userIDStr == "" {
		return "", fmt.Errorf("%s is empty in HTTP Header", common.UserInfoKey)
	}

	_, err := snowflake.ParseString(userIDStr)
	if err != nil {
		return "", fmt.Errorf("parse userId [%s] failed, %w", userIDStr, err)
	}

	return userIDStr, nil
}

func ToJsonString(c context.Context, v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		logging.GetLogger(c).Warnf("json parse error, err: %v", err)
		return ""
	}
	return string(bytes)
}

// ToStringIfSnowflake 如果是snowflake格式, 则返回snowflake格式base字符串
func ToStringIfSnowflake(str string) string {
	if IsBlank(str) {
		return ""
	}

	// 是否是snowflake 格式
	snowflakeID := snowflake.String2SnowflakeID(str)
	if snowflakeID != 0 {
		return snowflakeID.String()
	}

	return str
}

func IsBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func DefaultString(str, defaultStr string) string {
	if IsNotBlank(str) {
		return str
	} else if IsNotBlank(defaultStr) {
		return defaultStr
	} else {
		return ""
	}
}

func AllStringsNonEmpty(arr []string) bool {
	for _, str := range arr {
		if len(strings.TrimSpace(str)) == 0 {
			return false
		}
	}
	return true
}
