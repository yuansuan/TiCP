package dao

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	promodel "github.com/prometheus/common/model"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type ReportDaoImpl struct {
}

func NewReportDao() ReportDao {
	return &ReportDaoImpl{}
}

func (s *ReportDaoImpl) GetHostResourceMetricAvgUT(ctx context.Context, reportType, prefix string, timeRange *dto.TimeRange) ([]*dto.Value, error) {
	var queryString string
	if reportType == consts.LicenseAppUsedUtAvg {
		queryString = genLicenseAppUsedQueryString(prefix, timeRange.TimeStep)
	} else {
		queryString = genHostResourceQueryString(reportType, prefix, timeRange.TimeStep)
	}

	resp, err := util.PromRangeQuery(queryString, timeRange.StartTime/1000, timeRange.EndTime/1000, timeRange.TimeStep)
	if err != nil {
		return nil, err
	}

	return parsePromValToReportMetrics(resp.(promodel.Matrix)), nil
}

func parsePromValToReportMetrics(resAvgUTValues promodel.Matrix) []*dto.Value {
	hostResMetricAvgUTRecords := make([]*dto.Value, 0)
	if len(resAvgUTValues) <= 0 {
		return hostResMetricAvgUTRecords
	}

	timeValuesMap := make(map[int64][]float64)
	for _, promValue := range resAvgUTValues {
		if promValue.Values != nil && len(promValue.Values) > 0 {
			for _, value := range promValue.Values {
				timeStamp := value.Timestamp.Unix()
				float64Value := float64(value.Value)

				if _, ok := timeValuesMap[timeStamp]; ok {
					timeValuesMap[timeStamp] = append(timeValuesMap[timeStamp], float64Value)
				} else {
					timeValuesMap[timeStamp] = []float64{float64Value}
				}
			}
		}
	}

	for timeStamp, values := range timeValuesMap {
		hostResourceMetricAvgUT := &dto.Value{
			T: timeStamp * 1000,
			V: util.CalculateAvg(values),
		}
		hostResMetricAvgUTRecords = append(hostResMetricAvgUTRecords, hostResourceMetricAvgUT)
	}

	// 按照时间进行升序排序
	if len(hostResMetricAvgUTRecords) >= 2 {
		sort.Slice(hostResMetricAvgUTRecords, func(i, j int) bool {
			return hostResMetricAvgUTRecords[i].T < hostResMetricAvgUTRecords[j].T
		})
	}

	return hostResMetricAvgUTRecords
}

func genHostResourceQueryString(reportType string, nodePrefix string, timeStep int64) string {
	// 屏蔽主机配置
	excludeHost := config.GetConfig().HiddenNode
	var excludeHostLabels string
	if strutil.IsNotEmpty(excludeHost) {
		excludeHosts := strings.Split(excludeHost, ",")
		if len(excludeHosts) > 0 {
			labelList := make([]string, 0, len(excludeHosts))
			for _, host := range excludeHosts {
				labelList = append(labelList, fmt.Sprintf("%s!=\"%s\"", consts.HostNameLabel, host))
			}
			excludeHostLabels = strings.Join(labelList, ",") + ","
		}
	}

	// 节点名称过滤
	var nodePrefixQueryStr string
	if strutil.IsNotEmpty(nodePrefix) {
		nodePrefixQueryStr = fmt.Sprintf(" %v=~\"%v.*\",", consts.HostNameLabel, nodePrefix)
	}

	switch reportType {
	case consts.CPUUtAvg:
		return fmt.Sprintf("max_over_time(%s{%s %s %s=\"%s\"}[%vs])", metrics.CPUMetrics, excludeHostLabels, nodePrefixQueryStr, consts.NameLabel, metrics.CPUPercent, timeStep)
	case consts.MemUtAvg:
		return fmt.Sprintf("max_over_time(%s{%s %s %s=\"%s\"}[%vs])", metrics.MemoryMetrics, excludeHostLabels, nodePrefixQueryStr, consts.NameLabel, metrics.MemoryPercent, timeStep)
	case consts.TotalIoUtAvg:
		return fmt.Sprintf("max_over_time(%s{%s %s %s=\"%s\"}[%vs])", metrics.DiskMetrics, excludeHostLabels, nodePrefixQueryStr, consts.NameLabel, metrics.DiskTotalThroughput, timeStep)
	case consts.ReadIoUtAvg:
		return fmt.Sprintf("max_over_time(%s{%s %s %s=\"%s\"}[%vs])", metrics.DiskMetrics, excludeHostLabels, nodePrefixQueryStr, consts.NameLabel, metrics.DiskReadThroughput, timeStep)
	case consts.WriteIoUtAvg:
		return fmt.Sprintf("max_over_time(%s{%s %s %s=\"%s\"}[%vs])", metrics.DiskMetrics, excludeHostLabels, nodePrefixQueryStr, consts.NameLabel, metrics.DiskWriteThroughput, timeStep)
	default:
		return ""
	}
}

