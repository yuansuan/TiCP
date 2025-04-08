package add

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CashVoucherID string `uri:"CashVoucherID" json:"CashVoucherID" form:"CashVoucherID"`
	AccountIDs    string `json:"AccountIDs" form:"AccountIDs"`
}

type Data struct {
}

type Response struct {
	v20230530.Response `json:",inline"`
}
