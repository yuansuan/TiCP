package webrtc

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

func TestGenerateDesktopURL(t *testing.T) {
	config.SetConfig(config.CustomT{
		CloudApp: config.CloudApp{
			SignalHost: "cloudapp-signal.intern.yuansuan.cn",
			WebClientConf: config.WebClientConf{
				BaseURL: "https://cloudapp-web-client.intern.yuansuan.cn/index.html",
			},
		},
	})

	urlB64, err := GenerateDesktopURLBase64("roomIdA")
	assert.NoError(t, err)
	t.Logf("urlB64: %s", urlB64)

	url, err := base64.StdEncoding.DecodeString(urlB64)
	assert.NoError(t, err)
	t.Logf("url: %s", string(url))
}
