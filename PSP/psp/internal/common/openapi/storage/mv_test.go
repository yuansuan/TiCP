package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestMv(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	_, err := Mv(api, mv.Request{
		SrcPath:  "/4Afa3ivYikw/yskj/helloworld",
		DestPath: "/4Afa3ivYikw/yskj/helloworld2",
	})

	fmt.Printf("err:[%v]", err)
}