func genLicenseAppUsedQueryString(prefix string, timeStep int64) string {
	return fmt.Sprintf("avg_over_time(%s{%s=\"%s\", %s=\"%s\"}[%vs])", metrics.LicenseMonitorFeature, consts.ValueTypeLabel, metrics.FeatureUsagePercent, consts.AppTypeLabel, prefix, timeStep)
}

func (s *ReportDaoImpl) GetLicenseAppModuleUsedUtMetric(ctx context.Context, appType, featureName, licenseID string, timeRange *dto.TimeRange) ([]*dto.Value, error) {
	queryString := fmt.Sprintf("avg_over_time(%s{%s=\"%s\", %s=\"%s\", %s=\"%s\", %s=\"%s\"}[%vs])",
		metrics.LicenseMonitorFeature, consts.AppTypeLabel, appType, consts.FeatureName, featureName,
		consts.LicenseID, licenseID, consts.ValueTypeLabel, metrics.FeatureUsagePercent,
		timeRange.TimeStep)
	resp, err := util.PromRangeQuery(queryString, timeRange.StartTime/1000, timeRange.EndTime/1000, timeRange.TimeStep)
	if err != nil {
		return nil, err
	}

	return parsePromValToReportMetrics(resp.(promodel.Matrix)), nil
}

func (s *ReportDaoImpl) GetDiskUsageUtMetric(ctx context.Context, timeRange *dto.TimeRange) ([]*dto.UtAvgMetric, error) {
	queryString := fmt.Sprintf("avg_over_time(%s{%s=\"%s\"}[%vs])", metrics.DiskUsage, consts.NameLabel, metrics.DiskUsagePercent, timeRange.TimeStep)
	resp, err := util.PromRangeQuery(queryString, timeRange.StartTime/1000, timeRange.EndTime/1000, timeRange.TimeStep)
	if err != nil {
		return nil, err
	}

	resAvgUTValues := resp.(promodel.Matrix)
	metricAvgRecords := make([]*dto.UtAvgMetric, 0)
	if len(resAvgUTValues) <= 0 {
		return metricAvgRecords, nil
	}

	amountPathMap := make(map[string][]*promodel.SampleStream)
	for _, promValue := range resAvgUTValues {
		if promValue.Values != nil && len(promValue.Values) > 0 {
			amountPath := promValue.Metric[consts.MountPathLabel]
			amountPathStr := string(amountPath)
			// 无挂载路径记录为无效数据
			if strutil.IsEmpty(amountPathStr) {
				continue
			}

			if _, ok := amountPathMap[amountPathStr]; ok {
				amountPathMap[amountPathStr] = append(amountPathMap[amountPathStr], promValue)
			} else {
				amountPathMap[amountPathStr] = []*promodel.SampleStream{promValue}
			}
		}
	}

	for amountPath, promValueArray := range amountPathMap {
		timeValuesMap := make(map[int64][]float64)
		for _, promValue := range promValueArray {
			if promValue.Values != nil && len(promValue.Values) > 0 {
				for _, value := range promValue.Values {
					timeStamp := value.Timestamp.Unix()
					float64Value := float64(value.Value)

					if _, ok := timeValuesMap[timeStamp]; ok {
						timeValuesMap[timeStamp] = append(timeValuesMap[timeStamp], float64Value)
					} else {
						timeValuesMap[timeStamp] = []float64{float64Value}
					}
				}
			}
		}

		if len(timeValuesMap) > 0 {
			metricAvgUtValues := make([]*dto.MetricTV, 0)
			for timeStamp, values := range timeValuesMap {
				metricAvgUTValue := &dto.MetricTV{
					Timestamp: timeStamp * 1000,
					Value:     util.CalculateAvg(values),
				}
				metricAvgUtValues = append(metricAvgUtValues, metricAvgUTValue)
			}

			// 按照时间进行升序排序
			if len(metricAvgUtValues) >= 2 {
				sort.Slice(metricAvgUtValues, func(i, j int) bool {
					return metricAvgUtValues[i].Timestamp < metricAvgUtValues[j].Timestamp
				})
			}

			vs := &dto.UtAvgMetric{
				Name:    amountPath,
				Metrics: metricAvgUtValues,
			}

			metricAvgRecords = append(metricAvgRecords, vs)
		}
	}

	return metricAvgRecords, nil
}

