package util

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func PromQuery(queryCondition string) (model.Value, error) {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	newClient, err := api.NewClient(api.Config{Address: boot.Config.App.Middleware.Monitor.PrometheusServerEndpoint})
	if nil != err {
		logger.Errorf("creat prometheus client err: %v", err)
		return nil, err
	}

	promAPI := v1.NewAPI(newClient)
	resp, warnings, err := promAPI.Query(ctx, queryCondition, time.Time{})
	if nil != err {
		logger.Error("prometheus query err: %v, query condition: %v", err, queryCondition)
		return nil, err
	}

	if len(warnings) > 0 {
		logger.Warnf("prometheus query warnings: %+v", warnings)
	}

	return resp, nil
}

func PromRangeQuery(queryCondition string, startTime int64, endTime int64, step int64) (model.Value, error) {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	timeRange := v1.Range{
		Start: time.Unix(startTime, 0),
		End:   time.Unix(endTime, 0),
		Step:  time.Duration(step * 1000 * 1000 * 1000)}

	newClient, err := api.NewClient(api.Config{Address: boot.Config.App.Middleware.Monitor.PrometheusServerEndpoint})
	if nil != err {
		logger.Errorf("creat prometheus client err: %v", err)
		return nil, err
	}

	promAPI := v1.NewAPI(newClient)
	resp, warnings, err := promAPI.QueryRange(ctx, queryCondition, timeRange)
	if nil != err {
		logger.Error("prometheus query err: %v, condition: %v, start: %v, end: %v, step: %v", err, queryCondition, startTime, endTime, step)
		return nil, err
	}

	if len(warnings) > 0 {
		logger.Warnf("prometheus query warnings: %+v", warnings)
	}

	return resp, nil
}

func ParseHpcCommandOutFields(fields []string) (*dto.NodeHpcRunningInfo, error) {
	if len(fields) == 0 {
		return nil, errors.New("command fields is empty!")
	}
	hpcRunningInfo := &dto.NodeHpcRunningInfo{}

	hpcRunningInfo.NodeName = strings.TrimSpace(fields[0])
	hpcRunningInfo.State = strings.TrimSpace(fields[3])
	cpuInfos := strings.TrimSpace(fields[4])
	if !strutil.IsEmpty(cpuInfos) {
		infoArray := strings.Split(cpuInfos, "/")
		if len(infoArray) < 4 {
			return nil, errors.Errorf("cpus info [%s] format is invalid!", cpuInfos)
		}

		cpuAlloc, err := strconv.Atoi(infoArray[0])
		if err != nil {
			return nil, err
		}
		hpcRunningInfo.CPUAlloc = cpuAlloc

		cpuIdle, err := strconv.Atoi(infoArray[1])
		if err != nil {
			return nil, err
		}
		hpcRunningInfo.CPUIdle = cpuIdle

		cpuTotal, err := strconv.Atoi(infoArray[3])
		if err != nil {
			return nil, err
		}
		hpcRunningInfo.CPUTotal = cpuTotal
	}

	return hpcRunningInfo, nil
}

// CalculateAvg CalculateAvg
func CalculateAvg(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}

func ValidReportType(reportType string) bool {
	if strutil.IsEmpty(reportType) {
		return false
	}

	_, ok := consts.MetricTypeNameMap[reportType]
	return !ok
}
