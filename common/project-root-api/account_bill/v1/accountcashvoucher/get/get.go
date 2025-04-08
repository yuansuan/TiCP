package get

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountCashVoucherID string `form:"AccountCashVoucherID" xquery:"AccountCashVoucherID"`
}

type Data struct {
	*v20230530.AccountCashVoucher `json:",inline"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
