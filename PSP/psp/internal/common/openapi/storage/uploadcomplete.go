package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apiupcomplete "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/complete"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// UploadComplete 文件上传完成
func UploadComplete(api *openapi.OpenAPI, req complete.Request) (*complete.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.UploadComplete req:[%+v]", req)

	options := []apiupcomplete.Option{
		api.Client.Storage.UploadComplete.UploadID(req.UploadID),
		api.Client.Storage.UploadComplete.Path(req.Path),
	}
	response, err := api.Client.Storage.UploadComplete(options...)

	if err != nil {
		logger.Errorf("invoke openapi.UploadComplete failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi upload complete request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
