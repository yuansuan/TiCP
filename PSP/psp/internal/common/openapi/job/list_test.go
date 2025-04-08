package job

import (
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
)

func TestAdminListJobs(t *testing.T) {
	config.InitConfig()

	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	resp, err := AdminListJobs(api, "az-jinan", "", 1, 10)
	if err != nil {
		fmt.Println("err: ", err)
	}

	for _, job := range resp.Data.Jobs {
		println(fmt.Sprintf("==== job: %+v", job))
	}

	//assert.Equal(t, jobs.Total, int64(46))
}
