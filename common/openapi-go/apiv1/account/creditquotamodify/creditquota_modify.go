package creditquotamodify

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditquotamodify"
)

type API func(options ...Option) (*creditquotamodify.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*creditquotamodify.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s/creditquota", req.AccountID)).
			Json(req))

		ret := new(creditquotamodify.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*creditquotamodify.Request, error) {
	req := &creditquotamodify.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *creditquotamodify.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *creditquotamodify.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) CreditQuotaAmount(creditQuotaAmount int64) Option {
	return func(req *creditquotamodify.Request) error {
		req.CreditQuotaAmount = creditQuotaAmount
		return nil
	}
}
