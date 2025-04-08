package get

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"net/http"
)

type API func(options ...Option) (*licmanager.GetLicManagerResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licmanager.GetLicManagerResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/licenseManagers/" + string(*req)))

		ret := new(licmanager.GetLicManagerResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licmanager.GetLicManagerRequest, error) {
	req := new(licmanager.GetLicManagerRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licmanager.GetLicManagerRequest) error

func (api API) Id(id string) Option {
	return func(req *licmanager.GetLicManagerRequest) error {
		*req = licmanager.GetLicManagerRequest(id)
		return nil
	}
}
