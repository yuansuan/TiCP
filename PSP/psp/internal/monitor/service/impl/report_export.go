package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/excelutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

func (s *ReportServiceImpl) ExportReport(ctx *gin.Context, req *dto.UniteReportReq) error {
	reportType := req.ReportType
	switch reportType {
	case consts.CPUUtAvg, consts.MemUtAvg, consts.TotalIoUtAvg:
		return s.exportHostResourceMetricAvgUt(ctx, req)
	case consts.CPUTimeSum:
		return s.exportCPUTimeSumReport(ctx, req)

	default:
		return status.Error(errcode.ErrCommonReportTypeUnsupported, errcode.MsgCommonReportTypeUnsupported)
	}
}

// exportHostResourceMetricAvgUt 导出cpu、内存使用情况报告
func (s *ReportServiceImpl) exportHostResourceMetricAvgUt(ctx *gin.Context, req *dto.UniteReportReq) error {
	reports, err := s.getAllNodesTypeReport(ctx, req)
	if err != nil {
		return err
	}

	sheetName := consts.MetricTypeNameMap[req.ReportType]
	fileName := sheetName + "表"

	var percentSymbol string
	if req.ReportType == consts.CPUUtAvg || req.ReportType == consts.MemUtAvg {
		percentSymbol = fmt.Sprintf("(%s)", consts.PercentSymbol)
	} else if req.ReportType == consts.TotalIoUtAvg {
		percentSymbol = fmt.Sprintf("(%s)", consts.IoBandwidthUnit)
	}

	// 表头
	headers := make([]string, 0, len(reports))
	headers = append(headers, "时间")
	// 将数据组装成 timeStamp : []MetricValues
	metricMap := make(map[int64][]float64)
	for _, report := range reports {
		// 表头拼接
		headers = append(headers, report.Name+percentSymbol)
		for _, metric := range report.Metrics {
			metricMap[metric.Timestamp] = append(metricMap[metric.Timestamp], metric.Value)
		}
	}

	// 声明二维数组
	colsLen := len(headers)
	var rows [][]interface{}

	// 向数组添加数据
	dataCols := make([]interface{}, 0, colsLen)
	for k, v := range metricMap {
		// 将cols数组初始化为空
		dataCols = []interface{}{}
		// 格式化表格第一列时间戳
		unixTime := timeutil.ParseUnixTime(k / 1000)
		formatTimeStamp := timeutil.FormatTime(unixTime, common.DatetimeFormat)
		dataCols = append(dataCols, formatTimeStamp) // 时间
		// 增加指标列
		for _, floatVal := range v {
			dataCols = append(dataCols, fmt.Sprint(floatVal))
		}

		// 将数据添加到二维数组中
		rows = append(rows, dataCols)
	}

	// 生成sheet 表单数据
	sheetData := &excelutil.ExportExcelSheetData{
		SheetName: sheetName,
		Headers:   headers,
		Rows:      rows,
	}

	// 导出表单数据
	var dataList = []*excelutil.ExportExcelSheetData{sheetData}
	file, err := excelutil.ExportExcel(dataList)
	if err != nil {
		return err
	}

	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	disposition := fmt.Sprintf("attachment; filename=\"%s.xlsx\"", fileName)
	ctx.Writer.Header().Set("Content-Disposition", disposition)
	err = file.Write(ctx.Writer)
	if err != nil {
		return err
	}

	return nil
}

// ExportCPUTimeSumReport 导出核时使用情况
func (s *ReportServiceImpl) exportCPUTimeSumReport(ctx *gin.Context, req *dto.UniteReportReq) error {
	pbReq := &pb.GetJobMetricRequest{
		StartTime: req.Start,
		EndTime:   req.End,
		TopSize:   consts.CPUTimeSumTopSize,
	}

	cpuTimeMetric, err := client.GetInstance().Job.GetJobCPUTimeMetric(ctx, pbReq)
	if err != nil {
		return err
	}

	fileName := consts.MetricTypeNameMap[req.ReportType]

	// 按软件使用情况统计
	appHeaders, appRows := genCPUTimeSumSheet(cpuTimeMetric.AppMetrics)
	appSheetData := &excelutil.ExportExcelSheetData{
		SheetName: "核时使用情况(按软件)",
		Headers:   appHeaders,
		Rows:      appRows,
	}

	// 按核时使用情况统计
	userHeaders, userRows := genCPUTimeSumSheet(cpuTimeMetric.UserMetrics)
	userSheetData := &excelutil.ExportExcelSheetData{
		SheetName: "核时使用情况(按用户)",
		Headers:   userHeaders,
		Rows:      userRows,
	}

	// 导出表单数据
	var dataList = []*excelutil.ExportExcelSheetData{appSheetData, userSheetData}
	file, err := excelutil.ExportExcel(dataList)
	if err != nil {
		return err
	}

	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	disposition := fmt.Sprintf("attachment; filename=\"%s.xlsx\"", fileName)
	ctx.Writer.Header().Set("Content-Disposition", disposition)
	err = file.Write(ctx.Writer)
	if err != nil {
		return err
	}

	return nil
}

func genCPUTimeSumSheet(metrics []*pb.MetricKV) ([]string, [][]interface{}) {
	headers := make([]string, 0, len(metrics)+1)
	headers = append(headers, " ") // 第一行第一列空格

	colsLen := len(headers)
	var rows [][]interface{}
	dataCols := make([]interface{}, 0, colsLen)
	dataCols = append(dataCols, "核时(小时)")

	for _, metric := range metrics {
		headers = append(headers, metric.Key)
		dataCols = append(dataCols, metric.Value)
	}
	rows = append(rows, dataCols)

	return headers, rows
}
