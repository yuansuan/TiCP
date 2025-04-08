package dto

type UniteReportReq struct {
	ReportType string `form:"type"`  // 报表类型
	Start      int64  `form:"start"` // 开始时间，毫秒单位
	End        int64  `form:"end"`   // 结束时间, 毫秒单位
}

type ResourceUtAvgReportResp struct {
	CPUUtAvg     []*UtAvgMetric `json:"cpu_ut_avg"`      // cpu 平均利用率
	MemUtAvg     []*UtAvgMetric `json:"mem_ut_avg"`      // 内存平均利用率
	TotalIOUtAvg []*UtAvgMetric `json:"total_io_ut_avg"` // 磁盘吞吐率
}

type DiskUtAvgReportResp struct {
	DiskUtAvg []*UtAvgMetric `json:"disk_ut_avg"` // 磁盘吞吐率
}

type LicenseAppUsedUtAvgReportResp struct {
	LicenseAppUtAvg []*UtAvgMetric `json:"license_app_ut_avg"` // App license 使用率
}

type LicenseAppModuleUsedUtAvgReportResp struct {
	LicenseAppModuleUtAvg []*UtAvgMetric `json:"license_app_module_ut_avg"` // App module license 使用率
}

type UtAvgMetric struct {
	Name    string      `json:"n"` // 指标名称
	Metrics []*MetricTV `json:"d"` // 指标数值
}

type MetricTV struct {
	Timestamp int64   `json:"t"` // 时间戳
	Value     float64 `json:"v"` // 指标值
}

type NodeTypeReportRequest struct {
	NodeLabel  string
	NodePrefix string
	ReportType string
	StartTime  int64
	EndTime    int64
	TimeStep   int64
}

type TimeRange struct {
	StartTime int64
	EndTime   int64
	TimeStep  int64
}

type CPUTimeSumMetricsResp struct {
	CPUTimeByApp  *ReportMetricKVData `json:"cpu_time_by_app"`  // app 使用指标值
	CPUTimeByUser *ReportMetricKVData `json:"cpu_time_by_user"` // 用户指标使用值
}

type JobCountMetricResp struct {
	JobCountByApp  *ReportMetricKVData `json:"job_count_by_app"`  // app 使用指标
	JobCountByUser *ReportMetricKVData `json:"job_count_by_user"` // 用户使用指标
}

type JobDeliverCountResp struct {
	JobDeliverUserCount []*UtAvgMetric `json:"job_deliver_user_count"` // 用户提交作业数
	JobDeliverJobCount  []*UtAvgMetric `json:"job_deliver_job_count"`  // 作业提交数量
}

type JobWaitStatisticResp struct {
	JobWaitTimeStatistic []*UtAvgMetric `json:"job_wait_time_statistic"` // 作业等待时间情况（按天）
	JobWaitNumStatistic  []*UtAvgMetric `json:"job_wait_num_statistic"`  // 作业等待人次情况（按天）
}

type ReportMetricKVData struct {
	Name         string      `json:"name"`          // 名称
	OriginalData []*MetricKV `json:"original_data"` // 数据
}

type MetricKV struct {
	Key   string  `json:"key"`   // 指标key
	Value float64 `json:"value"` // 指标value
}

type LicenseAppModuleUsedUtAvgReq struct {
	LicenseId   string `form:"license_id"`   // module Config 对应的 license id
	LicenseType string `form:"license_type"` // license 对应的license type
	UniteReportReq
}

type NodeDownStatisticReportReq struct {
	UniteReportReq
}

type StatisticItem struct {
	Key   string  `json:"key"`   // 键
	Value float64 `json:"value"` // 值
}

type OriginStatisticData struct {
	Name         string           `json:"name"`          // 名称
	OriginalData []*StatisticItem `json:"original_data"` // 原始数据
}

type NodeDownStatisticReportResp struct {
	NodeDownNumberRate *UtAvgMetric         `json:"node_down_number_rate"` // 节点宕机率
	NodeDownNumber     *OriginStatisticData `json:"node_down_number"`      // 节点宕机数
}

type NodeDownInfo struct {
	NodeName  string `json:"node_name"`  // 节点名称
	DownTime  int64  `json:"down_time"`  // 宕机时间
	DownStart int64  `json:"down_start"` // 宕机开始时间
	DownEnd   int64  `json:"down_end"`   // 宕机结束时间
}

type ExportNodeDownStatisticsReq struct {
	Start int64 `form:"start"` // 开始时间
	End   int64 `form:"end"`   // 结束时间
}

type ExportNodeDownStatisticsResp struct{}

type ReportData struct {
	DateTime int64
	Metrics  []float64
}
