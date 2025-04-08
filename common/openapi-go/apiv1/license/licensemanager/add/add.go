package add

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	"net/http"
)

type API func(options ...Option) (*licmanager.AddLicManagerResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licmanager.AddLicManagerResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/licenseManagers").
			Json(req))

		ret := new(licmanager.AddLicManagerResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licmanager.AddLicManagerRequest, error) {
	req := new(licmanager.AddLicManagerRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licmanager.AddLicManagerRequest) error

func (api API) AppType(appType string) Option {
	return func(req *licmanager.AddLicManagerRequest) error {
		req.AppType = appType
		return nil
	}
}

func (api API) Os(os int) Option {
	return func(req *licmanager.AddLicManagerRequest) error {
		req.Os = int(os)
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *licmanager.AddLicManagerRequest) error {
		req.Desc = desc
		return nil
	}
}

func (api API) ComputeRule(computeRule string) Option {
	return func(req *licmanager.AddLicManagerRequest) error {
		req.ComputeRule = computeRule
		return nil
	}
}
