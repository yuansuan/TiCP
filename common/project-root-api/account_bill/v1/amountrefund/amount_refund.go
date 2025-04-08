package amountrefund

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID    string `json:"AccountID" uri:"AccountID"`
	Amount       int64  `json:"Amount"`
	RefundID     string `json:"RefundID"`
	ResourceID   string `json:"ResourceID"`
	IdempotentID string `json:"IdempotentID"`
	Comment      string `json:"Comment"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	*v20230530.Account `json:",inline"`
}
