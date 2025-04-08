package paymentreduce

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID             string  `json:"AccountID" uri:"AccountID" binding:"required" error:"AccountID is required"`
	Comment               string  `json:"Comment"`
	TradeID               string  `json:"TradeID"`
	MerchandiseID         string  `json:"MerchandiseID"`
	MerchandiseName       string  `json:"MerchandiseName"`
	UnitPrice             int64   `json:"UnitPrice"`
	PriceDes              string  `json:"PriceDes"`
	Quantity              float64 `json:"Quantity"`
	QuantityUnit          string  `json:"QuantityUnit"`
	ResourceID            string  `json:"ResourceID"`
	StartTime             string  `json:"StartTime"`
	EndTime               string  `json:"EndTime"`
	Ext                   string  `json:"Ext"`
	AccountCashVoucherIDs string  `json:"AccountCashVoucherIDs"`
	VoucherConsumeMode    int64   `json:"VoucherConsumeMode"`
}

type OperationAccountReply struct {
	AccountId     string `json:"AccountId"`
	NormalBalance int64  `json:"NormalBalance"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

type Data struct {
	AccountID         string `json:"AccountID"`
	AccountName       string `json:"AccountName"`
	AccountBalance    int64  `json:"AccountBalance"`
	NormalBalance     int64  `json:"NormalBalance"`
	AwardBalance      int64  `json:"AwardBalance"`
	FreezedAmount     int64  `json:"FreezedAmount"`
	CashVoucherAmount int64  `json:"CashVoucherAmount"`
	CreditQuotaAmount int64  `json:"CreditQuotaAmount"`
}
