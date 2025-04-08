package joblist

import (
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	jl "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"
)

// API 管理员作业列表
type API func(options ...Option) (*jl.Response, error)

// New 管理员作业列表
func New(hc *xhttp.Client) API {
	return func(options ...Option) (*jl.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		builder := xhttp.NewRequestBuilder().URI("/admin/jobs")
		builder.AddQuery("PageSize", utils.Stringify(*req.PageSize))
		builder.AddQuery("PageOffset", utils.Stringify(*req.PageOffset))
		if req.JobState != "" {
			builder.AddQuery("JobState", req.JobState)
		}
		if req.Zone != "" {
			builder.AddQuery("Zone", req.Zone)
		}
		if req.UserID != "" {
			builder.AddQuery("UserID", req.UserID)
		}
		if req.AppID != "" {
			builder.AddQuery("AppID", req.AppID)
		}
		builder.AddQuery("WithDelete", utils.Stringify(req.WithDelete))
		builder.AddQuery("IsSystemFailed", utils.Stringify(req.IsSystemFailed))
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
	req := &jl.Request{}
	req.PageOffset = &offset
	req.PageSize = &size

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

// UserID 用户ID
func (api API) UserID(userID string) Option {
	return func(req *jl.Request) error {
		req.UserID = userID
		return nil
	}
}

// AppID 应用ID
func (api API) AppID(appID string) Option {
	return func(req *jl.Request) error {
		req.AppID = appID
		return nil
	}
}

// WithDelete 是否包含已删除的作业
func (api API) WithDelete(withDelete bool) Option {
	return func(req *jl.Request) error {
		req.WithDelete = withDelete
		return nil
	}
}

// IsSystemFailed 查询系统失败的作业
func (api API) IsSystemFailed(isSystemFailed bool) Option {
	return func(req *jl.Request) error {
		req.IsSystemFailed = isSystemFailed
		return nil
	}
}
