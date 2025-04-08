package creditadd

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/creditadd"
)

type API func(options ...Option) (*creditadd.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*creditadd.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s/credit", req.AccountID)).
			Json(req))

		ret := new(creditadd.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*creditadd.Request, error) {
	req := &creditadd.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *creditadd.Request) error

func (api API) AccountID(accountId string) Option {
	return func(req *creditadd.Request) error {
		req.AccountID = accountId
		return nil
	}
}

func (api API) TradeID(tradeId string) Option {
	return func(req *creditadd.Request) error {
		req.TradeId = tradeId
		return nil
	}
}

func (api API) IdempotentID(idempotentID string) Option {
	return func(req *creditadd.Request) error {
		req.IdempotentID = idempotentID
		return nil
	}
}

func (api API) DeltaAwardBalance(deltaAwardBalance int64) Option {
	return func(req *creditadd.Request) error {
		req.DeltaAwardBalance = deltaAwardBalance
		return nil
	}
}

func (api API) DeltaNormalBalance(deltaNormalBalance int64) Option {
	return func(req *creditadd.Request) error {
		req.DeltaNormalBalance = deltaNormalBalance
		return nil
	}
}

func (api API) Comment(comment string) Option {
	return func(req *creditadd.Request) error {
		req.Comment = comment
		return nil
	}
}
