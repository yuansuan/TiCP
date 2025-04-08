package webrtc

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

func GenerateDesktopURLBase64(roomId string) (string, error) {
	cfg := config.GetConfig()
	baseUrl := cfg.CloudApp.WebClientConf.BaseURL
	params := cfg.CloudApp.WebClientConf.QueryParams

	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("parse url %s failed, %w", baseUrl, err)
	}

	queries := u.Query()
	queries.Add("room_id", roomId)
	queries.Add("signal", cfg.CloudApp.SignalHost)
	for k, v := range params {
		queries.Add(k, v)
	}

	u.RawQuery = queries.Encode()

	return base64.StdEncoding.EncodeToString([]byte(u.String())), nil
}
