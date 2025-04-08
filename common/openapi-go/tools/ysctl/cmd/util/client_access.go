package util

import (
	"fmt"
	"github.com/spf13/viper"
	openys "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/tools/ysctl/util/clientcmd"
)

type ApiClient struct {
}

var AccessKeyID string
var AccessKeySecret string
var Endpoint string
var StorageEndpoint string
var Proxy string
var Loglevel string

func NewApiClient() *ApiClient {
	return &ApiClient{}
}

func (c *ApiClient) RESTClient() (*openys.Client, error) {
	config := clientcmd.NewConfig()

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}
	akId, akSecret, endpoint := config.AccessKeyID, config.AccessKeySecret, config.Endpoint
	if AccessKeyID != "" {
		akId = AccessKeyID
	}
	if AccessKeySecret != "" {
		akSecret = AccessKeySecret
	}
	if Endpoint != "" {
		endpoint = Endpoint
	}
	return openys.NewClient(credential.NewCredential(akId, akSecret), openys.WithBaseURL(endpoint))
}

func (c *ApiClient) StorageClient() (*openys.Client, *clientcmd.Config, error) {
	config := clientcmd.NewConfig()

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}
	if AccessKeyID == "" {
		AccessKeyID = config.AccessKeyID
	}
	if AccessKeySecret == "" {
		AccessKeySecret = config.AccessKeySecret
	}
	if StorageEndpoint == "" {
		StorageEndpoint = config.StorageEndpoint
	}

	switch {
	case len(StorageEndpoint) == 0:
		return nil, nil, fmt.Errorf("empty storage endpoint")
	case len(AccessKeyID) == 0:
		return nil, nil, fmt.Errorf("empty access key id")
	case len(AccessKeySecret) == 0:
		return nil, nil, fmt.Errorf("emtpy access key secret")
	}

	client, err := openys.NewClient(credential.NewCredential(AccessKeyID, AccessKeySecret), openys.WithBaseURL(StorageEndpoint))

	return client, config, err
}
