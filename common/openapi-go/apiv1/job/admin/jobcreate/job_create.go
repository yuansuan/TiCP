package jobcreate

import (
	"net/http"
	"time"

	"github.com/yuansuan/ticp/common/openapi-go/utils/payby"

	js "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	job "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

// API 提交作业api
type API func(options ...Option) (*js.Response, error)

// New 新建
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*js.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/jobs").
			Json(req))

		ret := new(js.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 新建请求
func NewRequest(options []Option) (*js.Request, error) {
	req := &js.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *js.Request) error

// Name 作业名称
func (api API) Name(name string) Option {
	return func(req *js.Request) error {
		req.Name = name
		return nil
	}
}

// Comment 作业备注
func (api API) Comment(comment string) Option {
	return func(req *js.Request) error {
		req.Comment = comment
		return nil
	}
}

// Timeout 作业超时时间
func (api API) Timeout(timeout int64) Option {
	return func(req *js.Request) error {
		req.Timeout = timeout
		return nil
	}
}

// Zone 作业运行区域
func (api API) Zone(zone string) Option {
	return func(req *js.Request) error {
		req.Zone = zone
		return nil
	}
}

// Params 作业参数
func (api API) Params(params job.Params) Option {
	return func(req *js.Request) error {
		req.Params = params
		return nil
	}
}

// Queue 作业队列
func (api API) Queue(queue string) Option {
	return func(req *js.Request) error {
		req.Queue = queue
		return nil
	}
}

// ChargeParams 计费相关参数
func (api API) ChargeParams(chargeParams v20230530.ChargeParams) Option {
	return func(req *js.Request) error {
		req.ChargeParam = chargeParams
		return nil
	}
}

// NoRound 单节点是否不进行取整,仅限内部用户使用
func (api API) NoRound(noRound bool) Option {
	return func(req *js.Request) error {
		req.NoRound = noRound
		return nil
	}
}

// AllocType 作业CPU资源的分配方式
func (api API) AllocType(allocType string) Option {
	return func(req *js.Request) error {
		req.AllocType = allocType
		return nil
	}
}

// PreScheduleID 预调度ID，如果为不为空则从预调度信息中获取资源信息等
func (api API) PreScheduleID(preScheduleID string) Option {
	return func(req *js.Request) error {
		req.PreScheduleID = preScheduleID
		return nil
	}
}

// JobSchedulerSubmitFlags 自定义调度器提交参数
func (api API) JobSchedulerSubmitFlags(jobSchedulerSubmitFlags map[string]string) Option {
	return func(req *js.Request) error {
		req.JobSchedulerSubmitFlags = jobSchedulerSubmitFlags
		return nil
	}
}

func (api API) PayBy(payBy string) Option {
	return func(req *js.Request) error {
		req.PayBy = payBy
		return nil
	}
}

func (api API) PayByParams(payByAccessKeyID, payByAccessSecret string) Option {
	return func(req *js.Request) error {
		if req.PayBy != "" {
			return nil
		}

		timestamp := time.Now().UTC().UnixMilli()
		payBy, _ := payby.NewPayBy(payByAccessKeyID, payByAccessSecret, req.Name, timestamp)
		req.PayBy = payBy.Token()
		return nil
	}
}
