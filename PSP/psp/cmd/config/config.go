package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/osutil"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/cmd/consts"
	app "github.com/yuansuan/ticp/PSP/psp/internal/app/config"
	approve "github.com/yuansuan/ticp/PSP/psp/internal/approve/config"
	gateway "github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/config"
	openapi "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
	job "github.com/yuansuan/ticp/PSP/psp/internal/job/config"
	monitor "github.com/yuansuan/ticp/PSP/psp/internal/monitor/config"
	notice "github.com/yuansuan/ticp/PSP/psp/internal/notice/config"
	project "github.com/yuansuan/ticp/PSP/psp/internal/project/config"
	rbac "github.com/yuansuan/ticp/PSP/psp/internal/rbac/config"
	storage "github.com/yuansuan/ticp/PSP/psp/internal/storage/config"
	sysconfig "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/config"
	user "github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	visual "github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
)

const PSPName = "psp"

type CustomConfig struct {
	Main    Main                   `yaml:"main"`
	App     app.CustomConfig       `yaml:"app"`
	Rbac    rbac.CustomConfig      `yaml:"rbac"`
	Job     job.CustomConfig       `yaml:"job"`
	Notice  notice.CustomConfig    `yaml:"notice"`
	Visual  visual.CustomConfig    `yaml:"visual"`
	OpenApi openapi.CustomConfig   `yaml:"openapi"`
	Storage storage.CustomConfig   `yaml:"storage"`
	GateWay gateway.CustomConfig   `yaml:"gateway"`
	Monitor monitor.CustomConfig   `yaml:"monitor"`
	User    user.CustomConfig      `yaml:"user"`
	Project project.CustomConfig   `yaml:"project"`
	System  sysconfig.CustomConfig `yaml:"system"`
	Logger  *Logger                `yaml:"logger"`
	Approve approve.CustomConfig   `yaml:"approve"`
}

type Main struct {
	Swagger          Swagger           `yaml:"swagger"`
	EnableVisual     bool              `yaml:"enable_visual"`
	ComputeTypeNames map[string]string `yaml:"compute_type_names"`
}

type Logger struct {
	MaxSize     int    `yaml:"max_size"`
	BackupCount int    `yaml:"backup_count"`
	LogDir      string `yaml:"log_dir"`
	MaxAge      int    `yaml:"max_age"`
}

type Swagger struct {
	Enable bool   `yaml:"enable"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
}

var (
	mutex  sync.Mutex
	Custom CustomConfig
)

// GetConfig 获取配置
func GetConfig() CustomConfig {
	mutex.Lock()
	defer mutex.Unlock()
	return Custom
}

func InitConfig() error {

	newViper, err := getNewViper()
	if err != nil {
		return err
	}

	md := mapstructure.Metadata{}
	err = newViper.Unmarshal(&Custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		return err
	}

	// set Custom here
	rbac.Custom = Custom.Rbac
	app.SetConfig(Custom.App)
	job.SetConfig(Custom.Job)
	notice.SetConfig(Custom.Notice)
	visual.SetConfig(Custom.Visual)
	openapi.SetConfig(Custom.OpenApi)
	storage.SetConfig(Custom.Storage)
	gateway.SetConfig(Custom.GateWay)
	monitor.SetConfig(Custom.Monitor)
	user.SetConfig(Custom.User)
	project.SetConfig(Custom.Project)
	sysconfig.SetConfig(Custom.System)
	approve.SetConfig(Custom.Approve)
	logging.Default().Infof("%v", Custom)

	return nil
}

func getNewViper() (*viper.Viper, error) {

	newViper := viper.New()

	// Get the YS_TOP from env
	ysTop := os.Getenv(consts.EnvVarYSTop)
	if strings.TrimSpace(ysTop) == "" {
		panic("YS_TOP env is not set.")
	}

	logging.Default().Debugf("%v=%v", consts.EnvVarYSTop, ysTop)

	newViper.SetConfigType(consts.SysConfigType)

	logging.Default().Debugf("mode=%v", env.ModeName(env.Env.Mode))
	configName := fmt.Sprintf("%v_custom", env.ModeName(env.Env.Mode))

	newViper.SetConfigName(configName)

	configPath := filepath.Join(ysTop, PSPName, config.ConfigDir)
	newViper.AddConfigPath(configPath)

	err := newViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return newViper, nil
}

func GetLoaderLogLevel() string {

	// Get the log level from psp.conf
	ysTop := os.Getenv(consts.EnvVarYSTop)
	pspConf := fmt.Sprintf("%v/%v/%v/psp.conf", ysTop, PSPName, config.ConfigDir)
	cmd := fmt.Sprintf("cat %v | grep \"^LOG_LEVEL=.*\" | awk -F= '{print $2}'", pspConf)
	stdout, _, err := osutil.CommandHelper.BashWithCurrent(context.TODO(), cmd)
	logLevel := ""
	if err != nil {
		logLevel = "error"
	} else {
		logLevel = strings.ReplaceAll(string(stdout), "\n", "")
	}
	logLevel = strings.ToLower(logLevel)

	if logLevel != "info" && logLevel != "warn" && logLevel != "error" && logLevel != "debug" {
		logLevel = "error"
	}

	return logLevel
}
