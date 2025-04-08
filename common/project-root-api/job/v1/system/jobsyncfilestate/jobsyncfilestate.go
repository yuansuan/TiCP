package jobsyncfilestate

import schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// Request 请求
// swagger:model JobSyncFileRequest
type Request struct {
	JobID                   string `uri:"JobID" binding:"required"` // 作业ID
	DownloadFinished        bool   `json:"DownloadFinished"`        // 文件是否同步完成
	DownloadFileSizeCurrent int64  `json:"DownloadFileSizeCurrent"` // 已同步文件大小
	DownloadFileSizeTotal   int64  `json:"DownloadFileSizeTotal"`   // 需下载文件总大小
	DownloadFinishedTime    string `json:"DownloadFinishedTime"`    // 同步完成时间, RFC3339 格式
	TransmittingTime        string `json:"TransmittingTime"`        // 开始下载文件时间，RFC3339 格式
	FileSyncState           string `json:"FileSyncState"`           // 文件传输状态: Waiting,Syncing,Pausing,Paused,Resuming,Completed,Failed
}

// Response 返回
// swagger:model JobSyncFileRequest
type Response struct {
	schema.Response `json:",inline"`
}
