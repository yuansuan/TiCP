package impl

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config 配置文件
type Config struct {
	CurrentEnvironment string `yaml:"current_environment"`
	Environments       map[string]EnvironmentConfig
}

// EnvironmentConfig 环境配置
type EnvironmentConfig struct {
	Endpoint               string `yaml:"endpoint"`
	ComputeYsID            string `yaml:"compute_ys_id"`
	ComputeAccessKeyID     string `yaml:"compute_access_key_id"`
	ComputeAccessKeySecret string `yaml:"compute_access_key_secret"`

	StorageYsID            string `yaml:"storage_ys_id"`
	StorageAccessKeyID     string `yaml:"storage_access_key_id"`
	StorageAccessKeySecret string `yaml:"storage_access_key_secret"`

	IamAdminEndpoint        string `yaml:"iam_endpoint"`
	IamAdminAccessKeyID     string `yaml:"iam_admin_access_key_id"`
	IamAdminAccessKeySecret string `yaml:"iam_admin_access_key_secret"`
}

var (
	// Cfg 配置
	Cfg *Config
	// CurrentCfg 当前环境配置
	CurrentCfg *EnvironmentConfig
)

// SaveConfig saves the configuration to file
func SaveConfig() error {
	return Save(Cfg)
}

// Save 保存配置
func Save(config *Config) error {
	viper.Set("current_environment", config.CurrentEnvironment)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// 创建默认config.yaml
func createDefaultConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.Set("current_environment", "default")
	viper.Set("environments.default.endpoint", "")
	viper.Set("environments.default.compute_ys_id", "")
	viper.Set("environments.default.compute_access_key_id", "")
	viper.Set("environments.default.compute_access_key_secret", "")
	viper.Set("environments.default.storage_ys_id", "")
	viper.Set("environments.default.storage_access_key_id", "")
	viper.Set("environments.default.storage_access_key_secret", "")
	viper.Set("environments.default.iam_endpoint", "")
	viper.Set("environments.default.iam_admin_access_key_id", "")
	viper.Set("environments.default.iam_admin_access_key_secret", "")

	if err := viper.WriteConfigAs("config.yaml"); err != nil {
		return err
	}

	fmt.Println("Create default config.yaml success")

	md := mapstructure.Metadata{}
	if err := viper.Unmarshal(Cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	}); err != nil {
		return err
	}

	return nil
}

// InitConfig 初始化配置
func InitConfig() {
	Cfg = new(Config)
	CurrentCfg = new(EnvironmentConfig)

	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("LoadConfigFail: ", err.Error())
		if err := createDefaultConfig(); err != nil {
			cobra.CheckErr(fmt.Errorf("CreateDefaultConfigFail: %s", err.Error()))
		}
	}

	md := mapstructure.Metadata{}
	if err := viper.Unmarshal(Cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	}); err != nil {
		cobra.CheckErr(fmt.Errorf("Unmarshal config fail: %s", err.Error()))
	}

	if Cfg.CurrentEnvironment == "" || Cfg.Environments == nil {
		fmt.Println("Config is empty")
		createDefaultConfig()
	}

	if env, ok := Cfg.Environments[Cfg.CurrentEnvironment]; !ok {
		cobra.CheckErr(fmt.Errorf("CurrentEnvironment %s not found", Cfg.CurrentEnvironment))
	} else {
		CurrentCfg = &env
	}

}

// CheckComputeConfig 检查Compute配置
func CheckComputeConfig() error {
	if CurrentCfg.ComputeAccessKeyID == "" || CurrentCfg.ComputeAccessKeySecret == "" || CurrentCfg.Endpoint == "" {
		return fmt.Errorf("compute_access_key_id or compute_access_key_secret or endpoint is empty, please check config.yaml")
	}
	return nil
}

// CheckComputeConfigExit 检查Compute配置
func CheckComputeConfigExit() {
	cobra.CheckErr(CheckComputeConfig())
}

// CheckStorageConfig 检查Storage配置
func CheckStorageConfig() error {
	if CurrentCfg.StorageAccessKeyID == "" || CurrentCfg.StorageAccessKeySecret == "" {
		return fmt.Errorf("storage_access_key_id or storage_access_key_secret is empty, please check config.yaml")
	}
	return nil
}

// CheckStorageConfigExit 检查Storage配置
func CheckStorageConfigExit() {
	cobra.CheckErr(CheckStorageConfig())
}

// CheckIamAdminConfig 检查IamAdmin配置
func CheckIamAdminConfig() error {
	if CurrentCfg.IamAdminAccessKeyID == "" || CurrentCfg.IamAdminAccessKeySecret == "" || CurrentCfg.IamAdminEndpoint == "" {
		return fmt.Errorf("iam_admin_access_key_id or iam_admin_access_key_secret or iam_endpoint is empty, please check config.yaml")
	}
	return nil
}

// CheckIamAdminConfigExit 检查IamAdmin配置
func CheckIamAdminConfigExit() {
	cobra.CheckErr(CheckIamAdminConfig())
}
