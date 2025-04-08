package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apiupinit "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/init"
	upinit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// UploadInit 文件上传初始化
func UploadInit(api *openapi.OpenAPI, req upinit.Request) (*upinit.Data, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.UploadInit req:[%+v]", req)

	options := []apiupinit.Option{
		api.Client.Storage.UploadInit.Path(req.Path),
		api.Client.Storage.UploadInit.Size(req.Size),
	}

	response, err := api.Client.Storage.UploadInit(options...)

	if err != nil {
		logger.Errorf("invoke openapi.UploadInit failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi upload init request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response.Data, err
}
