package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apils "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/ls"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// Ls 查询文件列表
func Ls(api *openapi.OpenAPI, req ls.Request) (*ls.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.Ls req:[%+v]", req)

	options := []apils.Option{
		api.Client.Storage.LsWithPage.Path(req.Path),
		api.Client.Storage.LsWithPage.FilterRegexpList(req.FilterRegexpList),
		api.Client.Storage.LsWithPage.PageOffset(req.PageOffset),
		api.Client.Storage.LsWithPage.PageSize(req.PageSize),
	}

	response, err := api.Client.Storage.LsWithPage(options...)

	if err != nil {
		logger.Errorf("invoke openapi.Ls failed ,req:[%+v]，errMsg:[%v]", req, err)
	}

	if response != nil {
		logging.Default().Debugf("openapi ls request id: [%v], req: [%+v]", response.RequestID, req)
	}

	return response, err
}
