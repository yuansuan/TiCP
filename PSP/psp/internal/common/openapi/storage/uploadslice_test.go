package storage

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestUpload(t *testing.T) {
	api, _ := openapi.NewLocalHPCAPI()

	// 调用运算云的openapi
	open, _ := os.Open("/Users/yskj/Documents/test/Abaqus.yaml")
	byte, _ := io.ReadAll(open)
	defer open.Close()

	_, err := UploadSlice(api, slice.Request{
		UploadID: "eef5b2a0-27c3-462f-8448-47ee47ce6fa2",
		Offset:   0,
		Length:   3082,
		Slice:    byte,
	})

	if err != nil {
		fmt.Println(err)
	}

}
