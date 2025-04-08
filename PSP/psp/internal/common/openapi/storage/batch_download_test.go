package storage

import (
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestBatchDownload(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	data, err := BatchDownload(api, batchDownload.Request{
		Paths: []string{
			"/4TiSBX39DtN/yskj/kk",
		},
		FileName: "kk.zip",
	}, nil)

	if err != nil {
		panic("")
	}

	println(data)
}
