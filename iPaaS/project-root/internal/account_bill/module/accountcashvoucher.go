package module

import (
	"context"
	"errors"
	"strings"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/add"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/get"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/list"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/statusmodify"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
	"xorm.io/xorm"
)

// List ...
func List(ctx context.Context, req *list.Request, optUserID string) (*list.Data, error) {
	logging.GetLogger(ctx).Infof("Account.BillList params: %v, optUserID: %s", util.ToJsonString(ctx, req), optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	// check account exists
	var accountID snowflake.ID
	if req.AccountID != "" {
		accountID = snowflake.MustParseString(req.AccountID)
		_, err := IsAccountExists(ctx, accountID, sess)
		if err != nil {
			return nil, err
		}
	}

	listRequest := &dao.ListRequest{
		OptUserID: snowflake.MustParseString(optUserID),
		AccountID: accountID,
		StartTime: util.InvalidTime,
		EndTime:   util.InvalidTime,
		PageSize:  req.PageSize,
		PageIndex: req.PageIndex,
	}
	if req.StartTime != "" {
		startTime, _ := util.StringToTime(req.StartTime)
		listRequest.StartTime = *startTime
	}
	if req.EndTime != "" {
		endTIme, _ := util.StringToTime(req.EndTime)
		listRequest.EndTime = *endTIme
	}
	accountVouchers, total, err := dao.AccountCashVoucherRelation.SelectByAccountId(ctx, sess, listRequest)
	if err != nil {
		return nil, err
	}
	listResp := make([]*v20230530.AccountCashVoucher, 0, len(accountVouchers))
	for _, item := range accountVouchers {
		listData := util.ModelToOpenApiAccountVoucher(item)
		listResp = append(listResp, listData)
	}

	return &list.Data{
		Total:               total,
		AccountCashVouchers: listResp,
	}, nil
}

func AccountCashVoucherGetByID(ctx context.Context, req *get.Request, optUserID string) (*v20230530.AccountCashVoucher, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("AccountCashVoucher.AccountCashVoucherGetByID params: %v, optUserID: %s", util.ToJsonString(ctx, req), optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	accountCashVoucherID := snowflake.MustParseString(req.AccountCashVoucherID)

	relation, err := dao.AccountCashVoucherRelation.Get(ctx, sess, accountCashVoucherID, false)
	if err != nil {
		return nil, err
	}

	return &v20230530.AccountCashVoucher{
		AccountCashVoucherID: relation.Id.String(),
		AccountID:            relation.AccountId.String(),
		CashVoucherID:        relation.CashVoucherId.String(),
		Amount:               relation.RemainingAmount,
		UsedAmount:           relation.UsedAmount,
		RemainingAmount:      relation.RemainingAmount,
		Status:               relation.Status,
		ExpiredTime:          relation.ExpiredTime.String(),
		IsExpired:            relation.IsExpired,
		CreateTime:           relation.CreateTime.String(),
	}, nil
}

func StatusModify(ctx context.Context, req *statusmodify.Request, optUserID string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("AccountCashVoucher.AccountCashVoucherGetByID params: %v, optUserID: %s", util.ToJsonString(ctx, req), optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	accountCashVoucherID := snowflake.MustParseString(req.AccountCashVoucherID)

	voucherRelation, err := dao.AccountCashVoucherRelation.Get(ctx, sess, accountCashVoucherID, false)
	if err != nil {
		return err
	}

	_, status := voucherstatus.ValidAccountCashVoucherStatusString(req.Status)
	voucherRelation.Id = accountCashVoucherID
	voucherRelation.Status = int(status)
	_, err = dao.AccountCashVoucherRelation.Edit(ctx, sess, voucherRelation)
	if err != nil {
		return err
	}
	return nil
}

func Add(ctx context.Context, req *add.Request, optUserID string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("AccountCashVoucher.Add params: %v, optUserID: %s", util.ToJsonString(ctx, req), optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	cashVoucherID := snowflake.MustParseString(req.CashVoucherID)

	//校验账户
	accountIDs := strings.Split(req.AccountIDs, ",")
	for _, accountID := range accountIDs {
		_, err := IsAccountExists(ctx, snowflake.MustParseString(accountID), sess)
		if err != nil {
			return err
		}
	}

	//构建待新增对象
	accountVouchers, err := buildModel(ctx, sess, cashVoucherID, accountIDs, optUserID)
	if err != nil {
		return err
	}

	//批量新增
	_, err = dao.AccountCashVoucherRelation.BatchAdd(ctx, sess, accountVouchers)
	if err != nil {
		return err
	}

	return nil
}

func buildModel(ctx context.Context, sess *xorm.Session, cashVoucherID snowflake.ID, accountIDs []string, optUserID string) ([]*models.AccountCashVoucherRelation, error) {
	//查询代金券信息
	cashVoucher, err := dao.CashVoucher.Get(ctx, sess, cashVoucherID, false)
	if err != nil {
		return nil, err
	}

	//构建账户代金券结构
	ids, err := rpc.GetInstance().GenIDs(ctx, int64(len(accountIDs)))
	if err != nil {
		return nil, errors.New("error calling idgen service")
	}

	var expiredTime time.Time
	now := time.Now()
	if voucherexpiredtype.ExpiredType(cashVoucher.ExpiredType) == voucherexpiredtype.RelativeExpired {
		duration := time.Duration(cashVoucher.RelExpiredTime) * time.Second / time.Nanosecond
		expiredTime = now.Add(duration)
	} else {
		expiredTime = cashVoucher.AbsExpiredTime
	}

	var accountVouchers []*models.AccountCashVoucherRelation
	for index, accountID := range accountIDs {
		accountVouchers = append(accountVouchers, &models.AccountCashVoucherRelation{
			Id:                ids[index],
			AccountId:         snowflake.MustParseString(accountID),
			CashVoucherId:     cashVoucherID,
			CashVoucherAmount: cashVoucher.Amount,
			UsedAmount:        consts.AmountZero,
			RemainingAmount:   cashVoucher.Amount,
			Status:            int(voucherstatus.ENABLED),
			ExpiredTime:       expiredTime,
			IsExpired:         int(voucherexpiredtype.NORMAL),
			IsDeleted:         consts.NoDeleted,
			OptUserId:         snowflake.MustParseString(optUserID),
			CreateTime:        now,
			UpdateTime:        now,
		})
	}
	return accountVouchers, nil
}
