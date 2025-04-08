package handler_rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/account_bill/scheduler"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/leader"
	"time"
)

// InitGRPCServer InitGRPCServer
func InitGRPCServer(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	util.ChkErr(err)
	_ = s

	cashVoucherDaoImpl := dao.NewCashVoucherDaoImpl()
	accountCashVoucherDaoImpl := dao.NewAccountCashVoucherDaoImpl()
	cashVoucherLog := dao.NewAccountCashVoucherLogDaoImpl()
	logger := logging.Default()
	scheduler := scheduler.NewScheduler(logger, cashVoucherDaoImpl, accountCashVoucherDaoImpl, cashVoucherLog)

	leader.Runner(ctx, "cash_voucher_scheduler", scheduler.CashVoucherRun, leader.SetInterval(1*time.Second))

	leader.Runner(ctx, "account_cash_voucher_scheduler", scheduler.AccountCashVoucherRun, leader.SetInterval(1*time.Second))

}
