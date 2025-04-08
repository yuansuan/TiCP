package dto

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

type GetAppInfoServiceRequest struct {
	ID          snowflake.ID
	OutAppID    string
	Name        string
	AppType     string
	Version     string
	ComputeType string
}

type AddAppServiceRequest struct {
	NewType           string
	NewVersion        string
	BaseName          string
	ComputeType       string
	Description       string
	Icon              string
	Image             string
	ResidualLogParser string
	CloudOutAppId     string
	EnableResidual    bool
	EnableSnapshot    bool
	Queues            []*QueueInfo
	Licenses          []*LicenseInfo
	BinPath           []*KeyValue
	SchedulerParam    []*KeyValue
}
