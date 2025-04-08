package get

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CashVoucherID string `uri:"CashVoucherID" json:"CashVoucherID" form:"CashVoucherID"`
}

type Data struct {
	v20230530.CashVoucher `json:",inline"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
