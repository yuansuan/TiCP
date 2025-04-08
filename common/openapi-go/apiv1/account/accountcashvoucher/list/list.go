package list

import (
	"net/http"
	"strconv"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/accountcashvoucher/list"
)

type API func(options ...Option) (*list.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*list.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodGet).
			URI("/internal/accountcashvouchers").
			AddQuery("AccountID", req.AccountID).
			AddQuery("StartTime", req.StartTime).
			AddQuery("EndTime", req.EndTime).
			AddQuery("PageIndex", strconv.FormatInt(req.PageIndex, 10)).
			AddQuery("PageSize", strconv.FormatInt(req.PageSize, 10)),
		)

		ret := new(list.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*list.Request, error) {
	req := &list.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *list.Request) error

func (api API) AccountID(accountID string) Option {
	return func(req *list.Request) error {
		req.AccountID = accountID
		return nil
	}
}

func (api API) StartTime(startTime string) Option {
	return func(req *list.Request) error {
		req.StartTime = startTime
		return nil
	}
}

func (api API) EndTime(endTime string) Option {
	return func(req *list.Request) error {
		req.EndTime = endTime
		return nil
	}
}

func (api API) PageIndex(index int64) Option {
	return func(req *list.Request) error {
		req.PageIndex = index
		return nil
	}
}
func (api API) PageSize(size int64) Option {
	return func(req *list.Request) error {
		req.PageSize = size
		return nil
	}
}
