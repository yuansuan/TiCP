package add

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	"net/http"
)

type API func(options ...Option) (*licenseinfo.AddLicenseInfoResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*licenseinfo.AddLicenseInfoResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/licenses").
			Json(req))

		ret := new(licenseinfo.AddLicenseInfoResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*licenseinfo.AddLicenseInfoRequest, error) {
	req := new(licenseinfo.AddLicenseInfoRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *licenseinfo.AddLicenseInfoRequest) error

func (api API) ManagerId(id string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.ManagerId = id
		return nil
	}
}

func (api API) Provider(name string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.Provider = name
		return nil
	}
}

func (api API) MacAddr(addr string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.MacAddr = addr
		return nil
	}
}

func (api API) ToolPath(p string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.ToolPath = p
		return nil
	}
}

func (api API) LicenseUrl(url string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.LicenseUrl = url
		return nil
	}
}

func (api API) Port(port int) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.Port = port
		return nil
	}
}

func (api API) LicenseAddresses(licenseProxies map[string]licenseinfo.LicenseProxy) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.LicenseProxies = licenseProxies
		return nil
	}
}

func (api API) LicenseNum(num string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.LicenseNum = num
		return nil
	}
}

func (api API) Weight(w int) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.Weight = w
		return nil
	}
}

func (api API) BeginTime(t string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.BeginTime = t
		return nil
	}
}

func (api API) EndTime(t string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.EndTime = t
		return nil
	}
}

func (api API) Auth(auth bool) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.Auth = auth
		return nil
	}
}

func (api API) LicenseEnvVar(v string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.LicenseEnvVar = v
		return nil
	}
}

func (api API) LicenseType(t int) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.LicenseType = t
		return nil
	}
}

func (api API) HpcEndpoint(endpoint string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.HpcEndpoint = endpoint
		return nil
	}
}

func (api API) AllowableHpcEndpoints(allowableHpcEndpoints []string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.AllowableHpcEndpoints = allowableHpcEndpoints
		return nil
	}
}

func (api API) CollectorType(t string) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.CollectorType = t
		return nil
	}
}

func (api API) EnableScheduling(enableScheduling bool) Option {
	return func(req *licenseinfo.AddLicenseInfoRequest) error {
		req.EnableScheduling = enableScheduling
		return nil
	}
}
