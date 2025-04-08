package impl

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/common/model"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/service/client"
)

type Secret string

type MatchRegexps map[string]Regexp

type Regexp struct {
	*regexp.Regexp
	original string
}

type AlertManagerConfig struct {
	Global    *GlobalConfig `yaml:"global,omitempty" json:"global,omitempty"`
	Route     *Route        `yaml:"route,omitempty" json:"route,omitempty"`
	Receivers []*Receiver   `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates []string      `yaml:"templates" json:"templates"`
}

type HostPort struct {
	Host string
	Port string
}

type GlobalConfig struct {
	//ResolveTimeout   string `yaml:"resolve_timeout" json:"resolve_timeout"`
	SMTPSmarthost    string `yaml:"smtp_smarthost,omitempty" json:"smtp_smarthost,omitempty"`
	SMTPFrom         string `yaml:"smtp_from,omitempty" json:"smtp_from,omitempty"`
	SMTPAuthUsername string `yaml:"smtp_auth_username,omitempty" json:"smtp_auth_username,omitempty"`
	SMTPAuthPassword Secret `yaml:"smtp_auth_password,omitempty" json:"smtp_auth_password,omitempty"`
	WechatApiUrl     string `yaml:"wechat_api_url,omitempty" json:"wechat_api_url,omitempty"`
	SMTPRequireTLS   bool   `yaml:"smtp_require_tls" json:"smtp_require_tls"`
}

type Route struct {
	Receiver string `yaml:"receiver,omitempty" json:"receiver,omitempty"`

	GroupByStr []string          `yaml:"group_by,omitempty" json:"group_by,omitempty"`
	GroupBy    []model.LabelName `yaml:"-" json:"-"`
	GroupByAll bool              `yaml:"-" json:"-"`

	Match    map[string]string `yaml:"match,omitempty" json:"match,omitempty"`
	MatchRE  MatchRegexps      `yaml:"match_re,omitempty" json:"match_re,omitempty"`
	Continue bool              `yaml:"continue,omitempty" json:"continue,omitempty"`
	Routes   []*Route          `yaml:"routes,omitempty" json:"routes,omitempty"`

	GroupWait      string `yaml:"group_wait,omitempty" json:"group_wait,omitempty"`
	GroupInterval  string `yaml:"group_interval,omitempty" json:"group_interval,omitempty"`
	RepeatInterval string `yaml:"repeat_interval,omitempty" json:"repeat_interval,omitempty"`
}

type Receiver struct {
	Name string `yaml:"name" json:"name"`

	EmailConfigs []*AlertEmailConfig `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
}

type AlertEmailConfig struct {
	To      string            `yaml:"to,omitempty" json:"to,omitempty"`
	From    string            `yaml:"from,omitempty" json:"from,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	HTML    string            `yaml:"html,omitempty" json:"html,omitempty"`
}

type RuleConfig struct {
	Groups []*GroupsConfig `yaml:"groups,omitempty" json:"groups,omitempty"`
}

type GroupsConfig struct {
	Name  string  `yaml:"name,omitempty" json:"name,omitempty"`
	Rules []*Rule `yaml:"rules,omitempty" json:"rules,omitempty"`
}

type Rule struct {
	Alert       string            `yaml:"alert,omitempty" json:"alert,omitempty"`
	Expr        string            `yaml:"expr,omitempty" json:"expr,omitempty"`
	For         string            `yaml:"for,omitempty" json:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

func updateRuleConfig(emailData *dto.Notification) {
	// 更新 alertmanager 配置文件
	alertManagerMap := make(map[string]string)
	alertManagerMap[consts.NodeBreakdown] = boolToString(emailData.NodeBreakdown)
	alertManagerMap[consts.AgentBreakdown] = boolToString(emailData.AgentBreakdown)
	alertManagerMap[consts.DiskUsage] = boolToString(emailData.DiskUsage)
	alertManagerMap[consts.JobFailNum] = boolToString(emailData.JobFailNum)
	err := updateAlertManagerFile(alertManagerMap)
	if err != nil {
		logging.GetLogger(context.Background()).Errorf("Failed to update alertmanager.yml. %v", err.Error())
	}

	// 重新加载 alertmanager 配置
	client.AlertManagerReload()
}

func updateGlobalConfig(emailData *dto.EmailConfig) {

	// 更新 alertmanager 配置文件
	alertManagerMap := make(map[string]string)
	alertManagerMap[consts.KeyFrom] = emailData.From
	alertManagerMap[consts.KeyHost] = fmt.Sprintf("%v:%v", emailData.Host, emailData.Port)
	alertManagerMap[consts.KeyUsername] = emailData.UserName
	alertManagerMap[consts.KeyPassword] = emailData.Password
	alertManagerMap[consts.KeyUseTLS] = boolToString(emailData.UseTLS)
	alertManagerMap[consts.KeyAdminAddr] = emailData.AdminAddr
	err := updateAlertManagerFile(alertManagerMap)
	if err != nil {
		logging.GetLogger(context.Background()).Errorf("Failed to update alertmanager.yml. %v", err.Error())
	}

	// 重新加载 alertmanager 配置
	client.AlertManagerReload()
}

func updateAlertManagerFile(alertManagerMap map[string]string) error {
	ctx := context.Background()

	// 1.读取 alertmanager.yml
	newViper := getNewViper()
	alertManagerConfig := AlertManagerConfig{}
	md := mapstructure.Metadata{}
	err := newViper.Unmarshal(&alertManagerConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		logging.GetLogger(ctx).Errorf("Failed to unmarshal alertmanager.yml. %v", err.Error())
		return err
	}

	// 2.更新告警开关
	if alertManagerMap[consts.NodeBreakdown] != "" && alertManagerMap[consts.AgentBreakdown] != "" && alertManagerMap[consts.DiskUsage] != "" && alertManagerMap[consts.JobFailNum] != "" {
		for _, route := range alertManagerConfig.Route.Routes {
			route.Match[consts.Enable] = alertManagerMap[route.Receiver]
		}
		newViper.Set("route.routes", alertManagerConfig.Route.Routes)
	}

	// 3.设置邮件全局变量
	if alertManagerMap[consts.KeyFrom] != "" && alertManagerMap[consts.KeyHost] != "" {
		alertManagerConfig.Global.SMTPFrom = alertManagerMap[consts.KeyFrom]
		alertManagerConfig.Global.SMTPSmarthost = alertManagerMap[consts.KeyHost]
		alertManagerConfig.Global.SMTPAuthUsername = alertManagerMap[consts.KeyUsername]
		alertManagerConfig.Global.SMTPAuthPassword = Secret(alertManagerMap[consts.KeyPassword])
		alertManagerConfig.Global.SMTPRequireTLS = stringToBool(alertManagerMap[consts.KeyUseTLS])
		newViper.Set("global", alertManagerConfig.Global)
	}

	// 4.设置接受者邮箱
	if alertManagerMap[consts.KeyAdminAddr] != "" {
		for _, receiver := range alertManagerConfig.Receivers {
			if receiver.Name == consts.KeyCommon {
				continue
			}
			for _, emailConfig := range receiver.EmailConfigs {
				emailConfig.To = fmt.Sprintf("%v ", alertManagerMap[consts.KeyAdminAddr])
			}
		}
		newViper.Set("receivers", alertManagerConfig.Receivers)
	}

	// 5.写入配置文件
	err = newViper.WriteConfig()
	if err != nil {
		logging.GetLogger(ctx).Errorf("Failed to write mail config into alertmanager.yml. %v", err.Error())
		return err
	}
	return nil
}
