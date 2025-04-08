package admin

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type PutStorageQuotaRequest struct {
	UserID       string  `form:"UserID" json:"UserID" query:"UserID"`
	StorageLimit float64 `form:"StorageLimit" json:"StorageLimit" query:"StorageLimit"`
}

type PutStorageQuotaResponse struct {
	v20230530.Response `json:",inline"`
}
