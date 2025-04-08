package impl

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	promodel "github.com/prometheus/common/model"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type ReportServiceImpl struct {
	ReportDao dao.ReportDao
	nodeDao   dao.NodeDao
	localAPI  *openapi.OpenAPI
}

func NewReportService() (service.ReportService, error) {
	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	return &ReportServiceImpl{
		ReportDao: dao.NewReportDao(),
		nodeDao:   dao.NewNodeDao(),
		localAPI:  localAPI,
	}, nil
}

func (s *ReportServiceImpl) GetHostResourceMetricUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.ResourceUtAvgReportResp, error) {
	reports, err := s.getAllNodesTypeReport(ctx, req)
	if err != nil {
		return nil, err
	}

	switch req.ReportType {
	case consts.CPUUtAvg:
		return &dto.ResourceUtAvgReportResp{CPUUtAvg: reports}, nil
	case consts.MemUtAvg:
		return &dto.ResourceUtAvgReportResp{MemUtAvg: reports}, nil
	case consts.TotalIoUtAvg:
		return &dto.ResourceUtAvgReportResp{TotalIOUtAvg: reports}, nil
	default:
		return nil, status.Error(errcode.ErrCommonReportTypeUnsupported, errcode.MsgCommonReportTypeUnsupported)
	}
}

func (s *ReportServiceImpl) getAllNodesTypeReport(ctx *gin.Context, req *dto.UniteReportReq) ([]*dto.UtAvgMetric, error) {
	var reports []*dto.UtAvgMetric
	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  CalculateStep(req.Start, req.End),
	}

	// 读取节点规则, 如果有配置，则读取单个分类统计报告
	nodeClassification := config.GetConfig().NodeClassification
	if strutil.IsNotEmpty(nodeClassification.ClassificationRule) {
		for _, node := range nodeClassification.Nodes {
			// 单个节点数据
			results, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, req.ReportType, node.ClassifyTag, timeRange)
			if err != nil {
				return nil, err
			}

			if len(results) == 0 {
				continue
			}

			singleReport := util.Convert2ResourceUtAvgReport(results, node.Label)
			reports = append(reports, singleReport)
		}
	}

	// 获取所有节点统计报告
	results, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, req.ReportType, "", timeRange)
	if err != nil {
		return nil, err
	}

	allReports := util.Convert2ResourceUtAvgReport(results, consts.ResourceUtAvgAllNodesName)
	reports = append(reports, allReports)
	return reports, nil
}

func (s *ReportServiceImpl) GetLicenseAppUsedUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.LicenseAppUsedUtAvgReportResp, error) {
	licenseAppUtAvg, err := s.getLicenseAppUsedUtAvgData(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.LicenseAppUsedUtAvgReportResp{
		LicenseAppUtAvg: licenseAppUtAvg,
	}, nil
}

func (s *ReportServiceImpl) getLicenseAppUsedUtAvgData(ctx context.Context, req *dto.UniteReportReq) ([]*dto.UtAvgMetric, error) {
	appTypes, err := s.getLicenseAppTypes()
	if err != nil {
		return nil, err
	}

	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  CalculateStep(req.Start, req.End),
	}

	licenseAppUtAvg := make([]*dto.UtAvgMetric, 0, len(appTypes))
	for _, appType := range appTypes {
		metricAvgUT, err := s.ReportDao.GetHostResourceMetricAvgUT(ctx, consts.LicenseAppUsedUtAvg, appType, timeRange)
		if err != nil {
			return nil, err
		}
		utAvgMetric := util.Convert2ResourceUtAvgReport(metricAvgUT, appType)

		// 未配置好指标值的，不返回到前端
		if len(utAvgMetric.Metrics) > 0 {
			licenseAppUtAvg = append(licenseAppUtAvg, utAvgMetric)
		}
	}

	return licenseAppUtAvg, nil
}

