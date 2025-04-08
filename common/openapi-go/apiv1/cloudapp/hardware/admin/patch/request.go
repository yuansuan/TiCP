package patch

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*hardware.AdminPatchResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*hardware.AdminPatchResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/admin/hardwares/"
		if req.HardwareId != nil {
			uri += *req.HardwareId
		}

		resolver := hc.Prepare(xhttp.NewPatchRequestBuilder().
			URI(uri).
			Json(req))

		ret := new(hardware.AdminPatchResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*hardware.AdminPatchRequest, error) {
	req := &hardware.AdminPatchRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
