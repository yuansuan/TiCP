package delete

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*hardware.AdminDeleteResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*hardware.AdminDeleteResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/admin/hardwares/"
		if req.HardwareId != nil {
			uri += *req.HardwareId
		}

		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().URI(uri))

		ret := new(hardware.AdminDeleteResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*hardware.AdminDeleteRequest, error) {
	req := &hardware.AdminDeleteRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
