package ls

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"
	"net/http"
)

type API func(options ...Option) (*ls.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*ls.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/lsWithPage").
			AddQuery("Path", utils.Stringify(req.Path)).
			AddQuery("PageOffset", utils.Stringify(req.PageOffset)).
			AddQuery("PageSize", utils.Stringify(req.PageSize)).
			AddQuery("FilterRegexp", utils.Stringify(req.FilterRegexp)).
			JsonCond(req.FilterRegexpList != nil, ls.Request{FilterRegexpList: req.FilterRegexpList}))

		ret := new(ls.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*ls.Request, error) {
	req := &ls.Request{PageSize: 100}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *ls.Request) error

func (api API) Path(path string) Option {
	return func(req *ls.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) PageOffset(PageOffset int64) Option {
	return func(req *ls.Request) error {
		req.PageOffset = PageOffset
		return nil
	}
}

func (api API) PageSize(PageSize int64) Option {
	return func(req *ls.Request) error {
		req.PageSize = PageSize
		return nil
	}
}

func (api API) FilterRegexp(filterRegexp string) Option {
	return func(req *ls.Request) error {
		req.FilterRegexp = filterRegexp
		return nil
	}
}

func (api API) FilterRegexpList(filterRegexpList []string) Option {
	return func(req *ls.Request) error {
		req.FilterRegexpList = filterRegexpList
		return nil
	}
}
