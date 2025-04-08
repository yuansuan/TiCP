package accountbilltype

type AccountBillTradeType int32

const (
	AccountBillTradeUnknow AccountBillTradeType = 0
	// 支付
	AccountBillTradePay AccountBillTradeType = 1
	// 充值
	AccountBillTradeCredit AccountBillTradeType = 2
	// 退款
	AccountBillTradeRefund AccountBillTradeType = 3
	// 提现
	AccountBillTradeWithdraw AccountBillTradeType = 4
	// 加款（管理员）
	AccountBillTradeFundAdd AccountBillTradeType = 5
	// 扣款（管理员）
	AccountBillTradeFundSub AccountBillTradeType = 6
)

// 收支类型（其中 冻结、解冻属于中间态，不给用户展示）
type AccountBillSign int32

const (
	AccountBillUnknow AccountBillSign = 0
	// 收入
	AccountBillSignAdd AccountBillSign = 1
	// 支出
	AccountBillSignReduce AccountBillSign = 2
	// 冻结
	AccountBillSignFreeze AccountBillSign = 3
	// 解冻
	AccountBillSignUnfreeze AccountBillSign = 4
)
