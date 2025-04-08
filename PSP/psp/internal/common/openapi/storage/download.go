package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apidownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/download"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/download"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Download 文件下载
func Download(api *openapi.OpenAPI, req download.Request, resolver xhttp.ResponseResolver) (*download.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Download req:[%+v]", req)

	options := []apidownload.Option{
		api.Client.Storage.Download.Path(req.Path),
		api.Client.Storage.Download.WithResolver(resolver),
	}

	response, err := api.Client.Storage.Download(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Download failed ，req:[%+v]，errMsg:[%v]", req, err)
	}

	return response, err
}
