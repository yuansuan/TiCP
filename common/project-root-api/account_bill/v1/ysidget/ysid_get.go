package ysidget

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type GetAccountReply struct {
	// 账户ID
	AccountID string `json:"AccountID"`
	// 普通余额
	NormalBalance int64 `json:"NormalBalance"`
}

type Request struct {
	UserID string `form:"UserID" xquery:"UserID" json:"UserID" uri:"UserID"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	*v20230530.AccountDetail `json:",inline"`
}