func (s *ReportDaoImpl) GetNodeAvailableMetic(ctx context.Context, prefix string, timeRange *dto.TimeRange) (map[string][]promodel.SamplePair, error) {
	queryString := fmt.Sprintf("avg_over_time(%s{%s=\"%s\"}[%vs])", metrics.CollectorScrapeSuccess, consts.NameLabel, prefix, timeRange.TimeStep)
	resp, err := util.PromRangeQuery(queryString, timeRange.StartTime/1000, timeRange.EndTime/1000, timeRange.TimeStep)
	if err != nil {
		return nil, err
	}

	// 时间戳序列映射
	timeStepList, timeStepMap := getTiemStepListAndMap(timeRange.StartTime, timeRange.EndTime, timeRange.TimeStep)

	hostNameMetricMap := make(map[string][]promodel.SamplePair)
	for _, matrix := range resp.(promodel.Matrix) {
		hostname := matrix.Metric[consts.HostNameLabel]
		if hostname == "" {
			continue
		}

		timestampMetricMap := make(map[int64]promodel.SamplePair)
		for _, value := range matrix.Values {
			timestampMetricMap[int64(value.Timestamp)] = value
		}

		timestampMetric := make([]promodel.SamplePair, 0, len(timeStepMap))
		for _, timeStamp := range timeStepList {
			if _, ok := timestampMetricMap[timeStamp]; ok {
				timestampMetric = append(timestampMetric, timestampMetricMap[timeStamp])
			} else {
				timestampMetric = append(timestampMetric, promodel.SamplePair{
					Timestamp: promodel.Time(timeStamp),
					Value:     promodel.SampleValue(0),
				})
			}
		}
		hostNameMetricMap[string(hostname)] = timestampMetric
	}

	return hostNameMetricMap, nil
}

func getTiemStepListAndMap(startTime, endTime, timeStep int64) ([]int64, map[int64]struct{}) {
	timeStepList := make([]int64, 0)
	timeStepMap := make(map[int64]struct{})
	if startTime > endTime {
		return timeStepList, timeStepMap
	}

	startDateTime := time.UnixMilli(startTime).Truncate(1 * time.Second)
	endDateTime := time.UnixMilli(endTime).Truncate(1 * time.Second)
	for !startDateTime.After(endDateTime) {
		timeStepMap[startDateTime.UnixMilli()] = struct{}{}
		startDateTime = startDateTime.Add(time.Duration(timeStep) * time.Second)
	}

	for timeStamp := range timeStepMap {
		timeStepList = append(timeStepList, timeStamp)
	}
	sort.Slice(timeStepList, func(i, j int) bool {
		return timeStepList[i] < timeStepList[j]
	})

	return timeStepList, timeStepMap
}
