package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

func InsertBill(ctx context.Context, postPaidBill *models.Bill) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, e := db.Insert(postPaidBill)
		return e
	})
}

func UpdateBillTimeByOrderId(ctx context.Context, orderId snowflake.ID, billTime time.Time) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Where("order_id = ?", orderId).
			Cols("bill_time").
			Update(&models.Bill{
				BillTime: billTime,
			})
		return err
	})
}

func GetBill(ctx context.Context, sessionId, resourceId snowflake.ID) (*models.Bill, bool, error) {
	postPaidBill := new(models.Bill)
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		var e error
		exist, e = db.Where("session_id = ?", sessionId).
			Where("resource_id = ?", resourceId).
			Get(postPaidBill)
		return e
	})
	if err != nil {
		return nil, false, errors.Wrap(err, "dao")
	}
	if !exist {
		return nil, false, nil
	}

	return postPaidBill, true, nil
}
