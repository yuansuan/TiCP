package module

import (
	"context"
	"errors"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/add"
	availabilityModify "github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/availabilitymodify"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/get"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/cashvoucher/list"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherexpiredtype"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts/voucherstatus"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/util"
)

// AddCashVoucher ...
func AddCashVoucher(ctx context.Context, req *add.Request, optUserID string) (*add.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("cashVoucher.AddCashVoucher params: %+v, optUserID: %s", req, optUserID)
	newID, err := rpc.GetInstance().GenID(ctx)
	if err != nil {
		logger.Warnf("gen snowflake error, err: %v", err)
		return nil, errors.New("cashVoucher.addCashVoucher.gen_ID error!")
	}

	absExpiredTime, _ := util.StringToTime(req.AbsExpiredTime)
	if voucherexpiredtype.ExpiredType(req.ExpiredType) == voucherexpiredtype.RelativeExpired {
		absExpiredTime = &util.InvalidTime
	} else if voucherexpiredtype.ExpiredType(req.ExpiredType) == voucherexpiredtype.AbsExpired {
		req.RelExpiredTime = 0
	}

	now := time.Now()
	cashVoucher := &models.CashVoucher{
		Id:                 newID,
		Name:               req.CashVoucherName,
		Amount:             req.Amount,
		AvailabilityStatus: int(voucherstatus.UNAVAILABLE),
		OptUserId:          snowflake.MustParseString(optUserID),
		IsExpired:          0,
		ExpiredType:        int(req.ExpiredType),
		AbsExpiredTime:     *absExpiredTime,
		RelExpiredTime:     req.RelExpiredTime,
		Comment:            req.Comment,
		CreateTime:         now,
		UpdateTime:         now,
	}
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	cashVoucherID, err := dao.CashVoucher.Add(ctx, sess, cashVoucher)
	if err != nil {
		logger.Warnf("add cash voucher db error, err: %v", err)
		return nil, errors.New("add cash voucher db error")
	}
	return &add.Data{
		CashVoucherID: cashVoucherID.String(),
	}, nil
}

// GetCashVoucherByID
func GetCashVoucherByID(ctx context.Context, req *get.Request, optUserID string) (*get.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("cashVoucher.GetCashVoucherByID params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	ID := snowflake.MustParseString(req.CashVoucherID)
	cashVoucherResp, err := dao.CashVoucher.Get(ctx, sess, ID, false)
	if err != nil {
		logger.Warnf("cashVoucher.GetCashVoucherByID.query cash voucher err: %v", err)
		return nil, err
	}

	voucherDetail := util.ModelToOpenApiCashVoucherDetail(cashVoucherResp)

	resp := &get.Data{}
	resp.CashVoucher = *voucherDetail
	return resp, nil
}

// ListCashVoucher
func ListCashVoucher(ctx context.Context, req *list.Request, optUserID string) (*list.Data, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("cashVoucher.ListCashVoucher params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()

	selectReq := &models.SelectCashVouchersReq{
		Id:                 snowflake.MustParseString(req.CashVoucherID),
		Name:               req.CashVoucherName,
		IsExpired:          req.IsExpired,
		AvailabilityStatus: req.AvailabilityStatus,
		StartTime:          util.InvalidTime,
		EndTime:            util.InvalidTime,
		Index:              req.PageIndex,
		Size:               req.PageSize,
		OptUserId:          snowflake.MustParseString(optUserID),
	}

	if !util.IsBlank(req.StartTime) {
		startTime, b := util.StringToTime(req.StartTime)
		if b {
			selectReq.StartTime = *startTime
		}
	}

	if !util.IsBlank(req.EndTime) {
		endTime, b := util.StringToTime(req.EndTime)
		if b {
			selectReq.EndTime = *endTime
		}
	}
	logger.Infof("parsed Request: %v", selectReq)
	total, cashVouchers, err := dao.CashVoucher.SelectCashVouchers(ctx, sess, selectReq)
	if err != nil {
		logger.Warnf("cashVoucher.ListCashVoucher.query db err: %v", err)
		return nil, err
	}

	vouchersResp := make([]*v20230530.CashVoucher, 0, len(cashVouchers))
	for _, voucher := range cashVouchers {
		voucherData := util.ModelToOpenApiCashVoucherDetail(voucher)
		vouchersResp = append(vouchersResp, voucherData)
	}

	return &list.Data{
		Total:        total,
		CashVouchers: vouchersResp,
	}, nil
}

func AvailabilityModify(ctx context.Context, req *availabilityModify.Request, optUserID string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("cashVoucher.AvailabilityModify params: %+v, optUserID: %s", req, optUserID)
	sess := boot.MW.DefaultSession(ctx)
	defer sess.Close()
	ID := snowflake.MustParseString(req.CashVoucherID)
	cashVoucher, err := dao.CashVoucher.Get(ctx, sess, ID, false)
	if err != nil {
		logger.Warnf("cashVoucher.AvailabilityModify.queryCash.error, err: %v", err)
		return err
	}

	_, status := voucherstatus.ValidAvailabilityStatusString(req.AvailabilityStatus)

	cashVoucher.AvailabilityStatus = int(status)
	cashVoucher.UpdateTime = time.Now()
	_, err = dao.CashVoucher.Edit(ctx, sess, cashVoucher)

	if err != nil {
		logger.Warnf("cashVoucher.AvailabilityModify.update.error, err: %v", err)
		return err
	}

	return nil
}
