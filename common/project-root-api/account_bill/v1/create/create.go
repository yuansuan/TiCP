package create

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountName string `json:"AccountName" binding:"required"`
	UserID      string `json:"UserID"`
	AccountType int64  `json:"AccountType" binding:"required"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	// 资金账户ID
	AccountID string `json:"AccountID"`
}
