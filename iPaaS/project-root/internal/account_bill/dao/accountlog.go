package dao

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"xorm.io/xorm"
)

type AccountLogDaoImpl struct {
}

func NewAccountLogDaoImpl() *AccountLogDaoImpl {
	return &AccountLogDaoImpl{}
}

func (o *AccountLogDaoImpl) Add(ctx context.Context, session *xorm.Session, accountLog *models.AccountLog) (snowflake.ID, error) {
	_, err := session.Insert(accountLog)
	if err != nil {
		return consts.AmountZero, errors.Wrap(err, "account log add dao:")
	}

	return accountLog.Id, nil
}
