package writeAt

import (
	"encoding/json"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	common2 "github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/writeAt"
)

type API func(options ...Option) (*writeAt.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*writeAt.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		if req.Data == nil {
			res := v20230530.Response{ErrorCode: common2.InvalidData, ErrorMsg: "Data can not be empty"}
			jsonBytes, _ := json.Marshal(res)
			return &writeAt.Response{Response: res}, errors.Errorf("http: %v, body: %v", "400 Bad Request", string(jsonBytes))
		}
		data, err := utils.CompressData(req.Data, req.Compressor)
		if err != nil {
			logging.Default().Errorf("compress data err, err: %v", err)
			res := v20230530.Response{ErrorCode: common2.InvalidCompressor, ErrorMsg: "invalid compressor type"}
			jsonBytes, _ := json.Marshal(res)
			return &writeAt.Response{Response: res}, errors.Errorf("http: %v, body: %v", "400 Bad Request", string(jsonBytes))
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/writeAt").
			AddQuery("Path", utils.Stringify(req.Path)).
			AddQuery("Offset", utils.Stringify(req.Offset)).
			AddQuery("Length", utils.Stringify(req.Length)).
			AddQuery("Compressor", utils.Stringify(req.Compressor)).
			Body(data))

		ret := new(writeAt.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*writeAt.Request, error) {
	req := &writeAt.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *writeAt.Request) error

func (api API) Path(path string) Option {
	return func(req *writeAt.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Offset(Offset int64) Option {
	return func(req *writeAt.Request) error {
		req.Offset = Offset
		return nil
	}
}

func (api API) Length(length int64) Option {
	return func(req *writeAt.Request) error {
		req.Length = length
		return nil
	}
}

func (api API) Compressor(compressor string) Option {
	return func(req *writeAt.Request) error {
		req.Compressor = compressor
		return nil
	}
}

func (api API) Data(data io.Reader) Option {
	return func(req *writeAt.Request) error {
		req.Data = data
		return nil
	}
}
