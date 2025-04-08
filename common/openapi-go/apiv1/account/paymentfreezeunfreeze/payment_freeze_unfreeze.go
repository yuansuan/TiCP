package paymentfreezeunfreeze

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/paymentfreezeunfreeze"
)

type API func(options ...Option) (*paymentfreezeunfreeze.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*paymentfreezeunfreeze.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s/frozenamount", req.AccountID)).
			Json(req))

		ret := new(paymentfreezeunfreeze.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*paymentfreezeunfreeze.Request, error) {
	req := &paymentfreezeunfreeze.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *paymentfreezeunfreeze.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *paymentfreezeunfreeze.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) Amount(amount int64) Option {
	return func(req *paymentfreezeunfreeze.Request) error {
		req.Amount = amount
		return nil
	}
}

func (api API) Comment(comment string) Option {
	return func(req *paymentfreezeunfreeze.Request) error {
		req.Comment = comment
		return nil
	}
}

func (api API) TradeID(tradeID string) Option {
	return func(req *paymentfreezeunfreeze.Request) error {
		req.TradeID = tradeID
		return nil
	}
}

func (api API) IsFreezed(isFreezed bool) Option {
	return func(req *paymentfreezeunfreeze.Request) error {
		req.IsFreezed = isFreezed
		return nil
	}
}
