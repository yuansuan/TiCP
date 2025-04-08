package openapi

import (
	"errors"
	"sync"

	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
)

const (
	EndpointTypeAPI = iota // 用于可视化、计算应用等
	EndpointTypeHPC        // 用于存储调用
)

type OpenAPI struct {
	Client *openapi.Client
}

var once, storageOnce, cloudOnce, cloudStorageOnce sync.Once

var instance, storageInstance, cloudInstance, cloudStorageInstance *OpenAPI

func NewLocalAPI() (api *OpenAPI, err error) {
	once.Do(func() {
		instance, err = getInstance(EndpointTypeAPI)
	})

	if instance == nil {
		return nil, errors.New("failed to create local openapi instance")
	}

	return instance, nil
}

func NewLocalHPCAPI() (api *OpenAPI, err error) {
	storageOnce.Do(func() {
		storageInstance, err = getInstance(EndpointTypeHPC)
	})

	if storageInstance == nil {
		return nil, errors.New("failed to create local storage openapi instance")
	}

	return storageInstance, nil
}

// getInstance 根据配置信息初始化 openapi client
func getInstance(apiType uint8) (*OpenAPI, error) {
	config.InitConfig()
	cfg := config.GetConfig()

	if cfg.Local == nil {
		return nil, nil
	}

	conf := cfg.Local.Settings
	endpoint, err := getEndpoint(conf, apiType)
	if err != nil {
		return nil, err
	}

	client, err := openapi.NewClient(credential.NewCredential(conf.AppKey, conf.AppSecret), openapi.WithBaseURL(endpoint))

	return &OpenAPI{
		Client: client,
	}, err
}

func getEndpoint(settings *config.Settings, apiType uint8) (string, error) {
	switch apiType {
	case EndpointTypeAPI:
		return settings.Endpoint, nil
	case EndpointTypeHPC:
		return settings.HPCEndpoint, nil
	default:
		return "", errors.New("endpoint type [%v] not support")
	}
}
