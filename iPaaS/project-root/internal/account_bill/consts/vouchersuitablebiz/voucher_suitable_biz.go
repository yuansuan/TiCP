package vouchersuitablebiz

import "strings"

type VoucherSuitableBizType string

const (
	Computing VoucherSuitableBizType = "computing"
	CloudApp  VoucherSuitableBizType = "cloudApp"
)

func ValidVoucherSuitableBizType(bizType VoucherSuitableBizType) bool {
	s := string(bizType)
	if len(strings.TrimSpace(s)) == 0 {
		return true
	}

	splitString := strings.Split(s, ",")
	for _, str := range splitString {
		if len(strings.TrimSpace(str)) == 0 {
			continue
		}

		if Computing != VoucherSuitableBizType(str) && CloudApp != VoucherSuitableBizType(str) {
			return false
		}
	}

	return true
}
