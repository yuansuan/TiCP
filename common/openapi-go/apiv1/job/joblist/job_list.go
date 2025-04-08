package joblist

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jl "github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
)

// API 作业列表
type API func(options ...Option) (*jl.Response, error)

// New 作业列表
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jl.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		builder := xhttp.NewRequestBuilder().URI("/api/jobs")
		builder.AddQuery("PageSize", utils.Stringify(*req.PageSize))
		builder.AddQuery("PageOffset", utils.Stringify(*req.PageOffset))
		if req.JobState != "" {
			builder.AddQuery("JobState", req.JobState)
		}
		if req.Zone != "" {
			builder.AddQuery("Zone", req.Zone)
		}
		resolver := hc.Prepare(builder)

		ret := new(jl.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

// NewRequest 新建请求
func NewRequest(options []Option) (*jl.Request, error) {
	offset := int64(0)
	size := int64(100)
	req := &jl.Request{
		PageOffset: &offset,
		PageSize:   &size,
	}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option 选项
type Option func(req *jl.Request) error

// PageOffset 分页偏移量，0开始
func (api API) PageOffset(pageOffset int64) Option {
	return func(req *jl.Request) error {
		req.PageOffset = &pageOffset
		return nil
	}
}

// PageSize 分页大小，1~1000，默认100
func (api API) PageSize(size int64) Option {
	return func(req *jl.Request) error {
		req.PageSize = &size
		return nil
	}
}

// JobState 作业状态过滤
func (api API) JobState(jobState string) Option {
	return func(req *jl.Request) error {
		req.JobState = jobState
		return nil
	}
}

// Zone 分区，分区枚举从list zones接口获取
func (api API) Zone(zone string) Option {
	return func(req *jl.Request) error {
		req.Zone = zone
		return nil
	}
}
