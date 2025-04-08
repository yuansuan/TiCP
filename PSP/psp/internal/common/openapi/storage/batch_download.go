package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apibatchdownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/batchDownload"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// BatchDownload 批量文件下载
func BatchDownload(api *openapi.OpenAPI, req batchDownload.Request, resolver xhttp.ResponseResolver) (*batchDownload.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.BatchDownload req:[%+v]", req)

	options := []apibatchdownload.Option{
		api.Client.Storage.BatchDownload.Paths(req.Paths...),
		api.Client.Storage.BatchDownload.FileName(req.FileName),
		api.Client.Storage.BatchDownload.BasePath(req.BasePath),
		api.Client.Storage.BatchDownload.IsCompress(req.IsCompress),
		api.Client.Storage.BatchDownload.WithResolver(resolver),
	}

	response, err := api.Client.Storage.BatchDownload(options...)

	if err != nil {
		logger.Errorf("invoke openapi.BatchDownload failed ，req:[%+v]，errMsg:[%v]", req, err)
	}

	return response, err
}
