package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	pbcompressstart "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/compress/start"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/start"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func CompressStart(api *openapi.OpenAPI, req start.Request) (*start.Data, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.compress start req:[%+v]", req)

	options := []pbcompressstart.Option{
		api.Client.Storage.CompressStart.Paths(req.Paths...),
		api.Client.Storage.CompressStart.TargetPath(req.TargetPath),
		api.Client.Storage.CompressStart.BasePath(req.BasePath),
	}
	response, err := api.Client.Storage.CompressStart(options...)

	if err != nil {
		logger.Errorf("invoke openapi.compress start failed req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logging.Default().Debugf("openapi check compress start id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response.Data, err
}
