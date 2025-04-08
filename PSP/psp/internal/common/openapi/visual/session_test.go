package openapivisual

import (
	"testing"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func TestGetApp(t *testing.T) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Ready(
		client.Ready.Id("599bgw1nZto"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(response)

}
