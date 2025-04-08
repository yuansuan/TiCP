package voucherexpiredtype

type ExpiredType int64

const (
	AbsExpired      ExpiredType = 1
	RelativeExpired ExpiredType = 2
)

func Valid(expiredType ExpiredType) bool {
	return expiredType == AbsExpired || expiredType == RelativeExpired
}
