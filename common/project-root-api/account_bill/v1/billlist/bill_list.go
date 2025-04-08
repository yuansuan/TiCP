package billlist

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	AccountID   string `form:"AccountID" uri:"AccountID" xquery:"AccountID"`
	StartTime   string `form:"StartTime" json:"StartTime" xquery:"StartTime"`
	EndTime     string `form:"EndTime" json:"EndTime" xquery:"EndTime"`
	TradeType   int64  `form:"TradeType" xquery:"TradeType"`
	SignType    int64  `form:"SignType" xquery:"SignType"`
	ProductName string `json:"ProductName" xquery:"ProductName"`
	SortByAsc   bool   `json:"SortByAsc" xquery:"SortByAsc"`
	PageIndex   int64  `form:"PageIndex" xquery:"PageIndex"`
	PageSize    int64  `form:"PageSize" xquery:"PageSize"`
}

type Data struct {
	AccountBills []*v20230530.BillListData `json:"AccountBills"`
	Total        int64                     `json:"Total"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
