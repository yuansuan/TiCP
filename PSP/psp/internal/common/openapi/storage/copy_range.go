package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apicprange "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/copyRange"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copyRange"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func CopyRange(api *openapi.OpenAPI, req copyRange.Request) (*copyRange.Response, error) {

	logger := logging.Default()
	logger.Debugf("invoke openapi.CopyRange req:[%+v]", req)

	options := []apicprange.Option{
		api.Client.Storage.CopyRange.Src(req.SrcPath),
		api.Client.Storage.CopyRange.Dest(req.DestPath),
		api.Client.Storage.CopyRange.SrcOffset(req.SrcOffset),
		api.Client.Storage.CopyRange.DestOffset(req.DestOffset),
		api.Client.Storage.CopyRange.Length(req.Length),
	}

	response, err := api.Client.Storage.CopyRange(options...)

	if err != nil {
		logger.Errorf("invoke openapi.CopyRange failed ，req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi CopyRange request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err

}
