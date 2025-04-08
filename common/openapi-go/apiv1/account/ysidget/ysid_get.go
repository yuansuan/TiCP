package ysidget

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/ysidget"
)

type API func(options ...Option) (*ysidget.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*ysidget.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			Method(http.MethodGet).
			URI(fmt.Sprintf("/internal/users/:%s/account", req.UserID)))

		ret := new(ysidget.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*ysidget.Request, error) {
	req := &ysidget.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *ysidget.Request) error

func (api API) UserID(userID string) Option {
	return func(req *ysidget.Request) error {
		req.UserID = userID
		return nil
	}
}
