package scheduler

import (
	"context"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/module/rpc"
	"go.uber.org/zap"
)

// Scheduler 调度器
type Scheduler struct {
	logger                    *zap.SugaredLogger
	CashVoucherDaoImpl        *dao.CashVoucherDaoImpl
	AccountCashVoucherDaoImpl *dao.AccountCashVoucherDaoImpl
	CashVoucherLogDaoImpl     *dao.AccountCashVoucherLogDaoImpl
}

// NewScheduler 新建调度器
func NewScheduler(logger *zap.SugaredLogger, cashVoucherDao *dao.CashVoucherDaoImpl, accountCashVoucherDao *dao.AccountCashVoucherDaoImpl, cashVoucherLogDao *dao.AccountCashVoucherLogDaoImpl) *Scheduler {
	return &Scheduler{
		logger:                    logger,
		CashVoucherDaoImpl:        cashVoucherDao,
		AccountCashVoucherDaoImpl: accountCashVoucherDao,
		CashVoucherLogDaoImpl:     cashVoucherLogDao,
	}
}

func (s *Scheduler) CashVoucherRun(ctx context.Context) {
	logger := s.logger.With("func", "job.scheduler.run")
	logger.Infof("scheduler start...")

	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	//将超过有效期的代金券状态置为过期
	err := s.CashVoucherDaoImpl.UpdateStatusOfExpired(ctx, session)
	if err != nil {
		logger.Errorw("err_update_cash_voucher_status", "error", err)
		return
	}
	logger.Infof("scheduler end")
}

func (s *Scheduler) AccountCashVoucherRun(ctx context.Context) {
	logger := s.logger.With("func", "cashVoucher.scheduler.run")
	logger.Infof("scheduler start...")

	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	//1.查询超过有效期的代金券
	vouchers, err := s.AccountCashVoucherDaoImpl.QueryExpired(ctx, session)
	if err != nil {
		logger.Errorw("err  method SelectExpired ", "error", err)
		return
	}

	//2.将超过有效期的代金券状态置为过期
	expiredIds := make([]snowflake.ID, 0)
	for i := 0; i < len(vouchers); i++ {
		expiredIds = append(expiredIds, vouchers[i].Id)
	}
	_, err = s.AccountCashVoucherDaoImpl.UpdateStatusOfExpired(session, expiredIds, consts.Expired)
	if err != nil {
		logger.Errorw("err_update_cash_voucher_status", "error", err)
		return
	}

	//3.批量生成主键id
	ids, err := rpc.GetInstance().GenIDs(ctx, int64(len(vouchers)))
	if err != nil {
		logger.Warnf("idgen error! err: %v", err)
		return
	}

	//4.根据vouchers构建accountCashVoucherLogs
	var accountCashVoucherLogs []*models.AccountCashVoucherLog
	for i, voucher := range vouchers {
		nowVoucher := voucher
		nowVoucher.IsExpired = consts.Expired

		accountCashVoucherLogs = append(accountCashVoucherLogs, &models.AccountCashVoucherLog{
			Id:                   ids[i],
			AccountId:            voucher.AccountId,
			CashVoucherId:        voucher.CashVoucherId,
			AccountCashVoucherId: voucher.Id,
			SignType:             consts.VOUCHER_LOG_SIGN_EXPIRED,
			SourceInfo:           voucher.String(),
			TargetInfo:           nowVoucher.String(),
			Comment:              consts.EXPIRED_COMMENT,
			OptUserId:            voucher.OptUserId,
			CreateTime:           voucher.CreateTime,
			UpdateTime:           voucher.UpdateTime,
		})
	}

	//5.批量新增
	if len(accountCashVoucherLogs) > 0 {
		_, err = s.CashVoucherLogDaoImpl.BatchAdd(ctx, session, accountCashVoucherLogs)
		if err != nil {
			logger.Warnf("batch add error! err: %v", err)
			return
		}
	}

	logger.Infof("scheduler end")
}
