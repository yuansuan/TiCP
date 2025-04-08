package storage

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apiupslice "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/slice"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

// UploadSlice 文件切片上传
func UploadSlice(api *openapi.OpenAPI, req slice.Request) (*slice.Response, error) {
	logger := logging.Default()
	logger.Debugf("invoke openapi.UploadSlice req: UploadID:[%v], Length:[%v], Offset:[%v]", req.UploadID, req.Length, req.Offset)

	options := []apiupslice.Option{
		api.Client.Storage.UploadSlice.UploadID(req.UploadID),
		api.Client.Storage.UploadSlice.Offset(req.Offset),
		api.Client.Storage.UploadSlice.Slice(req.Slice),
		api.Client.Storage.UploadSlice.Length(req.Length),
	}
	response, err := api.Client.Storage.UploadSlice(options...)

	if err != nil {
		logger.Errorf("invoke openapi.UploadSlice failed req: UploadID:[%v], Lenth:[%v], Offset:[%v]，errMsg:[%v]", req.UploadID, req.Length, req.Offset, err)
	}

	if response != nil {
		logger.Debugf("openapi upload slice request id: [%v], uploadID:[%v], length:[%v], offset:[%v]",
			response.RequestID, req.UploadID, req.Length, req.Offset)
	}

	return response, err
}
