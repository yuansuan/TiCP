package ysidreduce

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	UserID                string  `json:"UserID" uri:"UserID"`
	IdempotentID          string  `json:"IdempotentID"`
	Amount                int64   `json:"Amount"`
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
	ProductName           string  `json:"ProductName"`
	Ext                   string  `json:"Ext"`
	Comment               string  `json:"Comment"`
	AccountCashVoucherIDs string  `json:"AccountCashVoucherIDs"`
	VoucherConsumeMode    int64   `json:"VoucherConsumeMode"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

type Data struct {
	*v20230530.Account `json:",inline"`
}
