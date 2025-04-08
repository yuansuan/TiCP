package voucherstatus

import (
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
)

type AvailabilityStatusType int64

const (
	UNAVAILABLE AvailabilityStatusType = 0
	AVAILABLE   AvailabilityStatusType = 1
)

var statusMap = map[string]AvailabilityStatusType{}

func init() {
	statusMap["UNAVAILABLE"] = UNAVAILABLE
	statusMap["AVAILABLE"] = AVAILABLE
}

func ValidAvailabilityStatusInt(status AvailabilityStatusType) bool {
	return status == UNAVAILABLE || status == AVAILABLE
}

func ValidAvailabilityStatusString(status string) (bool, AvailabilityStatusType) {
	if len(strings.TrimSpace(status)) == 0 {
		return false, -1
	}
	upper := strings.ToUpper(status)
	statusType, exists := statusMap[upper]
	if !exists {
		return false, 0
	}

	return true, statusType
}

func GetAvailabilityStatusType(status AvailabilityStatusType) (bool, AvailabilityStatusType) {
	if status == AVAILABLE {
		return true, AVAILABLE
	} else if status == UNAVAILABLE {
		return true, UNAVAILABLE
	} else {
		return false, consts.InvalidStatus
	}
}
func (p AvailabilityStatusType) String() string {
	switch p {
	case AVAILABLE:
		return "AVAILABLE"
	case UNAVAILABLE:
		return "UNAVAILABLE"
	default:
		return "UNKNOWN"
	}
}