func (s *ReportServiceImpl) NodeDownStatisticReport(ctx *gin.Context, req *dto.NodeDownStatisticReportReq) (*dto.NodeDownStatisticReportResp, error) {
	nodeList, err := s.nodeDao.NodeList(ctx)
	if err != nil {
		return nil, err
	}

	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  getTimeStep(req.Start, req.End),
	}
	hostNameMetricMap, err := s.ReportDao.GetNodeAvailableMetic(ctx, consts.CPU, timeRange)
	if err != nil {
		return nil, err
	}

	nodeNumber, timeTimesMap, hostNameDownNumberMap := resolveNodeAvailableMetic(nodeList, hostNameMetricMap)

	nodeDownStatusStatistic := make([]*dto.Value, 0, len(timeTimesMap))
	for k, v := range timeTimesMap {
		nodeDownStatusStatistic = append(nodeDownStatusStatistic, &dto.Value{T: k, V: float64(v) / float64(nodeNumber) * 100})
	}
	sort.Slice(nodeDownStatusStatistic, func(i, j int) bool {
		return nodeDownStatusStatistic[i].T < nodeDownStatusStatistic[j].T
	})

	nodeDownNumbers := make([]*dto.StatisticItem, 0, len(hostNameDownNumberMap))
	for k, v := range hostNameDownNumberMap {
		if v != 0 {
			nodeDownNumbers = append(nodeDownNumbers, &dto.StatisticItem{Key: k, Value: float64(v)})
		}
	}

	return &dto.NodeDownStatisticReportResp{
		NodeDownNumberRate: util.Convert2ResourceUtAvgReport(nodeDownStatusStatistic, "节点宕机率"),
		NodeDownNumber: &dto.OriginStatisticData{
			Name:         "节点宕机次数",
			OriginalData: nodeDownNumbers,
		},
	}, nil
}

func resolveNodeAvailableMetic(nodeList []*model.NodeInfo, hostNameMetricMap map[string][]promodel.SamplePair) (int, map[int64]int64, map[string]int64) {
	nodeNumber := 0
	timeTimesMap := make(map[int64]int64)
	hostNameDownNumberMap := make(map[string]int64)

	for _, node := range nodeList {
		if metrics, ok := hostNameMetricMap[node.NodeName]; ok {
			preValue := 0
			for i := 0; i < len(metrics); i++ {
				// 节点宕机率
				curValue := int(metrics[i].Value)
				timeKey := int64(metrics[i].Timestamp)
				if _, ok := timeTimesMap[timeKey]; !ok {
					timeTimesMap[timeKey] = 0
				}
				if curValue == 0 && i != 0 && i != len(metrics)-1 {
					timeTimesMap[timeKey] = timeTimesMap[timeKey] + 1
				}

				// 节点宕机次数
				if i == 0 {
					preValue = curValue
					continue
				}
				if _, ok := hostNameDownNumberMap[node.NodeName]; !ok {
					hostNameDownNumberMap[node.NodeName] = 0
				}
				// 特殊处理结尾仍处于 0 的情况
				if (preValue == 0 && curValue == 1) || (i == len(metrics)-1 && curValue == 0 && preValue == 0) {
					hostNameDownNumberMap[node.NodeName] = hostNameDownNumberMap[node.NodeName] + 1
				}
				preValue = curValue
			}

			nodeNumber++
		}
	}

	return nodeNumber, timeTimesMap, hostNameDownNumberMap
}

func (s *ReportServiceImpl) ExportResourceUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error {
	reports, err := s.getAllNodesTypeReport(ctx, req)
	if err != nil {
		return err
	}

	reportData := getTiemStepMap(req.Start, req.End, CalculateStep(req.Start, req.End))
	switch req.ReportType {
	case consts.CPUUtAvg:
		csvHeaders, reportDataList := getExportCSVDataForMetrics(reports, reportData, consts.PercentSymbol)
		return exportCSVFile(ctx, "CPU平均利用率报表", csvHeaders, reportDataList, common.FloatPrecision6)
	case consts.MemUtAvg:
		csvHeaders, reportDataList := getExportCSVDataForMetrics(reports, reportData, consts.PercentSymbol)
		return exportCSVFile(ctx, "内存平均利用率报表", csvHeaders, reportDataList, common.FloatPrecision6)
	case consts.TotalIoUtAvg:
		csvHeaders, reportDataList := getExportCSVDataForMetrics(reports, reportData, consts.IoBandwidthUnit)
		return exportCSVFile(ctx, "磁盘吞吐率报表", csvHeaders, reportDataList, common.FloatPrecision2)
	default:
		return fmt.Errorf("unsupported report type: [%v]", req.ReportType)
	}
}

