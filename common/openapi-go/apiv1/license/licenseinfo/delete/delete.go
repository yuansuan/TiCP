package delete

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	"net/http"
)

type API func(options ...Option) (*licenseinfo.DeleteLicenseInfoResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licenseinfo.DeleteLicenseInfoResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/admin/licenses/" + string(*req)))

		ret := new(licenseinfo.DeleteLicenseInfoResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licenseinfo.DeleteLicenseInfoRequest, error) {
	req := new(licenseinfo.DeleteLicenseInfoRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licenseinfo.DeleteLicenseInfoRequest) error

func (api API) Id(id string) Option {
	return func(req *licenseinfo.DeleteLicenseInfoRequest) error {
		*req = licenseinfo.DeleteLicenseInfoRequest(id)
		return nil
	}
}
