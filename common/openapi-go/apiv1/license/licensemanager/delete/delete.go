package delete

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"net/http"
)

type API func(options ...Option) (*licmanager.DeleteLicManagerResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licmanager.DeleteLicManagerResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/admin/licenseManagers/" + string(*req)))

		ret := new(licmanager.DeleteLicManagerResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licmanager.DeleteLicManagerRequest, error) {
	req := new(licmanager.DeleteLicManagerRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licmanager.DeleteLicManagerRequest) error

func (api API) Id(id string) Option {
	return func(req *licmanager.DeleteLicManagerRequest) error {
		*req = licmanager.DeleteLicManagerRequest(id)
		return nil
	}
}
