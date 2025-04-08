package v20230530

import "time"

type AccountBill struct {
	Id        string    `json:"id"`        // 收支记录ID
	AccountId string    `json:"AccountId"` // 资金账户ID
	TradeType int64     `json:"TradeType"` // 交易类型: 1支付；2充值；3退款；4提现
	Amount    int64     `json:"Amount"`    // 金额
	TradeTime time.Time `json:"TradeTime"` // 交易时间
	Comment   string    `json:"Comment"`   // 交易备注
}

type BillListData struct {
	ID                    string  `json:"ID"`                    // 收支记录ID
	AccountID             string  `json:"AccountID"`             // 资金账户ID
	TradeType             int64   `json:"TradeType"`             // 交易类型: 1支付；2充值；3退款；4提现
	SignType              int64   `json:"SignType"`              // 收支类型: 1收入；2支出；
	Amount                int64   `json:"Amount"`                // 金额
	AccountBalance        int64   `json:"AccountBalance"`        // 账户余额
	FreezedAmount         int64   `json:"FreezedAmount"`         // 冻结金额
	DeltaNormalBalance    int64   `json:"DeltaNormalBalance"`    // 普通余额操作金额
	DeltaAwardBalance     int64   `json:"DeltaAwardBalance"`     // 赠送余额操作金额
	DeltaDeductionBalance int64   `json:"DeltaDeductionBalance"` // 抵扣金额，比如代金券
	DeltaDiscountBalance  int64   `json:"DeltaDiscountBalance"`  // 折扣金额
	TradeID               string  `json:"TradeID"`               // 交易单ID
	TradeTime             string  `json:"TradeTime"`             // 交易时间
	Comment               string  `json:"Comment"`               // 交易备注
	MerchandiseID         string  `json:"MerchandiseID"`         // 商品ID
	MerchandiseName       string  `json:"MerchandiseName"`       // 商品名称
	UnitPrice             int64   `json:"UnitPrice"`             // 单价
	PriceDes              string  `json:"PriceDes"`              // 单价单位描述， 如 元
	Quantity              float64 `json:"Quantity"`              // 消耗数量
	QuantityUnit          string  `json:"QuantityUnit"`          // 消耗数量单位描述，如 核时
	ResourceID            string  `json:"ResourceID"`            // 资源id(作业计算为作业ID，其余业务类似为业务主键ID)
	StartTime             string  `json:"StartTime"`             // 扣费周期开始时间，按量付费使用
	EndTime               string  `json:"EndTime"`               // 扣费周期结束时间，按量付费使用
	ProductName           string  `json:"ProductName"`           // 产品类型 CloudCompute: pass求解作业，CloudApp : 云应用
}

type Account struct {
	AccountID         string `json:"AccountID"`
	AccountName       string `json:"AccountName"`
	AccountBalance    int64  `json:"AccountBalance"`
	NormalBalance     int64  `json:"NormalBalance"`
	AwardBalance      int64  `json:"AwardBalance"`
	FreezedAmount     int64  `json:"FreezedAmount"`
	CashVoucherAmount int64  `json:"CashVoucherAmount"`
	CreditQuotaAmount int64  `json:"CreditQuotaAmount"`
}

type AccountDetail struct {
	AccountID         string `json:"AccountID"`
	CustomerID        string `json:"CustomerID"`
	Currency          string `json:"Currency"`
	AccountName       string `json:"AccountName"`
	AccountBalance    int64  `json:"AccountBalance"`
	NormalBalance     int64  `json:"NormalBalance"`
	AwardBalance      int64  `json:"AwardBalance"`
	CreditQuotaAmount int64  `json:"CreditQuotaAmount"`
	FreezedAmount     int64  `json:"FreezedAmount"`
	CashVoucherAmount int64  `json:"CashVoucherAmount"`
	FrozenStatus      bool   `json:"FrozenStatus"`
	CreateTime        string `json:"CreateTime"`
	UpdateTime        string `json:"UpdateTime"`
	IsOverdrawn       bool   `json:"IsOverdrawn"`
}

type CashVoucher struct {
	CashVoucherID      string `json:"CashVoucherID"`
	CashVoucherName    string `json:"CashVoucherName"`
	AvailabilityStatus string `json:"AvailabilityStatus"`
	Amount             int64  `json:"Amount"`
	IsExpired          string `json:"IsExpired"`
	ExpiredType        int    `json:"ExpiredType"`
	AbsExpiredTime     string `json:"AbsExpiredTime"`
	RelExpiredTime     int64  `json:"RelExpiredTime"`
	Comment            string `json:"Comment"`
	CreateTime         string `json:"CreateTime"`
}

type AccountCashVoucher struct {
	AccountCashVoucherID string `json:"AccountCashVoucherID"`
	AccountID            string `json:"AccountID"`
	CashVoucherID        string `json:"CashVoucherID"`
	Amount               int64  `json:"Amount"`
	UsedAmount           int64  `json:"UsedAmount"`
	RemainingAmount      int64  `json:"RemainingAmount"`
	Status               int    `json:"Status"`
	ExpiredTime          string `json:"ExpiredTime"`
	IsExpired            int    `json:"IsExpired"`
	CreateTime           string `json:"CreateTime"`
}
