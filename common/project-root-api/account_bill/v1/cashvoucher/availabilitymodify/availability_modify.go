package availabilitymodify

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CashVoucherID      string `uri:"CashVoucherID" json:"CashVoucherID" form:"CashVoucherID"`
	AvailabilityStatus string `json:"AvailabilityStatus" form:"AvailabilityStatus"`
}

type Response struct {
	v20230530.Response `json:",inline"`
}
