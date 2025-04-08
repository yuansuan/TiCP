package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apidownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/download"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/download"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// DownloadRange 节选文件下载
func DownloadRange(api *openapi.OpenAPI, path string, beginOffset, endOffset int64) (*download.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.DownloadRange path:[%v]， beginOffset:[%v]，endOffset:[%v]，", path, beginOffset, endOffset)

	options := []apidownload.Option{
		api.Client.Storage.Download.Path(path),
	}
	if endOffset != 0 {
		options = append(options, api.Client.Storage.Download.Range(beginOffset, endOffset))
	}

	response, err := api.Client.Storage.Download(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Download failed,path:[%v]， beginOffset:[%v]，endOffset:[%v], err:[%v]", path, beginOffset, endOffset, err)
	}

	return response, err
}
