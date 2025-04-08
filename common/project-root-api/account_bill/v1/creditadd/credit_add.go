package creditadd

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID          string `json:"AccountID" uri:"AccountID" binding:"required"`
	DeltaAwardBalance  int64  `json:"DeltaAwardBalance"`
	DeltaNormalBalance int64  `json:"DeltaNormalBalance"`
	TradeId            string `json:"TradeId"`
	IdempotentID       string `json:"IdempotentID"`
	Comment            string `json:"Comment"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	*v20230530.Account `json:",inline"`
}

type OperationAccountReply struct {
	AccountId     string `json:"AccountId"`
	NormalBalance int64  `json:"NormalBalance"`
}
