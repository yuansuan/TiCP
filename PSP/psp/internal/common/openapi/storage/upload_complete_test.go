package storage

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestUploadCmplete(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	if data, err := UploadComplete(api, complete.Request{
		Path:     "1002/1111/Abaqus.yaml",
		UploadID: "eef5b2a0-27c3-462f-8448-47ee47ce6fa2",
	}); err != nil {
		fmt.Println(data)
		fmt.Println(err)
	}

}
