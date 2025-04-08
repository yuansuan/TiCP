package residual

import (
	"strings"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser/fluent"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser/starccm"
)

// NewParser ...
func NewParser(parserName string) parser.Parser {
	switch strings.ToLower(parserName) {
	case schema.ResidualLogParserTypeStarccm:
		return &starccm.Parser{}
	case schema.ResidualLogParserTypeFluent:
		return &fluent.Parser{}
	default:
		return nil
	}
}

// ConvertResidual ...
func ConvertResidual(in *parser.Result) *schema.Residual {
	res := &schema.Residual{
		Vars:          []*schema.ResidualVar{},
		AvailableXvar: in.XVar,
	}
	for k, v := range in.Series {
		rv := &schema.ResidualVar{
			Name:   k,
			Values: v,
		}

		res.Vars = append(res.Vars, rv)
	}
	return res
}
