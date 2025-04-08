package paymentfreezeunfreeze

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID string `json:"AccountID" uri:"AccountID"` // 资金账户ID
	Amount    int64  `json:"Amount"`                    // 扣减金额
	Comment   string `json:"Comment"`                   // 备注
	TradeID   string `json:"TradeID"`                   // 关联业务订单号
	IsFreezed bool   `json:"IsFreezed"`                 // 冻结状态操作
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	AccountID         string `json:"AccountID"`
	AccountName       string `json:"AccountName"`
	AccountBalance    int64  `json:"AccountBalance"`
	NormalBalance     int64  `json:"NormalBalance"`
	AwardBalance      int64  `json:"AwardBalance"`
	FreezedAmount     int64  `json:"FreezedAmount"`
	CreditQuotaAmount int64  `json:"CreditQuotaAmount"`
}
