package job

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
)

func TestAdminTerminate(t *testing.T) {
	config.InitConfig()

	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	resp, err := AdminTerminate(api, "4W98PNzzbWC")
	if err != nil {
		fmt.Println("err: ", err)
	}

	println(fmt.Sprintf("==== terminate job: %+v", resp))
}
