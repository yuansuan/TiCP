package pause

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/filesyncaction"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*filesyncaction.SystemPauseFileSyncResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*filesyncaction.SystemPauseFileSyncResponse, error) {
		req := NewRequest(opts)
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(fmt.Sprintf("/system/jobs/%s/pause_file_sync", req.JobID)))
		ret := new(filesyncaction.SystemPauseFileSyncResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *filesyncaction.SystemPauseFileSyncRequest {
	req := new(filesyncaction.SystemPauseFileSyncRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
