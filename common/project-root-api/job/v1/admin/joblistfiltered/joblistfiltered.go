package joblistfiltered

import (
	list "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"
	"time"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	list.Request  `json:",inline"`
	Name          string    `form:"Name"`
	JobID         string    `form:"JobID"`
	AccountID     string    `form:"AccountID"`
	FileSyncState string    `form:"FileSyncState"`
	StartTime     time.Time `form:"StartTime"`
	EndTime       time.Time `form:"EndTime"`
}

type Response struct {
	schema.Response `json:",inline"`
	Data            *Data `json:"Data,omitempty"`
}

type Data struct {
	Jobs  []*schema.AdminJobInfo `json:"Jobs"`
	Total int64                  `json:"Total"`
}
