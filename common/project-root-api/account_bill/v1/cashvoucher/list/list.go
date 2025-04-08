package list

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CashVoucherID      string `form:"CashVoucherID,omitempty" xquery:"CashVoucherID,omitempty"`
	CashVoucherName    string `form:"CashVoucherName,omitempty" xquery:"CashVoucherName,omitempty"`
	IsExpired          int64  `form:"IsExpired,omitempty" xquery:"IsExpired,omitempty"`
	AvailabilityStatus string `form:"AvailabilityStatus,omitempty" xquery:"AvailabilityStatus,omitempty"`
	StartTime          string `form:"StartTime" xquery:"StartTime"`
	EndTime            string `form:"EndTime" xquery:"EndTime"`
	PageIndex          int64  `form:"PageIndex" xquery:"PageIndex"`
	PageSize           int64  `form:"PageSize" xquery:"PageSize"`
}

type Data struct {
	CashVouchers []*v20230530.CashVoucher `json:"CashVouchers"`
	Total        int64                    `json:"Total"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
