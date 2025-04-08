package config

import (
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

const (
	RewriteFileName = "config_rewrite.yml"
)

// Config Config
type Config struct {
	App struct {
		Name                 string
		Ver                  string
		Host                 string
		Port                 int
		LoadRemoteConfig     bool   `yaml:"load_remote_config"`
		NetworkInterfaceName string `yaml:"network_if_name"`
		PprofHost            string `yaml:"pprof_host"`
		PprofPort            int    `yaml:"pprof_port"`
		DBMigration          struct {
			AutoMigrate  bool   `yaml:"auto_migrate"`
			ForceMigrate *bool  `yaml:"force_migrate"`
			Version      string `yaml:"version"`
		} `yaml:"db_migration"`

		Middleware struct {
			Logger  conf_type.Logger
			HTTP    conf_type.HTTP
			GRPC    conf_type.GRPC
			Monitor conf_type.Monitor
			Cache   conf_type.Caches
			APM     conf_type.APM
			Tracing conf_type.Tracing

			Mysql         conf_type.Mysqls
			Pgsql         conf_type.Pgsqls
			Sqlite        conf_type.Sqlites
			Redis         conf_type.Redises
			Etcd          conf_type.Etcd
			Kafka         conf_type.Kafka
			Elasticsearch conf_type.Elasticsearches
			Temporal      conf_type.Temporal
		}
	}
}

const (
	ConfigDir = "config"
)

var (
	Conf *Config
)

// InitConfig InitConfig
func InitConfig(filepath string) {
	f, e := os.Open(filepath)
	defer f.Close()
	util.ChkErr(e)
	Conf = &Config{}
	d := yaml.NewDecoder(f)

	if env.Env.LogLevel < env.LevelWarn {
		d.SetStrict(true)
	}

	util.ChkErr(d.Decode(Conf))

	for k, v := range Conf.App.Middleware.Mysql {
		// MYSQL_DEFAULT_DSN
		dsn := os.Getenv(fmt.Sprintf("YS_MYSQL_%v_DSN", strings.ToUpper(k)))
		if dsn != "" {
			Conf.App.Middleware.Mysql[k] = conf_type.Mysql{Dsn: dsn, Startup: v.Startup, Encrypt: v.Encrypt, HiddenSQL: v.HiddenSQL}
		}

		// Decrypt mysql username and password if "Mysql.Encrypt" is true
		if v.Encrypt {
			decryptDsn, err := decryptMysqlDsn(Conf.App.Middleware.Mysql[k].Dsn, path.Base(filepath))
			util.ChkErr(err)
			Conf.App.Middleware.Mysql[k] = conf_type.Mysql{Dsn: decryptDsn, Startup: v.Startup}
		}
	}
	for k := range Conf.App.Middleware.GRPC.Client {
		addr := os.Getenv(fmt.Sprintf("YS_GRPC_%v_ADDR", strings.ToUpper(k)))
		if addr != "" {
			Conf.App.Middleware.GRPC.Client[k].Addr = addr
		}
	}
	for k := range Conf.App.Middleware.Elasticsearch {
		esusername := os.Getenv(fmt.Sprintf("YS_ES_%v_USERNAME", strings.ToUpper(k)))
		espassword := os.Getenv(fmt.Sprintf("YS_ES_%v_PASSWORD", strings.ToUpper(k)))
		if esusername != "" {
			Conf.App.Middleware.Elasticsearch[k] = conf_type.Elasticsearch{Username: esusername, Password: espassword}
		}
	}
	for k, v := range Conf.App.Middleware.Pgsql {
		// PGSQL_DEFAULT_DSN
		url := os.Getenv(fmt.Sprintf("YS_PGSQL_%v_URL", strings.ToUpper(k)))
		if url != "" {
			Conf.App.Middleware.Pgsql[k] = conf_type.Pgsql{URL: url, Startup: v.Startup}
		}
	}

	storePerfix := os.Getenv("YS_STORAGE_PREFIX")
	if storePerfix != "" {
		for k, v := range Conf.App.Middleware.Redis {
			v.Addr = storePerfix + v.Addr
			Conf.App.Middleware.Redis[k] = v
		}
		for k, v := range Conf.App.Middleware.Etcd.Endpoints {
			Conf.App.Middleware.Etcd.Endpoints[k] = storePerfix + v
		}
		for k, v := range Conf.App.Middleware.Kafka.KafkaClusterURL {
			Conf.App.Middleware.Kafka.KafkaClusterURL[k] = storePerfix + v
		}
	}
}

// ReadConfig ReadConfig for custom
func ReadConfig(configName string, configPath string, configType string) error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// Decrypt mysql username and password
func decryptMysqlDsn(dsn, configName string) (string, error) {
	atSep, colonSep := "@", ":"
	var newDsnBuilder strings.Builder
	dbAccount := strings.Split(dsn, atSep)[0]
	dsnSuffix := strings.Replace(dsn, dbAccount, "", 1)
	errorMsg := "mysql username %v password is empty in %v config file, the invalid dsn is [%v]"

	if dbAccount == "" {
		return "", fmt.Errorf(errorMsg, "and", configName, dsn)
	}

	dbAccounts := strings.Split(dbAccount, colonSep)
	if len(dbAccounts) != 2 {
		return "", fmt.Errorf(errorMsg, "or", configName, dsn)
	}

	// Decrypt username and password
	for i, v := range dbAccounts {
		cipherText := strings.TrimSpace(v)
		if len(cipherText) == 0 {
			return "", fmt.Errorf(errorMsg, "or", configName, dsn)
		}

		plainText, err := util.Decrypt(v)
		if err != nil {
			return "", fmt.Errorf(errorMsg, "or", configName, dsn)
		}
		newDsnBuilder.WriteString(plainText)
		if i == 0 {
			newDsnBuilder.WriteString(colonSep)
		}
	}

	newDsnBuilder.WriteString(dsnSuffix)
	return newDsnBuilder.String(), nil
}
