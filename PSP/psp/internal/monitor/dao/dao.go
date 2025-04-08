package dao

import (
	"context"

	promodel "github.com/prometheus/common/model"
	"xorm.io/xorm"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type NodeDao interface {
	AddNodes(ctx context.Context, session *xorm.Session, nodes []*model.NodeInfo) error
	UpdateNode(ctx context.Context, session *xorm.Session, nodes *model.NodeInfo) error
	DeleteNotIds(ctx context.Context, session *xorm.Session, ids []snowflake.ID) error
	GetNodes(ctx context.Context, nodeName string, index int64, size int64) ([]*model.NodeInfo, int64, error)
	GetNodeByNames(ctx context.Context, nodeNames []string) ([]*model.NodeInfo, error)
	NodeList(ctx context.Context) ([]*model.NodeInfo, error)
	StatisticCoreNum(ctx context.Context, nodeNames []string) (*Statistics, error)
	QueueList(ctx context.Context) ([]string, error)
	GetQueueAvailableCores(ctx context.Context, queueNames []string) ([]*pb.QueueCore, error)
	GetPlatformCores(ctx context.Context) ([]*pb.PlatformCore, error)
	GetQueueCoreInfos(ctx context.Context) ([]*pb.QueueCoreInfo, error)
}

type ReportDao interface {
	GetHostResourceMetricAvgUT(ctx context.Context, reportType, prefix string, timeRange *dto.TimeRange) ([]*dto.Value, error)
	GetDiskUsageUtMetric(ctx context.Context, timeRange *dto.TimeRange) ([]*dto.UtAvgMetric, error)
	GetLicenseAppModuleUsedUtMetric(ctx context.Context, appType, featureName, licenseID string, timeRange *dto.TimeRange) ([]*dto.Value, error)
	GetNodeAvailableMetic(ctx context.Context, prefix string, timeRange *dto.TimeRange) (map[string][]promodel.SamplePair, error)
}
