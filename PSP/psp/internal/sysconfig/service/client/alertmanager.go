package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	sysconfig "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/config"
)

func AlertManagerReload() {
	ctx := context.Background()
	alertManagerUrl := sysconfig.GetConfig().AlertManager.AlertManagerUrl
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/-/reload", alertManagerUrl), nil)
	if err != nil {
		logging.GetLogger(ctx).Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	alertManagerClient := &http.Client{}
	resp, err := alertManagerClient.Do(req)
	if err != nil {
		logging.GetLogger(ctx).Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
}
