package config

import (
	"context"
	"sync"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type CustomT struct {
	OpenAPI                      OpenAPI       `yaml:"openapi"`
	BillEnabled                  bool          `yaml:"bill_enabled"`
	CloudApp                     CloudApp      `yaml:"cloud_app"`
	CreateSessionTimeoutForAlarm time.Duration `yaml:"create_session_timeout_for_alarm"`
}

type OpenAPI struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Endpoint        string `yaml:"endpoint"`
}

type CloudApp struct {
	SignalServerAddress string        `yaml:"signal_server_address"`
	SignalHost          string        `yaml:"signal_host"`
	WebClientConf       WebClientConf `yaml:"web_client"`
	Zones               Zones         `yaml:"zones"`
}

type WebClientConf struct {
	BaseURL     string            `yaml:"base_url"`
	QueryParams map[string]string `yaml:"query_params"`
}

type Zones map[Zone]*ZoneOption

type CloudType string

const (
	TencentCloudType   CloudType = "tencent"
	ShanheCloudType    CloudType = "shanhe"
	OpenstackCloudType CloudType = "openstack"
)

type ZoneOption struct {
	AccessOrigin string `yaml:"access_origin"`
	GuacdAddress string `yaml:"guacd_address"`

	Cloud CloudType `yaml:"cloud"`
	// 一个区域对应一个云
	Tencent   *Tencent   `yaml:"tencent,omitempty"`
	Shanhe    *ShanHe    `yaml:"shanhe,omitempty"`
	OpenStack *OpenStack `yaml:"openstack,omitempty"`
}

// Tencent 腾讯云配置
type Tencent struct {
	SecretID         string              `yaml:"secret_id"`
	SecretKey        string              `yaml:"secret_key"`
	Region           string              `yaml:"region"`
	Cluster          string              `yaml:"cluster"`
	VpcSelector      map[string]string   `yaml:"vpc_selector"`
	ZoneAffinity     map[string]Affinity `yaml:"zone_affinity"`
	SecurityGroupIds []*string           `yaml:"security_group_ids"`
	SystemDiskSize   int64               `yaml:"system_disk_size"`
}

// ShanHe 山河云配置
type ShanHe struct {
	SecretID  string    `yaml:"secret_id"`
	SecretKey string    `yaml:"secret_key"`
	Region    string    `yaml:"region"`
	Endpoint  string    `yaml:"endpoint"`
	VxNets    []*string `yaml:"vx_nets"`
	Cluster   string    `yaml:"cluster"`
}

// OpenStack openstack配置
type OpenStack struct {
	Auth                 Auth     `yaml:"auth"`
	Compute              Compute  `yaml:"compute"`
	Network              Network  `yaml:"network"`
	Tags                 []string `yaml:"tags"`
	CreateWithBootVolume bool     `yaml:"create_with_boot_volume"`
}

type Auth struct {
	IdentityEndpoint string `yaml:"identity_endpoint"`
	CredentialID     string `yaml:"credential_id"`
	CredentialSecret string `yaml:"credential_secret"`
}

type Compute struct {
	NovaEndpoint string `yaml:"nova_endpoint"`
	MicroVersion string `yaml:"micro_version"`
}

type Network struct {
	Name string `yaml:"name"`
	Uuid string `yaml:"uuid"`
}

// Affinity 每个可用区对应的亲和度
type Affinity map[string]int

// Storage 存储共享服务
type Storage struct {
	Server string `yaml:"server"`
}

// Mut Mut
var Mut = &sync.Mutex{}

// Custom Custom
var custom CustomT

// GetConfig GetConfig
func GetConfig() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return custom
}

// SetConfig SetConfig
func SetConfig(cfg CustomT) {
	Mut.Lock()
	defer Mut.Unlock()
	custom = cfg
}

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			custom = *data.(*CustomT)
			logger.Info(custom)
		}

		path := common.GetRemoteConfigPath()
		logger.Infof("start read remote config: %v", path)
		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true, path, &CustomT{}, switchFunc)
		return err
	}

	Mut.Lock()
	defer Mut.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		return err
	}

	return nil
}

type Zone string

func (z Zone) IsValid() bool {
	return z != AZInvalid
}

// IsEmpty 是否为空可用区
func (z Zone) IsEmpty() bool {
	return z == AZEmpty
}

// String 工具函数，支持语言的默认行为
func (z Zone) String() string {
	return string(z)
}

const (
	// AZInvalid 无效的区域
	AZInvalid Zone = "az-invalid"
	// AZEmpty 空区域
	AZEmpty Zone = "az-empty"
)

func Parse(s string) Zone {
	// foreach zone
	for z := range GetConfig().CloudApp.Zones {
		if z.String() == s {
			return z
		}
	}

	switch Zone(s) {
	case "", AZEmpty:
		return AZEmpty
	}
	return AZInvalid
}
