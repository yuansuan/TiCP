package common

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"time"
)

type contextKey string

const (
	IamAdminTag           string     = "IamAdmin"
	ContextTransactionKey contextKey = "Transaction"
)

func GetBigTime() time.Time {
	return time.Time{}.AddDate(2407, 1, 1)
}

var IdGen *snowflake.Node

func init() {
	var err error
	IdGen, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}
