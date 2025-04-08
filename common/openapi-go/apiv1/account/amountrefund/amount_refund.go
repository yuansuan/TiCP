package amountrefund

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/amountrefund"
)

type API func(options ...Option) (*amountrefund.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*amountrefund.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodPatch).
			URI(fmt.Sprintf("/internal/accounts/:%s/refundamount", req.AccountID)).
			Json(req))

		ret := new(amountrefund.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*amountrefund.Request, error) {
	req := &amountrefund.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *amountrefund.Request) error

func (api API) AccountID(accountId string) Option {
	return func(req *amountrefund.Request) error {
		req.AccountID = accountId
		return nil
	}
}

func (api API) RefundId(refundId string) Option {
	return func(req *amountrefund.Request) error {
		req.RefundID = refundId
		return nil
	}
}

func (api API) IdempotentID(idempotentID string) Option {
	return func(req *amountrefund.Request) error {
		req.IdempotentID = idempotentID
		return nil
	}
}

func (api API) ResourceID(resourceID string) Option {
	return func(req *amountrefund.Request) error {
		req.ResourceID = resourceID
		return nil
	}
}

func (api API) Amount(amount int64) Option {
	return func(req *amountrefund.Request) error {
		req.Amount = amount
		return nil
	}
}

func (api API) Comment(comment string) Option {
	return func(req *amountrefund.Request) error {
		req.Comment = comment
		return nil
	}
}
