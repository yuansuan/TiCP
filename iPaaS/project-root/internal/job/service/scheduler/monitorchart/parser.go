package monitorchart

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart/parser"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart/parser/fluent"
)

// NewParser ...
func NewParser(parserType string) parser.Parser {
	switch parserType {
	case "fluent":
		return &fluent.Parser{}
	default:
		return nil
	}
}

// ConvertMonitorChart ...
func ConvertMonitorChart(prs []*parser.Result) []*schema.MonitorChart {
	res := []*schema.MonitorChart{}
	for _, in := range prs {
		mc := &schema.MonitorChart{
			Key:   in.Key,
			Items: []*schema.MonitorChartItem{},
		}

		for _, item := range in.Items {
			mc.Items = append(mc.Items, &schema.MonitorChartItem{
				Kv: []float64{item.Iteration, item.Value},
			})
		}

		res = append(res, mc)
	}

	return res
}

func mergeResultMap(results []*parser.Result, resultMap map[string]*parser.Result) []*parser.Result {
	for _, result := range resultMap {
		if result == nil {
			continue
		}
		results = append(results, result)
	}
	return results
}
