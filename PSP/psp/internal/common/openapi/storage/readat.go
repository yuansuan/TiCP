package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apireadat "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/readAt"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// ReadAt 文件读取
func ReadAt(api *openapi.OpenAPI, req readAt.Request, resolver xhttp.ResponseResolver) (*readAt.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.readAt req:[%+v]", req)

	options := []apireadat.Option{
		api.Client.Storage.ReadAt.Path(req.Path),
		api.Client.Storage.ReadAt.Offset(req.Offset),
		api.Client.Storage.ReadAt.Length(req.Length),
		api.Client.Storage.ReadAt.WithResolver(resolver),
	}

	response, err := api.Client.Storage.ReadAt(options...)

	if err != nil {
		logger.Errorf("invoke openapi.readAt failed，req:[%+v]，errMsg:[%v]", req, err)
	}

	return response, err
}
