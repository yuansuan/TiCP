package storage

import (
	"fmt"
	"testing"

	apiInit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestUploadInit(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	rsp, err := UploadInit(api, apiInit.Request{
		Path: "/yskj//Abaqus.yaml",
		Size: 3082,
	})

	if err != nil {
		fmt.Printf("err:[%v]", err)
	}
	if rsp != nil {
		fmt.Printf("UPLOAD_ID:[%s]", rsp.UploadID)
	}

}
