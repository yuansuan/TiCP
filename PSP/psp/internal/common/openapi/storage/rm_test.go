package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/rm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestRm(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	_, err := Rm(api, rm.Request{
		Path: "/yskj/111",
	})

	fmt.Printf("err:[%v]", err)
}
