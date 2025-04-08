package accountlist

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID    string `form:"AccountID" uri:"AccountID" xquery:"AccountID"`
	AccountName  string `form:"AccountName" xquery:"AccountName"`
	CustomerID   string `form:"CustomerID" xquery:"CustomerID"`
	FrozenStatus *bool  `form:"FrozenStatus" xquery:"FrozenStatus"`
	PageIndex    int64  `form:"PageIndex" xquery:"PageIndex"`
	PageSize     int64  `form:"PageSize" xquery:"PageSize"`
}

type Data struct {
	Accounts []*v20230530.AccountDetail `json:"Account"`
	Total    int64                      `json:"Total"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
