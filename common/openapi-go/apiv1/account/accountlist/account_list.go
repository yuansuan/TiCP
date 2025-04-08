package accountlist

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountlist"
)

type API func(options ...Option) (*accountlist.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*accountlist.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		requestBuilder := xhttp.NewRequestBuilder().
			URI("/internal/accounts").
			AddQuery("AccountID", req.AccountID).
			AddQuery("AccountName", req.AccountName).
			AddQuery("CustomerID", req.CustomerID).
			AddQuery("PageSize", fmt.Sprintf("%v", req.PageSize)).
			AddQuery("PageIndex", fmt.Sprintf("%v", req.PageIndex))

		if req.FrozenStatus != nil {
			requestBuilder.AddQuery("FrozenStatus", fmt.Sprintf("%v", *req.FrozenStatus))
		}
		resolver := hc.Prepare(requestBuilder)

		ret := new(accountlist.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*accountlist.Request, error) {
	req := &accountlist.Request{
		PageIndex: 1,
		PageSize:  10,
	}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *accountlist.Request) error

func (api API) PageIndex(index int64) Option {
	return func(req *accountlist.Request) error {
		req.PageIndex = index
		return nil
	}
}
func (api API) PageSize(size int64) Option {
	return func(req *accountlist.Request) error {
		req.PageSize = size
		return nil
	}
}
func (api API) AccountID(accountID string) Option {
	return func(req *accountlist.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) AccountName(accountName string) Option {
	return func(req *accountlist.Request) error {
		req.AccountName = accountName
		return nil
	}
}

func (api API) CustomerID(customerID string) Option {
	return func(req *accountlist.Request) error {
		req.CustomerID = customerID
		return nil
	}
}

func (api API) FrozenStatus(frozenStatus bool) Option {
	return func(req *accountlist.Request) error {
		req.FrozenStatus = &frozenStatus
		return nil
	}
}
