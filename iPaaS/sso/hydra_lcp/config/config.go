package config

import (
	"context"
	"net/url"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/ory/x/urlx"
	"github.com/spf13/viper"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// CustomT CustomT
type CustomT struct {
	HydraAdminURL        string `yaml:"hydra_admin_url"`
	HydraPortalAdminURL  string `yaml:"hydra_portal_admin_url"`
	HydraAdminConf       *url.URL
	HydraPortalAdminConf *url.URL
	FrontendServer       string          `yaml:"frontend_server"`
	FrontendPortalServer string          `yaml:"frontend_portal_server"`
	SMTP                 SMTP            `yaml:"smtp"`
	Wechat               Wechat          `yaml:"wechat"`
	Offiaccount          Offiaccount     `yaml:"offiaccount"`
	SmsConfig            SmsConfig       `yaml:"sms_config"`
	CasiOauth            CasiOauthConfig `yaml:"casi_oauth"`
}

// CasiOauthConfig 航天云网Oauth
type CasiOauthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

// SmsConfig sms config
type SmsConfig struct {
	SmsType     string                 `yaml:"sms_type"`
	Aliyun      AliyunSmsConfig        `yaml:"aliyun"`
	Tencent     TencentSmsConfig       `yaml:"tencent"`
	DebugCode   string                 `yaml:"debug_code"`
	Templates   []SmsTemplate          `yaml:"templates"`
	TemplateMap map[string]SmsTemplate `yaml:"-"`
	DomainMap   map[string]SmsTemplate `yaml:"-"`
}

type AliyunSmsConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	SignName        string `yaml:"sign_name"`
}

type TencentSmsConfig struct {
	EndpointTencent    string `yaml:"endpoint_tencent"`
	AccessKeyID        string `yaml:"access_key_id"`
	AccessKeySecret    string `yaml:"access_key_secret"`
	SmsSdkAppID        string `yaml:"sms_sdk_app_id"`
	EndpointSmsTencent string `yaml:"endpoint_sms_tencent"`
	RegionSmsTencent   string `yaml:"region_sms_tencent"`
	MonyunEndpoint     string `yaml:"monyun_endpoint"`
	MonyunAccessYS     string `yaml:"monyun_access_ys"`
	MonyunAccessZS     string `yaml:"monyun_access_zs"`
}

// SmsTemplate sms template
type SmsTemplate struct {
	Key        string `yaml:"key"`
	MonYunTpl  string `yaml:"monyun_tpl"`
	TencentTID string `yaml:"tencent_tid"`
	AliyunTID  string `yaml:"aliyun_tid"`
	ParamCount int    `yaml:"param_count"`
	Domain     string `yaml:"domain"`
}

// SMTP SMTP
type SMTP struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// Wechat Wechat
type Wechat struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

// Offiaccount Wechat Official Account Platform
type Offiaccount struct {
	AppID          string      `yaml:"app_id"`
	AppSecret      string      `yaml:"app_secret"`
	CfgToken       string      `yaml:"cfg_token"`
	ExpireSeconds  int64       `yaml:"expire_seconds"`
	EncodingAESKey string      `yaml:"encoding_aes_key"`
	MsgTemplate    MsgTemplate `yaml:"msg_template"`
}

// MsgTemplate 公众号消息模板
type MsgTemplate struct {
	Job     string `yaml:"job"`
	Balance string `yaml:"balance"`
	Topup   string `yaml:"topup"`
	VisJob  string `yaml:"visjob"`
}

// Ldap Ldap
type Ldap struct {
	Dsn     string `yaml:"dsn"`
	Startup bool   `yaml:"startup"`
}

// switch storage
var switches map[string]bool

// Mut Mut
var Mut = &sync.Mutex{}

// Custom Custom
var Custom CustomT

func (c *CustomT) preload() {
	Custom.HydraAdminConf = urlx.ParseOrPanic(Custom.HydraAdminURL)
	Custom.HydraPortalAdminConf = urlx.ParseOrPanic(Custom.HydraPortalAdminURL)

	Custom.SmsConfig.TemplateMap = map[string]SmsTemplate{}
	Custom.SmsConfig.DomainMap = map[string]SmsTemplate{}
	for _, t := range Custom.SmsConfig.Templates {
		Custom.SmsConfig.TemplateMap[t.Key] = t
		if t.Domain != "" {
			Custom.SmsConfig.DomainMap[t.Domain] = t
		}
	}

}

// CheckSwitch CheckSwitch
func CheckSwitch(key string) bool {
	Mut.Lock()
	defer Mut.Unlock()
	if v, ok := switches[key]; v && ok {
		return true
	}
	return false
}

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			Custom.preload()
			logger.Info(Custom)
		}

		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true, boot.Config.App.Name+"/"+env.ModeName(boot.Env.Mode)+"_custom.yaml", &CustomT{}, switchFunc)
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
	Custom.preload()

	logger.Infof(">>>>>>>>>>>>>>>>>> config %#v", Custom)
	return nil
}