func getExportCSVDataForMetrics(reports []*dto.UtAvgMetric, reportData []*dto.ReportData, unit string) ([]string, []*dto.ReportData) {
	csvHeaders := []string{"时间"}
	for i := 0; i < len(reports); i++ {
		csvHeader := reports[i].Name
		if unit != "" {
			csvHeader = fmt.Sprintf("%s(%s)", csvHeader, unit)
		}
		csvHeaders = append(csvHeaders, csvHeader)

		metrics := reports[i].Metrics
		metricsMap := make(map[int64]*dto.MetricTV)
		for _, v := range metrics {
			metricsMap[v.Timestamp] = v
		}

		for _, v := range reportData {
			if metric, ok := metricsMap[v.DateTime]; ok {
				v.Metrics = append(v.Metrics, metric.Value)
			} else {
				v.Metrics = append(v.Metrics, 0)
			}
		}
	}

	return csvHeaders, reportData
}

func getTiemStepMap(startTime, endTime, timeStep int64) []*dto.ReportData {
	return getTiemStepMapWithDateType(startTime, endTime, timeStep, "")
}

func getTiemStepMapWithDateType(startTime, endTime, timeStep int64, dateType string) []*dto.ReportData {
	reportDataList := make([]*dto.ReportData, 0)
	reportDataMap := make(map[int64]*dto.ReportData)
	if startTime > endTime {
		return reportDataList
	}

	startDateTime := time.UnixMilli(startTime).Truncate(1 * time.Second)
	endDateTime := time.UnixMilli(endTime).Truncate(1 * time.Second)
	for !startDateTime.After(endDateTime) {
		reportDataMap[startDateTime.UnixMilli()] = &dto.ReportData{
			DateTime: startDateTime.UnixMilli(),
			Metrics:  []float64{},
		}
		startDateTime = startDateTime.Add(time.Duration(timeStep) * time.Second)
	}

	for _, v := range reportDataMap {
		switch dateType {
		case consts.DateTypeDay:
			crime := time.UnixMilli(v.DateTime)
			v.DateTime = time.Date(crime.Year(), crime.Month(), crime.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
		default:
			// don't need to do anything
		}
		reportDataList = append(reportDataList, v)
	}
	sort.Slice(reportDataList, func(i, j int) bool {
		return reportDataList[i].DateTime < reportDataList[j].DateTime
	})

	return reportDataList
}

func exportCSVFile(ctx *gin.Context, csvFileName string, csvHeaders []string, reportDataList []*dto.ReportData, prec int) error {
	return csvutil.ExportCSVFile(ctx, &csvutil.ExportCSVFileInfo{
		CSVFileName: csvFileName,
		CSVHeaders:  csvHeaders,
		FillCSVData: fillCSVDataForMetrics(csvHeaders, reportDataList, "", prec),
	})
}

func fillCSVDataForMetrics(csvHeaders []string, reportDataList []*dto.ReportData, dateType string, prec int) func(w *csv.Writer) error {
	return func(w *csv.Writer) error {
		for _, v := range reportDataList {
			rowData := make([]string, 0, len(csvHeaders))
			switch dateType {
			case consts.DateTypeDay:
				rowData = append(rowData, csvutil.CSVContentWithTab(time.UnixMilli(v.DateTime).Format(common.DateOnly)))
			default:
				rowData = append(rowData, csvutil.CSVContentWithTab(time.UnixMilli(v.DateTime).Format(common.DatetimeFormat)))
			}

			for _, metric := range v.Metrics {
				rowData = append(rowData, csvutil.CSVContentWithTab(strconv.FormatFloat(metric, 'f', prec, 64)))
			}
			_ = w.Write(rowData)
		}
		w.Flush()

		return nil
	}
}

func fillCSVDataForOriginData(metric []*job.MetricKV, prec int) func(w *csv.Writer) error {
	return func(w *csv.Writer) error {
		for _, v := range metric {
			rowData := make([]string, 0)
			rowData = append(rowData, csvutil.CSVContentWithTab(v.Key))
			rowData = append(rowData, csvutil.CSVContentWithTab(strconv.FormatFloat(v.Value, 'f', prec, 64)))
			_ = w.Write(rowData)
		}
		w.Flush()

		return nil
	}
}

func (s *ReportServiceImpl) ExportDiskUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error {
	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  CalculateStep(req.Start, req.End),
	}
	diskUsageUtMetric, err := s.ReportDao.GetDiskUsageUtMetric(ctx, timeRange)
	if err != nil {
		return err
	}

	reportData := getTiemStepMap(req.Start, req.End, CalculateStep(req.Start, req.End))
	csvHeaders, reportDataList := getExportCSVDataForMetrics(diskUsageUtMetric, reportData, consts.PercentSymbol)

	return exportCSVFile(ctx, "磁盘利用率报表", csvHeaders, reportDataList, common.FloatPrecision6)
}

