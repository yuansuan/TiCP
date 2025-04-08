package get

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(options ...Option) (*software.AdminGetResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*software.AdminGetResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		uri := "/admin/softwares/"
		if req.SoftwareId != nil {
			uri += *req.SoftwareId
		}

		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI(uri))

		ret := new(software.AdminGetResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*software.AdminGetRequest, error) {
	req := &software.AdminGetRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}
