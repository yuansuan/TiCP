package openapi

import (
	"fmt"

	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
)

func NewClient() (*openapi.Client, error) {
	conf := config.GetConfig()
	if conf.AK == "" {
		return nil, fmt.Errorf("openapi.access_key_id cannot be empty")
	}
	if conf.AS == "" {
		return nil, fmt.Errorf("openapi.access_key_secret cannot be empty")
	}
	if conf.OpenAPIEndpoint == "" {
		return nil, fmt.Errorf("openapi.endpoint cannot be empty")
	}
	opts := []openapi.Option{
		openapi.WithBaseURL(conf.OpenAPIEndpoint),
		openapi.WithTimeout(module.DefaultTimeout),
		openapi.WithRetryTimes(module.DefaultRetryTimes),
		openapi.WithRetryInterval(module.DefaultRetryInterval),
	}
	openapiClient, err := openapi.NewClient(credential.NewCredential(conf.AK, conf.AS), opts...)
	if err != nil {
		return nil, fmt.Errorf("new openapi client failed, %w", err)
	}
	return openapiClient, nil
}
