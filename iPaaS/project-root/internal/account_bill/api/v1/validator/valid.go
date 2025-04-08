package validator

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// ErrJudge ...
func ErrJudge(ctx *gin.Context, err error) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, common.ErrUserNotExists) {
		common.ErrorResp(ctx, http.StatusNotFound, consts.UserNotExistsErrorCode, "user not exists")
		return false
	}

	if errors.Is(err, common.ErrAccountNotExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountNotExists, "account not exists")
		return false
	}

	if errors.Is(err, common.ErrAccountExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountExists, "account exists")
		return false
	}

	if errors.Is(err, common.ErrAccountBillNotExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountBillNotExists, "account bill not exists")
		return false
	}

	if errors.Is(err, common.ErrCreditAddTradeExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.CreditAddTradeExists, "credit add trade id exists")
		return false
	}

	if errors.Is(err, common.ErrReduceTradeExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.ReduceTradeExists, "reduce trade id exists")
		return false
	}

	if errors.Is(err, common.ErrCreditQuotaExhausted) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.CreditQuotaExhausted, "credit quota exhausted, the amount can not less than current amount")
		return false
	}

	if errors.Is(err, common.ErrFreezedAccount) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.FreezedAccount, "the account has been frozen")
		return false
	}

	if errors.Is(err, common.ErrInsufficientBalance) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InsufficientBalance, "insufficient account balance")
		return false
	}

	if errors.Is(err, common.ErrCashVoucherNotExists) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.CashVoucherNotExists, err.Error())
		return false
	}

	if errors.Is(err, common.ErrAccountVoucherExpired) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountVoucherExpired, err.Error())
		return false
	}

	if errors.Is(err, common.ErrAccountVoucherDisabled) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountVoucherDisabled, err.Error())
		return false
	}

	if errors.Is(err, common.ErrAccountBillSignStatusInvalid) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountBillSignStatusInvalid, err.Error())
		return false
	}

	if errors.Is(err, common.ErrAccountBillIdempotentIDRepeat) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.AccountBillIdempotentIDRepeat, "account bill has repeat idempotentID ")
		return false
	}

	common.InternalServerError(ctx, err.Error())
	return false
}

