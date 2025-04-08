package frozenmodify

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type Request struct {
	AccountID   string `json:"AccountID" uri:"AccountID" binding:"required"` // 资金账户ID
	FrozenState bool   `json:"FrozenState"`                                  // 冻结状态  true:冻结   false: 解冻
}

type Response struct {
	v20230530.Response `json:",inline"`
}
