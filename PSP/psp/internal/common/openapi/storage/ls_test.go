package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestLs(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	data, err := Ls(api, ls.Request{
		Path:             "/4WMkNCQeWYQ/yskj",
		FilterRegexpList: []string{"^[.|_]"},
		PageSize:         1000,
		PageOffset:       0,
	})

	if err == nil && data != nil {
		for _, file := range data.Data.Files {
			fmt.Println(file)
		}
	}

}
