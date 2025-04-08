package db

import (
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	ErrDuplicatedEntry = errors.New("Duplicate entry")
)

const duplicatedEntryKeyWord = "Duplicate entry"

func IsDuplicatedError(err error) bool {
	if err == nil {
		return false
	}

	// 暂时只判断属于mysql的error类型
	var mysqlErr *mysql.MySQLError
	ok := errors.As(err, &mysqlErr)
	if ok {
		if mysqlErr.Number == 1062 || strings.Contains(mysqlErr.Message, duplicatedEntryKeyWord) {
			return true
		}
	}

	return false
}
