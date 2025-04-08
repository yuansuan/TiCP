package admin

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type GetStorageQuotaTotalResponse struct {
	v20230530.Response `json:",inline"`
	Data               *GetStorageQuotaTotalResponseData `json:"Data,omitempty"`
}

type GetStorageQuotaTotalResponseData struct {
	StorageUsage float64 `json:"StorageUsage,omitempty"`
	StorageLimit float64 `json:"StorageLimit,omitempty"`
}
