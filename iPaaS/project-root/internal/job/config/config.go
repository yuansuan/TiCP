package config

import (
	"context"
	"os"
	"sync"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/mongo"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// CustomT CustomT
type CustomT struct {
	ChangeLicense bool `yaml:"change_license"`

	Zones schema.Zones `yaml:"zones"`

	SelfYsID        string `yaml:"self_ys_id"`
	AK              string `yaml:"ak"`
	AS              string `yaml:"as"`
	OpenAPIEndpoint string `yaml:"openapi_endpoint"`

	// 计费开关
	BillEnabled bool `yaml:"bill_enabled"`

	// 残差解析读取文件大小上限
	ResidualMaxFileSize int64 `yaml:"residual_max_file_size"`
	// 残差解析读取文件大小上限
	MonitorChartMaxFileSize int64 `yaml:"monitor_chart_max_file_size"`

	SelectorWeights map[string]float64 `yaml:"selector_weights"`

	// mongoDB
	Mongo *mongo.Config `yaml:"mongo"`

	// 长时间运行作业告警阈值
	LongRunningJobThreshold int64 `yaml:"long_running_job_threshold"`
	// WebhookURL
	WebhookURL string `yaml:"webhook_url"`
}

// Mut Mut
var Mut = &sync.Mutex{}

// Custom Custom
var Custom CustomT

// String String
func (c CustomT) String() string {
	bs, _ := yaml.Marshal(&c)

	return string(bs)
}

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			logger.Info(Custom)
		}

		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true,
			common.GetRemoteConfigPath(), &CustomT{}, switchFunc)
		return err
	}

	Mut.Lock()
	defer Mut.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&Custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})

	if err != nil {
		return err
	}

	return nil
}

// GetConfig return the config
func GetConfig() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}

// GetAK 获取ak
func (c *CustomT) GetAK() string {
	ak := os.Getenv("AK")
	if ak != "" {
		return ak
	}
	return c.AK
}

// GetAS 获取as
func (c *CustomT) GetAS() string {
	as := os.Getenv("AS")
	if as != "" {
		return as
	}
	return c.AS
}
