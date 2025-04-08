package voucherexpiredtype

type IsExpiredType int64

const (
	NORMAL  IsExpiredType = 0
	EXPIRED IsExpiredType = 1
)

func GetIsExpiredType(expiredType IsExpiredType) (bool, IsExpiredType) {
	if expiredType == NORMAL {
		return true, NORMAL
	} else if expiredType == EXPIRED {
		return true, EXPIRED
	} else {
		return false, -1
	}
}

func (p IsExpiredType) String() string {
	switch p {
	case NORMAL:
		return "NORMAL"
	case EXPIRED:
		return "EXPIRED"
	default:
		return "UNKNOWN"
	}
}
