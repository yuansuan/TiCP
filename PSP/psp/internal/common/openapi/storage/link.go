package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apilink "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/link"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Link(api *openapi.OpenAPI, req link.Request) (*link.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Link req:[%+v]", req)

	options := []apilink.Option{
		api.Client.Storage.Link.Src(req.SrcPath),
		api.Client.Storage.Link.Dest(req.DestPath),
	}

	response, err := api.Client.Storage.Link(options...)

	if err != nil {
		logger.Errorf("invoke openapi.link failed ，req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logger.Debugf("openapi link request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
