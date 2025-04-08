package voucherstatus

import (
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
)

type AccountCashVoucherStatus int64

const (
	ENABLED  AccountCashVoucherStatus = 0
	DISABLED AccountCashVoucherStatus = 1
)

var accountCashStatusMap = map[string]AccountCashVoucherStatus{}

func init() {
	accountCashStatusMap["ENABLED"] = ENABLED
	accountCashStatusMap["DISABLED"] = DISABLED
}

func ValidAccountCashVoucherStatusInt(status AccountCashVoucherStatus) bool {
	return status == ENABLED || status == DISABLED
}

func ValidAccountCashVoucherStatusString(status string) (bool, AccountCashVoucherStatus) {
	if len(strings.TrimSpace(status)) == 0 {
		return false, -1
	}
	upper := strings.ToUpper(status)
	statusType, exists := accountCashStatusMap[upper]
	if !exists {
		return false, 0
	}

	return true, statusType
}

func GetAccountCashVoucherStatus(status AccountCashVoucherStatus) (bool, AccountCashVoucherStatus) {
	if status == ENABLED {
		return true, ENABLED
	} else if status == DISABLED {
		return true, DISABLED
	} else {
		return false, consts.InvalidStatus
	}
}

func (p AccountCashVoucherStatus) String() string {
	switch p {
	case ENABLED:
		return "ENABLED"
	case DISABLED:
		return "DISABLED"
	default:
		return "UNKNOWN"
	}
}
