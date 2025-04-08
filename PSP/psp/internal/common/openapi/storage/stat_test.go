package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestStat(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	if data, err := Stat(api, stat.Request{
		Path: "/1002/1111/Abaqus.yaml",
	}); err != nil {
		fmt.Println(data)
		fmt.Println(err)
	}

}
