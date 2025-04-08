package list

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID string `form:"AccountID" xquery:"AccountID"`
	StartTime string `form:"StartTime" xquery:"StartTime"`
	EndTime   string `form:"EndTime" xquery:"EndTime"`
	PageIndex int64  `form:"PageIndex" xquery:"PageIndex"`
	PageSize  int64  `form:"PageSize" xquery:"PageSize"`
}

type Data struct {
	AccountCashVouchers []*v20230530.AccountCashVoucher `json:"AccountCashVoucher"`
	Total               int64                           `json:"Total"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
