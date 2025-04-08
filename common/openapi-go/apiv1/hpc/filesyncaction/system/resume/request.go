package resume

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/filesyncaction"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*filesyncaction.SystemResumeFileSyncResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*filesyncaction.SystemResumeFileSyncResponse, error) {
		req := NewRequest(opts)
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(fmt.Sprintf("/system/jobs/%s/resume_file_sync", req.JobID)))
		ret := new(filesyncaction.SystemResumeFileSyncResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *filesyncaction.SystemResumeFileSyncRequest {
	req := new(filesyncaction.SystemResumeFileSyncRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
