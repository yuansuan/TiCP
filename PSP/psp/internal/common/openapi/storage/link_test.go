package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestLink(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	_, err := Link(api, link.Request{
		SrcPath:  "/4TpFFZDkFWy/yskj/Q3.sim",
		DestPath: "/4TpFFZDkFWy/.tmp_upload/4X8bvHrzgyd/Q3.sim",
	})

	fmt.Printf("err:[%v]", err)
}
