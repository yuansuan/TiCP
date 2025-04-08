package resourcebilllist

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	StartTime   string `form:"StartTime" json:"StartTime" xquery:"StartTime"`
	EndTime     string `form:"EndTime" json:"EndTime" xquery:"EndTime"`
	TradeType   int64  `form:"TradeType" xquery:"TradeType"`
	SignType    int64  `form:"SignType" xquery:"SignType"`
	ProductName string `json:"ProductName" xquery:"ProductName"`
	SortByAsc   bool   `json:"SortByAsc" xquery:"SortByAsc"`
	PageIndex   int64  `form:"PageIndex" xquery:"PageIndex"`
	PageSize    int64  `form:"PageSize" xquery:"PageSize"`
}

type AccountResourceBillListData struct {
	AccountID            string  `json:"AccountID"`            // 资金账户ID
	TotalAmount          int64   `json:"TotalAmount"`          // 金额
	TotalFreezedAmount   int64   `json:"TotalFreezedAmount"`   // 冻结金额
	TotalNormalAmount    int64   `json:"TotalNormalAmount"`    // 普通余额操作金额
	TotalAwardAmount     int64   `json:"TotalAwardAmount"`     // 赠送余额操作金额
	TotalDeductionAmount int64   `json:"TotalDeductionAmount"` // 抵扣金额，比如代金券
	TotalDiscountAmount  int64   `json:"TotalDiscountAmount"`  // 折扣金额
	TotalRefundAmount    int64   `json:"TotalRefundAmount"`    // 退款金额
	TradeID              string  `json:"TradeID"`              // 交易单ID
	LatestTradeTime      string  `json:"LatestTradeTime"`      // 交易时间
	MerchandiseID        string  `json:"MerchandiseID"`        // 商品ID
	MerchandiseName      string  `json:"MerchandiseName"`      // 商品名称
	UnitPrice            int64   `json:"UnitPrice"`            // 单价
	PriceDes             string  `json:"PriceDes"`             // 单价单位描述， 如 元
	Quantity             float64 `json:"Quantity"`             // 消耗数量
	QuantityUnit         string  `json:"QuantityUnit"`         // 消耗数量单位描述，如 核时
	ResourceID           string  `json:"ResourceID"`           // 资源id(作业计算为作业ID，其余业务类似为业务主键ID)
	StartTime            string  `json:"StartTime"`            // 扣费周期开始时间，按量付费使用
	EndTime              string  `json:"EndTime"`              // 扣费周期结束时间，按量付费使用
	ProductName          string  `json:"ProductName"`          // 产品类型 CloudCompute: pass求解作业，CloudApp : 云应用
}

type Data struct {
	AccountBills []*AccountResourceBillListData `json:"AccountBills"`
	Total        int64                          `json:"Total"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