func (s *ReportServiceImpl) ExportCPUTimeSum(ctx *gin.Context, req *dto.UniteReportReq) error {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
		TopSize:   consts.CPUTimeSumTopSize,
	}

	cpuTimeMetric, err := client.GetInstance().Job.GetJobCPUTimeMetric(ctx, pbReq)
	if err != nil {
		return err
	}

	err = csvutil.ExportCSVFilesToZip(ctx, "计算应用-核时使用情况", []*csvutil.ExportCSVFileInfo{
		{
			CSVFileName: "核时使用情况-按软件",
			CSVHeaders:  []string{"软件名称", "使用核时(小时)"},
			FillCSVData: fillCSVDataForOriginData(cpuTimeMetric.AppMetrics, common.FloatPrecision2),
		},
		{
			CSVFileName: "核时使用情况-按用户",
			CSVHeaders:  []string{"用户名称", "使用核时(小时)"},
			FillCSVData: fillCSVDataForOriginData(cpuTimeMetric.UserMetrics, common.FloatPrecision2),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ReportServiceImpl) ExportJobCount(ctx *gin.Context, req *dto.UniteReportReq) error {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobCountMetric, err := client.GetInstance().Job.GetJobCountMetric(ctx, pbReq)
	if err != nil {
		return err
	}

	err = csvutil.ExportCSVFilesToZip(ctx, "计算应用-作业投递数情况", []*csvutil.ExportCSVFileInfo{
		{
			CSVFileName: "作业投递数情况-按软件",
			CSVHeaders:  []string{"软件名称", "作业数(个)"},
			FillCSVData: fillCSVDataForOriginData(jobCountMetric.AppMetrics, common.FloatPrecision0),
		},
		{
			CSVFileName: "作业投递数情况-按用户",
			CSVHeaders:  []string{"用户名称", "作业数(个)"},
			FillCSVData: fillCSVDataForOriginData(jobCountMetric.UserMetrics, common.FloatPrecision0),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ReportServiceImpl) ExportJobDeliverCount(ctx *gin.Context, req *dto.UniteReportReq) error {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobDeliverCount, err := client.GetInstance().Job.GetJobDeliverCount(ctx, pbReq)
	if err != nil {
		return err
	}

	deliverCountResp := util.Convert2JobDeliverCountResp(jobDeliverCount)

	reportData := getTiemStepMapWithDateType(req.Start, req.End, consts.RangeOneDaySec, consts.DateTypeDay)
	csvHeadersForUser, reportDataListForUser := getExportCSVDataForMetrics(deliverCountResp.JobDeliverUserCount, reportData, "")
	reportData = getTiemStepMapWithDateType(req.Start, req.End, consts.RangeOneDaySec, consts.DateTypeDay)
	csvHeadersForJob, reportDataListForJob := getExportCSVDataForMetrics(deliverCountResp.JobDeliverJobCount, reportData, "")

	err = csvutil.ExportCSVFilesToZip(ctx, "用户数与作业数报表", []*csvutil.ExportCSVFileInfo{
		{
			CSVFileName: "提交作业的用户数-按天",
			CSVHeaders:  csvHeadersForUser,
			FillCSVData: fillCSVDataForMetrics(csvHeadersForUser, reportDataListForUser, consts.DateTypeDay, common.FloatPrecision0),
		},
		{
			CSVFileName: "投递作业数-按天",
			CSVHeaders:  csvHeadersForJob,
			FillCSVData: fillCSVDataForMetrics(csvHeadersForJob, reportDataListForJob, consts.DateTypeDay, common.FloatPrecision0),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ReportServiceImpl) ExportJobWaitStatistic(ctx *gin.Context, req *dto.UniteReportReq) error {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobWaitStatistic, err := client.GetInstance().Job.GetJobWaitStatistic(ctx, pbReq)
	if err != nil {
		return err
	}

	waitStatisticResp := util.Convert2JobWaitStatisticResp(jobWaitStatistic)

	reportData := getTiemStepMapWithDateType(req.Start, req.End, consts.RangeOneDaySec, consts.DateTypeDay)
	csvHeadersForWaitTime, reportDataListForWaitTime := getExportCSVDataForMetrics(waitStatisticResp.JobWaitTimeStatistic, reportData, "小时")
	reportData = getTiemStepMapWithDateType(req.Start, req.End, consts.RangeOneDaySec, consts.DateTypeDay)
	csvHeadersForWaitNum, reportDataListForWaitNum := getExportCSVDataForMetrics(waitStatisticResp.JobWaitNumStatistic, reportData, "")

	err = csvutil.ExportCSVFilesToZip(ctx, "作业等待情况报表", []*csvutil.ExportCSVFileInfo{
		{
			CSVFileName: "作业等待时间情况-按天",
			CSVHeaders:  csvHeadersForWaitTime,
			FillCSVData: fillCSVDataForMetrics(csvHeadersForWaitTime, reportDataListForWaitTime, consts.DateTypeDay, common.FloatPrecision2),
		},
		{
			CSVFileName: "作业等待人次情况-按天",
			CSVHeaders:  csvHeadersForWaitNum,
			FillCSVData: fillCSVDataForMetrics(csvHeadersForWaitNum, reportDataListForWaitNum, consts.DateTypeDay, common.FloatPrecision0),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ReportServiceImpl) ExportLicenseAppUsedUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error {
	licenseAppUtAvg, err := s.getLicenseAppUsedUtAvgData(ctx, req)
	if err != nil {
		return err
	}

	reportData := getTiemStepMap(req.Start, req.End, CalculateStep(req.Start, req.End))
	csvHeaders, reportDataList := getExportCSVDataForMetrics(licenseAppUtAvg, reportData, "")

	return exportCSVFile(ctx, "许可证应用使用情况报表", csvHeaders, reportDataList, common.FloatPrecision6)
}

func (s *ReportServiceImpl) ExportNodeDownStatistics(ctx *gin.Context, req *dto.ExportNodeDownStatisticsReq) error {
	start := time.Now()

	nodeDownInfos, err := s.getNodeDownInfos(ctx, req)
	if err != nil {
		return err
	}

	csvHeaders := []string{"机器名称", "宕机时长(小时)", "开始时间", "结束时间"}
	err = csvutil.ExportCSVFile(ctx, &csvutil.ExportCSVFileInfo{
		CSVFileName: "机器节点宕机统计表",
		CSVHeaders:  csvHeaders,
		FillCSVData: func(w *csv.Writer) error {
			for _, v := range nodeDownInfos {
				rowData := make([]string, 0, len(csvHeaders))
				rowData = append(rowData, csvutil.CSVContentWithTab(v.NodeName))
				rowData = append(rowData, csvutil.CSVContentWithTab(strconv.FormatFloat(float64(v.DownTime)/1000/60/60, 'f', common.FloatPrecision2, 64)))
				rowData = append(rowData, csvutil.CSVContentWithTab(time.UnixMilli(v.DownStart).Format(common.DatetimeFormat)))
				rowData = append(rowData, csvutil.CSVContentWithTab(time.UnixMilli(v.DownEnd).Format(common.DatetimeFormat)))

				_ = w.Write(rowData)
			}
			w.Flush()

			return nil
		},
	})
	if err != nil {
		return err
	}

	logging.GetLogger(ctx).Infof("export node down statistics, cost time: %v", time.Since(start))

	return nil
}

func (s *ReportServiceImpl) getNodeDownInfos(ctx *gin.Context, req *dto.ExportNodeDownStatisticsReq) ([]*dto.NodeDownInfo, error) {
	logger := logging.GetLogger(ctx)

	nodeList, err := s.nodeDao.NodeList(ctx)
	if err != nil {
		return nil, err
	}

	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  getTimeStep(req.Start, req.End),
	}
	hostNameMetricMap, err := s.ReportDao.GetNodeAvailableMetic(ctx, consts.CPU, timeRange)
	if err != nil {
		return nil, err
	}

	nodeDownInfos := make([]*dto.NodeDownInfo, 0)
	for _, node := range nodeList {
		if metrics, ok := hostNameMetricMap[node.NodeName]; ok {
			preValue := 0
			downStart := true
			var downStartTime int64
			for i := 0; i < len(metrics); i++ {
				curValue := int(metrics[i].Value)
				if i == 0 {
					preValue = curValue
					if curValue == 0 {
						downStart = false
						downStartTime = int64(metrics[i].Timestamp)
					}
					continue
				}

				if downStart {
					if preValue == 1 && curValue == 0 {
						downStart = false
						downStartTime = int64(metrics[i].Timestamp)
					}
				} else {
					// 特殊处理结尾仍处于 0 的情况
					if (preValue == 0 && curValue == 1) || (i == len(metrics)-1 && curValue == 0 && preValue == 0) {
						if downStartTime == 0 {
							logger.Infof("node %s down time is 0", node.NodeName)
						}
						nodeDownInfos = append(nodeDownInfos, &dto.NodeDownInfo{
							NodeName:  node.NodeName,
							DownTime:  int64(metrics[i].Timestamp) - downStartTime,
							DownStart: downStartTime,
							DownEnd:   int64(metrics[i].Timestamp),
						})
						downStart = true
						downStartTime = 0
					}
				}
				preValue = curValue
			}
		}
	}

	return nodeDownInfos, nil
}

func (s *ReportServiceImpl) getLicenseAppTypes() ([]string, error) {
	licenseManagerList, err := s.localAPI.Client.License.ListLicenseManager()
	if err != nil {
		return nil, err
	}

	if licenseManagerList.Data.Total == 0 {
		return []string{}, nil
	}

	appTypes := make([]string, 0, len(licenseManagerList.Data.Items))
	distinctMap := make(map[string]string)
	for _, item := range licenseManagerList.Data.Items {
		_, exist := distinctMap[item.AppType]
		if exist {
			continue
		} else {
			distinctMap[item.AppType] = item.AppType
			appTypes = append(appTypes, item.AppType)
		}
	}

	return appTypes, nil
}

func (s *ReportServiceImpl) GetLicenseAppModuleUsedUtAvg(ctx *gin.Context, req *dto.LicenseAppModuleUsedUtAvgReq) (*dto.LicenseAppModuleUsedUtAvgReportResp, error) {
	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  CalculateStep(req.Start, req.End),
	}

	licenseAppModules, err := s.getLicenseAppModules(req.LicenseId)
	if err != nil {
		return nil, err
	}

	licenseAppModuleUtAvg := make([]*dto.UtAvgMetric, 0, len(licenseAppModules))
	for _, appModule := range licenseAppModules {
		metricAvgUT, err := s.ReportDao.GetLicenseAppModuleUsedUtMetric(ctx, req.LicenseType, appModule, req.LicenseId, timeRange)
		if err != nil {
			return nil, err
		}
		utAvgMetric := util.Convert2ResourceUtAvgReport(metricAvgUT, appModule)
		licenseAppModuleUtAvg = append(licenseAppModuleUtAvg, utAvgMetric)
	}

	return &dto.LicenseAppModuleUsedUtAvgReportResp{
		LicenseAppModuleUtAvg: licenseAppModuleUtAvg,
	}, nil

}

func (s *ReportServiceImpl) getLicenseAppModules(licenseId string) ([]string, error) {
	resp, err := s.localAPI.Client.License.ListModuleConfig(s.localAPI.Client.License.ListModuleConfig.LicenseId(licenseId))
	if err != nil {
		return nil, err
	}

	if len(resp.Data.ModuleConfigs) == 0 {
		return []string{}, nil
	}

	moduleConfigs := make([]string, 0, len(resp.Data.ModuleConfigs))
	distinctMap := make(map[string]string)
	for _, moduleConfig := range resp.Data.ModuleConfigs {
		_, exist := distinctMap[moduleConfig.Id]
		if exist {
			continue
		} else {
			distinctMap[moduleConfig.Id] = moduleConfig.ModuleName
			moduleConfigs = append(moduleConfigs, moduleConfig.ModuleName)
		}
	}

	return moduleConfigs, nil
}

func CalculateStep(timeStart, timeEnd int64) int64 {
	diff := timeEnd/1000 - timeStart/1000
	if diff <= consts.RangeOneWeekSec { // 小于等于一月时间区间
		return consts.RangeFiveMinSec
	} else if diff > consts.RangeOneWeekSec && diff <= consts.RangeOneMonthSec {
		return consts.RangeTenMinSec
	} else if diff > consts.RangeOneMonthSec && diff <= consts.RangeHalfAYearSec { // 半年时间区间
		return consts.RangeFortyMinSec
	} else if diff > consts.RangeHalfAYearSec && diff <= consts.RangeOneYearSec { // 一年时间区间
		return consts.RangeOneHourSec
	} else { // 超过一年时间，按照 （年数 * 小时) 来返回
		return ((diff / consts.RangeOneYearSec) + 1) * consts.RangeOneHourSec
	}
}

func (s *ReportServiceImpl) GetDiskUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.DiskUtAvgReportResp, error) {
	timeRange := &dto.TimeRange{
		StartTime: req.Start,
		EndTime:   req.End,
		TimeStep:  CalculateStep(req.Start, req.End),
	}
	diskUsageUtMetric, err := s.ReportDao.GetDiskUsageUtMetric(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	return &dto.DiskUtAvgReportResp{
		DiskUtAvg: diskUsageUtMetric,
	}, nil
}

func (s *ReportServiceImpl) GetCPUTimeSum(ctx *gin.Context, req *dto.UniteReportReq) (*dto.CPUTimeSumMetricsResp, error) {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
		TopSize:   consts.CPUTimeSumTopSize,
	}

	cpuTimeMetric, err := client.GetInstance().Job.GetJobCPUTimeMetric(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	cpuTimeSumResp := util.Convert2CPUTimeSumResp(cpuTimeMetric)
	return cpuTimeSumResp, nil
}

func (s *ReportServiceImpl) GetJobCount(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobCountMetricResp, error) {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobCountMetric, err := client.GetInstance().Job.GetJobCountMetric(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	jobCountResp := util.Convert2JobCountResp(jobCountMetric)
	return jobCountResp, nil
}

func (s *ReportServiceImpl) GetJobDeliverCount(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobDeliverCountResp, error) {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobDeliverCount, err := client.GetInstance().Job.GetJobDeliverCount(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	deliverCountResp := util.Convert2JobDeliverCountResp(jobDeliverCount)
	return deliverCountResp, nil
}

func (s *ReportServiceImpl) GetJobWaitStatistic(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobWaitStatisticResp, error) {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
	}

	jobWaitStatistic, err := client.GetInstance().Job.GetJobWaitStatistic(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	waitStatisticResp := util.Convert2JobWaitStatisticResp(jobWaitStatistic)
	return waitStatisticResp, nil
}

func getTimeStep(startTime, endTime int64) int64 {
	timeStep := int64(math.Ceil(float64(endTime-startTime) / 1000 / 11000))
	if timeStep < 5*60 {
		timeStep = 5 * 60
	}
	return timeStep
}
