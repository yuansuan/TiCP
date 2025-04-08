package impl

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"google.golang.org/grpc/status"
)

var emailLock sync.RWMutex

type emailServiceImpl struct{}

func NewEmailService() (service.EmailService, error) {
	return &emailServiceImpl{}, nil
}

func (e *emailServiceImpl) GetEmail(ctx context.Context) (*dto.EmailConfig, error) {
	emailViper, err := getNewViper()
	if err != nil {
		return nil, err
	}

	systemConfig := dto.SystemConfig{}

	md := mapstructure.Metadata{}
	err = emailViper.Unmarshal(&systemConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = common.Yaml
		config.Metadata = &md
	})
	if err != nil {
		return nil, err
	}

	return systemConfig.Email, nil
}

func (e *emailServiceImpl) SetEmail(ctx context.Context, email *dto.EmailConfig) error {
	emailLock.Lock()
	defer emailLock.Unlock()

	emailViper, err := getNewViper()
	if err != nil {
		return err
	}

	systemConfig := dto.SystemConfig{}

	md := mapstructure.Metadata{}
	err = emailViper.Unmarshal(&systemConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = common.Yaml
		config.Metadata = &md
	})
	if err != nil {
		return err
	}

	emailViper.Set("email.enable", email.Enable)
	emailViper.Set("email.setting", email.Setting)
	err = emailViper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (e *emailServiceImpl) SendEmail(ctx context.Context, receiver, emailTemplate, jsonData string) error {
	emailViper, err := getNewViper()
	if err != nil {
		return err
	}

	systemConfig := dto.SystemConfig{}

	md := mapstructure.Metadata{}
	err = emailViper.Unmarshal(&systemConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = common.Yaml
		config.Metadata = &md
	})
	if err != nil {
		return err
	}

	if !systemConfig.Email.Enable {
		return status.Errorf(errcode.ErrNoticeSendEmailNotEnable, "not enable send email function")
	}

	template := systemConfig.Email.Template
	if len(template) == 0 {
		return status.Errorf(errcode.ErrNoticeNotSettingEmailTemplate, "not setting email template")
	}

	var matchedTemplate *dto.Template
	for _, v := range template {
		if v.Type == emailTemplate {
			matchedTemplate = v
			break
		}
	}

	if matchedTemplate == nil {
		return status.Errorf(errcode.ErrNoticeEmailTemplateTypeNotMatch, "email template type not match")
	}

	setting := systemConfig.Email.Setting

	mail := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	mail.SetHeader("From", setting.SendEmail)
	mail.SetHeader("To", receiver)

	emailSubject, err := resolveEmailParam(matchedTemplate.Subject, jsonData)
	if err != nil {
		return err
	}
	mail.SetHeader("Subject", emailSubject)

	emailBody, err := resolveEmailParam(matchedTemplate.Content, jsonData)
	if err != nil {
		return err
	}
	mail.SetBody("text/html", emailBody)

	dial := gomail.NewDialer(setting.Host, setting.Port, setting.SendEmail, setting.Password)
	dial.TLSConfig = &tls.Config{InsecureSkipVerify: setting.TLS}
	if err := dial.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}

func resolveEmailParam(templateStr, jsonData string) (string, error) {
	if jsonData == "" {
		return templateStr, nil
	}

	var jsonDataMap map[string]any
	err := json.Unmarshal([]byte(jsonData), &jsonDataMap)
	if err != nil {
		return "", err
	}

	t, err := template.New("resolve-email-param").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var bodyBuffer bytes.Buffer
	err = t.Execute(&bodyBuffer, jsonDataMap)
	if err != nil {
		return "", err
	}

	return bodyBuffer.String(), nil
}

func getNewViper() (*viper.Viper, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	newViper := viper.New()
	newViper.SetConfigType(common.Yaml)
	newViper.SetConfigName(consts.SysConfigName)
	configPath := filepath.Join(pwd, config.ConfigDir)
	newViper.AddConfigPath(configPath)

	err = newViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return newViper, nil
}
