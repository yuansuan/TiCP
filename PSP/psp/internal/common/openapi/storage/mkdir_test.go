package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestMkir(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	_, err := Mkdir(api, mkdir.Request{
		Path: "/4Afa3ivYikw/yskj/.starccm_test2",
	})

	fmt.Printf("err:[%v]", err)
}
