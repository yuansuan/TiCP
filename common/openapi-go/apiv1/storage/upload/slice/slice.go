package slice

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	common2 "github.com/yuansuan/ticp/common/project-root-api/common"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	slice2 "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"
	"net/http"
)

type API func(options ...Option) (*slice2.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*slice2.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}

		reader, err := utils.CompressData(req.Slice, req.Compressor)
		if err != nil {
			logging.Default().Errorf("compress data err, err: %v", err)
			res := v20230530.Response{ErrorCode: common2.InvalidCompressor, ErrorMsg: "invalid compressor type"}
			jsonBytes, _ := json.Marshal(res)
			return &slice2.Response{Response: res}, errors.Errorf("http: %v, body: %v", "400 Bad Request", string(jsonBytes))
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/upload/slice").
			AddQuery("UploadID", utils.Stringify(req.UploadID)).
			AddQuery("Offset", utils.Stringify(req.Offset)).
			AddQuery("Length", utils.Stringify(req.Length)).
			AddQuery("Compressor", utils.Stringify(req.Compressor)).
			Body(reader))

		ret := new(slice2.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*slice2.Request, error) {
	req := &slice2.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *slice2.Request) error

func (api API) UploadID(uploadID string) Option {
	return func(req *slice2.Request) error {
		req.UploadID = uploadID
		return nil
	}
}

func (api API) Offset(offset int64) Option {
	return func(req *slice2.Request) error {
		req.Offset = offset
		return nil
	}
}

func (api API) Length(length int64) Option {
	return func(req *slice2.Request) error {
		req.Length = length
		return nil
	}
}

func (api API) Slice(slice []byte) Option {
	return func(req *slice2.Request) error {
		req.Slice = slice
		return nil
	}
}

func (api API) Compressor(compressor string) Option {
	return func(req *slice2.Request) error {
		req.Compressor = compressor
		return nil
	}
}
