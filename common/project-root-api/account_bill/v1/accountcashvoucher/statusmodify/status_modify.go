package statusmodify

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountCashVoucherID string `uri:"AccountCashVoucherID" json:"AccountCashVoucherID" form:"AccountCashVoucherID"`
	Status               string `json:"Status" form:"Status"`
}

type Response struct {
	v20230530.Response `json:",inline"`
}
