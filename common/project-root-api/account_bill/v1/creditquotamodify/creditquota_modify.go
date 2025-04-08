package creditquotamodify

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	// 资金账户ID
	AccountID string `json:"AccountID" uri:"AccountID"`
	// 授信额度金额
	CreditQuotaAmount int64 `json:"CreditQuotaAmount"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}

type Data struct {
	// 资金账户ID
	AccountID string `json:"AccountID"`
	// 当前授信额度金额
	CreditQuotaAmount int64 `json:"CreditQuotaAmount"`
}
