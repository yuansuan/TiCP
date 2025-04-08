package collector

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
)

type nodeCollector struct {
	metric  GaugeDesc
	logger  log.Logger
	nodeDao dao.NodeDao
}

func init() {
	Register(consts.Scheduler, NewNodeCollector)
}

func NewNodeCollector(logger log.Logger) (Collector, error) {
	return &nodeCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.Scheduler, "scheduler state info", []string{"name", "node_name"}),
		},
		logger:  logger,
		nodeDao: dao.NewNodeDao(),
	}, nil
}

func (c *nodeCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	nodeList, err := c.nodeDao.NodeList(context.Background())
	if err != nil {
		logging.Default().Errorf("err: %v", err)
		return err
	}
	for _, item := range nodeList {
		ch <- c.metric.MustNewConstMetric(statusReflect(item.Status), metrics.SchedulerStatus, item.NodeName)
	}
	return nil
}
func statusReflect(status string) float64 {
	if status == consts.Downed {
		return consts.NodeAbnormal //节点异常
	}

	return consts.NodeNormal //节点正常
}
