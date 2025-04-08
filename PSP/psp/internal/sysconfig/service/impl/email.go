package impl

import (
	"context"
	"crypto/tls"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"gopkg.in/gomail.v2"
	"strings"
)

// GetEmailConfig GetEmailConfig
func (s *SysConfigServiceImpl) GetEmailConfig(ctx context.Context) (*dto.GetEmailConfigRes, error) {

	alertNotifications, err := s.alertManagerConfigDao.Get(ctx, consts.AlertManagerType)
	if err != nil {
		return nil, err
	}
	alertMap := make(map[string]string)
	for _, notification := range alertNotifications {
		alertMap[notification.Key] = notification.Value
	}

	return &dto.GetEmailConfigRes{
		Notification: dto.Notification{
			NodeBreakdown:  stringToBool(alertMap[consts.KeyNodeBreakdown]),
			DiskUsage:      stringToBool(alertMap[consts.KeyDiskUsage]),
			AgentBreakdown: stringToBool(alertMap[consts.KeyAgentBreakdown]),
			JobFailNum:     stringToBool(alertMap[consts.KeyJobFailNum]),
		},
	}, nil
}

// SetEmailConfig SetEmailConfig
func (s *SysConfigServiceImpl) SetEmailConfig(ctx context.Context, in *dto.SetEmailConfigReq) error {
	// 更新数据库
	alertNotification := buildAlertNotification(s.sid, in)
	err := s.alertManagerConfigDao.Set(ctx, alertNotification)
	if err != nil {
		return err
	}

	// 更新 alertmanager 配置文件
	go updateRuleConfig(&in.Notification)

	return nil
}

// SetGlobalEmail SetGlobalEmail
func (s *SysConfigServiceImpl) SetGlobalEmail(ctx context.Context, in *dto.EmailConfig) error {

	// 更新数据库
	alertNotification := buildGlobalEmail(s.sid, in)
	err := s.alertManagerConfigDao.Set(ctx, alertNotification)
	if err != nil {
		return err
	}

	// 更新 alertmanager 配置文件
	go updateGlobalConfig(in)

	return nil
}

// GetGlobalEmail GetGlobalEmail
func (s *SysConfigServiceImpl) GetGlobalEmail(ctx context.Context) (*dto.EmailConfig, error) {
	alertNotifications, err := s.alertManagerConfigDao.Get(ctx, consts.GlobalEmailType)
	if err != nil {
		return nil, err
	}

	alertMap := make(map[string]string)
	for _, notification := range alertNotifications {
		alertMap[notification.Key] = notification.Value
	}

	return &dto.EmailConfig{
		Host:      alertMap[consts.KeyHost],
		Port:      stringToInt(alertMap[consts.KeyPort]),
		UseTLS:    stringToBool(alertMap[consts.KeyUseTLS]),
		UserName:  alertMap[consts.KeyUsername],
		Password:  alertMap[consts.KeyPassword],
		From:      alertMap[consts.KeyFrom],
		AdminAddr: alertMap[consts.KeyAdminAddr],
	}, nil
}

func (s *SysConfigServiceImpl) SendEmail(ctx context.Context) error {

	alertNotifications, err := s.alertManagerConfigDao.Get(ctx, consts.GlobalEmailType)
	if err != nil {
		return err
	}

	alertMap := make(map[string]string)
	for _, notification := range alertNotifications {
		alertMap[notification.Key] = notification.Value
	}

	mail := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	mail.SetHeader("From", alertMap[consts.KeyFrom])
	mail.SetHeader("To", strings.Split(alertMap[consts.KeyAdminAddr], ",")...)

	mail.SetHeader("Subject", consts.TestEmailSubject)

	mail.SetBody("text/html", consts.TestEmailBody)

	dial := gomail.NewDialer(alertMap[consts.KeyHost], stringToInt(alertMap[consts.KeyPort]), alertMap[consts.KeyUsername], alertMap[consts.KeyPassword])
	if stringToBool(alertMap[consts.KeyUseTLS]) {
		dial.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := dial.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}
