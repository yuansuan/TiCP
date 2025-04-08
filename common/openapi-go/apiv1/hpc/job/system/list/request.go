package list

import (
	"net/http"
	"strings"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*job.SystemListResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*job.SystemListResponse, error) {
		req := NewRequest(opts)

		builder := xhttp.NewRequestBuilder().URI("/system/jobs")
		builder.AddQuery("PageSize", utils.Stringify(req.PageSize))
		builder.AddQuery("PageOffset", utils.Stringify(req.PageOffset))
		if req.Status != "" {
			builder.AddQuery("Status", req.Status)
		}
		if len(req.JobIDs) > 0 {
			builder.AddQuery("JobIDs", strings.Join(req.JobIDs, ","))
		}

		resolver := hc.Prepare(builder)
		ret := new(job.SystemListResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *job.SystemListRequest {
	req := new(job.SystemListRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
