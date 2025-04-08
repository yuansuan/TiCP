package util

import (
	"strconv"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
)

func Convert2ResourceUtAvgReport(values []*dto.Value, nodeLabel string) *dto.UtAvgMetric {
	if len(values) == 0 {
		return &dto.UtAvgMetric{Name: nodeLabel, Metrics: []*dto.MetricTV{}}
	}

	metrics := make([]*dto.MetricTV, 0, len(values))
	for _, result := range values {
		metric := &dto.MetricTV{
			Timestamp: result.T,
			Value:     result.V,
		}
		metrics = append(metrics, metric)
	}

	return &dto.UtAvgMetric{
		Name:    nodeLabel,
		Metrics: metrics,
	}
}

func Convert2CPUTimeSumResp(cpuTimeSumMetric *job.GetJobMetricResponse) *dto.CPUTimeSumMetricsResp {
	cpuTimeByApp := convert2JobMetricKV(cpuTimeSumMetric.AppMetrics, consts.CPUTimeSumUnit)
	userTimeByApp := convert2JobMetricKV(cpuTimeSumMetric.UserMetrics, consts.CPUTimeSumUnit)

	return &dto.CPUTimeSumMetricsResp{
		CPUTimeByApp:  cpuTimeByApp,
		CPUTimeByUser: userTimeByApp,
	}
}

func Convert2JobCountResp(cpuTimeSumMetric *job.GetJobMetricResponse) *dto.JobCountMetricResp {
	cpuTimeByApp := convert2JobMetricKV(cpuTimeSumMetric.AppMetrics, consts.JobCountUnit)
	userTimeByApp := convert2JobMetricKV(cpuTimeSumMetric.UserMetrics, consts.JobCountUnit)

	return &dto.JobCountMetricResp{
		JobCountByApp:  cpuTimeByApp,
		JobCountByUser: userTimeByApp,
	}
}

func Convert2JobDeliverCountResp(jobDeliverCountMetric *job.GetJobMetricResponse) *dto.JobDeliverCountResp {
	deliverUserCountMetric := convert2UtAvgMetric(jobDeliverCountMetric.GetUserMetrics(), consts.JobDeliverUserCountTitle)
	deliverJobCountMetric := convert2UtAvgMetric(jobDeliverCountMetric.GetAppMetrics(), consts.JobDeliverJobCountTitle)

	return &dto.JobDeliverCountResp{
		JobDeliverUserCount: []*dto.UtAvgMetric{deliverUserCountMetric},
		JobDeliverJobCount:  []*dto.UtAvgMetric{deliverJobCountMetric},
	}
}

func Convert2JobWaitStatisticResp(jobWaitStatistic *job.GetJobWaitTimeStatisticResponse) *dto.JobWaitStatisticResp {
	jobWaitTimeStatisticTotal := convert2UtAvgMetric(jobWaitStatistic.WaitTimeStatisticTotal, consts.JobWaitTimeStatisticTotalName)
	jobWaitTimeStatisticAvg := convert2UtAvgMetric(jobWaitStatistic.WaitTimeStatisticAvg, consts.JobWaitTimeStatisticAvgName)
	jobWaitTimeStatisticMax := convert2UtAvgMetric(jobWaitStatistic.WaitTimeStatisticMax, consts.JobWaitTimeStatisticMaxName)
	jobWaitNumStatisticTotal := convert2UtAvgMetric(jobWaitStatistic.WaitNumStatisticTotal, consts.JobWaitNumStatisticTotalName)

	return &dto.JobWaitStatisticResp{
		JobWaitTimeStatistic: []*dto.UtAvgMetric{
			jobWaitTimeStatisticAvg,
			jobWaitTimeStatisticMax,
			jobWaitTimeStatisticTotal,
		},
		JobWaitNumStatistic: []*dto.UtAvgMetric{
			jobWaitNumStatisticTotal,
		},
	}
}

func convert2UtAvgMetric(metrics []*job.MetricKV, name string) *dto.UtAvgMetric {
	metricsTVs := make([]*dto.MetricTV, 0, len(metrics))
	if len(metrics) > 0 {
		for _, metric := range metrics {
			timeStamp, err := strconv.ParseInt(metric.Key, 10, 64)
			if err != nil {
				timeStamp = 0
			}
			metricTV := &dto.MetricTV{
				Timestamp: timeStamp,
				Value:     metric.Value,
			}
			metricsTVs = append(metricsTVs, metricTV)
		}
	}

	return &dto.UtAvgMetric{
		Name:    name,
		Metrics: metricsTVs,
	}
}

func convert2JobMetricKV(jobMetricKVs []*job.MetricKV, name string) *dto.ReportMetricKVData {
	metricKVs := make([]*dto.MetricKV, 0, len(jobMetricKVs))
	if len(jobMetricKVs) > 0 {
		for _, metric := range jobMetricKVs {
			metricKV := &dto.MetricKV{
				Key:   metric.Key,
				Value: metric.Value,
			}
			metricKVs = append(metricKVs, metricKV)
		}
	}

	return &dto.ReportMetricKVData{
		Name:         name,
		OriginalData: metricKVs,
	}
}
