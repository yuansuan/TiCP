package frozenmodify

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/frozenmodify"
)

type API func(options ...Option) (*frozenmodify.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*frozenmodify.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s", req.AccountID)).
			Json(req))

		ret := new(frozenmodify.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*frozenmodify.Request, error) {
	req := &frozenmodify.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *frozenmodify.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *frozenmodify.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) FrozenState(frozenState bool) Option {
	return func(req *frozenmodify.Request) error {
		req.FrozenState = frozenState
		return nil
	}
}
