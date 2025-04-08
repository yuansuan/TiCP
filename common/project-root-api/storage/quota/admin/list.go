package admin

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type ListStorageQuotaRequest struct {
	PageOffset int `form:"PageOffset" json:"PageOffset" query:"PageOffset"`
	PageSize   int `form:"PageSize" json:"PageSize" query:"PageSize"`
}

type ListStorageQuotaResponse struct {
	v20230530.Response `json:",inline"`
	Data               *ListStorageQuotaResponseData `json:"Data,omitempty"`
}

type ListStorageQuotaResponseData struct {
	StorageQuota []*StorageQuota `json:"StorageQuota"`
	NextMarker   int64           `json:"NextMarker"`
	Total        int64           `json:"Total"`
}

type StorageQuota struct {
	UserID       string  `json:"UserID"`
	StorageUsage float64 `json:"StorageUsage"`
	StorageLimit float64 `json:"StorageLimit"`
}
