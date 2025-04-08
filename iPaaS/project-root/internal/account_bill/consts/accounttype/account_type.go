package accounttype

type AccountType int32

const (
	COMPANY  AccountType = 1
	PERSONAL AccountType = 2
)

func ValidAccountType(accountType AccountType) bool {
	if accountType == 0 {
		return false
	}

	if accountType == COMPANY || accountType == PERSONAL {
		return true
	}

	return false
}
