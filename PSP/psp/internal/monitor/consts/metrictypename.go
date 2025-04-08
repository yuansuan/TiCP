package consts

var MetricTypeNameMap = make(map[string]string)

func init() {
	MetricTypeNameMap[CPUUtAvg] = CPUUtAvgName
	MetricTypeNameMap[MemUtAvg] = MemUtAvgName
	MetricTypeNameMap[TotalIoUtAvg] = TotalIoUtAvgName
	MetricTypeNameMap[CPUTimeSum] = CPUTimeSumName
	MetricTypeNameMap[DiskUtAvg] = DiskUtAvgName
}
