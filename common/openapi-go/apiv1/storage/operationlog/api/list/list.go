package list

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	operationLogApi "github.com/yuansuan/ticp/common/project-root-api/storage/operationLog/api"
	"net/http"
)

type API func(options ...Option) (*operationLogApi.ListOperationLogResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*operationLogApi.ListOperationLogResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/operationLog").
			AddQuery("FileName", req.FileName).
			AddQuery("FileTypes", req.FileTypes).
			AddQuery("OperationTypes", req.OperationTypes).
			AddQuery("BeginTime", utils.Stringify(req.BeginTime)).
			AddQuery("EndTime", utils.Stringify(req.EndTime)).
			AddQuery("PageOffset", utils.Stringify(req.PageOffset)).
			AddQuery("PageSize", utils.Stringify(req.PageSize)))

		ret := new(operationLogApi.ListOperationLogResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*operationLogApi.ListOperationLogRequest, error) {
	req := &operationLogApi.ListOperationLogRequest{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *operationLogApi.ListOperationLogRequest) error

func (api API) PageSize(pageSize int64) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.PageSize = pageSize
		return nil
	}
}

func (api API) PageOffset(pageOffset int64) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.PageOffset = pageOffset
		return nil
	}
}

func (api API) FileName(fileName string) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.FileName = fileName
		return nil
	}
}

func (api API) FileTypes(fileTypes string) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.FileTypes = fileTypes
		return nil
	}
}

func (api API) OperationTypes(operationTypes string) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.OperationTypes = operationTypes
		return nil
	}
}

func (api API) BeginTime(beginTime int64) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.BeginTime = beginTime
		return nil
	}
}

func (api API) EndTime(endTime int64) Option {
	return func(req *operationLogApi.ListOperationLogRequest) error {
		req.EndTime = endTime
		return nil
	}
}
