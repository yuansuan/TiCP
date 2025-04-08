package cpu_usage

import (
	"fmt"
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/cpuusage"
	"net/http"
)

type API func(option ...Option) (*cpuusage.SystemGetResponse, error)

func New(hc *xhttp.Client) API {
	return func(option ...Option) (*cpuusage.SystemGetResponse, error) {
		req := NewRequest(option)
		builder := xhttp.NewRequestBuilder().URI(fmt.Sprintf("/system/jobs/%s/cpuusage", req.JobID))
		resolver := hc.Prepare(builder)

		ret := new(cpuusage.SystemGetResponse)

		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *cpuusage.SystemGetRequest {
	req := new(cpuusage.SystemGetRequest)
	for _, opt := range opts {
		opt(req)
	}
	return req
}
