package consts

// ProductNameType 云产品名称
type ProductNameType string

const (
	// ProductNameCloudCompute 云计算产品名称
	ProductNameCloudCompute ProductNameType = "CloudCompute"
	// ProductNameCloudApp 云应用产品名
	ProductNameCloudApp ProductNameType = "CloudApp"
)

// AccountBillTradeType 账单类型
type AccountBillTradeType int64

const (
	// AccountBillTradePay 账单支付
	AccountBillTradePay AccountBillTradeType = 1
	// AccountBillTradeCredit 账单充值
	AccountBillTradeCredit AccountBillTradeType = 2
	// AccountBillTradeRefund 账单退款
	AccountBillTradeRefund AccountBillTradeType = 3
)

// AccountBillSignType 账单收支类型
type AccountBillSignType int64

const (
	// AccountBillSignAdd 收入
	AccountBillSignAdd AccountBillSignType = 1
	// AccountBillSignReduce 支出
	AccountBillSignReduce AccountBillSignType = 2
)
