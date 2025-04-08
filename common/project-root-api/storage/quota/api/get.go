package api

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type GetStorageQuotaResponse struct {
	v20230530.Response `json:",inline"`
	Data               *GetStorageQuotaResponseData `json:"Data,omitempty"`
}

type GetStorageQuotaResponseData struct {
	StorageUsage float64 `json:"StorageUsage,omitempty"`
	StorageLimit float64 `json:"StorageLimit,omitempty"`
}
