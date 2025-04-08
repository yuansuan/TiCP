package collector

import (
	"context"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/client"
)

type jobStatusCollector struct {
	metric GaugeDesc
	logger log.Logger
	rpc    *client.GRPC
}

func init() {
	Register(consts.JobStatus, NewJobStatusCollector)
}

func NewJobStatusCollector(logger log.Logger) (Collector, error) {
	rpc := client.GetInstance()
	return &jobStatusCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.JobStatus, "job status info", []string{"name"}),
		},
		logger: logger,
		rpc:    rpc,
	}, nil
}

func (c *jobStatusCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	// 获取作业状态信息
	jobStatusMap, err := getJobStatus(c.rpc)
	if err != nil {
		return err
	}

	for key, value := range jobStatusMap {
		ch <- c.metric.MustNewConstMetric(float64(value), key)
	}
	return nil
}

// 获取作业状态信息
func getJobStatus(rpc *client.GRPC) (map[string]int64, error) {
	//1.获取作业信息
	resp, err := rpc.Job.GetJobStatus(context.Background(), &pb.GetJobStatusRequest{})
	if err != nil {
		return nil, err
	}
	return resp.JobStatusMap, nil
}