func ValidAccountID(ctx *gin.Context, accountID string) bool {
	if util.IsBlank(accountID) || strings.TrimSpace(accountID) == ":" {
		msg := fmt.Sprintf("AccountID is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	ID, err := snowflake.ParseString(accountID)
	if err != nil || ID == 0 {
		msg := fmt.Sprintf("invalid accountID: %v", accountID)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidTradeID(ctx *gin.Context, tradeID string) bool {
	if util.IsBlank(tradeID) {
		msg := fmt.Sprintf("TradeID is required")
		common.InvalidParams(ctx, msg)
		return false
	}
	if !util.IsBlank(tradeID) && len(tradeID) > 128 {
		msg := fmt.Sprintf("invalid params, TradeID length can not be exceed 128 characters")
		logging.GetLogger(ctx).Warnf(msg)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidRefundID(ctx *gin.Context, tradeID string) bool {
	if util.IsBlank(tradeID) {
		msg := fmt.Sprintf("RefundID is required")
		common.InvalidParams(ctx, msg)
		return false
	}
	if !util.IsBlank(tradeID) && len(tradeID) > 128 {
		msg := fmt.Sprintf("invalid params, RefundID length can not be exceed 128 characters")
		logging.GetLogger(ctx).Warnf(msg)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidIdempotentID(ctx *gin.Context, idempotentID string) bool {
	if !util.IsBlank(idempotentID) && len(idempotentID) > 256 {
		msg := fmt.Sprintf("invalid params, idempotent ID length can not be exceed 256 characters")
		logging.GetLogger(ctx).Error(msg)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidAmount(ctx *gin.Context, amount int64) bool {
	if amount < 0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "Amount can not be less than 0")
		return false
	}

	return true
}

func ValidComment(ctx *gin.Context, comment string) bool {
	if util.IsBlank(comment) {
		msg := fmt.Sprintf("comment is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	if len(comment) > 255 {
		msg := fmt.Sprintf("commnet length can not be exceed 256 characters")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidMerchandiseID(ctx *gin.Context, merchandiseID string) bool {
	if util.IsBlank(merchandiseID) {
		msg := fmt.Sprintf(" merchandise ID is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidMerchandiseName(ctx *gin.Context, merchandiseName string) bool {
	if util.IsBlank(merchandiseName) {
		msg := fmt.Sprintf(" merchandise name is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	if len(merchandiseName) > 256 {
		msg := fmt.Sprintf(" merchandise name length can not be exceed 256 characters")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidUnitPrice(ctx *gin.Context, unitPrice int64) bool {
	if unitPrice <= 0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "unit price can not be less than or equal to 0")
		return false
	}

	return true
}

func ValidQuantity(ctx *gin.Context, quantity float64) bool {
	if quantity <= 0.0 {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidAmount, "quantity can not be less than or equal to 0")
		return false
	}

	return true
}

func ValidQuantityUnit(ctx *gin.Context, quantityUnit string) bool {
	if util.IsBlank(quantityUnit) {
		msg := fmt.Sprintf("quantity unit is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	if len(quantityUnit) > 256 {
		msg := fmt.Sprintf("quantity unit length can not be exceed 256 characters")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidResourceID(ctx *gin.Context, resourceID string) bool {

	ID, err := snowflake.ParseString(resourceID)
	if err != nil || ID == 0 {
		msg := fmt.Sprintf("invalid resourceID")
		common.InvalidParams(ctx, msg)
		return false
	}

	if !util.IsBlank(resourceID) && len(resourceID) > 128 {
		msg := fmt.Sprintf("invalid params, resourceID length can not be exceed 128 characters")
		logging.GetLogger(ctx).Warnf(msg)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidStartTimeAndEndTime(ctx *gin.Context, startTimeStr string, endTimeStr string) bool {
	var startTime, endTime *time.Time
	var toTimeResult bool
	if startTimeStr != "" {
		startTime, toTimeResult = util.StringToTime(startTimeStr)
		if !toTimeResult {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "start time date format is incorrect，for example:2006-01-02 15:04:05")
			return false
		}
	}

	if endTimeStr != "" {
		endTime, toTimeResult = util.StringToTime(endTimeStr)
		if !toTimeResult {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "end time date format is incorrect，for example:2006-01-02 15:04:05")
			return false
		}
	}

	// 校验 end > start
	if startTimeStr != "" && endTimeStr != "" {
		if endTime.Before(*startTime) {
			common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidEndTime, "end time should not be before start time")
			return false
		}
	}

	return true
}

func ValidResourceStartTimeAndEndTime(ctx *gin.Context, startTimeStr string, endTimeStr string) bool {
	if (util.IsNotBlank(startTimeStr) && util.IsBlank(endTimeStr)) ||
		(util.IsBlank(startTimeStr) && util.IsNotBlank(endTimeStr)) {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidArgumentErrorCode, "start time and end time should not be empty at same time")
		return false
	}

	return true
}

func ValidUserID(ctx *gin.Context, userID string) bool {
	if util.IsBlank(userID) || strings.TrimSpace(userID) == ":" {
		msg := fmt.Sprintf("userID is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	ID, err := snowflake.ParseString(userID)
	if err != nil || ID == 0 {
		msg := fmt.Sprintf("invalid userID: %v", userID)
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidAccountCashVoucherID(ctx *gin.Context, accountCashVoucherID string) bool {
	if util.IsBlank(accountCashVoucherID) || strings.TrimSpace(accountCashVoucherID) == ":" {
		msg := fmt.Sprintf("AccountCashVoucherID is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	_, err := snowflake.ParseString(accountCashVoucherID)
	if err != nil {
		msg := fmt.Sprintf("invalid accountCashVoucherID")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidCashVoucherID(ctx *gin.Context, cashVoucherID string) bool {
	if util.IsBlank(cashVoucherID) || strings.TrimSpace(cashVoucherID) == ":" {
		msg := fmt.Sprintf("CashVoucherID is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	_, err := snowflake.ParseString(cashVoucherID)
	if err != nil {
		msg := fmt.Sprintf("invalid cashVoucherID")
		common.InvalidParams(ctx, msg)
		return false
	}

	return true
}

func ValidAccountIDs(ctx *gin.Context, accountIDs string) bool {
	if util.IsBlank(accountIDs) {
		msg := fmt.Sprintf("accountIDs is required")
		common.InvalidParams(ctx, msg)
		return false
	}

	for _, accountID := range strings.Split(accountIDs, ",") {
		_, err := snowflake.ParseString(accountID)
		if err != nil {
			msg := fmt.Sprintf(" invalid accountIDs")
			common.InvalidParams(ctx, msg)
			return false
		}
	}

	return true
}

func ValidAuthUserID(ctx *gin.Context) (string, bool) {
	userID, err := util.GetUserID(ctx)
	if err != nil {
		common.ErrorResp(ctx, http.StatusBadRequest, consts.InvalidUserId, err.Error())
		return "", false
	}

	return userID, true

}
