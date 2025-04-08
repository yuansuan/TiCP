package add

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	CashVoucherName string `json:"CashVoucherName" form:"CashVoucherName"`
	Amount          int64  `json:"Amount" form:"Amount"`
	ExpiredType     int64  `json:"ExpiredType" form:"ExpiredType"`
	AbsExpiredTime  string `json:"AbsExpiredTime,omitempty" form:"AbsExpiredTime,omitempty"`
	RelExpiredTime  int64  `json:"RelExpiredTime,omitempty" form:"RelExpiredTime,omitempty"`
	Comment         string `json:"Comment,omitempty" form:"Comment,omitempty"`
}

type Data struct {
	CashVoucherID string `json:"CashVoucherID"`
}

type Response struct {
	v20230530.Response `json:",inline"`

	Data *Data `json:"data,omitempty"`
}
