package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xurl"
)

const (
	defaultConfigPath = "./config/config.yaml"
)

var _config *Config

// GetConfig be care for nil pointer, _config initialized when NewConfig() called
func GetConfig() *Config {
	return _config
}

// Config 服务自定义配置项
type Config struct {
	HttpAddress        string           `yaml:"http_address"`
	PerformanceAddress string           `yaml:"performance_address"`
	HpcStorageAddress  string           `yaml:"hpc_storage_address"`
	OpenAPI            OpenAPIConfig    `yaml:"openapi"`
	Iam                IamConfig        `yaml:"iam"`
	AccessLog          AccessLog        `yaml:"access_log"`
	Log                Log              `yaml:"log"`
	Database           *Database        `yaml:"database"`
	Migrations         *Migrations      `yaml:"migrations"`
	BackendProvider    *BackendProvider `yaml:"backend-provider"`
	Singularity        *Singularity     `yaml:"singularity"`
	StateMachine       *StateMachine    `yaml:"state-machine"`
	Snowflake          Snowflake        `yaml:"snowflake"`
	PreparedFilePath   string           `yaml:"prepared_file_path"`
	Sync               SyncConfig       `yaml:"sync"`
}

type OpenAPIConfig struct {
	MaxRetryTimes int           `yaml:"max_retry_times"`
	RetryInterval time.Duration `yaml:"retry_interval"`
	Proxy         string        `yaml:"proxy"` //用于设置代理
}

type IamConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
	YsID      string `yaml:"ys_id"`
	Proxy     string `yaml:"proxy"`
}

type AccessLog struct {
	Path       string `yaml:"path"`
	UseConsole bool   `yaml:"use_console"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

type Log struct {
	Path         string `yaml:"path"`
	Level        string `yaml:"level"`
	ReleaseLevel string `yaml:"release_level"`
	UseConsole   bool   `yaml:"use_console"`
	MaxSize      int    `yaml:"max_size"`
	MaxAge       int    `yaml:"max_age"`
	MaxBackups   int    `yaml:"max_backups"`
}

// Database 数据库配置项
type Database struct {
	Type      string `yaml:"type"`
	DSN       string `yaml:"dsn"`
	HiddenSQL bool   `yaml:"hidden_sql"`
}

type StateMachine struct {
	Channel string `yaml:"channel"`
}

type BackendProvider struct {
	Type               string                 `yaml:"type"`
	SchedulerCommon    SchedulerCommon        `yaml:"scheduler-common"`
	Slurm              *SlurmBackendProvider  `yaml:"slurm"`
	PbsPro             *PbsProBackendProvider `yaml:"pbs-pro"`
	Mock               *MockBackendProvider   `yaml:"mock"`
	CheckAliveInterval int                    `yaml:"check-alive-interval"`
}

type SlurmBackendProvider struct {
	Submit        string `yaml:"submit"`
	SubmitAverage string `yaml:"submit-average"`
	Kill          string `yaml:"kill"`
	CheckAlive    string `yaml:"check-alive"`
	CheckHistory  string `yaml:"check-history"`
	GetResource   string `yaml:"get-resource"`
	JobIdRegex    string `yaml:"job-id-regex"`
}

type PbsProBackendProvider struct {
	Submit      string `yaml:"submit"`
	Kill        string `yaml:"kill"`
	CheckAlive  string `yaml:"check-alive"`
	GetResource string `yaml:"get-resource"`
}

type MockBackendProvider struct {
}

// Singularity 镜像相关配置
type Singularity struct {
	Storage  string                `yaml:"storage"`
	IsMock   bool                  `yaml:"is-mock"`
	Registry *ObjectStorageService `yaml:"registry"`
}

// ObjectStorageService 对象存储服务配置
type ObjectStorageService struct {
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Region       string `yaml:"region"`
	Endpoint     string `yaml:"endpoint"`
	Bucket       string `yaml:"bucket"`
	PathPrefix   string `yaml:"path_prefix"`
}
type CoreConfig struct {
	Name string `yaml:"name"`
	Core int    `yaml:"core"`
}

type SchedulerCommon struct {
	DefaultQueue      string         `yaml:"default-queue"`
	CandidateQueues   []string       `yaml:"candidate-queues"`
	Workspace         string         `yaml:"workspace"`
	CoresPerNode      map[string]int `yaml:"-"`
	ReservedCores     map[string]int `yaml:"-"` // 每个队列预留资源，未设置为0，表示不预留
	CoresPerNodeList  []CoreConfig   `yaml:"cores-per-node-list"`
	ReservedCoresList []CoreConfig   `yaml:"reserved-cores-list"`
	SubmitSysUser     string         `yaml:"submit-sys-user"`
	SubmitSysUserUid  int            `yaml:"submit-sys-user-uid"`
	SubmitSysUserGid  int            `yaml:"submit-sys-user-gid"`
}

// SyncConfig 传输参数配置
type SyncConfig struct {
	Compressor string `yaml:"compressor"`
}

// BaseURL 返回该仓库的URL地址
func (oss *ObjectStorageService) BaseURL() string {
	return "https://" + oss.Bucket + "." + oss.Endpoint + "/" + oss.PathPrefix
}

// RelativePath 获取仓库中文件的相对路径
func (oss *ObjectStorageService) RelativePath(paths ...string) string {
	return xurl.Join(append([]string{oss.PathPrefix}, paths...)...)
}

// Migrations 数据库自动升级配置
type Migrations struct {
	AutoMigration    bool   `yaml:"auto-migration"`
	MigrationVersion string `yaml:"migration-version"`
}

type Snowflake struct {
	Node int64 `yaml:"node"`
}

type CustomConfig struct {
	path string
}

type Option interface {
	apply(c *CustomConfig)
}

type optionFunc func(c *CustomConfig)

func (f optionFunc) apply(c *CustomConfig) {
	f(c)
}

func withDefaultOption() Option {
	return optionFunc(func(c *CustomConfig) {
		c.path = defaultConfigPath
	})
}

func WithPath(path string) Option {
	return optionFunc(func(c *CustomConfig) {
		if path == "" {
			return
		}

		c.path = path
	})
}

// NewConfig 从配置文件加载服务配置
func NewConfig(opts ...Option) (*Config, error) {
	customConf := new(CustomConfig)
	withDefaultOption().apply(customConf)
	for _, opt := range opts {
		opt.apply(customConf)
	}

	fmt.Printf("init config from config [%s]\n", customConf.path)

	dir, file := filepath.Split(customConf.path)
	fileName, fileExt := fileSplit(file)

	cfg := new(Config)
	viper.AddConfigPath(dir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileExt)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &mapstructure.Metadata{}
	})

	// 遍历 CoresPerNodeList 和 ReservedCoresList 填充 map
	cfg.BackendProvider.SchedulerCommon.CoresPerNode = make(map[string]int)
	for _, item := range cfg.BackendProvider.SchedulerCommon.CoresPerNodeList {
		cfg.BackendProvider.SchedulerCommon.CoresPerNode[item.Name] = item.Core
	}

	cfg.BackendProvider.SchedulerCommon.ReservedCores = make(map[string]int)
	for _, item := range cfg.BackendProvider.SchedulerCommon.ReservedCoresList {
		cfg.BackendProvider.SchedulerCommon.ReservedCores[item.Name] = item.Core
	}

	_config = cfg

	return cfg, err
}

func fileSplit(file string) (name, ext string) {
	ext = filepath.Ext(file)

	return strings.TrimSuffix(file, fmt.Sprintf(".%s", ext)), strings.TrimPrefix(ext, ".")
}
