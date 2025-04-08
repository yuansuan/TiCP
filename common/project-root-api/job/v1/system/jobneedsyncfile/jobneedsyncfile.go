package jobneedsyncfile

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// Request 请求
// swagger:model JobNeedSyncFilesRequest
type Request struct {
	Zone       string `form:"Zone"`
	PageOffset *int64 `form:"PageOffset"`
	PageSize   *int64 `form:"PageSize"`
}

type Response struct {
	schema.Response `json:",inline"`

	Data *Data `json:"Data,omitempty"`
}

type Data struct {
	Jobs  []*NeedSyncFileJobInfo `json:"Jobs"`
	Total int64                  `json:"Total"`
}

type NeedSyncFileJobInfo struct {
	ID                      string `json:"ID"`
	State                   string `json:"State"`
	Name                    string `json:"Name"`
	FileSyncState           string `json:"FileSyncState"`
	WorkDir                 string `json:"WorkDir"`
	OutputDir               string `json:"OutputDir"`
	NoNeededPaths           string `json:"NoNeededPaths"`
	NeededPaths             string `json:"NeededPaths"`
	FileOutputStorageZone   string `json:"FileOutputStorageZone"`
	DownloadFileSizeTotal   int64  `json:"DownloadFileSizeTotal"`
	DownloadFileSizeCurrent int64  `json:"DownloadFileSizeCurrent"`
}
