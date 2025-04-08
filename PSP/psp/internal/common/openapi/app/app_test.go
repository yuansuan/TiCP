package openapiapp

import (
	"strconv"
	"testing"
	"time"

	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestAddApp(t *testing.T) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	response, err := api.Client.Job.AdminAddAPP(
		api.Client.Job.AdminAddAPP.Name(strconv.FormatInt(time.Now().UnixMilli(), 10)),
		api.Client.Job.AdminAddAPP.Type("test"),
		api.Client.Job.AdminAddAPP.Version("0.0.1"),
		api.Client.Job.AdminAddAPP.Image("f421ca42-0fd9-462f-ba6a-03b800ca90de"),
		api.Client.Job.AdminAddAPP.BinPath(map[string]string{
			"az-jinan":    "https://www.baidu.com/img/bd_logo1.png",
			"az-shanghai": "https://www.baidu.com/img/bd_logo2.png",
		}),
	)

	if response != nil {
		t.Log(response)
	}
}

func TestUpdateApp(t *testing.T) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	response, err := api.Client.Job.AdminUpdateAPP(
		api.Client.Job.AdminUpdateAPP.AppID("4W76kK68q5Q"),
		api.Client.Job.AdminUpdateAPP.Name("sadfasdfas"),
		api.Client.Job.AdminUpdateAPP.Type("tessdft"),
		api.Client.Job.AdminUpdateAPP.Version("sadfas0.0.1"),
		api.Client.Job.AdminUpdateAPP.Image("f421ca42-0fd9-462f-ba6a-03b800ca90de"),
		api.Client.Job.AdminUpdateAPP.LicManagerId("f421ca42-0fd9-462f-ba6a-03b800ca90de"),
		api.Client.Job.AdminUpdateAPP.PublishStatus(update.Published),
		api.Client.Job.AdminUpdateAPP.BinPath(map[string]string{
			"az-jinan":    "https://www.baidu.com/img/bd_logo1.png",
			"az-shanghai": "https://www.baidu.com/img/bd_logo2.png",
		}),
	)

	if response != nil {
		t.Log(response)
	}
}

func TestListApp(t *testing.T) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	response, err := api.Client.Job.AdminListAPP()

	if response != nil {
		t.Log(response)
	}
}

func TestGetApp(t *testing.T) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	response, err := api.Client.Job.AdminGetAPP(
		api.Client.Job.AdminGetAPP.AppID("4W7gCJYm2go"),
	)

	if response != nil {
		t.Log(response)
	}
}
