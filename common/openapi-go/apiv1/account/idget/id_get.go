package idget

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/idget"
)

type API func(options ...Option) (*idget.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*idget.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodGet).
			URI(fmt.Sprintf("/internal/accounts/:%s", req.AccountID)))

		ret := new(idget.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*idget.Request, error) {
	req := &idget.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *idget.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *idget.Request) error {
		req.AccountID = accountID
		return nil
	}
}
