package service

import (
	"fmt"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"math/rand"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
)

const (
	Subject = "邮件验证码"
)

var randEmail *rand.Rand

// EmailService EmailService
type EmailService struct {
	userDao  *dao.UserDao
	host     string
	user     string
	password string
}

// NewEmailSrv NewEmailSrv
func NewEmailSrv(host string, user string, password string) *EmailService {
	return &EmailService{
		host:     host,
		user:     user,
		password: password,
		userDao:  dao.NewUserDao(),
	}
}

func (s *EmailService) concatMessage(content string, subject string, from string, to string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-version"] = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	return []byte(message)
}

// Signup Signup
func (s *EmailService) Signup(ctx context.Context, userID int64, email string, pwd string, activateURL string) error {
	pwdHash, err := util.PwdGenerate(pwd)
	if err != nil {
		return err
	}

	err = s.userDao.Add(ctx, &models.SsoUser{Ysid: userID, Email: email, PwdHash: pwdHash, IsActivated: false})
	if err != nil {
		return err
	}

	token, err := util.JWTGenerate(common.JWTSignup, strconv.FormatInt(userID, 10), time.Now().Add(time.Hour*48))
	if err != nil {
		return err
	}

	template := "<html><a href=http://%v?token=%v>activate</a></html>"

	err = s.Send(ctx, fmt.Sprintf(template, strings.TrimSuffix(activateURL, "/"), url.QueryEscape(token)), "activate account", email)
	if err != nil {
		return err
	}

	return nil
}

// Activate Activate
func (s *EmailService) Activate(ctx context.Context, token string, userID int64) error {
	t, err := url.QueryUnescape(token)
	if err != nil {
		return status.Error(consts.ErrHydraLcpJWTInvalid, err.Error())
	}
	userInToken, err := util.JWTGetSubject(common.JWTSignup, t)
	if err != nil {
		return err
	}

	if userInToken != strconv.FormatInt(userID, 10) {
		return status.Errorf(consts.ErrHydraLcpUserNotMatch, "userInToken in token is %v, while passed in userInToken is %v", userInToken, userID)
	}

	return nil
}

// Send Send
func (s *EmailService) Send(ctx context.Context, content string, subject string, recipient string) error {
	auth := LoginAuth(s.user, s.password)
	err := smtp.SendMail(s.host, auth, s.user, []string{recipient}, s.concatMessage(content, subject, s.user, recipient))
	if err != nil {
		return status.Error(consts.ErrHydraLcpSMTPError, err.Error())
	}
	return nil
}

// SendVerificationCode SendVerificationCode
func (s *EmailService) SendVerificationCode(ctx context.Context, email string) error {
	logger := logging.GetLogger(ctx)
	logger.Infof("[send code] send code to %v", email)
	cache := boot.MW.DefaultCache()

	var lastSendTime string
	// check whether get the captcha already within a minute
	_, ok := cache.Get(common.RedisPrefixEmailToLastSendTime, email, &lastSendTime)
	if ok {
		t, err := time.Parse(time.RFC3339, lastSendTime)
		if err != nil {
			logger.Warnf("[send code exception] parse last send time %v failed", lastSendTime)
			return status.Errorf(consts.ErrHydraLcpFailedToParseTime, err.Error())
		}
		if !time.Now().After(t.Add(time.Minute)) {
			logger.Warnf("[send code exception] code has been send to %v email within a minute", email)
			return status.Errorf(consts.ErrHydraLcpWait, "code has been send to email %v within a minute", email)
		}
	}

	// create a new captcha
	captcha := fmt.Sprintf("%06v", randEmail.Int31n(1000000))
	var err error
	// set redis key
	err = s.setVerificationCodeID(ctx, email, captcha)
	if err != nil {
		return err
	}

	//取短信模板
	tplConfig, exist := config.Custom.SmsConfig.TemplateMap[common.SMSTemplateEmailCode]
	if !exist {
		return status.Errorf(consts.ErrHydraLcpSendSmsTplNotFound, "sms template not found")
	}
	content := fmt.Sprintf(tplConfig.MonYunTpl, captcha)

	// send message
	err = s.Send(ctx, content, Subject, email)
	if err != nil {
		return err
	}
	logger.Infof("[send code] code sent to %v successful", email)
	return nil

}
func init() {
	randEmail = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (s *EmailService) setVerificationCodeID(ctx context.Context, email string, id string) error {
	cache := boot.MW.DefaultCache()
	logger := logging.GetLogger(ctx)
	err := cache.Put(common.RedisPrefixEmailToCodeID, email, id)

	if err != nil {
		logger.Warnf("[set code exception] save captcha to email failed, err is ", err.Error())
		return status.Error(consts.ErrHydraLcpRedisFailed, err.Error())
	}

	err = cache.Put(common.RedisPrefixEmailToLastSendTime, email, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.Warnf("[set code exception] save captcha generated time to email failed, err is ", err.Error())
		return status.Error(consts.ErrHydraLcpRedisFailed, err.Error())
	}

	return nil
}

// VerifyCode VerifyCode
func (s *EmailService) VerifyCode(ctx context.Context, email, code string) error {
	if debugCode := config.Custom.SmsConfig.DebugCode; debugCode != "" && debugCode == code {
		return nil
	}

	logger := logging.GetLogger(ctx)
	logger.Infof("[verify code] verify captcha for email %v", email)
	cache := boot.MW.DefaultCache()
	var captcha string
	_, ok := cache.Get(common.RedisPrefixEmailToCodeID, email, &captcha)
	if !ok {
		logger.Infof("[verify code exception] captcha for email %v is not set", email)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "captcha for email %v is not set", email)
	}

	if code != captcha {
		logger.Infof("[verify code exception] verify captcha for email %v fail", email)
		return status.Errorf(consts.ErrHydraLcpCaptchaVerifyFailed, "verify captcha for email %v fail", email)
	}
	// clear email code cache after validate success
	cache.Delete(common.RedisPrefixEmailToCodeID, email)
	cache.Delete(common.RedisPrefixEmailToLastSendTime, email)
	logger.Infof("code cache for email %v cleared", email)
	return nil
}

type loginAuth struct {
	username, password string
}

// LoginAuth LoginAuth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start Start
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next Next
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, status.Error(consts.ErrHydraLcpSMTPError, "Unkown fromServer")
		}
	}
	return nil, nil
}
