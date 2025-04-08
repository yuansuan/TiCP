package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestRealpath(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	data, _ := Realpath(api, &realpath.Request{
		RelativePath: "/4TiSBX39DtN/yskj/workspace/",
	})

	fmt.Println("real_path: ", data.Data.RealPath)
}
