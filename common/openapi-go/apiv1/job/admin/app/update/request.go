package update

import (
	"encoding/json"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
)

type API func(options ...Option) (*update.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*update.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPutRequestBuilder().
			URI("/admin/apps/" + req.AppID).
			Json(req))

		ret := new(update.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*update.Request, error) {
	req := &update.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *update.Request) error

func (api API) AppID(id string) Option {
	return func(req *update.Request) error {
		req.AppID = id
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *update.Request) error {
		req.Name = name
		return nil
	}
}

func (api API) Type(t string) Option {
	return func(req *update.Request) error {
		req.Type = t
		return nil
	}
}

func (api API) Version(version string) Option {
	return func(req *update.Request) error {
		req.Version = version
		return nil
	}
}

func (api API) AppParamsVersion(appParamsVersion int) Option {
	return func(req *update.Request) error {
		req.AppParamsVersion = appParamsVersion
		return nil
	}
}

func (api API) Image(image string) Option {
	return func(req *update.Request) error {
		req.Image = image
		return nil
	}
}

func (api API) Endpoint(endpoint string) Option {
	return func(req *update.Request) error {
		req.Endpoint = endpoint
		return nil
	}
}

func (api API) Command(command string) Option {
	return func(req *update.Request) error {
		req.Command = command
		return nil
	}
}

func (api API) Description(description string) Option {
	return func(req *update.Request) error {
		req.Description = description
		return nil
	}
}

func (api API) IconUrl(iconUrl string) Option {
	return func(req *update.Request) error {
		req.IconUrl = iconUrl
		return nil
	}
}

func (api API) CoresMaxLimit(coresMaxLimit int64) Option {
	return func(req *update.Request) error {
		req.CoresMaxLimit = coresMaxLimit
		return nil
	}
}

func (api API) CoresPlaceholder(coresPlaceholder string) Option {
	return func(req *update.Request) error {
		req.CoresPlaceholder = coresPlaceholder
		return nil
	}
}

func (api API) FileFilterRule(fileFilterRule string) Option {
	return func(req *update.Request) error {
		req.FileFilterRule = fileFilterRule
		return nil
	}
}

// ResidualEnable 残差图是否开启
func (api API) ResidualEnable(residualEnable bool) Option {
	return func(req *update.Request) error {
		req.ResidualEnable = residualEnable
		return nil
	}
}

// ResidualLogRegexp 残差图文件，默认为工作路径下的stdout.log
func (api API) ResidualLogRegexp(residualLogRegexp string) Option {
	return func(req *update.Request) error {
		req.ResidualLogRegexp = residualLogRegexp
		return nil
	}
}

// ResidualLogParser 残差图解析器类型，目前只支持以下枚举：["starccm","fluent"]
func (api API) ResidualLogParser(residualLogParser string) Option {
	return func(req *update.Request) error {
		req.ResidualLogParser = residualLogParser
		return nil
	}
}

// MonitorChartEnable 监控图表是否开启
func (api API) MonitorChartEnable(monitorChartEnable bool) Option {
	return func(req *update.Request) error {
		req.MonitorChartEnable = monitorChartEnable
		return nil
	}
}

// MonitorChartRegexp 监控图表文件规则，默认为'.*\.out'
func (api API) MonitorChartRegexp(monitorChartRegexp string) Option {
	return func(req *update.Request) error {
		req.MonitorChartRegexp = monitorChartRegexp
		return nil
	}
}

// MonitorChartParser 监控图表解析器类型，目前只支持以下枚举：["fluent","cfx"]
func (api API) MonitorChartParser(monitorChartParser string) Option {
	return func(req *update.Request) error {
		req.MonitorChartParser = monitorChartParser
		return nil
	}
}

func (api API) LicenseVars(licenseVars string) Option {
	return func(req *update.Request) error {
		req.LicenseVars = licenseVars
		return nil
	}
}

func (api API) SnapshotEnable(snapshotEnable bool) Option {
	return func(req *update.Request) error {
		req.SnapshotEnable = snapshotEnable
		return nil
	}
}

// BinPath 应用路径
func (api API) BinPath(binPath map[string]string) Option {
	return func(req *update.Request) error {
		if len(binPath) != 0 {
			b, _ := json.Marshal(binPath)
			req.BinPath = string(b)
			return nil
		}
		return nil
	}
}

// ExtentionParams 扩展参数
func (api API) ExtentionParams(extentionParams string) Option {
	return func(req *update.Request) error {
		req.ExtentionParams = extentionParams
		return nil
	}
}

func (api API) LicManagerId(id string) Option {
	return func(req *update.Request) error {
		req.LicManagerId = id
		return nil
	}
}

func (api API) PublishStatus(status update.Status) Option {
	return func(req *update.Request) error {
		req.PublishStatus = status
		return nil
	}
}

func (api API) NeedLimitCore(needLimitCore bool) Option {
	return func(req *update.Request) error {
		req.NeedLimitCore = needLimitCore
		return nil
	}
}

func (api API) SpecifyQueue(q map[string]string) Option {
	return func(req *update.Request) error {
		req.SpecifyQueue = q
		return nil
	}
}
